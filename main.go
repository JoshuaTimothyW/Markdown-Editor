package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type M map[string]interface{}

var tmpl, err = template.ParseGlob("views/*")

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		files, err := ioutil.ReadDir(".")

		var list_md = []string
		
		if err == nil {
			for _, file := range files {
				append(list_md,file.Name())
				// fmt.Println(file.Name(), file.IsDir())
			}
		}

		// fmt.Println(files)

		err := tmpl.ExecuteTemplate(w, "index.html", list_md)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	port := "9000"
	fmt.Println("server started at localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
