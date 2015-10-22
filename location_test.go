package metadata

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestLocation_File(t *testing.T) {
	var err error

	var l Location
	t.Log("Check loading location file.")
	{
		if _, err = toml.DecodeFile("testdata/location.toml", &l); err != nil {
			t.Error(err)
		}
	}

	var b []byte
	t.Log("Compare loaded location file.")
	{
		if b, err = ioutil.ReadFile("testdata/location.toml"); err != nil {
			t.Error(err)
		}
		if string(b) != l.String() {
			t.Errorf("location file text mismatch: [\n%s\n]", SimpleDiff(string(b), l.String()))
		}
	}

	var f Location
	t.Log("Compare loaded location details.")
	{
		var n *IPNetwork
		var a [4]*IPAddress

		if n, err = ParseIPNetwork("192.168.192.0/28"); err != nil {
			t.Error(err)
			return
		}

		for i := 0; i < 4; i++ {
			s := fmt.Sprintf("192.168.192.%d/28", i+1)
			if a[i], err = ParseIPAddress(s); err != nil {
				t.Error(err)
				return
			}
		}

		f = Location{
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
		if l.String() != f.String() {
			t.Errorf("location details mismatch: [\n%s\n]", SimpleDiff(l.String(), f.String()))
		}
	}
}
