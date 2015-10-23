package metadata

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

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

func (mak Maker) StoreMaker(dir string) error {

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(dir+"/maker.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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

	var l []string

	/*
		l = append(l, "# IP4 network allocation tables, for a given service maker or entity.")
		l = append(l, "")
		l = append(l, "#")
		l = append(l, "# RCF1918 private address space.")
		l = append(l, "#")
		l = append(l, "#    10.0.0.0/8")
		l = append(l, "#    192.168.0.0/16")
		l = append(l, "#    172.16.0.0/12")
		l = append(l, "#")
		l = append(l, "")
	*/

	l = append(l, "## The name of the device maker.")
	l = append(l, fmt.Sprintf("name = %s", strconv.Quote(mak.Name)))
	l = append(l, "")

	l = append(l, "## Optional maker specific notes and documentation.")
	if mak.Notes != nil {
		n := strings.Split(strings.Replace(strings.TrimSpace(*mak.Notes), "\\n", "\n", -1), "\n")
		l = append(l, fmt.Sprintf("notes = \"\"\"\\\n\t%s\\n\\\n\t\"\"\"", strings.Join(n, "\\n\\\n\t")))
	} else {
		l = append(l, fmt.Sprintf("#notes = \"\"\"\\\n#\t\\n\\\n#\t\"\"\""))
	}
	l = append(l, "")
	l = append(l, "## An array of device model types.")
	l = append(l, "")
	l = append(l, "#[[model]]")
	l = append(l, "#\t## The name of the device model.")
	l = append(l, "#\tname = \"\"")
	l = append(l, "#")
	l = append(l, "#\t## The generic type of the model.")
	l = append(l, "#\ttype = \"\"")
	l = append(l, "#")
	l = append(l, "#\t## An array of extra tags associated with this model.")
	l = append(l, "#\t#tags = []")
	l = append(l, "#")
	l = append(l, "#\t## Optional model specific notes and documentation.")
	l = append(l, fmt.Sprintf("#\t#notes = \"\"\"\\\n#\t#\t\\n\\\n#\t#\t\"\"\""))

	for i := 0; i < len(mak.Models); i++ {
		l = append(l, "")
		l = append(l, "[[model]]")
		l = append(l, "\t## The name of the device model.")
		l = append(l, fmt.Sprintf("\tname = %s", strconv.Quote(mak.Models[i].Name)))
		l = append(l, "")
		l = append(l, "\t## The generic type of the model.")
		l = append(l, fmt.Sprintf("\ttype = %s", strconv.Quote(mak.Models[i].Type)))
		l = append(l, "")
		l = append(l, "\t## An array of extra tags associated with this model.")
		if len(mak.Models[i].Tags) > 0 {
			l = append(l, fmt.Sprintf("\ttags = [\n\t\t\"%s\"\n\t]", strings.Join(mak.Models[i].Tags, "\",\n\t\t\"")))
		} else {
			l = append(l, fmt.Sprintf("\t#tags = []"))
		}

		l = append(l, "")
		l = append(l, "\t## Optional model specific notes and documentation.")
		if mak.Models[i].Notes != nil {
			n := strings.Split(strings.Replace(strings.TrimSpace(*mak.Models[i].Notes), "\\n", "\n", -1), "\n")
			l = append(l, fmt.Sprintf("\tnotes = \"\"\"\\\n\t\t%s\\n\\\n\t\t\"\"\"", strings.Join(n, "\\n\\\n\t\t")))
		} else {
			l = append(l, fmt.Sprintf("\t#notes = \"\"\"\\\n\t#\t\\n\\\n\t#\t\"\"\""))
		}
		/*
		 */

		/*
			l = append(l, "")
			l = append(l, "\t## An array of networks.")

			var networks []string
			for _, n := range mak.Models[i].Networks {
				networks = append(networks, strconv.Quote(n.String()))
			}
			if len(mak.Models[i].Networks) > 0 {
				l = append(l, fmt.Sprintf("\tnetworks = [\n\t\t%s\n\t]", strings.Join(networks, ",\n\t\t")))
			} else {
				l = append(l, fmt.Sprintf("\t#networks = []"))
			}
		*/
	}

	l = append(l, "")
	l = append(l, "# "+"vim:"+" tabstop=4 expandtab shiftwidth=4 softtabstop=4")
	l = append(l, "")

	return strings.Replace(strings.Join(l, "\n"), "\t", "    ", -1)
}
