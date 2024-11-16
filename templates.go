package hyper

import (
	"log"
	"path"
	"text/template"
)

func (w *responseWriter) HTML(fileName string, data interface{}) {
	defer func() {
		recover()
	}()

	templatePath := path.Join(w.templatesPath, fileName)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println("error in parsing template file ", err)
	}
	err = t.Execute(w, data)
	panic(err)
}
