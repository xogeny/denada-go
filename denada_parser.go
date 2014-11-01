package denada

import "io"

import "log"
import "fmt"
import "bytes"
import "strings"
import "io/ioutil"
import "encoding/json"

import "github.com/bitly/go-simplejson"

type Parser struct {
	src        *strings.Reader
	lineNumber int
	colNumber  int
	file       string
}

func NewParser(s io.Reader, file string) (p *Parser, err error) {
	str, err := ioutil.ReadAll(s)
	if err != nil {
		return
	}
	src := strings.NewReader(string(str))
	p = &Parser{src: src, lineNumber: 0, colNumber: 0, file: file}
	return
}
func (p *Parser) ParseFile() (ElementList, error) {
	ret := ElementList{}
	for {
		// Try to parse an Element
		elem, err := p.ParseElement(true)

		// If any other error
		if err != nil {
			return nil, err
		}

		// If elem is nil, that means there are no more elements to parse
		if elem == nil {
			return ret, nil
		} else {
			// Add element and continue
			ret = append(ret, elem)
		}
	}
}

// Returns an element if one found, nil on "}"...otherwise an error
func (p *Parser) ParseElement(parsingFile bool) (ret *Element, err error) {
	ret = &Element{}

	// Get first token of the element
	t, err := p.nextNonWhiteToken()
	if err != nil {
		log.Printf("ParseElement is returning with error %v", err)
		return
	}

	// Depending on the context, the element list is terminated by
	// either an EOF or a }.  Check for these...
	if t.Type == T_EOF {
		if parsingFile {
			// Expected, indicate no more elements
			return nil, nil
		} else {
			// Unexpected
			err = t.Expected("definition, declaration or '}'")
			return
		}
	}

	if t.Type == T_RBRACE {
		if parsingFile {
			// Unexpected
			err = t.Expected("definition, declaration or EOF")
			return
		} else {
			// Expected
			return nil, nil
		}
	}

	// Assuming there are more elements, the first thing should always
	// be an identifier
	if t.Type != T_IDENTIFIER {
		err = t.Expected("definition, declaration or EOF")
		return
	}
	ret.Name = t.String

	// Get next token
	t, err = p.nextNonWhiteToken()
	if err != nil {
		return
	}

	// Parse all remaining identifiers and white space
	for {
		if t.Type == T_IDENTIFIER {
			// More qualifiers
			ret.Qualifiers = append(ret.Qualifiers, ret.Name)
			ret.Name = t.String
		} else if t.Type == T_LPAREN || t.Type == T_QUOTE || t.Type == T_EQUALS ||
			t.Type == T_SEMI || t.Type == T_LBRACE {
			// Expected next tokens
			break
		} else {
			err = UnexpectedToken{
				Found:    t,
				Expected: "identifier, whitespace, (, \", =, ; or {",
			}
			return
		}
		t, err = p.nextNonWhiteToken()
		if err != nil {
			return
		}
	}

	// Check if there is a modification (declaration or definition)
	if t.Type == T_LPAREN {
		ret.Modifications, err = p.ParseModifications()
		if err != nil {
			return
		}
		// Now get next token
		t, err = p.nextNonWhiteToken()
		if err != nil {
			return
		}
	}

	foundString := false

	// Read the description, if present
	if t.Type == T_QUOTE {
		ret.Description, err = p.ParseString()
		if err != nil {
			return
		}
		foundString = true

		// Now get next token
		t, err = p.nextNonWhiteToken()
		if err != nil {
			return
		}
	}

	// Is this a definition?
	if t.Type == T_LBRACE {
		// Definitely a definition, finish reading and return
		ret.definition = true
		for {
			e, terr := p.ParseElement(false)
			if terr != nil {
				err = terr
				return
			}
			// This means we are done
			if e == nil {
				err = nil
				return
			}
			ret.Contents = append(ret.Contents, e)
		}
	}

	// At this point, we know we have a declaration
	ret.definition = false

	// Check to see if it has a value
	if t.Type == T_EQUALS {
		// If we already parsed a string, then we shouldn't find an '=', we should
		// find a ';'
		if foundString {
			err = t.Expected(";")
			return
		}

		// Otherwise, parse the expression
		ret.Value, err = p.ParseExpr(false)
		if err != nil {
			return
		}

		// Grab the next token
		t, err = p.nextNonWhiteToken()
		if err != nil {
			return
		}

		// Check to see if there is a description after the value
		if t.Type == T_QUOTE {
			ret.Description, err = p.ParseString()
			if err != nil {
				return
			}

			t, err = p.nextNonWhiteToken()
			if err != nil {
				return
			}
		}
	}

	// This should a SEMI that terminates the declaration
	if t.Type != T_SEMI {
		err = t.Expected(";")
	}
	return
}

