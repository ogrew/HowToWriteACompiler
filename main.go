package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

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
