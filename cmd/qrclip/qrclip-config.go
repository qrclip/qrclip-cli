package main

import (
	"encoding/json"
	"errors"
	"github.com/kirsle/configdir"
	"os"
	"path/filepath"
)

// getQRClipConfigFilePath /////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getQRClipConfigFilePath() string {
	tConfigPath := configdir.LocalConfig("qrclip")
	tErr := configdir.MakePath(tConfigPath)
	if tErr != nil {
		ExitWithError("Creating config file path " + tConfigPath)
	}
	return filepath.Join(tConfigPath, "qrclip.json")
}

// SaveQRClipConfigFile ////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SaveQRClipConfigFile(pQRClipConfigDto QRClipConfigDto) {
	tConfigFile := getQRClipConfigFilePath()
	tFh, tErr := os.Create(tConfigFile)
	if tErr != nil {
		ExitWithError("Creating config file at " + tConfigFile)
	}
	defer tFh.Close()

	tEncoder := json.NewEncoder(tFh)
	tErr = tEncoder.Encode(&pQRClipConfigDto)
	if tErr != nil {
		ExitWithError("Error saving config")
	}

}

// GetQRClipConfig /////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GetQRClipConfig() (QRClipConfigDto, error) {
	tConfigFile := getQRClipConfigFilePath()
	var tQRClipConfigDto QRClipConfigDto
	if _, tErr := os.Stat(tConfigFile); os.IsNotExist(tErr) {
		return tQRClipConfigDto, errors.New("NO CONFIG FOUND")
	} else {
		tFile, tErr := os.Open(tConfigFile)
		if tErr != nil {
			ExitWithError("Opening config file at " + tConfigFile)
		}

		tDecoder := json.NewDecoder(tFile)
		tErr = tDecoder.Decode(&tQRClipConfigDto)
		if tErr != nil {
			ExitWithError("Error decoding config file!")
		}

		return tQRClipConfigDto, nil
	}
}
