package metadata

import (
	"net"
)

type IPAddress struct {
	net.IPNet
}

func (a *IPAddress) UnmarshalText(text []byte) error {

	aa, nn, err := net.ParseCIDR(string(text))
	if err != nil {
		return nil
	}
	*a = IPAddress{net.IPNet{IP: aa, Mask: nn.Mask}}

	return nil
}
func (a IPAddress) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}
func ParseIPAddress(cidr string) (*IPAddress, error) {
	a, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return &IPAddress{net.IPNet{IP: a, Mask: n.Mask}}, nil
}

func MustParseIPAddress(cidr string) *IPAddress {
	a, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return &IPAddress{}
	}
	return &IPAddress{net.IPNet{IP: a, Mask: n.Mask}}
}

type IPNetwork struct {
	net.IPNet
}

func (n *IPNetwork) UnmarshalText(text []byte) error {

	_, nn, err := net.ParseCIDR(string(text))
	if err != nil {
		return nil
	}
	*n = IPNetwork{*nn}

	return nil
}
func (n IPNetwork) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

func ParseIPNetwork(cidr string) (*IPNetwork, error) {
	_, nn, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return &IPNetwork{*nn}, nil
}

func MustParseIPNetwork(cidr string) *IPNetwork {
	_, nn, err := net.ParseCIDR(cidr)
	if err != nil {
		return &IPNetwork{}
	}
	return &IPNetwork{*nn}
}
