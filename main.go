package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/davidrbourke/pricechecker/webapp/viewmodel"
)

// type myHandler struct {
// 	greeting string
// }

// func (mh myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte(fmt.Sprintf("%v world", mh.greeting)))
// }

// func main() {
// 	http.Handle("/", &myHandler{greeting: "Hello"})
// 	http.Handle("/foo", &myHandler{greeting: "Foo"})
// 	http.ListenAndServe(":8000", nil)
// }

func main() {
	// fmt.Println(http.Dir("public"))
	// http.ListenAndServe(":8000", http.FileServer(http.Dir("e:\\dev\\goworkspace\\public")))
	templates := populateTemplates()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestedFile := r.URL.Path[1:]
		template := templates[requestedFile+".html"]
		context := viewmodel.NewBase()
		if template != nil {
			err := template.Execute(w, context)
			if err != nil {
				log.Println(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.Handle("/img/", http.FileServer(http.Dir("public")))
	http.Handle("/css/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":8000", nil)
}

// func populateTemplates() *template.Template {
// 	result := template.New("templates")
// 	const basePath = "templates"
// 	template.Must(result.ParseGlob(basePath + "/*.html"))
// 	return result
// }

func populateTemplates() map[string]*template.Template {
	result := make(map[string]*template.Template)
	const basePath = "templates"
	layout := template.Must(template.ParseFiles(basePath + "/_layout.html"))
	template.Must(layout.ParseFiles(basePath+"/_header.html", basePath+"/_footer.html"))
	dir, err := os.Open(basePath + "/content")
	if err != nil {
		panic("Failed to open template blocks directory: " + err.Error())
	}
	fis, err := dir.Readdir(-1)
	if err != nil {
		panic("Failed to read contents of content directory: " + err.Error())
	}
	for _, fi := range fis {
		f, err := os.Open(basePath + "/content/" + fi.Name())
		if err != nil {
			panic("Failed to open template: '" + fi.Name() + "'")
		}
		content, err := ioutil.ReadAll(f)
		if err != nil {
			panic("Failed to open content from the file: '" + fi.Name() + "'")
		}
		f.Close()
		tmpl := template.Must(layout.Clone())
		_, err = tmpl.Parse(string(content))
		if err != nil {
			panic("Failed to parse contents of '" + fi.Name() + "' as template")
		}
		result[fi.Name()] = tmpl
	}
	return result
}
