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

// getToken token列のtokenIdx番目のものを返す君
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
	kind     string // "intliteral", "unary"
	intval   int    // for intliteral
	operator string // "+", "-", ...
	operand  *Expr  // for unary expression
	left     *Expr  // for binary expr
	right    *Expr  // for binary expr
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

// tokenize tokenを解析する君
func tokenize() []*Token {
	var tokens []*Token
	fmt.Printf("#Tokens:")

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
			fmt.Printf(" '%s'", token.value)
		case ';', '+', '-', '*', '/':
			token := &Token{
				kind:  "punct",
				value: string([]byte{char}),
			}
			tokens = append(tokens, token)
			fmt.Printf(" '%s'", token.value)
		case ' ', '\t', 'n':
			continue
		default:
			panic(fmt.Sprintf("tokenize error! Invalid Input: '%c'", char))
		}
	}

	fmt.Printf("\n\n")
	return tokens
}

// parseUnaryExpr 単項式の解析をしてくれる君
func parseUnaryExpr() *Expr {
	token := getToken()
	switch token.kind {
	case "intliteral":
		intval, _ := strconv.Atoi(token.value)
		return &Expr{
			kind:   "intliteral",
			intval: intval,
		}
	case "punct":
		operator := token.value
		operand := parseUnaryExpr()
		return &Expr{
			kind:     "unary",
			operator: operator,
			operand:  operand,
		}
	default:
		panic("Unexpected token kind:" + token.kind)
	}
}

// parse パーサー本体
func parse() *Expr {
	expr := parseUnaryExpr()

	for {
		token := getToken()
		if token == nil || token.value == ";" {
			return expr
		}

		switch token.value {
		case "+", "-", "*", "/":
			return &Expr{
				kind:     "binary",
				operator: token.value,
				left:     expr,
				right:    parseUnaryExpr(),
			}
		default:
			panic("Unexpected token value:" + token.value)
		}
	}
}

// generateExpr アセンブリで計算式部分を吐く君
func generateExpr(expr *Expr) {
	switch expr.kind {
	case "intliteral":
		fmt.Printf("  mov $%d, %%rax\n", expr.intval)
	case "unary":
		switch expr.operator {
		case "-":
			fmt.Printf("  mov $-%d, %%rax\n", expr.operand.intval)
		case "+":
			fmt.Printf("  mov $+%d, %%rax\n", expr.operand.intval)
		default:
			panic("generator error! Unknown unary operator:" + expr.operator)
		}
	case "binary":
		fmt.Printf("  mov $%d, %%rax\n", expr.left.intval)
		fmt.Printf("  mov $%d, %%rcx\n", expr.right.intval)

		switch expr.operator {
		case "+":
			fmt.Printf("  add %%rcx, %%rax\n")
		case "-":
			fmt.Printf("  sub %%rcx, %%rax\n")
		case "*":
			fmt.Printf("  imul %%rcx, %%rax\n")
		case "/":
			fmt.Printf("  mov $0, %%rdx\n")
			fmt.Printf("  idiv %%rcx\n")
		default:
			panic("generator error! Unknown binary operator:" + expr.operator)
		}
	default:
		panic("generator error! Unknown expr.kind:" + expr.kind)
	}
}

// generateAssembly アセンブリコードを吐く君
func generateAssembly(expr *Expr) {
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	generateExpr(expr)
	fmt.Printf("  ret\n")
}

func main() {
	// 標準入力を受け取る
	src, _ = ioutil.ReadAll(os.Stdin)

	tokens = tokenize()
	expr := parse()
	generateAssembly(expr)
}
