package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/mdp/qrterminal"
	rand2 "math/rand"
	"net/http"
	"os"
)

// GenerateEncryptionKey ///////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateEncryptionKey() string {
	tBytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(tBytes); err != nil {
		ExitWithError("Error generating encryption key:" + err.Error())
	}
	return base64.StdEncoding.EncodeToString(tBytes)
}

// GenerateRandomIV ////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateRandomIV() string {
	//goland:noinspection SpellCheckingInspection
	var tValidChars = []rune("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	tIV := make([]rune, 16)
	for i := range tIV {
		tIV[i] = tValidChars[rand2.Intn(len(tValidChars))]
	}
	return string(tIV)
}

// GenerateEncryptionKeyWithPhrase /////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateEncryptionKeyWithPhrase(pPhrase string) string {
	if len(pPhrase) < 32 {
		tLength := len(pPhrase)
		tCharactersNeeded := 32 - tLength
		for i := 0; i < tCharactersNeeded; i++ {
			pPhrase = pPhrase + "X"
		}
	}

	if len(pPhrase) > 32 {
		pPhrase = pPhrase[0:32]
	}

	return base64.StdEncoding.EncodeToString([]byte(pPhrase))
}

// CreateQRClip ////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func CreateQRClip(pReceiveMode bool) ClipDto {
	tErrorPrefix := "Creating QRClip, "
	var tUrl = gApiUrl + "/clips/create"

	tJwt := CheckJwtToken()

	if tJwt != "" {
		tUrl = tUrl + "/user"
	}

	var tCreateClipDto CreateClipDto
	tCreateClipDto.ReceivingMode = pReceiveMode

	tJsonPayload, tErr := json.Marshal(tCreateClipDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	// CREATE REQUEST
	tRequest, tErr := http.NewRequest("POST", tUrl, bytes.NewBuffer(tJsonPayload))
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	tRequest.Header.Set("Content-Type", "application/json")

	// IF THERES A JWT TOKEN USE IT
	if tJwt != "" {
		tRequest.Header.Add("Authorization", "Bearer "+tJwt)
	}

	// SEND REQUEST
	tClient := &http.Client{}
	tResponse, tErr := tClient.Do(tRequest)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer tResponse.Body.Close()

	// PARSE RESPONSE
	var tClipDto ClipDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tClipDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	return tClipDto
}

// GetQRCodeTerminalConfig /////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DisplayQRClipQRCode(pClipDto ClipDto, pKey string) {
	tConfig := GetQRCodeTerminalConfig()

	tUrl := gSpaUrl + "/receive/open?id=" + pClipDto.Id + "&subId=" + pClipDto.SubId + "#" + pKey

	ShowSuccess("----------------------------------------")
	ShowSuccess("ID     : " + pClipDto.Id)
	ShowSuccess("SUB ID : " + pClipDto.SubId)
	ShowSuccess("KEY    : " + pKey)
	ShowSuccess(" ")
	qrterminal.GenerateWithConfig(tUrl, tConfig)
}

// GetUserInput ////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ExitWithError(pInfo string) {
	ShowError(pInfo)
	color.Set(color.FgRed)
	fmt.Println("QRClip EXIT")
	color.Unset()
	os.Exit(1)
}

// ShowError ///////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowError(pInfo string) {
	color.Set(color.FgRed)
	fmt.Println("ERROR: " + pInfo)
	color.Unset()
}

// ShowWarning /////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowWarning(pInfo string) {
	color.Set(color.FgYellow)
	fmt.Println("WARNING: " + pInfo)
	color.Unset()
}

// ShowInfo ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowInfo(pInfo string) {
	fmt.Println(pInfo)
}

// ShowInfoCyan ////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowInfoCyan(pInfo string) {
	color.Set(color.FgCyan)
	fmt.Println(pInfo)
	color.Unset()
}

// ShowInfoYellow //////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowInfoYellow(pInfo string) {
	color.Set(color.FgYellow)
	fmt.Println(pInfo)
	color.Unset()
}

// ShowSuccess /////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShowSuccess(pInfo string) {
	color.Set(color.FgGreen)
	fmt.Println(pInfo)
	color.Unset()
}

// EncodeBase64 ////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func EncodeBase64(pBytes []byte) string {
	return base64.StdEncoding.EncodeToString(pBytes)
}

// DecodeBase64 ////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DecodeBase64(pString string) []byte {
	tBytes, tErr := base64.StdEncoding.DecodeString(pString)
	if tErr != nil {
		ExitWithError("Error decoding from base64, " + tErr.Error())
	}
	return tBytes
}
