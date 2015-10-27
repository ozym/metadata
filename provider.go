package metadata

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
)

const providerTemplate = `# IP4 network allocation tables, for a given service provider or entity.

#
# RCF1918 private address space.
#
#    10.0.0.0/8
#    192.168.0.0/16
#    172.16.0.0/12
#

## The name of the network provider.
name = "{{.Name}}"

## Povider notes and documentation.
{{if .Notes}}notes = """\
{{$lines := Lines .Notes}}{{range $k, $v := $lines}}    {{$v}}\n\
{{end}}    """{{else}}#notes = """\
#    \n\
#    """{{end}}

## The provided services.

#[[service]]
#    ## The provided service.
#    name = ""
#
#    ## Service reference.
#    #reference = ""
#
#    ## Service contact details.
#    #contact = ""
#
#    ## Service specific notes.
#    #notes = """\
#    #    \n\
#    #    """{{range .Services}}

[[service]]
    ## The provided service.
    name = "{{.Name}}"

    ## Service reference.
{{if .Reference}}    reference = "{{.Reference}}"{{else}}    #reference = ""{{end}}

    ## Service contact details.
{{if .Contact}}    contact = "{{.Contact}}"{{else}}    #contact = ""{{end}}

    ## Service specific notes.
{{if .Notes}}    notes = """\
{{$lines := Lines .Notes}}{{range $k, $v := $lines}}        {{$v}}\n\
{{end}}        """{{else}}    #notes = """\
    #    \n\
    #    """{{end}}{{end}}

## The provided network ranges.

#[[range]]
#    ## The name of the network range.
#    name = ""
#
#    ## The network area identification.
#    area = ""
#
#    ## An array of networks.
#    networks = []
#
#    ## Optional model specific notes and documentation.
#    #notes = """\
#    #    \n\
#    #    """{{range .Ranges}}

[[range]]
    ## The name of the network range.
    name = "{{.Name}}"

    ## The network area identification.
    area = "{{.Area}}"

    ## An array of networks.
{{if .Networks}}    networks = [{{range $n, $t := .Networks}}{{if gt $n 0}},{{end}}
        "{{$t}}"{{end}}
    ]{{else}}    #networks = []{{end}}

    ## Optional model specific notes and documentation.
{{if .Notes}}    notes = """\
{{$lines := Lines .Notes}}{{range $k, $v := $lines}}        {{$v}}\n\
{{end}}        """{{else}}    #notes = """\
    #    \n\
    #    """{{end}}{{end}}

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
`

type Service struct {
	Name      string  `json:"name"`
	Reference *string `json:"reference"`
	Contact   *string `json:"contact"`
	Notes     *string `json:"notes"`
}

type Range struct {
	Name     string      `json:"name"`
	Area     string      `json:"area"`
	Networks []IPNetwork `json:"networks,omitempty"`
	Notes    *string     `json:"notes"`
}

type Provider struct {
	Name     string    `json:"name"`
	Services []Service `json:"services,omitempty" toml:"service"`
	Ranges   []Range   `json:"ranges,omitempty" toml:"range"`
	Notes    *string   `json:"notes"`
}

func LoadProvider(filename string) (*Provider, error) {
	var p Provider

	if _, err := toml.DecodeFile(filename, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func LoadProviders(dirname, filename string) ([]Provider, error) {
	var pp []Provider

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			p, e := LoadProvider(path)
			if e != nil {
				return e
			}
			pp = append(pp, *p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return pp, nil
}

func (pro Provider) StoreProvider(path string) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(pro.String()))
	if err != nil {
		return err
	}

	return nil
}

func (pro Provider) String() string {
	tplFuncMap := make(template.FuncMap)
	tplFuncMap["Lines"] = Lines

	tmpl, err := template.New("").Funcs(tplFuncMap).Parse(providerTemplate)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, pro)
	if err != nil {
		panic(err)
	}

	return doc.String()

}
