package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/fatih/color"
	"github.com/mdp/qrterminal"
	"os"
	"strings"
)

// GenerateEncryptionKey ///////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateEncryptionKey() string {
	tBytes := make([]byte, 32) //generate a random 32 byte key
	if _, err := rand.Read(tBytes); err != nil {
		ExitWithError("Error generating encryption key:" + err.Error())
	}
	return EncodeBase64UrlNoPad(tBytes)
}

// CreateQRClip ////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func CreateQRClip(pReceiveMode bool) (ClipDto, error) {
	var tUrlPath = "/clips/create"

	tJwt := CheckJwtToken()
	if tJwt != "" {
		tUrlPath = tUrlPath + "/user"
	}

	var tCreateClipDto CreateClipDto
	tCreateClipDto.ReceivingMode = pReceiveMode

	// REQUEST
	tResponse, tErr := HttpDoPost(tUrlPath, tJwt, tCreateClipDto)
	if tErr != nil {
		return ClipDto{}, tErr
	}

	// PARSE RESPONSE
	var tClipDto ClipDto
	tErr = DecodeJSONResponse(tResponse, &tClipDto)
	if tErr != nil {
		return ClipDto{}, tErr
	}

	return tClipDto, nil
}

// GetQRCodeTerminalConfig /////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GetQRCodeTerminalConfig() qrterminal.Config {
	//goland:noinspection ALL
	return qrterminal.Config{
		Level:      qrterminal.M,
		Writer:     os.Stdout,
		HalfBlocks: gHalfBlocks,
		BlackChar:  qrterminal.BLACK,
		WhiteChar:  qrterminal.WHITE,
		QuietZone:  2,
	}
}

// DisplayQRClipQRCode /////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DisplayQRClipQRCode(pClipDto ClipDto, pKey string) {
	tConfig := GetQRCodeTerminalConfig()

	tUrl := gSpaUrl + "/receive/open?id=" + pClipDto.Id + "&subId=" + pClipDto.SubId + "#" + pKey

	ShowSuccess("----------------------------------------")
	ShowSuccess("ID     : " + pClipDto.Id)
	ShowSuccess("SUB ID : " + pClipDto.SubId)
	ShowSuccess("KEY    : " + pKey)
	ShowSuccess("LINK   : " + tUrl)
	ShowSuccess(" ")
	qrterminal.GenerateWithConfig(tUrl, tConfig)
}

// GetUserInput ////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GetUserInput() {
	tReader := bufio.NewReader(os.Stdin)
	tChar, _, tErr := tReader.ReadRune()
	if tErr != nil {
		ShowError(tErr.Error())
	}

	if tChar == 'a' {
		ExitWithError("ABORTED")
	}
}

// ExitWithError ///////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ExitWithError(pInfo string) {
	ShowError(pInfo)
	color.Set(color.FgRed)
	fmt.Println("QRClip EXIT")
	color.Unset()
	os.Exit(1)
}

// ShowError ///////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowError(pInfo string) {
	color.Set(color.FgRed)
	fmt.Println("ERROR: " + pInfo)
	color.Unset()
}

// ShowWarning /////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowWarning(pInfo string) {
	color.Set(color.FgYellow)
	fmt.Println("WARNING: " + pInfo)
	color.Unset()
}

// ShowInfo ////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowInfo(pInfo string) {
	fmt.Println(pInfo)
}

// ShowInfoCyan ////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowInfoCyan(pInfo string) {
	color.Set(color.FgCyan)
	fmt.Println(pInfo)
	color.Unset()
}

// ShowInfoYellow //////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowInfoYellow(pInfo string) {
	color.Set(color.FgYellow)
	fmt.Println(pInfo)
	color.Unset()
}

// ShowSuccess /////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowSuccess(pInfo string) {
	color.Set(color.FgGreen)
	fmt.Println(pInfo)
	color.Unset()
}

// EncodeBase64UrlNoPad ////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func EncodeBase64UrlNoPad(pBytes []byte) string {
	tBase64UrlWithPad := base64.URLEncoding.EncodeToString(pBytes)
	return strings.Replace(tBase64UrlWithPad, "=", "", -1)
}

// DecodeBase64Url ////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DecodeBase64Url(pString string) []byte {
	// ADD THE PADDING IN CASE DOES NOT EXIST
	l := len(pString) % 4
	if l != 0 {
		pString += strings.Repeat("=", 4-l)
	}

	tBytes, tErr := base64.URLEncoding.DecodeString(pString)
	if tErr != nil {
		ExitWithError("Error decoding from base64, " + tErr.Error())
	}
	return tBytes
}
