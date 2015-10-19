package metadata

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestLocation_File(t *testing.T) {

	t.Log("Check loading location files.")
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
					Code:        &[]string{"CODE"}[0],
					Uninstalled: &[]bool{false}[0],
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
		if _, err := toml.DecodeFile("testdata/location.toml", &l); err != nil {
			t.Error(err)
		}
		if !l.Equal(f) {
			t.Errorf("location file entry mismatch: %s [\n%s\n]", "testdata/location.toml", Diff(f, l))
		}
	}
	t.Log("Finished checking location files.")
}
