package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type M map[string]interface{}

var data = M{}

func list_dir(path string) int {
	files, err := ioutil.ReadDir("./content")

	var list_md []string

	data = M{
		"title":      "Markdown Editor",
		"list_files": list_md,
	}

	if err == nil {
		for _, file := range files {
			list_md = append(list_md, file.Name())
		}
	} else {
		return 0
	}

	return 1
}

func main() {

	// laod static files
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("views/static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("views/index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		is_dir := list_dir("./content")

		if is_dir != 1 {
			list_dir(".")
		}

		err = tmpl.Execute(w, data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	port := "9000"
	fmt.Println("server started at localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
