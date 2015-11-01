package metadata

import (
	//	"fmt"
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
)

var testNetwork Network

func init() {

	testNetwork = Network{
		Location: "network",
		Name:     &[]string{"A Network Name"}[0],
		Notes:    &[]string{"Some Notes\nSome More Notes"}[0],
		Runnet:   MustParseIPNetwork("192.168.192.0/28"),
		Linknets: []Linknet{
			Linknet{Name: "From A to B"},
			Linknet{},
			Linknet{Name: "From A to C"},
		},
		Devices: []Device{
			Device{
				Name:        "rf2somewhere-network",
				Address:     MustParseIPAddress("192.168.192.5/28"),
				Model:       "Test Radio",
				Links:       []string{"rf2network-somewhere"},
				Uninstalled: &[]bool{true}[0],
			},
			Device{
				Name:    "test1-network",
				Address: MustParseIPAddress("192.168.192.1/28"),
				Aliases: []IPAddress{
					*MustParseIPAddress("192.168.192.2/28"),
					*MustParseIPAddress("192.168.192.3/28"),
				},
				Model:       "Test Model 1",
				Tags:        []string{"ABCD", "EFG", "HIJ"},
				Uninstalled: &[]bool{false}[0],
				Notes:       &[]string{"Some Notes\nSome More Notes"}[0],
			},
			Device{
				Name:        "test2-network",
				Address:     MustParseIPAddress("192.168.192.4/28"),
				Model:       "Test Model 2",
				Uninstalled: &[]bool{true}[0],
			},
		},
	}
}

func TestNetwork_IPAddress(t *testing.T) {

	t.Log("Check valid ip addresses.")
	{
		tests := []string{
			"192.168.0.1/16",
			"10.0.0.1/8",
			"10.0.0.1/24",
			"10.0.0.1/28",
			"10.0.0.1/32",
		}

		for _, s := range tests {
			a, err := ParseIPAddress(s)
			if err != nil {
				t.Error(err)
			}
			if s != a.String() {
				t.Errorf("address mismatch: \"%s\" != \"%s\"", s, a.String())
			}
		}
	}

	t.Log("Check invalid ip addresses.")
	{
		tests := []string{
			"a.b.c.d",
			"192.168.0.1",
			"192.168.0",
			"192.168",
			"192",
			"",
			"256.256.256.256/256",
			"256.256.256.256/8",
			"256.256.256.1/8",
			"256.256.0.1/8",
			"256.0.0.1/8",
			"10.0.0.0/-1",
		}

		for _, s := range tests {
			_, err := ParseIPAddress(s)
			if err == nil {
				t.Errorf("address should not be valid: \"%s\"", s)
			}
		}
	}
}

func TestNetwork_IPNetwork(t *testing.T) {

	t.Log("Check valid ip networks.")
	{
		tests := []string{
			"192.168.0.0/16",
			"10.0.0.0/8",
			"10.0.1.0/24",
			"10.0.2.16/28",
			"10.0.3.20/30",
		}

		for _, s := range tests {
			a, err := ParseIPNetwork(s)
			if err != nil {
				t.Error(err)
			}
			if s != a.String() {
				t.Errorf("network mismatch: \"%s\" != \"%s\"", s, a.String())
			}
		}
	}

	t.Log("Check invalid ip networks.")
	{
		tests := []string{
			"a.b.c.d",
			"192.168.0.1",
			"192.168.0",
			"192.168",
			"192",
			"",
			"256.256.256.256/256",
			"256.256.256.256/8",
			"256.256.256.1/8",
			"256.256.0.1/8",
			"256.0.0.1/8",
		}
		for _, s := range tests {
			_, err := ParseIPNetwork(s)
			if err == nil {
				t.Errorf("network should be invalid: \"%s\"", s)
			}
		}
	}

	t.Log("Check mapped ip networks.")
	{
		tests := []struct{ a, b string }{
			{"10.0.0.1/8", "10.0.0.0/8"},
			{"10.0.0.1/24", "10.0.0.0/24"},
			{"10.0.1.1/28", "10.0.1.0/28"},
			{"10.0.2.1/32", "10.0.2.1/32"},
			{"10.0.3.18/30", "10.0.3.16/30"},
		}

		for _, s := range tests {
			n, err := ParseIPNetwork(s.a)
			if err != nil {
				t.Errorf("network should be valid: \"%s\"", s.a)
			} else if s.b != n.String() {
				t.Errorf("network mismatch: \"%s\" != \"%s\"", s.b, n.String())
			}
		}
	}
}

func TestNetwork_DecodeFile(t *testing.T) {
	t.Log("Check decoding network file.")
	{
		var l Network
		if _, err := toml.DecodeFile("testdata/network.toml", &l); err != nil {
			t.Fatal(err)
		}
		if l.String() != testNetwork.String() {
			t.Errorf("network file text mismatch: [\n%s\n]", SimpleDiff(l.String(), testNetwork.String()))
		}
	}
}

func TestNetwork_ReadFile(t *testing.T) {
	t.Log("Compare loaded network file.")
	{
		b, err := ioutil.ReadFile("testdata/network.toml")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != testNetwork.String() {
			t.Errorf("network file text mismatch: [\n%s\n]", SimpleDiff(string(b), testNetwork.String()))
		}
	}
}

func TestNetwork_LoadFile(t *testing.T) {
	t.Log("Check loading network file.")
	{
		l, err := LoadNetwork("testdata/network.toml")
		if err != nil {
			t.Fatal(err)
		}
		if l == nil {
			t.Fatalf("network file load problem")
		}
		if l.String() != testNetwork.String() {
			t.Errorf("network file decode mismatch: [\n%s\n]", SimpleDiff(l.String(), testNetwork.String()))
		}
	}
}

func TestNetwork_LoadFiles(t *testing.T) {
	t.Log("Check loading network files.")
	{
		l, err := LoadNetworks("testdata", "network.toml")
		if err != nil {
			t.Fatal(err)
		}
		if len(l) != 1 {
			t.Fatalf("network files load problem")
		}
		if l[0].String() != testNetwork.String() {
			t.Errorf("network file decode mismatch: [\n%s\n]", SimpleDiff(l[0].String(), testNetwork.String()))
		}
	}
}
