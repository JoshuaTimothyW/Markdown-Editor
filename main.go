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
	Title      string
	Content    string
	List_files []Files
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
			println("No markdown files")
		}
	}
}

// read file by path
func readFile(path string) {
	b, err := ioutil.ReadFile("./content/markdown-syntax.md")
	if err != nil {
		return
	}

	data.Content = string(b)
}

func main() {

	r := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	r.Renderer = renderer

	r.Static("./static", "views/static")

	r.GET("/", func(ctx echo.Context) error {
		data.Title = "whatever"
		readFile("./content/markdown-syntax.md")
		check_dir()
		return ctx.Render(http.StatusOK, "index.html", M{
			"Title":      data.Title,
			"Content":    data.Content,
			"List_files": data.List_files,
		})
	})

	r.GET("/read", func(ctx echo.Context) error {
		b, err := ioutil.ReadFile("./content/markdown-syntax.md")
		if err != nil {
			return nil
		}
		data := M{"content": string(b)}
		return ctx.JSON(http.StatusOK, data)
	})

	r.Start(":9000")

}
