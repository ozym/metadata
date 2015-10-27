package metadata

import (
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
)

var testProvider Provider

func init() {

	testProvider = Provider{
		Name:  "Example Provider",
		Notes: &[]string{"Some Notes\nSome More Notes"}[0],
		Services: []Service{Service{
			Name:      "Test Service",
			Notes:     &[]string{"Some Notes\nSome More Notes"}[0],
			Reference: &[]string{"ABC1234"}[0],
			Contact:   &[]string{"0800 123123"}[0],
		}},
		Ranges: []Range{
			Range{
				Name: "Private Networks",
				Area: "0.0.0.1",
				Networks: []IPNetwork{
					*MustParseIPNetwork("10.100.41.0/24"),
					*MustParseIPNetwork("10.100.45.0/24"),
				},
			},
			Range{
				Name: "More Private Networks",
				Area: "0.0.0.2",
				Networks: []IPNetwork{
					*MustParseIPNetwork("10.51.0.0/16"),
					*MustParseIPNetwork("10.52.0.0/16"),
				},
			},
			Range{
				Name: "An Empty Range",
				Area: "0.0.0.3",
			},
		},
	}
}

func TestProvider_DecodeFile(t *testing.T) {
	t.Log("Check decoding provider file.")
	{
		var p Provider
		if _, err := toml.DecodeFile("testdata/provider.toml", &p); err != nil {
			t.Fatal(err)
		}
		if p.String() != testProvider.String() {
			t.Errorf("provider file text mismatch: [\n%s\n]", SimpleDiff(p.String(), testProvider.String()))
		}
	}
}

func TestProvider_ReadFile(t *testing.T) {
	t.Log("Compare loaded provider file.")
	{
		b, err := ioutil.ReadFile("testdata/provider.toml")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != testProvider.String() {
			t.Errorf("provider file text mismatch: [\n%s\n]", SimpleDiff(string(b), testProvider.String()))
		}
	}
}

func TestProvider_LoadFile(t *testing.T) {
	t.Log("Check loading provider file.")
	{
		p, err := LoadProvider("testdata/provider.toml")
		if err != nil {
			t.Fatal(err)
		}
		if p == nil {
			t.Fatalf("provider file load problem")
		}
		if p.String() != testProvider.String() {
			t.Errorf("provider file decode mismatch: [\n%s\n]", SimpleDiff(p.String(), testProvider.String()))
		}
	}
}

func TestProvider_LoadFiles(t *testing.T) {
	t.Log("Check loading provider files.")
	{
		p, err := LoadProviders("testdata", "provider.toml")
		if err != nil {
			t.Fatal(err)
		}
		if len(p) != 1 {
			t.Fatalf("provider files load problem")
		}
		if p[0].String() != testProvider.String() {
			t.Errorf("provider file decode mismatch: [\n%s\n]", SimpleDiff(p[0].String(), testProvider.String()))
		}
		//t.Log(p[0].String())
	}
}
