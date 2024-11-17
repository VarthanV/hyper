# hyper

A toy web framework written in Go inspired by [gin-gonic/gin](https://gin-gonic.com/)


# Installing

```sh
go get github.com/VarthanV/hyper
```

# Getting Started

```go
package main

import "github.com/VarthanV/hyper"

func main() {
	h := hyper.New()
	h.GET("/ping", func(w *hyper.ResponseWriter, request *hyper.Request) {
		w.WriteString(200, "PONG")
	})

	h.ListenAndServe("localhost", "6060", `
                               
 _   ___   ______  _____ ____  
| | | \ \ / |  _ \| ____|  _ \ 
| |_| |\ V /| |_) |  _| | |_) |
|  _  | | | |  __/| |___|  _ < 
|_| |_| |_| |_|   |_____|_| \_\
                               
`)
}
```

## Sample Blog application 

[hyper-blog](https://github.com/VarthanV/hyper-blog)

```


