package main

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
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

// initial views
func index(c *fiber.Ctx) error {
	check_dir()

	return c.Render("views/edit.html", fiber.Map{
		"title":      data.Title,
		"list_files": data.List_files,
		"content":    data.Content,
	})

}

func main() {

	app := fiber.New()

	// Static files : css and js
	app.Static("/", "./views")

	// main page
	app.Get("/", index)

	port := ":9000"

	println("Server started at ", port)
	app.Listen(port)

}
