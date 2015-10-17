package metadata

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestLocation_File(t *testing.T) {

	res := `{"tag":"location","name":"A Location Name","runnet":{"IP":"192.168.192.0","Mask":"////8A=="},"locnet":true,"linknets":[{"name":"From A to B"},{"name":""},{"name":"From A to C"}],"equipment":{"test1":{"name":"test1-location","model":"Test Model 1","tags":["ABCD","EFG","HIJ"],"code":"CODE"},"test2":{"name":"test2-location","model":"Test Model 2","uninstalled":true}}}`

	t.Log("Check loading location files.")
	{
		var l Location
		if _, err := toml.DecodeFile("testdata/location.toml", &l); err != nil {
			t.Error(err)
		}
		if l.String() != res {
			t.Errorf("decoded location file entry mismatch: %s", "testdata/location.toml")
		}
	}
	t.Log("Finished checking location files.")
}
