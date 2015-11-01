package metadata

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
)

const equipmentTemplate = `## Equipment details and deployment history.

## The equipment serial number.
serial = "{{.Serial}}"

## The equipment model name.
model = "{{.Model}}"

## Optional equipment asset number.
{{if .Asset}}asset = "{{.Asset}}"{{else}}#asset = ""{{end}}

## Optional equipment specific notes.
notes = """\
{{$lines := Lines .Notes}}{{range $k, $v := $lines}}    {{$v}}\n\
{{end}}    """

## An array of equipment deployments.

#[[deploy]]
#    ## The location of the equipment deployment.
#    location = ""
#
#    ## Deployment start time.
#    start = 2000-01-01T00:00:00Z
#
#    ## Optional deployment stop time.
#    #stop = datetime
#
#    ## Optional deployment specific notes.
#    #notes = """\
#    #    \n\
#    #    """{{range .Deploys}}

[[deploy]]
    ## The location of the equipment deployment.
    location = "{{.Location}}"

    ## Deployment start time.
    start = {{DateTime .Start}}

    ## Optional deployment stop time.
{{if .Stop}}    stop = {{DateTimePtr .Stop}}{{else}}    #stop = datetime{{end}}

    ## Optional deployment specific notes.
{{if .Notes}}    notes = """\
{{range $k, $v := $lines}}        {{$v}}\n\
{{end}}        """{{else}}    #notes = """\
    #    \n\
    #    """{{end}}{{end}}

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
`

type Deploy struct {
	Location string     `json:"location"`
	Start    time.Time  `json:"start"`
	Stop     *time.Time `json:"stop,omitempty" toml:"stop"`
	Notes    *string    `json:"notes,omitempty"`
}

type Equipment struct {
	Serial  string   `json:"serial"`
	Model   string   `json:"model"`
	Asset   *string  `json:"asset"`
	Notes   *string  `json:"notes"`
	Deploys []Deploy `json:"deploys,omitempty" toml:"deploy"`
}

func LoadEquipment(filename string) (*Equipment, error) {
	var m Equipment

	if _, err := toml.DecodeFile(filename, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func LoadEquipments(dirname, filename string) ([]Equipment, error) {
	var mm []Equipment

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			m, e := LoadEquipment(path)
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

func (eq Equipment) StoreEquipment(path string) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(eq.String()))
	if err != nil {
		return err
	}

	return nil
}

func (eq Equipment) String() string {

	tplFuncMap := make(template.FuncMap)
	tplFuncMap["Lines"] = Lines
	tplFuncMap["DateTime"] = DateTime
	tplFuncMap["DateTimePtr"] = DateTimePtr

	tmpl, err := template.New("").Funcs(tplFuncMap).Parse(equipmentTemplate)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, eq)
	if err != nil {
		panic(err)
	}

	return doc.String()
}
