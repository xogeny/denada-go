package denada

import "io"
import "fmt"
import "log"
import "strings"
import "io/ioutil"

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

func EOF(err error) bool {
	return err == io.EOF
}

func (p *Parser) ParseFile() (ElementList, error) {
	log.Printf(">> File")
	ret := ElementList{}
	for {
		elem, err := p.ParseElement()
		log.Printf("  -> Got element %v", elem)
		if EOF(err) {
			log.Printf("<< File")
			return ret, nil
		}
		if err != nil {
			log.Printf("<< File (Error: %v)", err)
			return nil, err
		} else {
			ret = append(ret, elem)
		}
	}
}

func (p *Parser) ParseContents() (ElementList, error) {
	ret := ElementList{}
	for {
		elem, err := p.ParseElement()
		if EOF(err) {
			return ret, nil
		}
		if err != nil {
			return nil, err
		} else {
			ret = append(ret, elem)
		}
	}
}

type TokenType int

const (
	T_IDENTIFIER TokenType = iota
	T_RBRACE
	T_LBRACE
	T_LPAREN
	T_RPAREN
	T_QUOTE
	T_EQUALS
	T_SEMI
	T_COMMA
	T_WHITE
	T_UNKNOWN
)

func (tt TokenType) String() string {
	switch tt {
	case T_IDENTIFIER:
		return "<identifier>"
	case T_RBRACE:
		return "}"
	case T_LBRACE:
		return "{"
	case T_LPAREN:
		return "("
	case T_RPAREN:
		return ")"
	case T_QUOTE:
		return "\""
	case T_EQUALS:
		return "="
	case T_SEMI:
		return ";"
	case T_COMMA:
		return ","
	case T_WHITE:
		return "<whitespace>"
	case T_UNKNOWN:
		fallthrough
	default:
		return "<???>"
	}
}

type UnexpectedToken struct {
	Found    Token
	Expected string
}

func (u UnexpectedToken) Error() string {
	return fmt.Sprintf("Expecting %s, found '%v' @ (%d, %d)", u.Expected, u.Found.Type,
		u.Found.Line, u.Found.Column)
}

type Token struct {
	Type   TokenType
	String string
	Line   int
	Column int
}

func (t Token) Expected(expected string) UnexpectedToken {
	return UnexpectedToken{
		Found:    t,
		Expected: expected,
	}
}

func (p *Parser) nextNonWhiteToken() (t Token, err error) {
	for {
		t, err = p.nextToken()
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
	if err != nil {
		log.Printf("    Token -> Error: %v", err)
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
		log.Printf("    Token -> %v", t)
		return
	case '}':
		t = Token{Type: T_RBRACE, String: "}", Line: line, Column: col}
		log.Printf("    Token -> %v", t)
		return
	case '(':
		t = Token{Type: T_LPAREN, String: "(", Line: line, Column: col}
		log.Printf("    Token -> %v", t)
		return
	case ')':
		t = Token{Type: T_RPAREN, String: ")", Line: line, Column: col}
		log.Printf("    Token -> %v", t)
		return
	case '"':
		t = Token{Type: T_QUOTE, String: "\"", Line: line, Column: col}
		log.Printf("    Token -> %v", t)
		return
	case '=':
		t = Token{Type: T_EQUALS, String: "=", Line: line, Column: col}
		log.Printf("    Token -> %v", t)
		return
	case ';':
		t = Token{Type: T_SEMI, String: ";", Line: line, Column: col}
		log.Printf("    Token -> %v", t)
		return
	case ',':
		t = Token{Type: T_COMMA, String: ",", Line: line, Column: col}
		log.Printf("    Token -> %v", t)
		return
	}
	if white {
		// If this was white space, keep reading white space
		for {
			ch, _, err = p.src.ReadRune()
			if err == io.EOF {
				t = Token{Type: T_WHITE, String: " ", Line: line, Column: col}
				log.Printf("    Token -> %v", t)
				err = nil
				return
			}
			if err != nil {
				log.Printf("    Token -> Error: %v", err)
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
				log.Printf("    Token -> %v", t)
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
				log.Printf("    Token -> %v", t)
				err = nil
				return
			}
			if err != nil {
				log.Printf("    Token -> Error: %v", err)
				return
			}
			if ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r' || ch == '=' ||
				ch == '{' || ch == '}' || ch == '(' || ch == ')' || ch == ',' || ch == ';' {
				p.src.UnreadRune()
				t = Token{Type: T_IDENTIFIER, String: string(ret), Line: line, Column: col}
				log.Printf("    Token -> %v", t)
				return
			} else {
				ret = append(ret, ch)
			}
		}
	}
}

func (p *Parser) ParseExpr() (Expr, error) {
	panic("Unimplemented")
	return nil, nil
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
func (p *Parser) ParseElement() (ret *Element, err error) {
	ret = &Element{}

	t, err := p.nextNonWhiteToken()
	if err != nil {
		return
	}
	if t.Type == T_RBRACE {
		return nil, nil
	}

	log.Printf("First token: %v", t)
	// Parse all identifiers and white space
	for {
		if t.Type == T_IDENTIFIER {
			ret.Qualifiers = append(ret.Qualifiers, t.String)
		} else if t.Type == T_LPAREN || t.Type == T_QUOTE || t.Type == T_EQUALS ||
			t.Type == T_SEMI || t.Type == T_LBRACE {
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
		log.Printf("Next token: %v", t)
	}

	// Check if there is a modification
	if t.Type == T_LPAREN {
		log.Printf("Parsing modification")
		ret.Modifications, err = p.ParseModifications()
		log.Printf("  err = %v", err)
		if err != nil {
			return
		}
		// Now get next token
		t, err = p.nextNonWhiteToken()
		log.Printf("  err = %v", err)
		if err != nil {
			return
		}
	}

	log.Printf("After modification parsed, next token is = %v (%v)", t, t.Type)

	foundString := false

	if t.Type == T_QUOTE {
		ret.Description, err = p.ParseString()
		log.Printf("    Added description: '%s'", ret.Description)
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
		log.Printf("Definitely a definition")
		// Definition a definition, finish reading and return
		ret.definition = true
		for {
			e, terr := p.ParseElement()
			if terr != nil {
				err = terr
				return
			}
			if e != nil {
				log.Printf("  Nested element: %v", e)
				ret.Contents = append(ret.Contents, e)
			} else {
				err = nil
				return
			}
		}
	}

	if t.Type == T_EQUALS {
		// If we already parsed a string, then we shouldn't find an '=', we should
		// find a ';'
		if foundString {
			err = t.Expected(";")
			return
		}
		expr, terr := p.ParseExpr()
		if terr != nil {
			err = terr
			return
		}
		ret.Value = expr
		t, err = p.nextNonWhiteToken()
		if err != nil {
			return
		}
		if t.Type != T_SEMI {
			err = t.Expected(";")
		}
		return
	}

	if t.Type == T_SEMI {
		ret.definition = false
		return
	}

	// TODO: Need to handle any remaining contingencies...error if we get here?
	err = fmt.Errorf("Unhandled case in Element parser")
	return
}
