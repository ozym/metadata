package metadata

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var InstallsHeader []string = []string{
	"Seismic Station",
	"Site Location",
	"Sensor Model",
	"Sensor Serial Number",
	"Azimuth",
	"Dip",
	"Depth",
	"Installation Start",
	"Installation Stop",
}

func IsInstallsHeader(line []string) bool {
	if len(line) != len(InstallsHeader) {
		return false
	}
	for i, n := 0, len(InstallsHeader); i < n; i++ {
		if InstallsHeader[i] != line[i] {
			return false
		}
	}
	return true
}

type Install struct {
	Station string
	Site    string
	Model   string
	Serial  string
	Azimuth float64
	Dip     float64
	Depth   float64
	Start   time.Time
	Stop    time.Time
}

func ParseInstall(line []string) (*Install, error) {
	if len(line) != 9 {
		return nil, fmt.Errorf("not enough installation values")
	}
	if IsInstallsHeader(line) {
		return nil, nil
	}

	azimuth, err := strconv.ParseFloat(line[4], 64)
	if err != nil {
		return nil, err
	}

	dip, err := strconv.ParseFloat(line[5], 64)
	if err != nil {
		return nil, err
	}

	depth, err := strconv.ParseFloat(line[6], 64)
	if err != nil {
		return nil, err
	}

	start, err := time.Parse(DateTimeFormat, line[7])
	if err != nil {
		return nil, err
	}

	stop, err := time.Parse(DateTimeFormat, line[8])
	if err != nil {
		return nil, err
	}

	return &Install{
		Station: line[0],
		Site:    line[1],
		Model:   line[2],
		Serial:  line[3],
		Azimuth: azimuth,
		Dip:     dip,
		Depth:   depth,
		Start:   start,
		Stop:    stop,
	}, nil
}

func (in Install) Strings() []string {
	var result []string

	result = append(result, in.Station, in.Site)
	result = append(result, in.Model, in.Serial)
	result = append(result, strconv.FormatFloat(in.Azimuth, 'g', -1, 64))
	result = append(result, strconv.FormatFloat(in.Dip, 'g', -1, 64))
	result = append(result, strconv.FormatFloat(in.Depth, 'g', -1, 64))
	result = append(result, in.Start.UTC().Format(DateTimeFormat))
	result = append(result, in.Stop.UTC().Format(DateTimeFormat))

	return result
}

type Installs []Install

func (in Installs) String() string {
	var lines [][]string
	lines = append(lines, InstallsHeader)
	for _, i := range in {
		lines = append(lines, i.Strings())
	}

	var b bytes.Buffer
	if err := csv.NewWriter(&b).WriteAll(lines); err != nil {
		return ""
	}

	return b.String()
}

func (in Installs) StoreInstalls(path string) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines [][]string

	lines = append(lines, InstallsHeader)

	w := csv.NewWriter(file)
	w.WriteAll(lines)

	if err := w.Error(); err != nil {
		return err
	}

	return nil

}

func LoadInstalls(path string) (Installs, error) {

	var installs []Install
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, line := range records {
		l, err := ParseInstall(line)
		if err != nil {
			return nil, err
		}
		if l != nil {
			installs = append(installs, *l)
		}
	}

	return installs, nil
}

func LoadInstallsDir(dirname, filename string) (Installs, error) {
	var installs Installs

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			i, e := LoadInstalls(path)
			if e != nil {
				return e
			}
			installs = append(installs, i...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return installs, nil
}

/*
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
*/
