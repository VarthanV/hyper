package main

import (
	"github.com/VarthanV/hyper/core"
)

func handleFoo(req *core.Request, res *core.Response) {
	res.WriteStatus(200)
}

func main() {
	h := core.New()

	h.POST("/foo", handleFoo)
	h.GET("/foo", handleFoo)

	h.ListenAndServe("localhost", "9000", `                               
	_   ___   ______  _____ ____  
	| | | \ \ / |  _ \| ____|  _ \ 
	| |_| |\ V /| |_) |  _| | |_) |
	|  _  | | | |  __/| |___|  _ < 
	|_| |_| |_| |_|   |_____|_| \_\
								   
`)
}
