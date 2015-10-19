package metadata

import (
	//	"encoding/json"
	//	"fmt"
	"net"
	"reflect"
	//	"strconv"
	//"strings"
	"io"

	"github.com/BurntSushi/toml"
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

func ParseIPAddress(cidr string) *IPAddress {
	a, n, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}
	return &IPAddress{net.IPNet{IP: a, Mask: n.Mask}}
}

type IPNetwork struct {
	net.IPNet
}

func ParseIPNetwork(cidr string) *IPNetwork {
	_, nn, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}
	return &IPNetwork{*nn}
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

type Linknet struct {
	Name string `json:"name,omitempty" comment:"Linknet name."`
}

func (l Linknet) Encode(w io.Writer, prefix string) error {
	return EncodeField(w, l, prefix)
}
func (l Linknet) Default() string {
	return "default"
}

type Equipment struct {
	Name        string      `json:"name" comment:"The unique equipment name."`
	Model       string      `json:"model" comment:"A generic equipment model name."`
	Address     *IPAddress  `json:"address,omitempty" comment:"The primary equipment ip address used to connect."`
	Aliases     []IPAddress `json:"aliases,omitempty" comment:"Any extra ip addresses assigned to the equipment."`
	Code        *string     `json:"code,omitempty" comment:"Optional equipment code or label."`
	Tags        []string    `json:"tags,omitempty" comment:"Any extra equipment tags or labels."`
	Notes       *string     `json:"notes,omitempty" comment:"Any equipment specific notes or comments."`
	Uninstalled *bool       `json:"uninstalled,omitempty" comment:"Indicate whether the equipment is not present."`
}

func (e Equipment) Encode(w io.Writer, prefix string) error {
	return EncodeField(w, e, prefix)
}

type Location struct {
	Tag       string               `json:"tag" comment:"Location specific tag."`
	Name      string               `json:"name" comment:"Location place name."`
	Latitude  *float32             `json:"latitude,omitempty" comment:"Optional location latitude."`
	Longitude *float32             `json:"longitude,omitempty" comment:"Optional location longitude."`
	Runnet    *IPNetwork           `json:"runnet,omitempty" comment:"Optional location runnet."`
	Locnet    *bool                `json:"locnet,omitempty" comment:"Should a locnet be required."`
	Linknets  []Linknet            `json:"linknets,omitempty" comment:"A list of required location linknets."`
	Equipment map[string]Equipment `json:"equipment,omitempty" comment:"The equipment installed at the location."`
}

func (l Location) Encode(w io.Writer, prefix string) error {
	return EncodeField(w, l, prefix)
}

func (l Location) String() string {
	j, _ := EncodeString(l)
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
