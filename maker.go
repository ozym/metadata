package metadata

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
)

const makerTemplate = `## The name of the device maker.
name = "{{.Name}}"

## Optional maker specific notes and documentation.
notes = """\
{{$lines := Lines .Notes}}{{range $k, $v := $lines}}    {{$v}}\n\
{{end}}    """

## An array of device model types.

#[[model]]
#    ## The name of the device model.
#    name = ""
#
#    ## The generic type of the model.
#    type = ""
#
#    ## An array of extra tags associated with this model.
#    #tags = []
#
#    ## Optional model specific notes and documentation.
#    #notes = """\
#    #    \n\
#    #    """{{range .Models}}

[[model]]
    ## The name of the device model.
    name = "{{.Name}}"

    ## The generic type of the model.
    type = "{{.Type}}"

    ## An array of extra tags associated with this model.
{{if .Tags}}    tags = [{{range $n, $t := .Tags}}{{if gt $n 0}},{{end}}
        "{{$t}}"{{end}}
    ]{{else}}    #tags = []{{end}}

    ## Optional model specific notes and documentation.
{{if .Notes}}    notes = """\
{{range $k, $v := $lines}}        {{$v}}\n\
{{end}}        """{{else}}    #notes = """\
    #    \n\
    #    """{{end}}{{end}}

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
`

type Model struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Notes    *string     `json:"notes,omitempty"`
	Tags     []string    `json:"tags,omitempty"`
	Networks []IPNetwork `json:"networks,omitempty"`
}

type Maker struct {
	Name   string  `json:"name"`
	Notes  *string `json:"notes"`
	Models []Model `json:"ranges,omitempty" toml:"model"`
}

func LoadMaker(filename string) (*Maker, error) {
	var m Maker

	if _, err := toml.DecodeFile(filename, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func LoadMakers(dirname, filename string) ([]Maker, error) {
	var mm []Maker

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			m, e := LoadMaker(path)
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

func (mak Maker) StoreMaker(path string) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(mak.String()))
	if err != nil {
		return err
	}

	return nil
}

func (mak Maker) String() string {

	tplFuncMap := make(template.FuncMap)
	tplFuncMap["Lines"] = Lines

	tmpl, err := template.New("").Funcs(tplFuncMap).Parse(makerTemplate)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, mak)
	if err != nil {
		panic(err)
	}

	return doc.String()
}
