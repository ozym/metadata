package metadata

import (
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	//	"github.com/ozym/metadata"
)

func TestLocation_File(t *testing.T) {

	t.Log("Check loading location files.")
	{

		n, err := ParseIPNetwork("192.168.192.0/28")
		if err != nil {
			t.Error(err)
			return
		}

		a1, err := ParseIPAddress("192.168.192.1/28")
		if err != nil {
			t.Error(err)
			return
		}

		a2, err := ParseIPAddress("192.168.192.2/28")
		if err != nil {
			t.Error(err)
			return
		}

		a3, err := ParseIPAddress("192.168.192.3/28")
		if err != nil {
			t.Error(err)
			return
		}

		a4, err := ParseIPAddress("192.168.192.4/28")
		if err != nil {
			t.Error(err)
			return
		}

		f := Location{
			Tag:      "location",
			Name:     "A Location Name",
			Latitude: &[]float32{-41.5}[0],
			//			Longitude: &[]float32{174.5}[0],
			Runnet: n,
			Locnet: &[]bool{true}[0],
			Linknets: []Linknet{
				Linknet{Name: "From A to B"},
				Linknet{},
				Linknet{Name: "From A to C"},
			},
			Equipment: map[string]Located{
				"test1": Located{
					Name:        "test1-location",
					Address:     a1,
					Aliases:     []IPAddress{*a2, *a3},
					Model:       "Test Model 1",
					Tags:        []string{"ABCD", "EFG", "HIJ"},
					Uninstalled: &[]bool{false}[0],
					Notes:       &[]string{"Some Notes"}[0],
				},
				"test2": Located{
					Name:        "test2-location",
					Address:     a4,
					Model:       "Test Model 2",
					Uninstalled: &[]bool{true}[0],
				},
			},
		}
		var l Location
		if _, err := toml.DecodeFile("testdata/location.toml", &l); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(l, f) {
			t.Errorf("location file entry mismatch: %s", "testdata/location.toml")
		}
	}
}
