package metadata

import (
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
)

var testModel Model

func init() {

	testModel = Model{
		Name:         "An Example Model",
		Manufacturer: "An Example Model Manufacturer",
		Notes:        &[]string{"Some Notes\nSome More Notes"}[0],
		Versions: map[string]Version{
			"model_a": Version{
				Name: "Model A",
				Type: "Model Type A",
			},
			"model_b": Version{
				Name: "Model B",
				Type: "Model Type B",
				Tags: []string{"A", "B", "C"},
			},
			"model_c": Version{
				Name:  "Model C",
				Notes: &[]string{"Some Notes\nSome More Notes"}[0],
			},
		},
	}
}

func TestModel_DecodeFile(t *testing.T) {
	t.Log("Check decoding model file.")
	{
		var m Model
		if _, err := toml.DecodeFile("testdata/model.toml", &m); err != nil {
			t.Fatal(err)
		}
		if m.String() != testModel.String() {
			t.Errorf("model file text mismatch: [\n%s\n]", SimpleDiff(m.String(), testModel.String()))
		}
	}
}

func TestModel_ReadFile(t *testing.T) {
	t.Log("Compare loaded model file.")
	{
		b, err := ioutil.ReadFile("testdata/model.toml")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != testModel.String() {
			t.Errorf("model file text mismatch: [\n%s\n]", SimpleDiff(string(b), testModel.String()))
		}
	}
}

func TestModel_LoadFile(t *testing.T) {
	t.Log("Check loading model file.")
	{
		m, err := LoadModel("testdata/model.toml")
		if err != nil {
			t.Fatal(err)
		}
		if m == nil {
			t.Fatalf("model file load problem")
		}
		if m.String() != testModel.String() {
			t.Errorf("model file decode mismatch: [\n%s\n]", SimpleDiff(m.String(), testModel.String()))
		}
	}
}

func TestModel_LoadFiles(t *testing.T) {
	t.Log("Check loading model files.")
	{
		m, err := LoadModels("testdata", "model.toml")
		if err != nil {
			t.Fatal(err)
		}
		if len(m) != 1 {
			t.Fatalf("model files load problem")
		}
		if m[0].String() != testModel.String() {
			t.Errorf("model file decode mismatch: [\n%s\n]", SimpleDiff(m[0].String(), testModel.String()))
		}
	}
}
