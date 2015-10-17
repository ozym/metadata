package metadata

import (
	"encoding/json"
	"net"
	"reflect"

	"github.com/BurntSushi/toml"
)

type ipaddress struct {
	net.IP
}

func (a *ipaddress) UnmarshalText(text []byte) error {

	aa, _, err := net.ParseCIDR(string(text))
	if err != nil {
		return nil
	}
	*a = ipaddress{aa}

	return nil
}

func ParseIPAddress(cidr string) *ipaddress {
	a, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}
	return &ipaddress{a}
}

type ipnetwork struct {
	net.IPNet
}

func ParseIPNetwork(cidr string) *ipnetwork {
	_, nn, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}
	return &ipnetwork{*nn}
}

func (n *ipnetwork) UnmarshalText(text []byte) error {

	_, nn, err := net.ParseCIDR(string(text))
	if err != nil {
		return nil
	}
	*n = ipnetwork{*nn}

	return nil
}

type Linknet struct {
	Name string `json:"name"`
}

type Equipment struct {
	Name        string      `json:"name"`
	Model       string      `json:"model"`
	Address     *ipaddress  `json:"address,omitempty"`
	Addresses   []ipaddress `json:"addresses,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	Code        *string     `json:"code,omitempty"`
	Uninstalled bool        `json:"uninstalled,omitempty"`
}

type Location struct {
	Tag       string               `json:"tag"`
	Name      string               `json:"name"`
	Runnet    *ipnetwork           `json:"runnet,omitempty"`
	Locnet    bool                 `json:"locnet,omitempty"`
	Linknets  []Linknet            `json:"linknets,omitempty"`
	Equipment map[string]Equipment `json:"equipment,omitempty"`
}

func (l Location) String() string {
	j, _ := json.MarshalIndent(l, "", "\t")
	return string(j)
}
func (l Location) Equal(location Location) bool {
	return reflect.DeepEqual(l, location)
}

func LoadLocation(file string) (Location, error) {
	var l Location

	if _, err := toml.DecodeFile(file, &l); err != nil {
		return l, err
	}

	return l, nil
}
