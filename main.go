package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
)

type M map[string]interface{}

type Files struct {
	Filename string
	Filepath string
}

type Data struct {
	Title       string
	Content     string
	CurrentPath string
	List_files  []Files
}

var data Data

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

// list all directories
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
			println("No markdown files")
		}
	}
}

// read file by path
func readFile(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	data.Content = string(b)
	data.CurrentPath = path
	data.Title = filepath.Base(path)
}

func writeFile() {
	ioutil.WriteFile(data.CurrentPath, []byte(data.Content), 0)

}

func main() {

	r := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	r.Renderer = renderer

	r.Static("/static", "views/static")

	r.GET("/", func(ctx echo.Context) error {

		check_dir()

		return ctx.Render(http.StatusOK, "index.html", M{})
	})

	r.GET("/edit", func(ctx echo.Context) error {

		path := ctx.QueryParam("path")

		// list all file directory
		check_dir()

		if len(path) > 0 {
			// read file to fetch content
			readFile(path)
		}

		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
		// return ctx.Render(http.StatusOK, "index.html", M{})
	})

	r.GET("/read", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, data)
	})

	r.POST("/", func(ctx echo.Context) error {
		data.CurrentPath = ctx.FormValue("Filepath")
		data.Content = ctx.FormValue("Content")

		writeFile()

		return ctx.JSON(http.StatusOK, M{
			"status": "OK",
		})
	})

	r.Start(":9000")

}
