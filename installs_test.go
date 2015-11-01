package metadata

import (
	"io/ioutil"
	"testing"
)

var testSensorInstalls SensorInstalls
var testDataloggerInstalls DataloggerInstalls
var testEquipmentInstalls EquipmentInstalls

func init() {

	testSensorInstalls = SensorInstalls{
		SensorInstall{
			Station: "ABCD",
			Site:    "10",
			Model:   "Model",
			Serial:  "Serial #1",
			Azimuth: 10.0,
			Dip:     10.0,
			Depth:   10.0,
			Start:   MustParseTime("2010-01-01T00:00:00Z"),
			Stop:    MustParseTime("2011-01-01T00:00:00Z"),
		},
		SensorInstall{
			Station: "ABCD",
			Site:    "20",
			Model:   "Model",
			Serial:  "Serial #2",
			Azimuth: 20.0,
			Dip:     20.0,
			Depth:   20.0,
			Start:   MustParseTime("2010-01-01T00:00:00Z"),
			Stop:    MustParseTime("9999-01-01T00:00:00Z"),
		},
		SensorInstall{
			Station: "EFGH",
			Site:    "10",
			Model:   "Model",
			Serial:  "Serial #3",
			Azimuth: 10.0,
			Dip:     10.0,
			Depth:   10.0,
			Start:   MustParseTime("2010-01-01T00:00:00Z"),
			Stop:    MustParseTime("2011-01-01T00:00:00Z"),
		},
		SensorInstall{
			Station: "EFGH",
			Site:    "20",
			Model:   "Model",
			Serial:  "Serial #4",
			Azimuth: 20.0,
			Dip:     20.0,
			Depth:   20.0,
			Start:   MustParseTime("2010-01-01T00:00:00Z"),
			Stop:    MustParseTime("2011-01-01T00:00:00Z"),
		},
		SensorInstall{
			Station: "EFGH",
			Site:    "20",
			Model:   "Model",
			Serial:  "Serial #5",
			Azimuth: 20.0,
			Dip:     20.0,
			Depth:   20.0,
			Start:   MustParseTime("2012-01-01T00:00:00Z"),
			Stop:    MustParseTime("2013-01-01T00:00:00Z"),
		},
	}

	testDataloggerInstalls = DataloggerInstalls{
		DataloggerInstall{
			Station: "ABCD",
			Site:    "01",
			Model:   "Model",
			Serial:  "Serial #1",
			Start:   MustParseTime("2010-01-01T00:00:00Z"),
			Stop:    MustParseTime("2011-01-01T00:00:00Z"),
		},
		DataloggerInstall{
			Station: "EFGH",
			Site:    "02",
			Model:   "Model",
			Serial:  "Serial #2",
			Start:   MustParseTime("2010-01-01T00:00:00Z"),
			Stop:    MustParseTime("9999-01-01T00:00:00Z"),
		},
	}

	testEquipmentInstalls = EquipmentInstalls{
		EquipmentInstall{
			Location: "Somewhere",
			Model:    "Model #1",
			Serial:   "Serial #1",
			Start:    MustParseTime("2010-01-01T00:00:00Z"),
			Stop:     MustParseTime("2011-01-01T00:00:00Z"),
		},
		EquipmentInstall{
			Location: "Somewhere",
			Model:    "Model #2",
			Serial:   "Serial #2",
			Start:    MustParseTime("2010-01-01T00:00:00Z"),
			Stop:     MustParseTime("2011-01-01T00:00:00Z"),
		},
		EquipmentInstall{
			Location: "Somewhere Else",
			Model:    "Model #2",
			Serial:   "Serial #2",
			Start:    MustParseTime("2012-01-01T00:00:00Z"),
			Stop:     MustParseTime("2013-01-01T00:00:00Z"),
		},
		EquipmentInstall{
			Location: "Somewhere Else",
			Model:    "Model #1",
			Serial:   "Serial #1",
			Start:    MustParseTime("2012-01-01T00:00:00Z"),
			Stop:     MustParseTime("2013-01-01T00:00:00Z"),
		},
	}
}

