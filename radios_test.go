package metadata

import (
	"io/ioutil"
	"testing"
)

var testRadioInstalls RadioInstalls

func init() {

	testRadioInstalls = RadioInstalls{
		RadioInstall{
			Location:  "Somewhere",
			Target:    "Somewhere Else",
			Role:      "Master",
			Model:     "Radio Model #1",
			Serial:    "Radio Serial #1",
			Polarity:  "V",
			Frequency: 10,
		},
		RadioInstall{
			Location:  "Somewhere Else",
			Target:    "Somewhere",
			Role:      "Slave",
			Model:     "Radio Model #1",
			Serial:    "Radio Serial #2",
			Polarity:  "V",
			Frequency: 10,
		},
	}
}

func TestRadioInstalls_ReadFile(t *testing.T) {
	t.Log("Compare loaded radio installs file.")
	{
		b, err := ioutil.ReadFile("testdata/radios.csv")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != Strings(testRadioInstalls) {
			t.Errorf("radio installs file text mismatch: [\n%s\n]", SimpleDiff(string(b), Strings(testRadioInstalls)))
		}
	}
}

func TestRadioInstalls_LoadFile(t *testing.T) {
	t.Log("Check loading radio installs file.")
	{
		var installs RadioInstalls
		if err := LoadList("testdata/radios.csv", &installs); err != nil {
			t.Fatal(err)
		}
		if Strings(installs) != Strings(testRadioInstalls) {
			t.Errorf("radio installs file decode mismatch: [\n%s\n]", SimpleDiff(Strings(installs), Strings(testRadioInstalls)))
		}
	}
}

func TestRadioInstalls_LoadFiles(t *testing.T) {
	t.Log("Check loading radio installs files.")
	{
		var installs RadioInstalls
		if err := LoadLists("testdata", "radios.csv", &installs); err != nil {
			t.Fatal(err)
		}
		if Strings(installs) != Strings(testRadioInstalls) {
			t.Errorf("radio installs file decode mismatch: [\n%s\n]", SimpleDiff(Strings(installs), Strings(testRadioInstalls)))
		}
	}
}
