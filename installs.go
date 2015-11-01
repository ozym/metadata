package metadata

import (
	"time"
)

type RadioInstall struct {
	Location  string  `csv:"Radio Location",`
	Target    string  `csv:"Radio Target Location",`
	Role      string  `csv:"Radio Role",`
	Model     string  `csv:"Radio Model",`
	Serial    string  `csv:"Radio Serial Number",`
	Polarity  string  `csv:"Antenna Polarity",`
	Frequency float64 `csv:"Frequency Key",`
}

type Asset struct {
	Model  string `csv:"Model Name",`
	Serial string `csv:"Serial Number",`
	Asset  string `csv:"Asset Number",`
}

type EquipmentInstall struct {
	Location string    `csv:"Equipment Location",`
	Model    string    `csv:"Equipment Model",`
	Serial   string    `csv:"Equipment Serial Number",`
	Start    time.Time `csv:"Installation Start",`
	Stop     time.Time `csv:"Installation Stop",`
}

type SensorInstall struct {
	Station string    `csv:"Seismic Station",`
	Site    string    `csv:"Sensor Location",`
	Model   string    `csv:"Sensor Model",`
	Serial  string    `csv:"Sensor Serial Number",`
	Azimuth float64   `csv:"Azimuth",`
	Dip     float64   `csv:"Dip",`
	Depth   float64   `csv:"Depth",`
	Start   time.Time `csv:"Installation Start",`
	Stop    time.Time `csv:"Installation Stop",`
}

type DataloggerInstall struct {
	Station string    `csv:"Seismic Station",`
	Site    string    `csv:"Datalogger Location",`
	Model   string    `csv:"Datalogger Model",`
	Serial  string    `csv:"Datalogger Serial Number",`
	Start   time.Time `csv:"Installation Start",`
	Stop    time.Time `csv:"Installation Stop",`
}

type AssetList []Asset
type RadioInstalls []RadioInstall
type EquipmentInstalls []EquipmentInstall
type SensorInstalls []SensorInstall
type DataloggerInstalls []DataloggerInstall

func (a AssetList) List()          {}
func (r RadioInstalls) List()      {}
func (e EquipmentInstalls) List()  {}
func (s SensorInstalls) List()     {}
func (d DataloggerInstalls) List() {}
