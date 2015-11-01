package metadata

import (
	"io/ioutil"
	"testing"
	//	"time"

	//	"github.com/BurntSushi/toml"
)

var testInstalls Installs

func init() {

	testInstalls = Installs{
		Install{
			Station: "ABCD",
			Site:    "10",
			Model:   "Model",
			Serial:  "Serial #1",
			Azimuth: 10.0,
			Dip:     10.0,
			Depth:   10.0,
			Start:   MustParseTime("2010-01-01T00:00:00Z"),
			Stop:    MustParseTime("2011-01-01T00:00:00Z"),
		},
	}
}

/*
func TestInstalls_DecodeFile(t *testing.T) {
	t.Log("Check decoding installs file.")
	{
		var i Installs
		if _, err := installs.DecodeFile("testdata/installs.csv", &i); err != nil {
			t.Fatal(err)
		}
		if i.String() != testInstalls.String() {
			t.Errorf("installs file text mismatch: [\n%s\n]", SimpleDiff(i.String(), testInstalls.String()))
		}
	}
}
*/

func TestInstalls_ReadFile(t *testing.T) {
	t.Log("Compare loaded installs file.")
	{
		b, err := ioutil.ReadFile("testdata/installs.csv")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != testInstalls.String() {
			t.Errorf("installs file text mismatch: [\n%s\n]", SimpleDiff(string(b), testInstalls.String()))
		}
	}
}

func TestInstalls_LoadFile(t *testing.T) {
	t.Log("Check loading installs file.")
	{
		m, err := LoadInstalls("testdata/installs.csv")
		if err != nil {
			t.Fatal(err)
		}
		if m == nil {
			t.Fatalf("installs file load problem")
		}
		if m.String() != testInstalls.String() {
			t.Errorf("installs file decode mismatch: [\n%s\n]", SimpleDiff(m.String(), testInstalls.String()))
		}
	}
}

/*
func TestInstalls_LoadFiles(t *testing.T) {
	t.Log("Check loading installs files.")
	{
		m, err := LoadInstallsDir("testdata", "installs.csv")
		if err != nil {
			t.Fatal(err)
		}
		if len(m) != 1 {
			t.Fatalf("installs files load problem")
		}
		if m[0].String() != testInstalls.String() {
			t.Errorf("installs file decode mismatch: [\n%s\n]", SimpleDiff(m[0].String(), testInstalls.String()))
		}
	}
}
*/
