package main

import (
	"bufio"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

type M map[string]interface{}

//go:embed views/index.html
var indexTemplate embed.FS

//go:embed views
var assets embed.FS

type Files struct {
	Filename string
	Filepath string
}

type Data struct {
	Title      string
	Content    []string
	List_files []Files
}

var data Data

func list_dir(path string) int {

	data.List_files = nil

	var file Files

	err := filepath.Walk("./content",
		func(name string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if len(filepath.Ext(name)) > 0 {
				file.Filename = filepath.Base(name)
				file.Filepath = name
				data.List_files = append(data.List_files, file)
			}

			return nil
		})

	data.Title = "Markdown Editor"

	if err != nil {
		return 0
	}

	return 1
}

// get list of files
func check_dir() {
	is_dir := list_dir("./content")

	if is_dir != 1 {
		list_dir(".")

		if is_dir != 1 {
			println("No views called index.html")
		}
	}
}

func read(path string) {
	file, err := os.Open("./" + path)

	if err != nil {
		fmt.Println("The path doesn't exist")
		return
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		data.Content = append(data.Content, scanner.Text())
	}

	file.Close()
}

func main() {

	// load static files

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.FS(assets))))

	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.RawPath)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(u.Query()["path"][0])

		http.Redirect(w, r, "/", 200)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get list of files
		check_dir()

		// tmpl, err := template.ParseFS(indexTemplate, "views/index.html")
		tmpl, err := template.ParseFiles("views/index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	port := "9000"

	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:", port).Start()
	// err := exec.Command("start", "http://localhost:", port).Start()

	if err != nil {
		fmt.Println("Cannot open browser automaticly")
	}

	urlStr := "http://localhost:9000/edit?path=content%5cposts%5cpost-1.md"

	u, err := url.Parse(urlStr)

	if err != nil {
		return
	}

	fmt.Println(u.Query()["path"][0])

	read(u.Query()["path"][0])

	fmt.Println(data.Content)

	fmt.Println("server started at localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