func (p *Parser) ParseContents() (ElementList, error) {
	ret := ElementList{}
	for {
		// Parse elements until there aren't any more
		elem, err := p.ParseElement(false)
		if err != nil {
			return nil, err
		} else {
			ret = append(ret, elem)
		}
	}
}

// Silly function to work around the fact that you can't
// build a simplejson.Json object from an existing interface{}
func makeJson(data interface{}) *simplejson.Json {
	tmp := simplejson.New()
	tmp.Set("tmp", data)
	return tmp.Get("tmp")
}

func (p *Parser) ParseExpr(modification bool) (expr *simplejson.Json, err error) {
	line := p.lineNumber
	col := p.colNumber
	// Read input stream keeping track of nesting of {}s and "s.  The
	// next unquoted ',', ')' or ';' outside of quotes and outside an
	// object definition is the end of the JSON string.
	objcount := 0
	arraycount := 0
	quote := false
	escaped := false
	empty := true

	w := bytes.NewBuffer([]byte{})

	for {
		ch, _, terr := p.src.ReadRune()
		err = terr
		if err == io.EOF {
			err = fmt.Errorf("Reached EOF while trying to read expression at (L%d, C%d)",
				p.lineNumber+1, p.colNumber+1)
			return
		}
		if err != nil {
			return
		}

		l := p.lineNumber
		c := p.colNumber
		white := p.updatePosition(ch)
		if !white {
			empty = false
		}

		if quote {
			if escaped {
				escaped = false
			} else {
				if ch == '"' {
					quote = false
				}
				if ch == '\\' {
					escaped = true
				}
			}
		} else {
			if objcount == 0 && arraycount == 0 &&
				((white && !empty) || ch == '/' || ch == ',' || ch == ')' || ch == ';') {
				p.src.UnreadRune()
				p.lineNumber = l
				p.colNumber = c

				// First, I use Go's native json encoding
				var data interface{}
				err = json.Unmarshal(w.Bytes(), &data)
				if err != nil {
					err = fmt.Errorf("Error parsing expression starting @ (L%d, C%d): %v",
						line+1, col+1, err)
				}

				// Now, convert it into a simplejson.Json object for
				// convenient access later.  I don't use the
				// simplejson parsing routines because they turn on
				// "UseNumber" which stores integers as strings.
				// This, in turn, messes up JSON schema
				// representation.
				expr = makeJson(data)
				return
			}
			if ch == '"' {
				quote = true
			}
			if ch == '{' {
				objcount++
			}
			if ch == '[' {
				arraycount++
			}
			if ch == '}' {
				objcount--
			}
			if ch == ']' {
				arraycount--
			}
		}
		w.WriteRune(ch)
	}
}

// This is called if we've already parsed a "(" after a name
func (p *Parser) ParseModifications() (mods Modifications, err error) {
	mods = Modifications{}
	first := true
	for {
		// Identifier
		nt, terr := p.nextNonWhiteToken()
		if terr != nil {
			err = terr
			return
		}
		if first && nt.Type == T_RPAREN {
			return
		}
		if nt.Type != T_IDENTIFIER {
			err = UnexpectedToken{
				Found:    nt,
				Expected: "identifier",
			}
			return
		}

		// =
		et, terr := p.nextNonWhiteToken()
		if terr != nil {
			err = terr
			return
		}
		if et.Type != T_EQUALS {
			err = UnexpectedToken{
				Found:    et,
				Expected: "=",
			}
			return
		}

		expr, terr := p.ParseExpr(true)
		if terr != nil {
			err = terr
			return
		}

		// , or )
		tt, terr := p.nextNonWhiteToken()
		if terr != nil {
			err = terr
			return
		}
		if tt.Type != T_COMMA && tt.Type != T_RPAREN {
			err = UnexpectedToken{
				Found:    tt,
				Expected: ") or ,",
			}
			return
		}

		mods[nt.String] = expr
		first = false
		if tt.Type == T_RPAREN {
			break
		}
	}
	return
}

// Called when we've already read a '"'...read until we get to the closing '"'
func (p *Parser) ParseString() (ret string, err error) {
	runes := []rune{}
	for {
		ch, _, terr := p.src.ReadRune()
		if terr != nil {
			err = terr
			return
		}
		p.updatePosition(ch)

		if ch == '"' {
			break
		}

		runes = append(runes, ch)
	}
	ret = string(runes)
	return
}

func (p *Parser) nextNonWhiteToken() (t Token, err error) {
	for {
		t, err = p.nextToken()
		if err != nil {
			return
		}
		// Ignore white space and comments
		if t.Type != T_WHITE && t.Type != T_SLCOMMENT && t.Type != T_MLCOMMENT {
			return
		}
	}
}

