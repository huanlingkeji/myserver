package util

import "fmt"

func ShowBytes(data []byte) {
	for _, v := range data {
		fmt.Print(v, " ")
	}
	fmt.Println()
}
