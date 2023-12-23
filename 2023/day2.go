package solution

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/ryanpdenoux/advent-of-code/utils"
)

func Day2(file *os.File) {
	var sum int = 0
	var id int = 1

	scanner := bufio.NewScanner(file)
	current_rules := Rules{12, 13, 14}

	for scanner.Scan() {
		line := scanner.Text()
		game := newGame(id, line)
		game.solveGame(current_rules)
		if game.Valid {
			sum += id
		}
		id++
	}

	fmt.Printf("Sum of game ids: %d\n", sum)
}

type parser interface {
	parse() Rules
}

const (
	ILLEGAL = "ILLEGAL"
	EOL     = "EOL"
	GAME    = "GAME"
	ID      = "ID"
	COLON   = "COLON"
	COMMA   = "COMMA"
	COMMENT = "COMMENT"
	SEMIC   = "SEMIC"
	RED     = "RED"
	GREEN   = "GREEN"
	BLUE    = "BLUE"
	VALUE   = "VALUE"
)

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

var colors = map[string]TokenType{
	"red":   RED,
	"green": GREEN,
	"blue":  BLUE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := colors[ident]; ok {
		return tok
	}
	return GAME
}

func newToken(tokenType TokenType, ch byte) *Token {
	return &Token{Type: tokenType, Literal: string(ch)}
}

type Rules struct {
	Red   int
	Green int
	Blue  int
}

func (r *Rules) Update(other Rules) {
	if other.Red > r.Red {
		r.Red = other.Red
	}
	if other.Blue > r.Blue {
		r.Blue = other.Blue
	}
	if other.Green > r.Green {
		r.Green = other.Green
	}
}

func (r *Rules) Compare(other Rules) bool {
	if other.Red > r.Red {
		return false
	}
	if other.Blue > r.Blue {
		return false
	}
	if other.Green > r.Green {
		return false
	}
	return true
}

type Game struct {
	Id     int
	Valid  bool
	parser *Parser
}

func newGame(id int, gameData string) *Game {
	game := &Game{Id: id, parser: newParser(gameData)}
	return game
}

func (g *Game) solveGame(rules Rules) {
	outcome := g.parser.Parse()
	g.Valid = rules.Compare(outcome)
}

// line based Game parser
type Parser struct {
	lexer *Lexer
	curr  *Token
	peek  *Token
}

func newParser(input string) *Parser {
	p := &Parser{}
	p.lexer = newLexer(input)
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curr = p.peek
	p.peek = p.lexer.NextToken()
	if p.curr != nil {
		slog.Debug("Current token", "token", p.curr)
	}
}

func (p *Parser) Parse() Rules {
	rules := Rules{}

	sets := p.parseSets()
	slog.Debug("Sets of current game", "sets", sets)
	for _, set := range(sets) {
		rules.Update(set)
	}

	return rules
}

func (p *Parser) parseSets() []Rules {
	results := []Rules{}

	for !p.currTokenIs(EOL) {
		set := p.parseSet()
		if set != nil {
			results = append(results, *set)
		}
	}

	return results
}

func (p *Parser) parseSet() *Rules {
	switch p.curr.Type {
	case GAME:
		return p.parseHeader()
	case COLON:
		return p.parseResult()
	case SEMIC:
		return p.parseResult()
	default:
		p.nextToken()
		return nil
	}
}

func (p *Parser) currTokenIs(token TokenType) bool {
	return p.curr.Type == token
}

func (p *Parser) peekTokenIs(token TokenType) bool {
	return p.peek.Type == token
}

func (p *Parser) expectPeek(token TokenType) bool {
	if p.peekTokenIs(token) {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) parseHeader() *Rules {
	p.nextToken()
	return nil
}

func (p *Parser) parseResult() *Rules {
	rules := &Rules{}

	if !p.expectPeek(VALUE) {
		return nil
	}

	for !p.currTokenIs(SEMIC) {

		if p.currTokenIs(EOL) {
			return rules
		}

		num, err := strconv.Atoi(p.curr.Literal)
		if err != nil {
			return nil
		}
		p.nextToken()

		switch p.curr.Type {
		case RED:
			rules.Red = num
		case BLUE:
			rules.Blue = num
		case GREEN:
			rules.Green = num
		}
		p.nextToken()

		if p.currTokenIs(COMMA) {
			p.nextToken()
		}
	}

	return rules
}

// Line based Lexer that codifies game information
type Lexer struct {
	input    string
	position int
	readPos  int
	ch       byte
}

func newLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.position = l.readPos
	l.readPos += 1
}

func (l *Lexer) NextToken() *Token {
	tok := &Token{}

	l.skipWhitespace()

	switch l.ch {
	case ':':
		tok = newToken(COLON, l.ch)
	case ';':
		tok = newToken(SEMIC, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '#':
		tok = newToken(COMMENT, l.ch)
	case 0:
		tok.Type = EOL
		tok.Literal = ""
	default:
		if utils.IsLetter(l.ch) {  // Only returns Game or Colors
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if utils.IsDigit(l.ch) {
			tok.Type = VALUE
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for utils.IsLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for utils.IsDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
