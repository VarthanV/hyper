package main

import (
	"fmt"

	"github.com/VarthanV/hyper/core"
)

func main() {
	fmt.Println("hello world")

	h := core.New()

	h.ListenAndServe("localhost", "9000", `                               
	_   ___   ______  _____ ____  
	| | | \ \ / |  _ \| ____|  _ \ 
	| |_| |\ V /| |_) |  _| | |_) |
	|  _  | | | |  __/| |___|  _ < 
	|_| |_| |_| |_|   |_____|_| \_\
								   
`)
}
