package metadata

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type EquipmentInstall struct {
	Location string    `csv:"Equipment Location",`
	Model    string    `csv:"Equipment Model",`
	Serial   string    `csv:"Equipment Serial Number",`
	Start    time.Time `csv:"Installation Start",`
	Stop     time.Time `csv:"Installation Stop",`
}

type EquipmentInstalls []EquipmentInstall

func (e EquipmentInstalls) CSV() {}

type SensorInstall struct {
	Station string    `csv:"Seismic Station",`
	Site    string    `csv:"Sensor Location",`
	Model   string    `csv:"Sensor Model",`
	Serial  string    `csv:"Sensor Serial Number",`
	Azimuth float64   `csv:"Azimuth",`
	Dip     float64   `csv:"Dip",`
	Depth   float64   `csv:"Depth",`
	Start   time.Time `csv:"Installation Start",`
	Stop    time.Time `csv:"Installation Stop",`
}

type SensorInstalls []SensorInstall

func (s SensorInstalls) CSV() {}

type DataloggerInstall struct {
	Station string    `csv:"Seismic Station",`
	Site    string    `csv:"Datalogger Location",`
	Model   string    `csv:"Datalogger Model",`
	Serial  string    `csv:"Datalogger Serial Number",`
	Start   time.Time `csv:"Installation Start",`
	Stop    time.Time `csv:"Installation Stop",`
}

type DataloggerInstalls []DataloggerInstall

func (d DataloggerInstalls) CSV() {}

type Install interface {
	CSV()
}

func Strings(install Install) string {
	var b bytes.Buffer

	if lines, err := Encode(install); err == nil {
		csv.NewWriter(&b).WriteAll(lines)
	}

	return b.String()
}

func Decode(data [][]string, install Install) error {
	rv := reflect.ValueOf(install)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("install decode requires a pointer")
	}
	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("install decode requires a pointer to a slice")
	}

	// no data ...
	if !(len(data) > 1) {
		return nil
	}

	// where to start ...
	offset := rv.Len()

	// make space ...
	if rv.IsNil() {
		rv.Set(reflect.MakeSlice(rv.Type(), len(data)-1, len(data)-1))
	} else {
		rv.Set(reflect.AppendSlice(rv, reflect.MakeSlice(rv.Type(), len(data)-1, len(data)-1)))
	}

	// skip the header line ...
	for i, n := 1, len(data); i < n; i++ {
		// gather the slot to store the data
		rvv := reflect.Indirect(rv.Index(offset + i - 1))
		if rvv.Kind() != reflect.Struct {
			return fmt.Errorf("decode requires a pointer to a slice of structs")
		}
		if rvv.NumField() != len(data[i]) {
			return fmt.Errorf("decode incorrect number of fields")
		}
		// decode each field, in order ...
		for j := 0; j < rvv.NumField(); j++ {
			vv := rvv.Field(j)
			switch vv.Kind() {
			case reflect.String:
				vv.SetString(data[i][j])
			case reflect.Int32:
				i, err := strconv.ParseInt(data[i][j], 10, 32)
				if err != nil {
					return err
				}
				vv.SetInt(i)
			case reflect.Float32:
				f, err := strconv.ParseFloat(data[i][j], 32)
				if err != nil {
					return err
				}
				vv.SetFloat(f)
			case reflect.Float64:
				f, err := strconv.ParseFloat(data[i][j], 64)
				if err != nil {
					return err
				}
				vv.SetFloat(f)
			default:
				t, err := time.Parse(DateTimeFormat, data[i][j])
				if err != nil {
					return err
				}
				vv.Set(reflect.ValueOf(t))
			}
		}
	}

	return nil
}

func Encode(install Install) ([][]string, error) {
	var data [][]string

	rv := reflect.ValueOf(install)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("decode requires a pointer to a slice")
	}
	if rv.IsNil() {
		return data, nil
	}
	for i, n := 0, rv.Len(); i < n; i++ {
		v := rv.Index(i).Interface()
		rvv := reflect.ValueOf(v)
		if i == 0 {
			var header []string
			t := reflect.TypeOf(v)
			for j := 0; j < t.NumField(); j++ {
				vv := t.Field(j)
				tags := strings.SplitAfter(strings.TrimSpace(vv.Tag.Get("csv")), ",")
				if len(tags) > 0 && len(tags[0]) > 0 {
					header = append(header, tags[0])
				} else {
					header = append(header, strings.Title(vv.Name))
				}
			}

			data = append(data, header)
		}
		var line []string
		for j := 0; j < rvv.NumField(); j++ {
			vv := rvv.Field(j)
			switch vv.Kind() {
			case reflect.String:
				line = append(line, vv.String())
			case reflect.Int32:
				line = append(line, strconv.FormatInt(vv.Int(), 10))
			case reflect.Float64:
				line = append(line, strconv.FormatFloat(vv.Float(), 'g', -1, 64))
			default:
				line = append(line, vv.Interface().(time.Time).UTC().Format(DateTimeFormat))
			}
		}
		data = append(data, line)
	}

	return data, nil
}

func LoadInstall(path string, install Install) error {

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}

	if err := Decode(data, install); err != nil {
		return err
	}

	return nil
}

func LoadInstalls(dirname, filename string, install Install) error {

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			if err := LoadInstall(path, install); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
