package httprequest

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Lexer tokenizes HTTP request files according to the specification
type Lexer struct {
	input    string
	position int // current position in input (points to current char)
	line     int // current line number
	column   int // current column number
	start    int // start position of current token
	tokens   []Token
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		line:   1,
		column: 1,
		tokens: make([]Token, 0),
	}
}

// Tokenize processes the input and returns all tokens
func (l *Lexer) Tokenize() ([]Token, error) {
	for l.position < len(l.input) {
		if err := l.nextToken(); err != nil {
			return nil, err
		}
	}

	// Add EOF token
	l.emit(TokenEOF, "")
	return l.tokens, nil
}

// nextToken identifies and emits the next token
func (l *Lexer) nextToken() error {
	l.skipWhitespace()

	if l.position >= len(l.input) {
		return nil
	}

	l.start = l.position
	char := l.current()

	switch {
	case char == '\n':
		l.emit(TokenNewline, "\n")
		l.advance()

	case char == '\r':
		if l.peek() == '\n' {
			l.advance() // skip \r
			l.emit(TokenNewline, "\n")
			l.advance() // skip \n
		} else {
			l.emit(TokenNewline, "\r")
			l.advance()
		}

	case char == '#' && l.peek() == '#' && l.peekN(2) == '#':
		return l.scanRequestSeparator()

	case char == '#':
		return l.scanComment()

	case char == '/' && l.peek() == '/':
		return l.scanComment()

	case char == '<' && l.peek() == '>':
		return l.scanResponseReference()

	case char == '<' && l.peek() == ' ':
		return l.scanFileReference()

	case char == '>' && l.peek() == ' ':
		return l.scanResponseHandler()

	case char == '{' && l.peek() == '{':
		return l.scanVariable()

	case char == ':':
		l.emit(TokenColon, ":")
		l.advance()

	case char == '-' && l.peek() == '-':
		return l.scanBoundary()

	case l.isHTTPMethod():
		return l.scanMethod()

	case l.isHTTPVersion():
		return l.scanHTTPVersion()

	case l.isURL():
		return l.scanURL()

	default:
		return l.scanText()
	}

	return nil
}

// scanComment scans a line comment (# or //)
func (l *Lexer) scanComment() error {
	start := l.position

	// Skip comment start (# or //)
	if l.current() == '#' {
		l.advance()
	} else {
		l.advance() // first /
		l.advance() // second /
	}

	// Skip whitespace after comment marker
	l.skipWhitespace()

	// Read until end of line
	for l.position < len(l.input) && l.current() != '\n' && l.current() != '\r' {
		l.advance()
	}

	comment := strings.TrimSpace(l.input[start:l.position])
	l.emit(TokenComment, comment)
	return nil
}

// scanRequestSeparator scans ### separator
func (l *Lexer) scanRequestSeparator() error {
	l.advance() // first #
	l.advance() // second #
	l.advance() // third #

	// Read optional comment after ###
	l.skipWhitespace()
	start := l.position

	for l.position < len(l.input) && l.current() != '\n' && l.current() != '\r' {
		l.advance()
	}

	separator := "###"
	if l.position > start {
		separator += " " + strings.TrimSpace(l.input[start:l.position])
	}

	l.emit(TokenRequestSeparator, separator)
	return nil
}

// scanResponseReference scans <> response.json
func (l *Lexer) scanResponseReference() error {
	l.advance() // <
	l.advance() // >
	l.skipWhitespace()

	start := l.position
	for l.position < len(l.input) && l.current() != '\n' && l.current() != '\r' {
		l.advance()
	}

	l.emit(TokenResponseRefStart, "<>")
	if l.position > start {
		path := strings.TrimSpace(l.input[start:l.position])
		l.emit(TokenResponseRefPath, path)
	}

	return nil
}

// scanFileReference scans < ./file.json
func (l *Lexer) scanFileReference() error {
	l.advance() // <
	l.skipWhitespace()

	start := l.position
	for l.position < len(l.input) && l.current() != '\n' && l.current() != '\r' {
		l.advance()
	}

	path := strings.TrimSpace(l.input[start:l.position])
	l.emit(TokenFileReference, path)
	return nil
}

