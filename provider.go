package metadata

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type Range struct {
	Name     string      `json:"name"`
	Area     string      `json:"area"`
	Networks []IPNetwork `json:"networks,omitempty"`
}

type Provider struct {
	Name   string  `json:"name"`
	Notes  *string `json:"notes"`
	Ranges []Range `json:"ranges,omitempty" toml:"range"`
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

func (pro Provider) StoreProvider(dir string) error {

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(dir+"/provider.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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

	var l []string

	l = append(l, "# IP4 network allocation tables, for a given service provider or entity.")
	l = append(l, "")
	l = append(l, "#")
	l = append(l, "# RCF1918 private address space.")
	l = append(l, "#")
	l = append(l, "#    10.0.0.0/8")
	l = append(l, "#    192.168.0.0/16")
	l = append(l, "#    172.16.0.0/12")
	l = append(l, "#")
	l = append(l, "")

	l = append(l, "## The name of the network provider.")
	l = append(l, fmt.Sprintf("name = %s", strconv.Quote(pro.Name)))
	l = append(l, "")

	l = append(l, "## Optional provider notes and documentation.")
	if pro.Notes != nil {
		n := strings.Split(strings.Replace(strings.TrimSpace(*pro.Notes), "\\n", "\n", -1), "\n")
		l = append(l, fmt.Sprintf("notes = \"\"\"\\\n\t%s\\n\\\n\t\"\"\"", strings.Join(n, "\\n\\\n\t")))
	} else {
		l = append(l, fmt.Sprintf("#notes = \"\"\"\\\n#\t\\n\\\n#\t\"\"\""))
	}
	l = append(l, "")
	l = append(l, "## An array of provided network ranges.")
	l = append(l, "")
	l = append(l, "#[[range]]")
	l = append(l, "#\t## The name of the network range.")
	l = append(l, "#\tname = \"\"")
	l = append(l, "#")
	l = append(l, "#\t## The network area identification.")
	l = append(l, "#\tarea = \"\"")
	l = append(l, "#")
	l = append(l, "#\t## An array of networks.")
	l = append(l, "#\tnetworks = []")
	for i := 0; i < len(pro.Ranges); i++ {
		l = append(l, "")
		l = append(l, "[[range]]")
		l = append(l, "\t## The name of the network range.")
		l = append(l, fmt.Sprintf("\tname = %s", strconv.Quote(pro.Ranges[i].Name)))
		l = append(l, "")
		l = append(l, "\t## The network area identification.")
		l = append(l, fmt.Sprintf("\tarea = %s", strconv.Quote(pro.Ranges[i].Area)))
		l = append(l, "")
		l = append(l, "\t## An array of networks.")

		var networks []string
		for _, n := range pro.Ranges[i].Networks {
			networks = append(networks, strconv.Quote(n.String()))
		}
		if len(pro.Ranges[i].Networks) > 0 {
			l = append(l, fmt.Sprintf("\tnetworks = [\n\t\t%s\n\t]", strings.Join(networks, ",\n\t\t")))
		} else {
			l = append(l, fmt.Sprintf("\t#networks = []"))
		}
	}

	l = append(l, "")
	l = append(l, "# "+"vim:"+" tabstop=4 expandtab shiftwidth=4 softtabstop=4")
	l = append(l, "")

	return strings.Replace(strings.Join(l, "\n"), "\t", "    ", -1)
}