func TestSensorInstalls_ReadFile(t *testing.T) {
	t.Log("Compare loaded sensor installs file.")
	{
		b, err := ioutil.ReadFile("testdata/sensors.csv")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != Strings(testSensorInstalls) {
			t.Errorf("sensor installs file text mismatch: [\n%s\n]", SimpleDiff(string(b), Strings(testSensorInstalls)))
		}
	}
}

func TestSensorInstalls_LoadFile(t *testing.T) {
	t.Log("Check loading sensor installs file.")
	{
		var installs SensorInstalls
		if err := LoadList("testdata/sensors.csv", &installs); err != nil {
			t.Fatal(err)
		}
		if Strings(installs) != Strings(testSensorInstalls) {
			t.Errorf("sensor installs file decode mismatch: [\n%s\n]", SimpleDiff(Strings(installs), Strings(testSensorInstalls)))
		}
	}
}

func TestSensorInstalls_LoadFiles(t *testing.T) {
	t.Log("Check loading sensor installs files.")
	{
		var installs SensorInstalls
		if err := LoadLists("testdata", "sensors.csv", &installs); err != nil {
			t.Fatal(err)
		}
		if Strings(installs) != Strings(testSensorInstalls) {
			t.Errorf("sensor installs file decode mismatch: [\n%s\n]", SimpleDiff(Strings(installs), Strings(testSensorInstalls)))
		}
	}
}

func TestDataloggerInstalls_ReadFile(t *testing.T) {
	t.Log("Compare loaded datalogger installs file.")
	{
		b, err := ioutil.ReadFile("testdata/dataloggers.csv")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != Strings(testDataloggerInstalls) {
			t.Errorf("datalogger installs file text mismatch: [\n%s\n]", SimpleDiff(string(b), Strings(testDataloggerInstalls)))
		}
	}
}

func TestDataloggerInstalls_LoadFile(t *testing.T) {
	t.Log("Check loading datalogger installs file.")
	{
		var installs DataloggerInstalls
		if err := LoadList("testdata/dataloggers.csv", &installs); err != nil {
			t.Fatal(err)
		}
		if Strings(installs) != Strings(testDataloggerInstalls) {
			t.Errorf("datalogger installs file decode mismatch: [\n%s\n]", SimpleDiff(Strings(installs), Strings(testDataloggerInstalls)))
		}
	}
}

func TestDataloggerInstalls_LoadFiles(t *testing.T) {
	t.Log("Check loading datalogger installs files.")
	{
		var installs DataloggerInstalls
		if err := LoadLists("testdata", "dataloggers.csv", &installs); err != nil {
			t.Fatal(err)
		}
		if Strings(installs) != Strings(testDataloggerInstalls) {
			t.Errorf("datalogger installs file decode mismatch: [\n%s\n]", SimpleDiff(Strings(installs), Strings(testDataloggerInstalls)))
		}
	}
}

func TestEquipmentInstalls_ReadFile(t *testing.T) {
	t.Log("Compare loaded equipment installs file.")
	{
		b, err := ioutil.ReadFile("testdata/equipment.csv")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != Strings(testEquipmentInstalls) {
			t.Errorf("equipment installs file text mismatch: [\n%s\n]", SimpleDiff(string(b), Strings(testEquipmentInstalls)))
		}
	}
}

func TestEquipmentInstalls_LoadFile(t *testing.T) {
	t.Log("Check loading equipment installs file.")
	{
		var installs EquipmentInstalls
		if err := LoadList("testdata/equipment.csv", &installs); err != nil {
			t.Fatal(err)
		}
		if Strings(installs) != Strings(testEquipmentInstalls) {
			t.Errorf("equipment installs file decode mismatch: [\n%s\n]", SimpleDiff(Strings(installs), Strings(testEquipmentInstalls)))
		}
	}
}

func TestEquipmentInstalls_LoadFiles(t *testing.T) {
	t.Log("Check loading equipment installs files.")
	{
		var installs EquipmentInstalls
		if err := LoadLists("testdata", "equipment.csv", &installs); err != nil {
			t.Fatal(err)
		}
		if Strings(installs) != Strings(testEquipmentInstalls) {
			t.Errorf("equipment installs file decode mismatch: [\n%s\n]", SimpleDiff(Strings(installs), Strings(testEquipmentInstalls)))
		}
	}
}