// scanResponseHandler scans > {% script %}
func (l *Lexer) scanResponseHandler() error {
	l.advance() // >
	l.skipWhitespace()

	// Check for inline handler {%
	if l.current() == '{' && l.peek() == '%' {
		l.advance() // {
		l.advance() // %
		l.emit(TokenResponseHandlerStart, "> {%")

		// Scan until %}
		start := l.position
		depth := 1

		for l.position < len(l.input) && depth > 0 {
			if l.current() == '%' && l.peek() == '}' {
				depth--
				if depth == 0 {
					break
				}
			} else if l.current() == '{' && l.peek() == '%' {
				depth++
			}
			l.advance()
		}

		if depth > 0 {
			return fmt.Errorf("unclosed response handler at line %d", l.line)
		}

		script := l.input[start:l.position]
		l.emit(TokenResponseHandlerCode, script)

		l.advance() // %
		l.advance() // }
		l.emit(TokenResponseHandlerEnd, "%}")

	} else {
		// File reference handler
		start := l.position
		for l.position < len(l.input) && l.current() != '\n' && l.current() != '\r' {
			l.advance()
		}

		path := strings.TrimSpace(l.input[start:l.position])
		l.emit(TokenResponseRefPath, path)
	}

	return nil
}

// scanVariable scans {{variableName}}
func (l *Lexer) scanVariable() error {
	l.advance() // first {
	l.advance() // second {
	l.emit(TokenVariableStart, "{{")

	l.skipWhitespace()

	start := l.position
	for l.position < len(l.input) {
		char := l.current()
		if char == '}' && l.peek() == '}' {
			break
		}
		if !l.isIdentifierChar(char) && !unicode.IsSpace(rune(char)) {
			return fmt.Errorf("invalid character in variable name at line %d, column %d", l.line, l.column)
		}
		l.advance()
	}

	if l.position >= len(l.input) || l.current() != '}' {
		return fmt.Errorf("unclosed variable at line %d", l.line)
	}

	varName := strings.TrimSpace(l.input[start:l.position])
	if varName == "" {
		return fmt.Errorf("empty variable name at line %d", l.line)
	}

	l.emit(TokenVariableName, varName)

	l.advance() // first }
	l.advance() // second }
	l.emit(TokenVariableEnd, "}}")

	return nil
}

// scanBoundary scans multipart boundary --boundary
func (l *Lexer) scanBoundary() error {
	start := l.position

	// Must start with --
	if l.current() != '-' || l.peek() != '-' {
		return l.scanText()
	}

	l.advance() // first -
	l.advance() // second -

	// Read boundary name
	for l.position < len(l.input) {
		char := l.current()
		if char == '\n' || char == '\r' {
			break
		}
		l.advance()
	}

	boundary := l.input[start:l.position]
	l.emit(TokenBoundary, boundary)
	return nil
}

// scanMethod scans HTTP method
func (l *Lexer) scanMethod() error {
	start := l.position

	for l.position < len(l.input) && l.isMethodChar(l.current()) {
		l.advance()
	}

	method := l.input[start:l.position]
	if ValidHTTPMethods[strings.ToUpper(method)] {
		l.emit(TokenMethod, strings.ToUpper(method))
	} else {
		l.emit(TokenIdentifier, method)
	}

	return nil
}

// scanHTTPVersion scans HTTP version like HTTP/1.1
func (l *Lexer) scanHTTPVersion() error {
	start := l.position

	// HTTP/
	for i := 0; i < 5; i++ {
		l.advance()
	}

	// Version number
	for l.position < len(l.input) && (unicode.IsDigit(rune(l.current())) || l.current() == '.') {
		l.advance()
	}

	version := l.input[start:l.position]
	l.emit(TokenHTTPVersion, version)
	return nil
}

