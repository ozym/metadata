package metadata

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestLocation_File(t *testing.T) {

	t.Log("Check loading location files.")
	{

		f := Location{
			Tag:    "location",
			Name:   "A Location Name",
			Runnet: ParseIPNetwork("192.168.192.0/28"),
			Locnet: true,
			Linknets: []Linknet{
				Linknet{Name: "From A to B"},
				Linknet{},
				Linknet{Name: "From A to C"},
			},
			Equipment: map[string]Equipment{
				"test1": Equipment{
					Name:    "test1-location",
					Address: ParseIPAddress("192.168.192.1/28"),
					Model:   "Test Model 1",
					Tags:    []string{"ABCD", "EFG", "HIJ"},
					Code:    &[]string{"CODE"}[0],
				},
				"test2": Equipment{
					Name:        "test2-location",
					Address:     ParseIPAddress("192.168.192.4/28"),
					Model:       "Test Model 2",
					Uninstalled: true,
				},
			},
		}

		var l Location
		if _, err := toml.DecodeFile("testdata/location.toml", &l); err != nil {
			t.Error(err)
		}
		if !l.Equal(f) {
			t.Errorf("location file entry mismatch: %s [\n%s\n]", "testdata/location.toml", diff(f, l))
		}
	}
	t.Log("Finished checking location files.")
}
