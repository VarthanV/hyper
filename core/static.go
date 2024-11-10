package core

import (
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	staticRoute = "^/static/*"
)

func (h *hyper) ConfigureStaticPath(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal("error in configuring static path ", err)
	}
	h.staticPath = absPath
	log.Println("Configured static file path to ", absPath)
	h.GET(staticRoute, h.staticHandler)
}

func (h *hyper) staticHandler(r *Request, w *ResponseWriter) {
	re := regexp.MustCompile(`^/static/(.*)`)
	match := re.FindStringSubmatch(r.Path)
	if len(match) < 1 {
		w.WriteString(200, "Invalid request")
		return
	}

	pathToOpen := path.Join(append([]string{h.staticPath}, match[1:]...)...)
	log.Println("path to open ", pathToOpen)

	b, err := os.ReadFile(pathToOpen)
	if err != nil {
		log.Println("error in opening file ", err)
		w.WriteString(200, "Invalid request")
		return
	}
	w.WriteStatus(200)
	w.WriteHeader("Content-Type", mime.TypeByExtension(filepath.Ext(pathToOpen)))
	w.WriteHeader("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)

}
