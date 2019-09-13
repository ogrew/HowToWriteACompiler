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

type Token struct {
	kind  string // intliteral
	value string
}

func main() {
	var src []byte
	// 標準入力を受け取る
	src, _ = ioutil.ReadAll(os.Stdin)
	input := string(src)
	num, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  movq $%d, %%rax\n", num)
	fmt.Printf("  ret\n")
}
