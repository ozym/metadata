package metadata

import (
	"encoding/json"
	"net"

	"github.com/BurntSushi/toml"
)

type address struct {
	net.IP
}

func (a *address) UnmarshalText(text []byte) error {

	aa, _, err := net.ParseCIDR(string(text))
	if err != nil {
		return nil
	}
	*a = address{aa}

	return nil
}

type network struct {
	net.IPNet
}

func (n *network) UnmarshalText(text []byte) error {

	_, nn, err := net.ParseCIDR(string(text))
	if err != nil {
		return nil
	}
	*n = network{*nn}

	return nil
}

type Linknet struct {
	Name string `json:"name"`
}

type Equipment struct {
	Name        string    `json:"name"`
	Model       string    `json:"model"`
	Addresses   []address `json:"addresses,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Code        *string   `json:"code,omitempty"`
	Uninstalled bool      `json:"uninstalled,omitempty"`
}

type Location struct {
	Tag       string               `json:"tag"`
	Name      string               `json:"name"`
	Runnet    *network             `json:"runnet,omitempty"`
	Locnet    bool                 `json:"locnet,omitempty"`
	Linknets  []Linknet            `json:"linknets,omitempty"`
	Equipment map[string]Equipment `json:"equipment,omitempty"`
}

func (l Location) String() string {
	j, _ := json.Marshal(l)
	return string(j)
}

func LoadLocation(file string) (Location, error) {
	var l Location

	if _, err := toml.DecodeFile(file, &l); err != nil {
		return l, err
	}

	return l, nil
}
