package metadata

import (
	"bytes"
	"net"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
)

const networkTemplate = `# Network and device IP address information.

## The network location ID tag.
location = "{{.Location}}"

## Name of the network, defaults to location name.
{{if .Name}}name = "{{.Name}}"{{else}}#name = ""{{end}}

## Network notes and documentation.
{{if .Notes}}notes = """\
{{$lines := Lines .Notes}}{{range $k, $v := $lines}}    {{$v}}\n\
{{end}}    """{{else}}#notes = """\
#    \n\
#    """{{end}}

## Site specific IP 192.168.X.Y/28 equipment range.
{{if .Runnet}}runnet = "{{.Runnet}}"{{else}}#runnet = ""{{end}}

## An array of 10.X.Y.N/28 linking networks, the order dictates the network offset.

#[[linknet]]
#    ## The name of the link, usually of the form "Remote Site to Local Site".
#    name = ""{{range .Linknets}}

[[linknet]]
    ## The name of the link, usually of the form "Remote Site to Local Site".
    name = "{{.Name}}"{{end}}

## Local devices.

#[[device]]
#    ## Device name, generally an equipment tag plus the site network tag.
#    name = ""
#
#    ## Model name, a generic term useful for monitoring or configuration.
#    model = ""
#
#    ## Primary IP address of the device.
#    #address = ""
#
#    ## Extra addresses associated with this device.
#    #aliases = []
#
#    ## Extra tags associated with this device.
#    #tags = []
#
#    ## Linked devices.
#    #links = []
#
#    ## Device specific notes and documentation.
#    #notes = """\
#    #    \n\
#    #    """
#
#    ## Is the device uninstalled or unreachable.
#    #uninstalled = false{{range .Devices}}

[[device]]
    ## Device name, generally an equipment tag plus the site network tag.
    name = "{{.Name}}"

    ## Model name, a generic term useful for monitoring or configuration.
    model = "{{.Model}}"

    ## Primary IP address of the device.
{{if .Address}}    address = "{{.Address}}"{{else}}    #address=""{{end}}

    ## Extra addresses associated with this device.
{{if .Aliases}}    aliases = [{{range $n, $t := .Aliases}}{{if gt $n 0}},{{end}}
        "{{$t}}"{{end}}
    ]{{else}}    #aliases = []{{end}}

    ## Extra tags associated with this device.
{{if .Tags}}    tags = [{{range $n, $t := .Tags}}{{if gt $n 0}},{{end}}
        "{{$t}}"{{end}}
    ]{{else}}    #tags = []{{end}}

    ## Linked devices.
{{if .Links}}    links = [{{range $n, $t := .Links}}{{if gt $n 0}},{{end}}
        "{{$t}}"{{end}}
    ]{{else}}    #links = []{{end}}

    ## Device specific notes and documentation.
{{if .Notes}}    notes = """\
{{$lines := Lines .Notes}}{{range $k, $v := $lines}}        {{$v}}\n\
{{end}}        """{{else}}    #notes = """\
    #    \n\
    #    """{{end}}

    ## Is the device uninstalled or unreachable.
{{if .Uninstalled}}    uninstalled = {{.Uninstalled}}{{else}}    #unistalled = true|false{{end}}{{end}}

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
`

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

type Linknet struct {
	Name string `json:"name,omitempty"`
}

type Device struct {
	Name        string      `json:"name"`
	Model       string      `json:"model"`
	Address     *IPAddress  `json:"address,omitempty"`
	Aliases     []IPAddress `json:"aliases,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	Links       []string    `json:"links,omitempty"`
	Notes       *string     `json:"notes,omitempty"`
	Uninstalled *bool       `json:"uninstalled,omitempty"`
}

type Network struct {
	Location string     `json:"location"`
	Name     *string    `json:"name,omitempty"`
	Notes    *string    `json:"notes,omitempty"`
	Runnet   *IPNetwork `json:"runnet,omitempty"`
	Linknets []Linknet  `json:"linknets,omitempty" toml:"linknet"`
	Devices  []Device   `json:"devices,omitempty" toml:"device"`
}

func LoadNetwork(filename string) (*Network, error) {
	var l Network

	if _, err := toml.DecodeFile(filename, &l); err != nil {
		return nil, err
	}

	return &l, nil
}

func LoadNetworks(dirname, filename string) ([]Network, error) {
	var ll []Network

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			l, e := LoadNetwork(path)
			if e != nil {
				return e
			}
			ll = append(ll, *l)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ll, nil
}

func (net Network) StoreNetwork(path string) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(net.String()))
	if err != nil {
		return err
	}

	return err
}

func (net Network) String() string {
	tplFuncMap := make(template.FuncMap)
	tplFuncMap["Lines"] = Lines
	tplFuncMap["LatLon"] = LatLon

	tmpl, err := template.New("").Funcs(tplFuncMap).Parse(networkTemplate)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, net)
	if err != nil {
		panic(err)
	}

	return doc.String()

}
