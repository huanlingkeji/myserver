package main

import "fmt"

func main() {
	var a interface{} = nil
	Func(&a)
	fmt.Println(a)
}

func Func(a *interface{}) {
	*a = "dsfsfs"
}
