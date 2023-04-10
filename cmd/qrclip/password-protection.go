package main

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
	"syscall"
)

// GetUserPassword /////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GetUserPassword(pKey string, pSubId string) (string, string) {
	var tAccessKey string
	var tPasswordKey []byte
	// GET PASSWORD FROM USER
	color.Set(color.FgYellow)
	fmt.Println("This QRClip is password protected")
	fmt.Println("Enter Password: ")
	color.Unset()
	tBytePassword, tErr := term.ReadPassword(int(syscall.Stdin))
	if tErr != nil {
		ExitWithError("Error getting password")
	}
	tPassword := string(tBytePassword)

	// GET THE ACCESS KEY AND PASSWORD KEY
	tAccessKey, tPasswordKey, tErr = calculatePasswordKey(tPassword, pSubId)
	if tErr != nil {
		ExitWithError("Error using the password to calculate the password key")
	}

	// NOW GENERATE THE COMBINE KEY
	tCombinedKey, tErr := generateCombinedKey(pKey, tPasswordKey)
	if tErr != nil {
		ExitWithError("Error using the password to generate the combined key")
	}

	return tCombinedKey, tAccessKey
}

// CalculatePasswordKey ////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func calculatePasswordKey(pPassword string, pSubId string) (string, []byte, error) {
	tSalt := []byte(pSubId[0:16])
	tPasswordKey := argon2.IDKey([]byte(pPassword), tSalt, 3, 62500, 1, 64)
	return EncodeBase64UrlNoPad(tPasswordKey[32:48]), tPasswordKey, nil
}

// GenerateCombinedKey /////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func generateCombinedKey(pKey string, pPasswordKey []byte) (string, error) {
	tSalt := pPasswordKey[48:64]
	tCombinedArray := make([]byte, 64)
	copy(tCombinedArray[0:32], pPasswordKey[0:32])

	tByteKey, tErr := DecodeBase64Url(pKey)
	if tErr != nil {
		ExitWithError("Decoding Base64")
	}

	copy(tCombinedArray[32:64], tByteKey)
	tCombinedKey := argon2.IDKey(tCombinedArray, tSalt, 3, 62500, 1, 32)
	return EncodeBase64UrlNoPad(tCombinedKey), nil
}
