package metadata

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type Linknet struct {
	Name string `json:"name,omitempty" comment:"Linknet name."`
}

type Located struct {
	Name        string      `json:"name" comment:"The unique equipment name."`
	Model       string      `json:"model" comment:"A generic equipment model name."`
	Address     *IPAddress  `json:"address,omitempty" comment:"The primary equipment ip address used to connect."`
	Aliases     []IPAddress `json:"aliases,omitempty" comment:"Any extra ip addresses assigned to the equipment."`
	Tags        []string    `json:"tags,omitempty" comment:"Any extra equipment tags or labels."`
	Notes       *string     `json:"notes,omitempty" comment:"Any equipment specific notes or comments."`
	Uninstalled *bool       `json:"uninstalled,omitempty" comment:"Indicate whether the equipment is not present."`
}

type Location struct {
	Tag       string             `json:"tag" comment:"Location specific tag."`
	Name      string             `json:"name" comment:"Location place name."`
	Latitude  *float32           `json:"latitude,omitempty" comment:"Optional location latitude."`
	Longitude *float32           `json:"longitude,omitempty" comment:"Optional location longitude."`
	Runnet    *IPNetwork         `json:"runnet,omitempty" comment:"Optional location runnet."`
	Locnet    *bool              `json:"locnet,omitempty" comment:"Should a locnet be required."`
	Linknets  []Linknet          `json:"linknets,omitempty" comment:"A list of required location linknets."`
	Equipment map[string]Located `json:"equipment,omitempty" comment:"The equipment installed at the location."`
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

	l = append(l, "## The site specific single word tag. This is used for DNS entries related to the site and for equipment linkages.")
	l = append(l, fmt.Sprintf("tag = %s", strconv.Quote(loc.Tag)))
	l = append(l, "")
	l = append(l, "## The name of the location, which is used for labels and general identification, punctuation should be avoided.")
	l = append(l, fmt.Sprintf("name = %s", strconv.Quote(loc.Name)))
	l = append(l, "")
	l = append(l, "## The site specific IP 192.168.X.Y/28 equipment range, which also provides a site specific supernet 10.X.Y.0/20.")
	l = append(l, "## which are generally split into /28 linknets, with the first being defined as the local network (locnet).")
	if loc.Runnet != nil {
		l = append(l, fmt.Sprintf("runnet = %s", strconv.Quote(loc.Runnet.String())))
	} else {
		l = append(l, fmt.Sprintf("#runnet = %s", strconv.Quote("192.168.0.0/28")))
	}
	l = append(l, "")
	l = append(l, "## Should a locknet be defined, useful for network planning and network link association.")
	if loc.Locnet != nil {
		l = append(l, fmt.Sprintf("locnet = %s", strconv.FormatBool(*loc.Locnet)))
	} else {
		l = append(l, fmt.Sprintf("#locnet = true|false"))
	}
	l = append(l, "")
	l = append(l, "## An array of /28 linking networks, the order dictates the network offset. Linknets with no names can be assigned,")
	l = append(l, "## these will be skipped, although they will still be used as offset placeholders.")
	l = append(l, "")
	l = append(l, "#[[linknets]]")
	l = append(l, "#\t## The name of the link, usually of the form \"Remote Site to Local Site\"")
	l = append(l, fmt.Sprintf("#\tname=%s", strconv.Quote("")))
	for i := 0; i < len(loc.Linknets); i++ {
		l = append(l, "")
		l = append(l, "[[linknets]]")
		l = append(l, "\t## The name of the link, usually of the form \"Remote Site to Local Site\"")
		l = append(l, fmt.Sprintf("\tname = %s", strconv.Quote(loc.Linknets[i].Name)))
	}
	l = append(l, "")

	l = append(l, "## A list of location equipment.")
	l = append(l, "")
	l = append(l, "#[equipment.label]")
	l = append(l, "#\t## The name of the equipment, generally a equipment tag plus the site location tag.")
	l = append(l, fmt.Sprintf("#\tname = %s", strconv.Quote("")))
	l = append(l, "#")
	l = append(l, "#\t## The model name, a generic term useful for monitoring or configuration.")
	l = append(l, fmt.Sprintf("#\t#model = %s", strconv.Quote("")))
	l = append(l, "#")
	l = append(l, "#\t## The assigned primary IP address of the equipment, this is expected to be returned for any DNS name lookups.")
	l = append(l, fmt.Sprintf("#\t#address = %s", strconv.Quote("")))
	l = append(l, "#")
	l = append(l, "#\t## An array of optional extra addresses associated with this equipment, these will provide the equipment name.")
	l = append(l, fmt.Sprintf("#\t#aliases = []"))
	l = append(l, "#")
	l = append(l, "#\t## Equipment code, useful for monitoring or configuration")
	l = append(l, fmt.Sprintf("#\t#code = %s", strconv.Quote("")))
	l = append(l, "#")
	l = append(l, "#\t## A string array of optional extra tags associated with this equipment.")
	l = append(l, fmt.Sprintf("#\t#tags = []"))
	l = append(l, "#")
	l = append(l, "#\t## Equipment specific notes and documentation")
	l = append(l, fmt.Sprintf("#\t#notes = \"\"\"\\\n#\t#\t\\n\\\n#\t#\t\"\"\""))
	l = append(l, "#")
	l = append(l, "#\t## Whether the equipment is active, the default is false which represents an installed device.")
	l = append(l, fmt.Sprintf("#\t#uninstalled = %s", strconv.FormatBool(false)))
	l = append(l, "")

	for e, equipment := range loc.Equipment {
		p := strings.Split(e, "-")
		if len(p) != 2 {
			continue
		}
		l = append(l, fmt.Sprintf("[equipment.%s]", p[0]))
		l = append(l, "")
		l = append(l, "\t## The name of the equipment, generally a equipment tag plus the site location tag.")
		l = append(l, fmt.Sprintf("\tname = %s", strconv.Quote(e)))
		l = append(l, "")
		l = append(l, "\t## The model name, a generic term useful for monitoring or configuration.")
		l = append(l, fmt.Sprintf("\tmodel = %s", strconv.Quote(equipment.Model)))
		l = append(l, "")
		l = append(l, "\t## The assigned primary IP address of the equipment, this is expected to be returned for any DNS name lookups.")

		if equipment.Address != nil {
			l = append(l, fmt.Sprintf("\taddress = %s", strconv.Quote(equipment.Address.String())))
		} else {
			l = append(l, fmt.Sprintf("\t#address = %s", strconv.Quote("")))
		}
		l = append(l, "")
		l = append(l, "\t## An array of optional extra addresses associated with this equipment, these will provide the equipment name.")
		var aliases []string
		for _, a := range equipment.Aliases {
			aliases = append(aliases, a.String())
		}
		if len(aliases) > 0 {
			l = append(l, fmt.Sprintf("\taliases = [\n\t\t%s\n\t]", strings.Join(aliases, ",\n\t\t")))
		} else {
			l = append(l, fmt.Sprintf("\t#aliases = []"))
		}
		l = append(l, "")
		l = append(l, "\t## An array of optional extra addresses associated with this equipment, these will provide the equipment name.")
		if len(equipment.Tags) > 0 {
			l = append(l, fmt.Sprintf("\ttags = [\n\t\t\"%s\"\n\t]", strings.Join(equipment.Tags, "\",\n\t\t\"")))
		} else {
			l = append(l, fmt.Sprintf("\t#tags = []"))
		}
		l = append(l, "")
		l = append(l, fmt.Sprintf("\t## Equipment specific notes and documentation"))
		if equipment.Notes != nil {
			n := strings.Split(strings.Replace(strings.TrimSpace(*equipment.Notes), "\\n", "\n", -1), "\n")
			l = append(l, fmt.Sprintf("\tnotes = \"\"\"\\\n\t\t%s\\\n\t\t\"\"\"", strings.Join(n, "\\n\\\n\t\t")))
		} else {
			l = append(l, fmt.Sprintf("\t#notes = \"\"\"\\\n\t#\t\\n\\\n\t#\t\"\"\""))
		}
		l = append(l, "")
		l = append(l, fmt.Sprintf("\t## Whether the equipment is active, the default is false which represents an installed device."))
		if equipment.Uninstalled != nil {
			l = append(l, fmt.Sprintf("\tuninstalled = %s", strconv.FormatBool(*equipment.Uninstalled)))
		} else {
			l = append(l, fmt.Sprintf("\t#uninstalled = %s", strconv.FormatBool(false)))
		}
		l = append(l, "")
	}

	return strings.Join(l, "\n")
}
