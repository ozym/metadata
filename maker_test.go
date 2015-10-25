package metadata

import (
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
)

var testMaker Maker

func init() {

	testMaker = Maker{
		Name:  "An Example Maker",
		Notes: &[]string{"Some Notes\nSome More Notes"}[0],
		Models: []Model{
			Model{
				Name: "Model A",
				Type: "Model Type A",
			},
			Model{
				Name: "Model B",
				Type: "Model Type B",
				Tags: []string{"A", "B", "C"},
			},
			Model{
				Name:  "Model C",
				Notes: &[]string{"Some Notes\nSome More Notes"}[0],
			},
		},
	}
}

func TestMaker_DecodeFile(t *testing.T) {
	t.Log("Check decoding maker file.")
	{
		var m Maker
		if _, err := toml.DecodeFile("testdata/maker.toml", &m); err != nil {
			t.Fatal(err)
		}
		if m.String() != testMaker.String() {
			t.Errorf("maker file text mismatch: [\n%s\n]", SimpleDiff(m.String(), testMaker.String()))
		}
	}
}

func TestMaker_ReadFile(t *testing.T) {
	t.Log("Compare loaded maker file.")
	{
		b, err := ioutil.ReadFile("testdata/maker.toml")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != testMaker.String() {
			t.Errorf("maker file text mismatch: [\n%s\n]", SimpleDiff(string(b), testMaker.String()))
		}
	}
}

func TestMaker_LoadFile(t *testing.T) {
	t.Log("Check loading maker file.")
	{
		m, err := LoadMaker("testdata/maker.toml")
		if err != nil {
			t.Fatal(err)
		}
		if m == nil {
			t.Fatalf("maker file load problem")
		}
		if m.String() != testMaker.String() {
			t.Errorf("maker file decode mismatch: [\n%s\n]", SimpleDiff(m.String(), testMaker.String()))
		}
	}
}

func TestMaker_LoadFiles(t *testing.T) {
	t.Log("Check loading maker files.")
	{
		m, err := LoadMakers("testdata", "maker.toml")
		if err != nil {
			t.Fatal(err)
		}
		if len(m) != 1 {
			t.Fatalf("maker files load problem")
		}
		if m[0].String() != testMaker.String() {
			t.Errorf("maker file decode mismatch: [\n%s\n]", SimpleDiff(m[0].String(), testMaker.String()))
		}
	}
}
