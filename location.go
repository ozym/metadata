package metadata

import (
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
