package metadata

import (
	"io/ioutil"
	"testing"
)

var testAssetList AssetList

func init() {

	testAssetList = AssetList{
		Asset{
			Model:  "Model #1",
			Serial: "Serial #1",
			Asset:  "Asset #1",
		},
		Asset{
			Model:  "Model #2",
			Serial: "Serial #2",
			Asset:  "Asset #2",
		},
	}
}

func TestAssetList_ReadFile(t *testing.T) {
	t.Log("Compare loaded asset list file.")
	{
		b, err := ioutil.ReadFile("testdata/assets.csv")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != Strings(testAssetList) {
			t.Errorf("asset list file text mismatch: [\n%s\n]", SimpleDiff(string(b), Strings(testAssetList)))
		}
	}
}

func TestAssetList_LoadFile(t *testing.T) {
	t.Log("Check loading asset list file.")
	{
		var list AssetList
		if err := LoadList("testdata/assets.csv", &list); err != nil {
			t.Fatal(err)
		}
		if Strings(list) != Strings(testAssetList) {
			t.Errorf("asset list file decode mismatch: [\n%s\n]", SimpleDiff(Strings(list), Strings(testAssetList)))
		}
	}
}

func TestAssetList_LoadFiles(t *testing.T) {
	t.Log("Check loading asset list files.")
	{
		var list AssetList
		if err := LoadLists("testdata", "assets.csv", &list); err != nil {
			t.Fatal(err)
		}
		if Strings(list) != Strings(testAssetList) {
			t.Errorf("asset list file decode mismatch: [\n%s\n]", SimpleDiff(Strings(list), Strings(testAssetList)))
		}
	}
}
