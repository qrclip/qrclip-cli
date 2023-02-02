package main

import (
	"strconv"
)

// CheckLimits /////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func CheckLimits() {
	tLimits, tErr := GetLimits()
	if tErr != nil {
		ExitWithError("Error checking limits: " + tErr.Error())
	} else {
		tLimits.printClipLimits()
	}
}

// GetLimits ///////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GetLimits() (ClipLimitsDto, error) {
	tJwt := CheckJwtToken()
	if tJwt == "" {
		return getClipLimits()
	} else {
		return getClipLimitsUser(tJwt)
	}
}

// printClipLimits /////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (pClipLimitsDto ClipLimitsDto) printClipLimits() {
	ShowSuccess("QRClip LIMITS")
	ShowSuccess(" Max Characters: " + strconv.Itoa(pClipLimitsDto.Text))
	ShowSuccess(" Max Expiration Minutes: " + strconv.Itoa(pClipLimitsDto.ExpiresInMinutes))
	ShowSuccess(" File Size(Mb): " + strconv.Itoa(pClipLimitsDto.FileMb))
	if pClipLimitsDto.MaxTransfers == 0 { // IF ZERO ITS UNLIMITED
		ShowSuccess(" Max Transfers: Unlimited")
	} else {
		ShowSuccess(" Max Transfers: " + strconv.Itoa(pClipLimitsDto.MaxTransfers))
	}
}

// getClipLimits ///////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getClipLimits() (ClipLimitsDto, error) {
	// REQUEST
	tResponse, tErr := HttpDoGet("/clips/limits", "")
	if tErr != nil {
		return ClipLimitsDto{}, tErr
	}

	// RESPONSE
	var tClipLimitsDto ClipLimitsDto
	tErr = DecodeJSONResponse(tResponse, &tClipLimitsDto)
	if tErr != nil {
		return ClipLimitsDto{}, tErr
	}

	return tClipLimitsDto, nil
}

// getClipLimitsUser ///////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getClipLimitsUser(pJwt string) (ClipLimitsDto, error) {
	// REQUEST
	tResponse, tErr := HttpDoGet("/clips/limits/user", pJwt)
	if tErr != nil {
		return ClipLimitsDto{}, tErr
	}

	// RESPONSE
	var tClipLimitsDto ClipLimitsDto
	tErr = DecodeJSONResponse(tResponse, &tClipLimitsDto)
	if tErr != nil {
		return ClipLimitsDto{}, tErr
	}
	return tClipLimitsDto, nil
}

// CheckIfCanBeSent ////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func CheckIfCanBeSent(pUpdateClipDto *UpdateClipDto) {
	tLimits, tErr := GetLimits()
	if tErr != nil {
		return
	}
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
