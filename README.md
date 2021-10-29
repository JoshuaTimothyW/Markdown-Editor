# Markdown Editor

Create my own version of markdown editor with Go and Toast UI Editor.

By default it detect all markdown files in current directory, also fits with hugo project to detect content directory :

```
hugo serve

mdeditor
```

You can download binary [here](https://github.com/JoshuaTimothyW/Markdown-Editor/releases)

## Targeted Features

* [x] List files
* [x] Make hierarchy from list of directories and files
* [x] Embedded views and assets
* [x] Create New File
* [ ] Drag and drop markdown file
* [ ] Rename and Delete Files 
* [ ] Autocomplete

## Run local server

```
go run main.go
```

or

```
npm run dev
```

## Build binary

```
go build -o /bin
```

or

```
npm run build
```