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

func (p *Parser) ParseContents() (ElementList, error) {
	ret := ElementList{}
	for {
		elem, err := p.ParseElement(false)
		if err != nil {
			return nil, err
		} else {
			ret = append(ret, elem)
		}
	}
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

func (p *Parser) nextToken() (t Token, err error) {
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

	// Increment column number
	p.colNumber++

	// Assume this isn't white space
	white := false

	switch ch {
	case '\n':
		p.colNumber = 0
		p.lineNumber++
		white = true
	case '\t':
		p.colNumber += 3
		white = true
	case ' ':
		white = true
	case '\r':
		p.colNumber = 0
		white = true
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
			switch ch {
			case ' ':
				p.colNumber++
				continue
			case '\n':
				p.colNumber = 0
				p.lineNumber++
				continue
			case '\t':
				p.colNumber += 4
				continue
			case '\r':
				p.colNumber = 0
				continue
			default:
				p.src.UnreadRune()
				t = Token{Type: T_WHITE, String: " ", Line: line, Column: col}
				return
			}
		}
	} else {
		ret := []rune{ch}
		// If not white space, read until we hit a special character
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
			if ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r' || ch == '=' ||
				ch == '{' || ch == '}' || ch == '(' || ch == ')' || ch == ',' || ch == ';' {
				p.src.UnreadRune()
				t = Token{Type: T_IDENTIFIER, String: string(ret), Line: line, Column: col}
				return
			} else {
				ret = append(ret, ch)
			}
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

func (p *Parser) ParseExpr() (expr Expr, err error) {
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

//  This is called if we've already parsed a "("
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

		expr, terr := p.ParseExpr()
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
		p.colNumber++
		if ch == '"' {
			break
		}
		runes = append(runes, ch)
	}
	ret = string(runes)
	return
}

// Returns an element if one found, nil on "}"...otherwise an error
func (p *Parser) ParseElement(parsingFile bool) (ret *Element, err error) {
	ret = &Element{}

	t, err := p.nextNonWhiteToken()
	if err != nil {
		return
	}

	if t.Type == T_EOF {
		if parsingFile {
			return nil, nil
		} else {
			err = t.Expected("definition, declaration or '}'")
			return
		}
	}

	if t.Type == T_RBRACE {
		if parsingFile {
			err = t.Expected("definition, declaration or EOF")
			return
		} else {
			return nil, nil
		}
	}

	// First thing should always be an identifier
	if t.Type != T_IDENTIFIER {
		err = t.Expected("definition, declaration or EOF")
		return
	} else {
		ret.Qualifiers = append(ret.Qualifiers, t.String)
	}

	// Get next token
	t, err = p.nextNonWhiteToken()
	if err != nil {
		return
	}

	// Parse all remaining identifiers and white space
	for {
		if t.Type == T_IDENTIFIER {
			ret.Qualifiers = append(ret.Qualifiers, t.String)
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

	log.Printf("After modification parsed, next token is = %v (%v)", t, t.Type)

	foundString := false

	if t.Type == T_QUOTE {
		ret.Description, err = p.ParseString()
		if err != nil {
			return
		}
		foundString = true

		// Now get next token
		t, err = p.nextNonWhiteToken()
		log.Printf("  err = %v", err)
		if err != nil {
			return
		}
	}

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

	if t.Type == T_EQUALS {
		// If we already parsed a string, then we shouldn't find an '=', we should
		// find a ';'
		if foundString {
			err = t.Expected(";")
			return
		}

		// Otherwise, this is definitely a declaration.  Parse the expression
		// and the semicolon
		ret.Value, err = p.ParseExpr()
		if err != nil {
			return
		}

		t, err = p.nextNonWhiteToken()
		if err != nil {
			return
		}
	}

	// This should a SEMI
	if t.Type != T_SEMI {
		err = t.Expected(";")
	}
	return
}
