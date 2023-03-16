package main

import (
	"fmt"
	"tempconv0/tempconv"
)

func main() {
	fmt.Printf("%g\n", tempconv.BoilingC-tempconv.FreezingC)
	boilingF := tempconv.CToF(tempconv.BoilingC)
	fmt.Printf("%g\n", boilingF-tempconv.CToF(tempconv.FreezingC))
	//fmt.Printf("%g\n", boilingF-FreezingC) //Diferent Types
}