// scanURL scans URL
func (l *Lexer) scanURL() error {
	start := l.position

	// Read until whitespace or newline
	for l.position < len(l.input) {
		char := l.current()
		if unicode.IsSpace(rune(char)) || char == '\n' || char == '\r' {
			break
		}
		l.advance()
	}

	url := l.input[start:l.position]
	l.emit(TokenURL, url)
	return nil
}

// scanText scans general text content
func (l *Lexer) scanText() error {
	start := l.position

	for l.position < len(l.input) {
		char := l.current()

		// Stop at special characters or constructs
		if char == '\n' || char == '\r' ||
			char == '#' ||
			(char == '/' && l.peek() == '/') ||
			(char == '<' && (l.peek() == '>' || l.peek() == ' ')) ||
			(char == '>' && l.peek() == ' ') ||
			(char == '{' && l.peek() == '{') ||
			(char == '-' && l.peek() == '-') {
			break
		}

		l.advance()
	}

	if l.position > start {
		text := l.input[start:l.position]

		// Determine if this looks like a header name (ends with :)
		if strings.HasSuffix(strings.TrimSpace(text), ":") {
			headerName := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(text), ":"))
			l.emit(TokenHeaderName, headerName)
			// Also emit the colon token
			l.emit(TokenColon, ":")
		} else {
			l.emit(TokenText, text)
		}
	}

	return nil
}

// Helper methods

// current returns the current character
func (l *Lexer) current() byte {
	if l.position >= len(l.input) {
		return 0
	}
	return l.input[l.position]
}

// peek returns the next character without advancing
func (l *Lexer) peek() byte {
	if l.position+1 >= len(l.input) {
		return 0
	}
	return l.input[l.position+1]
}

// peekN returns the character at position+n without advancing
func (l *Lexer) peekN(n int) byte {
	if l.position+n >= len(l.input) {
		return 0
	}
	return l.input[l.position+n]
}

// advance moves to the next character
func (l *Lexer) advance() {
	if l.position < len(l.input) {
		if l.input[l.position] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.position++
	}
}

// skipWhitespace skips whitespace characters except newlines
func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) {
		char := l.current()
		if char == ' ' || char == '\t' || char == '\f' {
			l.advance()
		} else {
			break
		}
	}
}

// emit creates and adds a token
func (l *Lexer) emit(tokenType TokenType, value string) {
	token := Token{
		Type:     tokenType,
		Value:    value,
		Line:     l.line,
		Column:   l.column - len(value),
		Position: l.start,
	}
	l.tokens = append(l.tokens, token)
}

// isHTTPMethod checks if current position starts with an HTTP method
func (l *Lexer) isHTTPMethod() bool {
	remaining := l.input[l.position:]

	for method := range ValidHTTPMethods {
		if strings.HasPrefix(strings.ToUpper(remaining), method) {
			// Check that it's followed by whitespace or end of input
			if len(remaining) == len(method) {
				return true
			}
			next := remaining[len(method)]
			if unicode.IsSpace(rune(next)) {
				return true
			}
		}
	}

	return false
}

// isHTTPVersion checks if current position starts with HTTP version
func (l *Lexer) isHTTPVersion() bool {
	remaining := l.input[l.position:]
	httpVersionRegex := regexp.MustCompile(`^HTTP/\d+\.\d+`)
	return httpVersionRegex.MatchString(remaining)
}

// isURL checks if current position starts with a URL
func (l *Lexer) isURL() bool {
	remaining := l.input[l.position:]

	// Simple URL detection
	urlRegex := regexp.MustCompile(`^(https?://|/|\*|{{)`)
	return urlRegex.MatchString(remaining)
}

// isMethodChar checks if character is valid in HTTP method
func (l *Lexer) isMethodChar(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
}

// isIdentifierChar checks if character is valid in identifier
func (l *Lexer) isIdentifierChar(char byte) bool {
	return unicode.IsLetter(rune(char)) || unicode.IsDigit(rune(char)) || char == '-' || char == '_'
}
