package main

import (
	"fmt"
	"runtime"
)

//global var init func
func init() {

}

//useage in README.txt

func main() {
	runtime.GOMAXPROCS(1)
	var a struct{ name string }
	fmt.Print(a)
}
