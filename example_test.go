package metadata

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestEncoder_File(t *testing.T) {

	t.Log("Check loading toml files.")
	{

		f := Location{
			Tag:      "location",
			Name:     "A Location Name",
			Latitude: &[]float32{-41.5}[0],
			//			Longitude: &[]float32{174.5}[0],
			Runnet: ParseIPNetwork("192.168.192.0/28"),
			Locnet: &[]bool{true}[0],
			Linknets: []Linknet{
				Linknet{Name: "From A to B"},
				Linknet{},
				Linknet{Name: "From A to C"},
			},
			Equipment: map[string]Equipment{
				"test1": Equipment{
					Name:    "test1-location",
					Address: ParseIPAddress("192.168.192.1/28"),
					Aliases: []IPAddress{
						*ParseIPAddress("192.168.192.2/28"),
						*ParseIPAddress("192.168.192.3/28"),
					},
					Model:       "Test Model 1",
					Tags:        []string{"ABCD", "EFG", "HIJ"},
					Uninstalled: &[]bool{false}[0],
					Notes:       &[]string{"Some Notes"}[0],
				},
				"test2": Equipment{
					Name:        "test2-location",
					Address:     ParseIPAddress("192.168.192.4/28"),
					Model:       "Test Model 2",
					Uninstalled: &[]bool{true}[0],
				},
			},
		}

		var l Location
		if _, err := toml.DecodeFile("testdata/example.toml", &l); err != nil {
			t.Error(err)
		}
		if !l.Equal(f) {
			t.Errorf("file entry mismatch: \n%s\n", SimpleDiff(f, l))
		}
		if err := f.Validate(); err != nil {
			t.Error(err)
		}
	}
	t.Log("Finished checking location files.")
}