func (p *Parser) updatePosition(ch rune) bool {
	switch ch {
	case '\n':
		p.colNumber = 0
		p.lineNumber++
		return true
	case '\t':
		p.colNumber += 4 // Assume tabs as four spaces
		return true
	case ' ':
		p.colNumber++
		return true
	case '\r':
		p.colNumber = 0
		return true
	default:
		p.colNumber++
	}
	return false
}

func (p *Parser) peek() (r rune, err error) {
	r, _, err = p.src.ReadRune()
	p.src.UnreadRune()
	if err != nil {
		return
	}
	return
}

func (p *Parser) nextToken() (t Token, err error) {
	// Record line number and column number at the start of this token
	line := p.lineNumber
	col := p.colNumber

	// Read the first character of the token
	ch, _, err := p.src.ReadRune()
	if err == io.EOF {
		t = Token{Type: T_EOF, String: "", Line: line, Column: col, File: p.file}
		err = nil
		return
	}
	if err != nil {
		return
	}

	// Assume this isn't white space
	white := p.updatePosition(ch)

	switch ch {
	case '/':
		nch, perr := p.peek()
		// Check if the next character is also a '/'
		if perr == nil && nch == '/' {
			// If so, read until end of line
			comment := []rune{}
			for {
				nch, _, err = p.src.ReadRune()
				if err != nil && err != io.EOF {
					return
				}
				p.updatePosition(nch)
				if err == io.EOF || nch == '\n' {
					t = Token{
						Type:   T_SLCOMMENT,
						String: string(comment),
						Line:   line,
						Column: col,
						File:   p.file}
					err = nil
					return
				}
				comment = append(comment, nch)
			}
		}
		// Check if the next character is also a '*'
		if perr == nil && nch == '*' {
			// If so, read until matching */
			comment := []rune{'/'}
			star := false
			for {
				nch, _, err = p.src.ReadRune()
				if err != nil {
					err = fmt.Errorf("Error while reading multi-line comment: %v", err)
					return
				}
				p.updatePosition(nch)
				comment = append(comment, nch)
				if nch == '/' && star {
					t = Token{
						Type:   T_MLCOMMENT,
						String: string(comment),
						Line:   line,
						Column: col,
						File:   p.file}
					return
				}
				if nch == '*' {
					star = true
				} else {
					star = false
				}
			}
		}
	case '{':
		t = Token{Type: T_LBRACE, String: "{", Line: line, Column: col, File: p.file}
		return
	case '}':
		t = Token{Type: T_RBRACE, String: "}", Line: line, Column: col, File: p.file}
		return
	case '(':
		t = Token{Type: T_LPAREN, String: "(", Line: line, Column: col, File: p.file}
		return
	case ')':
		t = Token{Type: T_RPAREN, String: ")", Line: line, Column: col, File: p.file}
		return
	case '"':
		t = Token{Type: T_QUOTE, String: "\"", Line: line, Column: col, File: p.file}
		return
	case '=':
		t = Token{Type: T_EQUALS, String: "=", Line: line, Column: col, File: p.file}
		return
	case ';':
		t = Token{Type: T_SEMI, String: ";", Line: line, Column: col, File: p.file}
		return
	case ',':
		t = Token{Type: T_COMMA, String: ",", Line: line, Column: col, File: p.file}
		return
	}

	if white {
		// If this was white space, keep reading white space
		for {
			ch, _, err = p.src.ReadRune()
			if err == io.EOF {
				t = Token{Type: T_WHITE, String: " ", Line: line, Column: col, File: p.file}
				err = nil
				return
			}
			if err != nil {
				return
			}

			l := p.lineNumber
			c := p.colNumber
			// If not white space, we are done
			if !p.updatePosition(ch) {
				p.src.UnreadRune()
				p.lineNumber = l
				p.colNumber = c
				t = Token{Type: T_WHITE, String: " ", Line: line, Column: col, File: p.file}
				return
			}
		}
	} else {
		ret := []rune{ch}
		// If not white space, this is the start of an identifier
		// read until we hit a special character
		for {
			ch, _, err = p.src.ReadRune()
			if err == io.EOF {
				t = Token{Type: T_IDENTIFIER, String: string(ret), Line: line, Column: col, File: p.file}
				err = nil
				return
			}
			if err != nil {
				return
			}

			l := p.lineNumber
			c := p.colNumber

			if p.updatePosition(ch) || ch == '=' || ch == '/' ||
				ch == '{' || ch == '}' || ch == '(' ||
				ch == ')' || ch == ',' || ch == ';' {
				p.src.UnreadRune()
				p.lineNumber = l
				p.colNumber = c
				t = Token{Type: T_IDENTIFIER, String: string(ret), Line: line, Column: col, File: p.file}
				return
			}
			ret = append(ret, ch)
		}
	}
}
