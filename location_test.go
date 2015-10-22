package metadata

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
)

var testLocation Location

func init() {

	n, _ := ParseIPNetwork("192.168.192.0/28")

	var a [4]*IPAddress
	for i := 0; i < 4; i++ {
		a[i], _ = ParseIPAddress(fmt.Sprintf("192.168.192.%d/28", i+1))
	}

	testLocation = Location{
		Tag:      "location",
		Name:     "A Location Name",
		Latitude: &[]float32{-41.5}[0],
		Runnet:   n,
		Locnet:   &[]bool{true}[0],
		Linknets: []Linknet{
			Linknet{Name: "From A to B"},
			Linknet{},
			Linknet{Name: "From A to C"},
		},
		Devices: map[string]Device{
			"test1": Device{
				Name:        "test1-location",
				Address:     a[0],
				Aliases:     []IPAddress{*a[1], *a[2]},
				Model:       "Test Model 1",
				Tags:        []string{"ABCD", "EFG", "HIJ"},
				Uninstalled: &[]bool{false}[0],
				Notes:       &[]string{"Some Notes\nSome More Notes"}[0],
			},
			"test2": Device{
				Name:        "test2-location",
				Address:     a[3],
				Model:       "Test Model 2",
				Uninstalled: &[]bool{true}[0],
			},
		},
	}
}

func TestLocation_DecodeFile(t *testing.T) {
	t.Log("Check decoding location file.")
	{
		var l Location
		if _, err := toml.DecodeFile("testdata/location.toml", &l); err != nil {
			t.Fatal(err)
		}
		if l.String() != testLocation.String() {
			t.Errorf("location file text mismatch: [\n%s\n]", SimpleDiff(l.String(), testLocation.String()))
		}
	}
}

func TestLocation_ReadFile(t *testing.T) {
	t.Log("Compare loaded location file.")
	{
		b, err := ioutil.ReadFile("testdata/location.toml")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != testLocation.String() {
			t.Errorf("location file text mismatch: [\n%s\n]", SimpleDiff(string(b), testLocation.String()))
		}
	}
}

func TestLocation_LoadFile(t *testing.T) {
	t.Log("Check loading location file.")
	{
		l, err := LoadLocation("testdata/location.toml")
		if err != nil {
			t.Fatal(err)
		}
		if l == nil {
			t.Fatalf("location file load problem")
		}
		if l.String() != testLocation.String() {
			t.Errorf("location file decode mismatch: [\n%s\n]", SimpleDiff(l.String(), testLocation.String()))
		}
	}
}

func TestLocation_LoadFiles(t *testing.T) {
	t.Log("Check loading location files.")
	{
		l, err := LoadLocations("testdata", "location.toml")
		if err != nil {
			t.Fatal(err)
		}
		if len(l) != 1 {
			t.Fatalf("location files load problem")
		}
		if l[0].String() != testLocation.String() {
			t.Errorf("location file decode mismatch: [\n%s\n]", SimpleDiff(l[0].String(), testLocation.String()))
		}
	}
}
