package main

import (
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	layoutCopy := make([]string, len(layouts))
	for _, include := range includes {
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFilesFuncs(filepath.Base(include), nil, files...)
	}
	return r
}

func renderHTML(ctx *gin.Context, code int, templateName string, data gin.H) {
	ctx.HTML(code, templateName, data)
}
