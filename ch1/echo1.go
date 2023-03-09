package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	var s, sep string
	for i := 0; i < len(os.Args); i++ {
		sep = " "
		s += strconv.Itoa(i) + "->" + os.Args[i] + sep
	}
	fmt.Println(s)
}
