package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type QRCStorageLocation struct {
	Index int
	Code  string
	Name  string
}

// getLocations ////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getLocations() []QRCStorageLocation {
	var tLocations []QRCStorageLocation

	// FRANKFURT
	var tLoc0 QRCStorageLocation
	tLoc0.Index = 0
	tLoc0.Name = "Frankfurt (Europe)"
	tLoc0.Code = "storage05"
	tLocations = append(tLocations, tLoc0)

	// New York
	var tLoc1 QRCStorageLocation
	tLoc1.Index = 1
	tLoc1.Name = "New York (North America)"
	tLoc1.Code = "storage01"
	tLocations = append(tLocations, tLoc1)

	// San Francisco
	var tLoc2 QRCStorageLocation
	tLoc2.Index = 2
	tLoc2.Name = "San Francisco (North America)"
	tLoc2.Code = "storage04"
	tLocations = append(tLocations, tLoc2)

	// San Francisco
	var tLoc3 QRCStorageLocation
	tLoc3.Index = 3
	tLoc3.Name = "Singapore (Asia)"
	tLoc3.Code = "storage03"
	tLocations = append(tLocations, tLoc3)

	return tLocations
}

// SelectStorage ///////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SelectStorage() {
	tSelectedStorage := ""
	tConfig, tErr := GetQRClipConfig()
	if tErr == nil {
		tSelectedStorage = tConfig.Storage
	}

	tLocations := getLocations()

	ShowInfoCyan("Available Storage Locations")

	for _, tLocation := range tLocations {
		tLine := fmt.Sprintf("%d - %s", tLocation.Index, tLocation.Name)
		if tSelectedStorage == tLocation.Code {
			tLine = tLine + " <-----SELECTED"
			ShowSuccess(tLine)
		} else {
			fmt.Println(tLine)
		}
	}

	ShowInfoYellow("Select the storage:")
	tReader := bufio.NewReader(os.Stdin)
	tChar, _, tErr := tReader.ReadRune()
	if tErr != nil {
		ShowError(tErr.Error())
	}

	tStorageLocationFound := false
	for _, tLocation := range tLocations {
		if rune(strconv.Itoa(tLocation.Index)[0]) == tChar {
			tStorageLocationFound = true
			ShowSuccess(tLocation.Name)
			setStorageLocation(tLocation.Code)
		}
	}

	if !tStorageLocationFound {
		ExitWithError("Storage Location not found!")
	}
}

// setStorageLocation //////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func setStorageLocation(tCode string) {
	tConfig, tErr := GetQRClipConfig()
	if tErr != nil {
		var tConf QRClipConfigDto
		tConf.Storage = tCode
		tConfig = tConf
	}
	SaveQRClipConfigFile(tConfig)
}
