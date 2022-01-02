package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// CheckLimits /////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func CheckLimits() {
	tLimits := GetLimits()
	printClipLimits(tLimits)
}

// GetLimits ///////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GetLimits() ClipLimitsDto {
	tJwt := CheckJwtToken()
	if tJwt == "" {
		return getClipLimits()
	} else {
		return getClipLimitsUser(tJwt)
	}
}

// printClipLimits /////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func printClipLimits(tClipLimitsDto ClipLimitsDto) {
	ShowSuccess("QRClip LIMITS")
	ShowSuccess(" Max Characters: " + strconv.Itoa(tClipLimitsDto.Text))
	ShowSuccess(" Max Expiration Minutes: " + strconv.Itoa(tClipLimitsDto.ExpiresInMinutes))
	ShowSuccess(" File Size(Mb): " + strconv.Itoa(tClipLimitsDto.FileMb))
	if tClipLimitsDto.MaxTransfers == 0 { // IF ZERO ITS UNLIMITED
		ShowSuccess(" Max Transfers: Unlimited")
	} else {
		ShowSuccess(" Max Transfers: " + strconv.Itoa(tClipLimitsDto.MaxTransfers))
	}
}

// getClipLimits ///////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getClipLimits() ClipLimitsDto {
	tResponse, tErr := http.Get(gApiUrl + "/clips/limits")
	if tErr != nil {
		ExitWithError("Checking QRClip limits failed, " + tErr.Error())
	}
	defer tResponse.Body.Close()

	var tClipLimitsDto ClipLimitsDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tClipLimitsDto)
	if tErr != nil {
		ExitWithError("Decoding ClipLimitsDto, " + tErr.Error())
	}

	return tClipLimitsDto
}

// getClipLimitsUser ///////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getClipLimitsUser(pJwt string) ClipLimitsDto {
	tErrorPrefix := "Checking QRClip limits for user, "
	var tUrl = gApiUrl + "/clips/limits/user"

	// CREATE THE REQUEST
	tRequest, tErr := http.NewRequest("GET", tUrl, nil)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	var tBearer = "Bearer " + pJwt
	tRequest.Header.Add("Authorization", tBearer)
	tRequest.Header.Set("Content-Type", "application/json")

	// SEND REQUEST
	tClient := &http.Client{}
	tResponse, tErr := tClient.Do(tRequest)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	// PARSE RESPONSE
	var tClipLimitsDto ClipLimitsDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tClipLimitsDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	return tClipLimitsDto
}

// CheckIfCanBeSent ////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func CheckIfCanBeSent(pUpdateClipDto *UpdateClipDto) {
	tLimits := GetLimits()

	// TEXT LIMIT
	if tLimits.Text < len(pUpdateClipDto.EncryptedText) {
		ExitWithError("Text message to big, your limit is " + strconv.Itoa(tLimits.Text))
	}

	// MAX TRANSFERS
	if tLimits.MaxTransfers > 0 { // ZERO IS UNLIMITED
		if tLimits.MaxTransfers < pUpdateClipDto.MaxTransfers {
			ShowWarning("Max transfers value to big, your max value is " + strconv.Itoa(tLimits.MaxTransfers))
			pUpdateClipDto.MaxTransfers = tLimits.MaxTransfers
		}
	}

	// EXPIRATION
	if tLimits.ExpiresInMinutes < pUpdateClipDto.ExpiresInMinutes {
		ShowWarning("Max expiration value to big, your max value is " + strconv.Itoa(tLimits.ExpiresInMinutes))
		pUpdateClipDto.ExpiresInMinutes = tLimits.ExpiresInMinutes
	}

	// FILE SIZE
	if tLimits.FileMb < int(pUpdateClipDto.FileSize/1000000) {
		ExitWithError("File to big, your limit is " + strconv.Itoa(tLimits.FileMb) + "Mb")
	}

}
