package skins

import (
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
)

const tplFile = "template.html"

var (
	baseTemplate *template.Template
)

var args struct {
	BasePath string
}

func Init(basePath string) error {
	args.BasePath = basePath

	tpl, err := compileTemplate(tplFile)
	if err != nil {
		return err
	}
	baseTemplate = tpl

	return nil
}

func Render(fileName string, wr io.Writer, data interface{}) error {
	tpl, err := compileTemplate(fileName)
	if err != nil {
		return err
	}

	return tpl.Execute(wr, data)
}

func compileTemplate(fileName string) (*template.Template, error) {
	filePath := filepath.Join(args.BasePath, fileName)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tpl *template.Template

	if baseTemplate == nil {
		tpl = template.New(fileName)
	} else {
		tpl, err = baseTemplate.Clone()
		if err != nil {
			return nil, err
		}
	}

	tpl, err = tpl.Parse(string(data))
	if err != nil {
		return nil, err
	}

	return tpl, nil
}
