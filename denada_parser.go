package denada

import "io"
import "fmt"
import "log"
import "bytes"
import "strings"
import "io/ioutil"
import "encoding/json"

import "github.com/bitly/go-simplejson"

type Parser struct {
	src        *strings.Reader
	lineNumber int
	colNumber  int
	log        *log.Logger
}

func NewParser(s io.Reader, l *log.Logger) (p *Parser, err error) {
	str, err := ioutil.ReadAll(s)
	if err != nil {
		return
	}
	src := strings.NewReader(string(str))
	p = &Parser{src: src, lineNumber: 0, colNumber: 0, log: l}
	return
}
func (p *Parser) ParseFile() (ElementList, error) {
	log.Printf(">> File")
	ret := ElementList{}
	for {
		// Try to parse an Element
		elem, err := p.ParseElement(true)

		// If any other error
		if err != nil {
			log.Printf("<< File (Error: %v)", err)
			return nil, err
		}

		// If elem is nil, that means there are no more elements to parse
		if elem == nil {
			log.Printf("<< File")
			return ret, nil
		} else {
			log.Printf("  -> Got element %v", elem)
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

		log.Printf("Token after expression %v is %v", ret.Value, t)

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

func (p *Parser) ParseExpr2() (expr Expr, err error) {
	// Grab whatever is left of the input stream
	left, err := ioutil.ReadAll(p.src)
	log.Printf("left = %s", left)
	if err != nil {
		return
	}

	// Create a JSON decoder on the remaining bytes
	obj := strings.NewReader(string(left))
	jexpr, err := simplejson.NewFromReader(obj)
	if err == nil {
		// Read what is left and reset input stream to that
		left, err = ioutil.ReadAll(obj)
		log.Printf("SUCCESS: %v, left = %s", jexpr, left)

		if err != nil {
			return
		}
		p.src = strings.NewReader(string(left))
	}

	// TODO: Check for float, int, string
	err = fmt.Errorf("Unrecognized expression: %s", left[0:20])

	return
}

func (p *Parser) ParseExpr(modification bool) (expr Expr, err error) {
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
				((white && !empty) || ch == ',' || ch == ')' || ch == ';') {
				p.src.UnreadRune()
				p.lineNumber = l
				p.colNumber = c
				log.Printf("Expr string = %s", w.String())
				decoder := json.NewDecoder(w)
				decoder.UseNumber()
				err = decoder.Decode(&expr)
				if err != nil {
					err = fmt.Errorf("Error parsing expression starting @ (L%d, C%d): %v",
						line+1, col+1, err)
				}
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

// Parse an expression (argument indicates whether we are parsing
// it within a modification or not
func (p *Parser) ParseExpr3(modification bool) (expr Expr, err error) {
	// Grab whatever is left of the input stream
	left, err := ioutil.ReadAll(p.src)
	log.Printf("left = %s", left)
	if err != nil {
		return
	}

	// Create a JSON decoder on the remaining bytes
	obj := strings.NewReader(string(left))
	w := bytes.NewBuffer([]byte{})
	tee := io.TeeReader(obj, w)
	decoder := json.NewDecoder(tee)
	decoder.UseNumber()

	var data interface{}
	// Try to extract a JSON object
	err = decoder.Decode(&data)
	if err == nil {
		log.Printf("SUCCESS: %v", data)
		expr = data
		// Read what is left and reset input stream to that
		left, err = ioutil.ReadAll(tee)
		log.Printf("left = %s", left)
		log.Printf("w = %s", w.String())
		if err != nil {
			return
		}
		p.src = strings.NewReader(string(left))
	}

	// TODO: Check for float, int, string
	err = fmt.Errorf("Unrecognized expression: %s", left[0:20])

	return
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
			log.Printf("--Error while reading tokens: %v", err)
			return
		}
		log.Printf("--Got Token: %v", t)
		if t.Type != T_WHITE {
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

func (p *Parser) nextToken() (t Token, err error) {
	// Record line number and column number at the start of this token
	line := p.lineNumber
	col := p.colNumber

	// Read the first character of the token
	ch, _, err := p.src.ReadRune()
	if err == io.EOF {
		t = Token{Type: T_EOF, String: "", Line: line, Column: col}
		err = nil
		return
	}
	if err != nil {
		return
	}

	// Assume this isn't white space
	white := p.updatePosition(ch)

	switch ch {
	case '{':
		t = Token{Type: T_LBRACE, String: "{", Line: line, Column: col}
		return
	case '}':
		t = Token{Type: T_RBRACE, String: "}", Line: line, Column: col}
		return
	case '(':
		t = Token{Type: T_LPAREN, String: "(", Line: line, Column: col}
		return
	case ')':
		t = Token{Type: T_RPAREN, String: ")", Line: line, Column: col}
		return
	case '"':
		t = Token{Type: T_QUOTE, String: "\"", Line: line, Column: col}
		return
	case '=':
		t = Token{Type: T_EQUALS, String: "=", Line: line, Column: col}
		return
	case ';':
		t = Token{Type: T_SEMI, String: ";", Line: line, Column: col}
		return
	case ',':
		t = Token{Type: T_COMMA, String: ",", Line: line, Column: col}
		return
	}

	if white {
		// If this was white space, keep reading white space
		for {
			ch, _, err = p.src.ReadRune()
			if err == io.EOF {
				t = Token{Type: T_WHITE, String: " ", Line: line, Column: col}
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
				t = Token{Type: T_WHITE, String: " ", Line: line, Column: col}
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
				t = Token{Type: T_IDENTIFIER, String: string(ret), Line: line, Column: col}
				err = nil
				return
			}
			if err != nil {
				return
			}

			l := p.lineNumber
			c := p.colNumber

			if p.updatePosition(ch) || ch == '=' ||
				ch == '{' || ch == '}' || ch == '(' ||
				ch == ')' || ch == ',' || ch == ';' {
				p.src.UnreadRune()
				p.lineNumber = l
				p.colNumber = c
				t = Token{Type: T_IDENTIFIER, String: string(ret), Line: line, Column: col}
				return
			}
			ret = append(ret, ch)
		}
	}
}
