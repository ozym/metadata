package metadata

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
)

var testEquipment Equipment

func init() {

	testEquipment = Equipment{
		Serial: "serial",
		Model:  "Model",
		Asset:  &[]string{"asset"}[0],
		Notes:  &[]string{"Some Notes\nSome More Notes"}[0],
		Deploys: []Deploy{
			Deploy{
				Location: "location1",
				Start:    MustParseTime("2010-01-01T00:00:00Z"),
				Stop:     &[]time.Time{MustParseTime("2011-01-01T00:00:00Z")}[0],
			},
			Deploy{
				Location: "location2",
				Start:    MustParseTime("2012-01-01T00:00:00Z"),
				Notes:    &[]string{"Some Notes\nSome More Notes"}[0],
			},
		},
	}
}

func TestEquipment_DecodeFile(t *testing.T) {
	t.Log("Check decoding equipment file.")
	{
		var m Equipment
		if _, err := toml.DecodeFile("testdata/equipment.toml", &m); err != nil {
			t.Fatal(err)
		}
		if m.String() != testEquipment.String() {
			t.Errorf("equipment file text mismatch: [\n%s\n]", SimpleDiff(m.String(), testEquipment.String()))
		}
	}
}

func TestEquipment_ReadFile(t *testing.T) {
	t.Log("Compare loaded equipment file.")
	{
		b, err := ioutil.ReadFile("testdata/equipment.toml")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != testEquipment.String() {
			t.Errorf("equipment file text mismatch: [\n%s\n]", SimpleDiff(string(b), testEquipment.String()))
		}
	}
}

func TestEquipment_LoadFile(t *testing.T) {
	t.Log("Check loading equipment file.")
	{
		m, err := LoadEquipment("testdata/equipment.toml")
		if err != nil {
			t.Fatal(err)
		}
		if m == nil {
			t.Fatalf("equipment file load problem")
		}
		if m.String() != testEquipment.String() {
			t.Errorf("equipment file decode mismatch: [\n%s\n]", SimpleDiff(m.String(), testEquipment.String()))
		}
	}
}

func TestEquipment_LoadFiles(t *testing.T) {
	t.Log("Check loading equipment files.")
	{
		m, err := LoadEquipments("testdata", "equipment.toml")
		if err != nil {
			t.Fatal(err)
		}
		if len(m) != 1 {
			t.Fatalf("equipment files load problem")
		}
		if m[0].String() != testEquipment.String() {
			t.Errorf("equipment file decode mismatch: [\n%s\n]", SimpleDiff(m[0].String(), testEquipment.String()))
		}
	}
}
