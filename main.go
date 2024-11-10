package main

import (
	"fmt"
	"log"

	"github.com/VarthanV/hyper/core"
)

func handleFoo(req *core.Request, res *core.ResponseWriter) {
	res.WriteString(200, "hello foo")
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

	h.POST("/post", func(req *core.Request, res *core.ResponseWriter) {
		type s struct {
			Foo string `json:"foo"`
		}

		val := s{}
		err := req.Bind(&val)
		if err != nil {
			log.Println("unable to bind req ", err)
			res.WriteJSON(400, map[string]string{"error": "invalid request"})
			return
		}
		log.Printf("%+v\n", val)

		res.WriteJSON(200, val)
	})

	h.GET("/:id", func(req *core.Request, res *core.ResponseWriter) {
		res.WriteStatus(200)
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
