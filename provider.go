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

## Optional provider notes and documentation.
notes = """\
{{$lines := Lines .Notes}}{{range $k, $v := $lines}}    {{$v}}\n\
{{end}}    """

## A list of provided services.

#[service.label]
#    ## The name of the provided service.
#    name = ""
#
#    ## An optional service reference.
#    #reference = ""
#
#    ## Optional contact details.
#    #contact = ""
#
#    ## Optional service specific notes and documentation.
#    #notes = """\
#    #    \n\
#    #    """{{range $k, $v := .Services}}

[service.{{$k}}]
    ## The name of the provided service.
    name = "{{$v.Name}}"

    ## An optional service reference.
{{if $v.Reference}}    reference = "{{$v.Reference}}"{{else}}    #reference = ""{{end}}

    ## Optional contact details.
{{if $v.Contact}}    contact = "{{$v.Contact}}"{{else}}    #contact = ""{{end}}

    ## Optional service specific notes and documentation.
{{if $v.Notes}}    notes = """\
{{$lines := Lines $v.Notes}}{{range $k, $v := $lines}}        {{$v}}\n\
{{end}}        """{{else}}    #notes = """\
    #    \n\
    #    """{{end}}{{end}}

## An array of provided network ranges.

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
	Notes    *string     `json:"notes"`
	Networks []IPNetwork `json:"networks,omitempty"`
}

type Provider struct {
	Name     string             `json:"name"`
	Notes    *string            `json:"notes"`
	Services map[string]Service `json:"services,omitempty" toml:"service"`
	Ranges   []Range            `json:"ranges,omitempty" toml:"range"`
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
