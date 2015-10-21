package metadata

import (
	"testing"
)

func TestNet_IPAddress(t *testing.T) {

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

func TestNet_IPNetwork(t *testing.T) {

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
