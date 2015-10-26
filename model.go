package metadata

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
)

const modelTemplate = `## The name of the equipment model.
name = "{{.Name}}"

## Primary device manufacturer.
manufacturer = "{{.Manufacturer}}"

## Optional model specific notes and documentation.
notes = """\
{{$lines := Lines .Notes}}{{range $k, $v := $lines}}    {{$v}}\n\
{{end}}    """

## A list of device model versions.

#[version.id]
#    ## The name of the model version.
#    name = ""
#
#    ## The generic type of the model version.
#    type = ""
#
#    ## An array of extra tags associated with this version.
#    #tags = []
#
#    ## Optional model version specific notes and documentation.
#    #notes = """\
#    #    \n\
#    #    """{{range $k, $v := .Versions}}

[version.{{$k}}]
    ## The name of the device model.
    name = "{{$v.Name}}"

    ## The generic type of the model version.
    type = "{{$v.Type}}"

    ## An array of extra tags associated with this version.
{{if $v.Tags}}    tags = [{{range $n, $t := $v.Tags}}{{if gt $n 0}},{{end}}
        "{{$t}}"{{end}}
    ]{{else}}    #tags = []{{end}}

    ## Optional model version specific notes and documentation.
{{if $v.Notes}}    notes = """\
{{$lines := Lines $v.Notes}}{{range $k, $v := $lines}}        {{$v}}\n\
{{end}}        """{{else}}    #notes = """\
    #    \n\
    #    """{{end}}{{end}}

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
`

type Version struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Notes *string  `json:"notes,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

type Model struct {
	Name         string             `json:"name"`
	Manufacturer string             `json:"manufacturer"`
	Notes        *string            `json:"notes"`
	Versions     map[string]Version `json:"versions,omitempty" toml:"version"`
}

func LoadModel(filename string) (*Model, error) {
	var m Model

	if _, err := toml.DecodeFile(filename, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func LoadModels(dirname, filename string) ([]Model, error) {
	var mm []Model

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			m, e := LoadModel(path)
			if e != nil {
				return e
			}
			mm = append(mm, *m)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return mm, nil
}

func (mod Model) StoreModel(path string) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(mod.String()))
	if err != nil {
		return err
	}

	return nil
}

func (mod Model) String() string {

	tplFuncMap := make(template.FuncMap)
	tplFuncMap["Lines"] = Lines

	tmpl, err := template.New("").Funcs(tplFuncMap).Parse(modelTemplate)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, mod)
	if err != nil {
		panic(err)
	}

	return doc.String()
}
