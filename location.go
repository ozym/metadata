package metadata

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type Linknet struct {
	Name string `json:"name,omitempty"`
}

type Device struct {
	Name        string      `json:"name"`
	Model       string      `json:"model"`
	Address     *IPAddress  `json:"address,omitempty"`
	Aliases     []IPAddress `json:"aliases,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	Notes       *string     `json:"notes,omitempty"`
	Uninstalled *bool       `json:"uninstalled,omitempty"`
}

type Location struct {
	Tag       string            `json:"tag"`
	Name      string            `json:"name"`
	Latitude  *float32          `json:"latitude,omitempty"`
	Longitude *float32          `json:"longitude,omitempty"`
	Runnet    *IPNetwork        `json:"runnet,omitempty"`
	Locnet    *bool             `json:"locnet,omitempty"`
	Linknets  []Linknet         `json:"linknets,omitempty" toml:"linknet"`
	Devices   map[string]Device `json:"devices,omitempty" toml:"device"`
}

func LoadLocation(filename string) (*Location, error) {
	var l Location

	if _, err := toml.DecodeFile(filename, &l); err != nil {
		return nil, err
	}

	return &l, nil
}

func (loc Location) StoreLocation(base string) error {

	// per site directory
	dir := base + "/" + loc.Tag
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// output file name
	filename := dir + "/" + "location.toml"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(loc.String()))
	if err != nil {
		return err
	}

	return err
}

func (loc Location) String() string {

	var l []string

	l = append(l, "## The unique site specific single word tag.")
	l = append(l, fmt.Sprintf("tag = %s", strconv.Quote(loc.Tag)))
	l = append(l, "")
	l = append(l, "## The general name of the location.")
	l = append(l, fmt.Sprintf("name = %s", strconv.Quote(loc.Name)))
	l = append(l, "")
	l = append(l, "## Optional site geographical position.")
	if loc.Latitude != nil {
		l = append(l, fmt.Sprintf("latitude = %.4f", *loc.Latitude))
	} else {
		l = append(l, fmt.Sprintf("#latitude = degrees"))
	}
	if loc.Longitude != nil {
		l = append(l, fmt.Sprintf("longitude = %.4f", *loc.Longitude))
	} else {
		l = append(l, fmt.Sprintf("#longitude = degrees"))
	}
	l = append(l, "")
	l = append(l, "## An optional site specific IP 192.168.X.Y/28 equipment range.")
	if loc.Runnet != nil {
		l = append(l, fmt.Sprintf("runnet = %s", strconv.Quote(loc.Runnet.String())))
	} else {
		l = append(l, fmt.Sprintf("#runnet = %s", strconv.Quote("192.168.0.0/28")))
	}
	l = append(l, "")
	l = append(l, "## Should a local IP 10.X.Y.0/28 range be assigned based on the runnet.")
	if loc.Locnet != nil {
		l = append(l, fmt.Sprintf("locnet = %s", strconv.FormatBool(*loc.Locnet)))
	} else {
		l = append(l, fmt.Sprintf("#locnet = true|false"))
	}
	l = append(l, "")
	l = append(l, "## An array of 10.X.Y.N/28 linking networks, the order dictates the network offset.")
	l = append(l, "")
	l = append(l, "#[[linknet]]")
	l = append(l, "#\t## The name of the link, usually of the form \"Remote Site to Local Site\".")
	l = append(l, fmt.Sprintf("#\tname=%s", strconv.Quote("")))
	for i := 0; i < len(loc.Linknets); i++ {
		l = append(l, "")
		l = append(l, "[[linknet]]")
		l = append(l, "\t## The name of the link, usually of the form \"Remote Site to Local Site\".")
		l = append(l, fmt.Sprintf("\tname = %s", strconv.Quote(loc.Linknets[i].Name)))
	}
	l = append(l, "")

	l = append(l, "## A list of local devices.")
	l = append(l, "")
	l = append(l, "#[device.label]")
	l = append(l, "#\t## The name of the device, generally an equipment tag plus the site location tag.")
	l = append(l, fmt.Sprintf("#\tname = %s", strconv.Quote("")))
	l = append(l, "#")
	l = append(l, "#\t## The model name, a generic term useful for monitoring or configuration.")
	l = append(l, fmt.Sprintf("#\t#model = %s", strconv.Quote("")))
	l = append(l, "#")
	l = append(l, "#\t## The assigned primary IP address of the device.")
	l = append(l, fmt.Sprintf("#\t#address = %s", strconv.Quote("")))
	l = append(l, "#")
	l = append(l, "#\t## An array of extra addresses associated with this device.")
	l = append(l, fmt.Sprintf("#\t#aliases = []"))
	l = append(l, "#")
	l = append(l, "#\t## An array of extra tags associated with this device.")
	l = append(l, fmt.Sprintf("#\t#tags = []"))
	l = append(l, "#")
	l = append(l, "#\t## Optional device specific notes and documentation.")
	l = append(l, fmt.Sprintf("#\t#notes = \"\"\"\\\n#\t#\t\\n\\\n#\t#\t\"\"\""))
	l = append(l, "#")
	l = append(l, "#\t## Whether the device is not currently installed or active.")
	l = append(l, fmt.Sprintf("#\t#uninstalled = %s", strconv.FormatBool(false)))

	var keys Keys
	for d, _ := range loc.Devices {
		keys = append(keys, d)
	}
	sort.Sort(keys)

	for _, d := range keys {
		device, ok := loc.Devices[d]
		if !ok {
			continue
		}
		l = append(l, "")
		l = append(l, fmt.Sprintf("[device.%s]", d))
		l = append(l, "\t## The name of the device, generally an equipment tag plus the site location tag.")
		l = append(l, fmt.Sprintf("\tname = %s", strconv.Quote(device.Name)))
		l = append(l, "")
		l = append(l, "\t## The model name, a generic term useful for monitoring or configuration.")
		l = append(l, fmt.Sprintf("\tmodel = %s", strconv.Quote(device.Model)))
		l = append(l, "")
		l = append(l, "\t## The assigned primary IP address of the device.")

		if device.Address != nil {
			l = append(l, fmt.Sprintf("\taddress = %s", strconv.Quote(device.Address.String())))
		} else {
			l = append(l, fmt.Sprintf("\t#address = %s", strconv.Quote("")))
		}
		l = append(l, "")
		l = append(l, "\t## An array of extra addresses associated with this device.")
		var aliases []string
		for _, a := range device.Aliases {
			aliases = append(aliases, strconv.Quote(a.String()))
		}
		if len(aliases) > 0 {
			l = append(l, fmt.Sprintf("\taliases = [\n\t\t\t%s\n\t\t]", strings.Join(aliases, ",\n\t\t\t")))
		} else {
			l = append(l, fmt.Sprintf("\t#aliases = []"))
		}
		l = append(l, "")
		l = append(l, "\t## An array of extra tags associated with this device.")
		if len(device.Tags) > 0 {
			l = append(l, fmt.Sprintf("\ttags = [\n\t\t\t\"%s\"\n\t\t]", strings.Join(device.Tags, "\",\n\t\t\t\"")))
		} else {
			l = append(l, fmt.Sprintf("\t#tags = []"))
		}
		l = append(l, "")
		l = append(l, "\t## Optional device specific notes and documentation.")
		if device.Notes != nil {
			n := strings.Split(strings.Replace(strings.TrimSpace(*device.Notes), "\\n", "\n", -1), "\n")
			l = append(l, fmt.Sprintf("\tnotes = \"\"\"\\\n\t\t%s\\n\\\n\t\t\"\"\"", strings.Join(n, "\\n\\\n\t\t")))
		} else {
			l = append(l, fmt.Sprintf("\t#notes = \"\"\"\\\n\t#\t\\n\\\n\t#\t\"\"\""))
		}
		l = append(l, "")
		l = append(l, "\t## Whether the device is not currently installed or active.")
		if device.Uninstalled != nil {
			l = append(l, fmt.Sprintf("\tuninstalled = %s", strconv.FormatBool(*device.Uninstalled)))
		} else {
			l = append(l, fmt.Sprintf("\t#uninstalled = %s", strconv.FormatBool(false)))
		}
	}
	l = append(l, "")
	l = append(l, "# "+"vim:"+" tabstop=4 expandtab shiftwidth=4 softtabstop=4")
	l = append(l, "")

	return strings.Replace(strings.Join(l, "\n"), "\t", "    ", -1)
}
