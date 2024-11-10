package main

import (
	"fmt"

	"github.com/VarthanV/hyper/core"
)

func handleFoo(req *core.Request, res *core.ResponseWriter) {
	res.WriteString(200, "hello nigga")
}

func main() {
	h := core.New()

	h.POST("/foo", handleFoo)
	h.GET("/foo", handleFoo)
	h.GET("/ping", func(req *core.Request, res *core.ResponseWriter) {
		res.WriteString(200, "PONG")
	})

	h.GET("/json", func(req *core.Request, res *core.ResponseWriter) {
		res.WriteJSON(200, map[string]string{
			"foo": "bar",
		})
	})

	h.GET("/pong", func(req *core.Request, res *core.ResponseWriter) {
		res.WriteString(200, "PONG")
	})

	h.GET("/q", func(req *core.Request, res *core.ResponseWriter) {
		res.WriteString(200, fmt.Sprintf("got query %s", req.Query("search")))
	})

	h.ConfigureStaticPath("./static")

	h.ListenAndServe("localhost", "9000", `                               
	_   ___   ______  _____ ____  
	| | | \ \ / |  _ \| ____|  _ \ 
	| |_| |\ V /| |_) |  _| | |_) |
	|  _  | | | |  __/| |___|  _ < 
	|_| |_| |_| |_|   |_____|_| \_\
								   
`)
}
