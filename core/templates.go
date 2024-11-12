package core

import (
	"log"
	"path"
	"text/template"
)

func (h *hyper) HTML(w ResponseWriter, fileName string, data interface{}) {
	templatePath := path.Join(h.templatesPath, fileName)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println("error in parsing template file ", err)
	}

	err = t.Execute(w, data)
}
