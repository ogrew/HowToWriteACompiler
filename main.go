package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var src []byte
var srcIdx int

// getChar srcのsrcIdx番目のものを返す君
func getChar() (byte, error) {
	if srcIdx == len(src) {
		return 0, errors.New("EOF")
	}
	char := src[srcIdx]
	srcIdx++
	return char, nil
}

func ungetChar() {
	srcIdx--
}

var tokens []*Token
var tokenIdx int

func getToken() *Token {
	if tokenIdx == len(tokens) {
		return nil
	}
	token := tokens[tokenIdx]
	tokenIdx++
	return token
}

// Token is XXX
type Token struct {
	kind  string // "intliteral", "punct"
	value string
}

// Expr is XXX
type Expr struct {
	kind   string // "intliteral"
	intval int    // for intliteral
}

// readNumber 数値列を最後まで読み取る君
func readNumber(char byte) string {
	num := []byte{char}
	for {
		char, err := getChar()
		if err != nil {
			// 最後まで読み取りきった場合は終了
			break
		}
		if '0' <= char && char <= '9' {
			num = append(num, char)
		} else {
			// 数値以外が来た場合は終了
			ungetChar()
			break
		}
	}

	return string(num)
}

func tokenize() []*Token {
	var tokens []*Token
	fmt.Printf("#Tokens : ")

	for {
		char, err := getChar()
		// もう終わり
		if err != nil {
			break
		}

		switch char {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			intLiteral := readNumber(char)
			token := &Token{
				kind:  "intliteral",
				value: intLiteral,
			}
			tokens = append(tokens, token)
			fmt.Printf("' %s ", token.value)
		case ';':
			token := &Token{
				kind:  "punct",
				value: string([]byte{char}),
			}
			tokens = append(tokens, token)
			fmt.Printf("' %s ", token.value)
		case ' ', '\t', 'n':
			continue
		default:
			panic(fmt.Sprintf("tokenize error: Invalid Input: '%c'", char))
		}
	}

	fmt.Printf("\n\n")
	return tokens
}

func parse() *Expr {
	token0 := getToken()

	intval, _ := strconv.Atoi(token0.value)
	expr := &Expr{
		kind:   "intliteral",
		intval: intval,
	}
	return expr
}

func generateAssembly(expr *Expr) {
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  movq $%d, %%rax\n", expr.intval)
	fmt.Printf("  ret\n")
}

func main() {
	// 標準入力を受け取る
	src, _ = ioutil.ReadAll(os.Stdin)

	tokens = tokenize()
	expr := parse()
	generateAssembly(expr)

}
