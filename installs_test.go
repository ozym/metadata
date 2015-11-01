package metadata

import (
	"io/ioutil"
	"testing"
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
		Install{
			Station: "ABCD",
			Site:    "20",
			Model:   "Model",
			Serial:  "Serial #2",
			Azimuth: 20.0,
			Dip:     20.0,
			Depth:   20.0,
			Start:   MustParseTime("2010-01-01T00:00:00Z"),
			Stop:    MustParseTime("9999-01-01T00:00:00Z"),
		},
		Install{
			Station: "EFGH",
			Site:    "10",
			Model:   "Model",
			Serial:  "Serial #3",
			Azimuth: 10.0,
			Dip:     10.0,
			Depth:   10.0,
			Start:   MustParseTime("2010-01-01T00:00:00Z"),
			Stop:    MustParseTime("2011-01-01T00:00:00Z"),
		},
		Install{
			Station: "EFGH",
			Site:    "20",
			Model:   "Model",
			Serial:  "Serial #4",
			Azimuth: 20.0,
			Dip:     20.0,
			Depth:   20.0,
			Start:   MustParseTime("2010-01-01T00:00:00Z"),
			Stop:    MustParseTime("2011-01-01T00:00:00Z"),
		},
		Install{
			Station: "EFGH",
			Site:    "20",
			Model:   "Model",
			Serial:  "Serial #5",
			Azimuth: 20.0,
			Dip:     20.0,
			Depth:   20.0,
			Start:   MustParseTime("2012-01-01T00:00:00Z"),
			Stop:    MustParseTime("2013-01-01T00:00:00Z"),
		},
	}

}

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

func TestInstalls_LoadFiles(t *testing.T) {
	t.Log("Check loading installs files.")
	{
		m, err := LoadInstallsDir("testdata", "installs.csv")
		if err != nil {
			t.Fatal(err)
		}
		if m.String() != testInstalls.String() {
			t.Errorf("installs file decode mismatch: [\n%s\n]", SimpleDiff(m.String(), testInstalls.String()))
		}
	}
}
