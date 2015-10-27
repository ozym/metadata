package metadata

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
)

const locationTemplate = `# Equipment location site information.

## The unique site specific single word id.
id = "{{.Id}}"

## The general name of the location.
name = "{{.Name}}"

## Geographical position.
{{if .Latitude}}latitude = {{LatLon .Latitude}}{{else}}#latitude = degrees{{end}}
{{if .Longitude}}longitude = {{LatLon .Longitude}}{{else}}#longitude = degrees{{end}}

## An array of service providers associated with this location.
{{if .Services}}services = [{{range $n, $t := .Services}}{{if gt $n 0}},{{end}}
    "{{$t}}"{{end}}
]{{else}}#services = []{{end}}

## An array of tags associated with this location.
{{if .Tags}}tags = [{{range $n, $t := .Tags}}{{if gt $n 0}},{{end}}
    "{{$t}}"{{end}}
]{{else}}#tags = []{{end}}

## Access notes and documentation.
{{if .Access}}access = """\
{{$lines := Lines .Access}}{{range $k, $v := $lines}}    {{$v}}\n\
{{end}}    """{{else}}#access = """\
#    \n\
#    """{{end}}

## Location notes and documentation.
{{if .Notes}}notes = """\
{{$lines := Lines .Notes}}{{range $k, $v := $lines}}    {{$v}}\n\
{{end}}    """{{else}}#notes = """\
#    \n\
#    """{{end}}

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
`

type Location struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Latitude  *float32 `json:"latitude,omitempty"`
	Longitude *float32 `json:"longitude,omitempty"`
	Services  []string `json:"services,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Access    *string  `json:"access,omitempty"`
	Notes     *string  `json:"notes,omitempty"`
}

func LoadLocation(filename string) (*Location, error) {
	var l Location

	if _, err := toml.DecodeFile(filename, &l); err != nil {
		return nil, err
	}

	return &l, nil
}

func LoadLocations(dirname, filename string) ([]Location, error) {
	var ll []Location

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			l, e := LoadLocation(path)
			if e != nil {
				return e
			}
			ll = append(ll, *l)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ll, nil
}

func (loc Location) StoreLocation(path string) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(loc.String()))
	if err != nil {
		return err
	}

	return err
}

func (loc Location) String() string {
	tplFuncMap := make(template.FuncMap)
	tplFuncMap["Lines"] = Lines
	tplFuncMap["LatLon"] = LatLon

	tmpl, err := template.New("").Funcs(tplFuncMap).Parse(locationTemplate)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, loc)
	if err != nil {
		panic(err)
	}

	return doc.String()

}
