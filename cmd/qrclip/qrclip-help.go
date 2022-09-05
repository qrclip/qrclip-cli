package main

import (
	"fmt"

	"github.com/fatih/color"
)

// PrintHelp ///////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelp() {

	// LOGIN
	PrintHelpLogin()

	fmt.Println("")

	// LOGOUT
	PrintHelpLogout()

	fmt.Println("")

	// CHECK LIMITS
	PrintHelpCheckLimits()

	fmt.Println("")

	// SEND
	PrintHelpSend()

	fmt.Println("")

	// RECEIVE
	PrintHelpReceive()

	fmt.Println("")

	// SELECT STORAGE
	PrintHelpSelectStorage()

	fmt.Println("")

	// ENCRYPT
	PrintHelpEncrypt()

	fmt.Println("")

	// DECRYPT
	PrintHelpDecrypt()

	fmt.Println("")

	// GENERATE KEY
	PrintHelpGenerateKey()

	fmt.Println("")

	// HELP
	ShowInfoCyan("HELP")
	fmt.Println(" qrclip h")

	fmt.Println("")
}

// PrintHelpLogin //////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpLogin() {
	ShowInfoCyan("LOGIN")
	ShowInfoYellow(" To login using qr code")
	fmt.Println("  qrclip l")
	ShowInfoYellow(" With username and password")
	fmt.Println("  qrclip l -u myemail@email.com -p \"MySecretPassword\"")
	ShowInfoYellow(" With username (password will be asked)")
	fmt.Println("  qrclip l -u myemail@email.com")
}

// PrintHelpLogout /////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpLogout() {
	ShowInfoCyan("LOGOUT")
	ShowInfoYellow(" To clear the credentials")
	fmt.Println("  qrclip logout")
}

// PrintHelpCheckLimits ////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpCheckLimits() {
	ShowInfoCyan("CHECK LIMITS")
	ShowInfoYellow(" Check QRClip limits of the current user")
	fmt.Println("  qrclip c")
}

// PrintHelpSend ///////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpSend() {
	ShowInfoCyan("SEND")
	fmt.Println(" qrclip s -m \"Message to Send\" -f fileToSend")
	fmt.Println(" qrclip s -m \"Message to Send\"")
	fmt.Println(" qrclip s -f fileToSend")
	ShowInfoYellow(" Other Options:")

	printHelpOption("  -e 15     ", "( Expiration Time in minutes - default 15 )")
	printHelpOption("  -mt 2     ", "( Max transfers - default 2 )")
	printHelpOption("  -ad true  ", "( Allow delete - default true )")
}

// printHelpOption /////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func printHelpOption(pCommand string, pDescription string) {
	fmt.Print(pCommand + " ")
	color.Set(color.FgYellow)
	fmt.Println(pDescription)
	color.Unset()
}

// PrintHelpReceive ////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpReceive() {
	ShowInfoCyan("RECEIVE")
	ShowInfoYellow(" Receive Mode:")
	fmt.Println("  qrclip r")
	ShowInfoYellow(" Get QRClip:")
	fmt.Println("  qrclip r -i QRClipID -s QRClipSubID -k 32CharactersEncryptionKeyEncodedInBase64")
	fmt.Println("  qrclip r -u \"QRClipURL\"")
}

// PrintHelpEncrypt ////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpEncrypt() {
	ShowInfoCyan("ENCRYPT (OFFLINE)")
	ShowInfoYellow(" Encrypt with automatic generated key:")
	fmt.Println("  qrclip e -f fileToEncrypt")
	ShowInfoYellow(" Encrypt with a specified key:")
	fmt.Println("  qrclip e -f fileToEncrypt -k 32CharactersEncryptionKeyEncodedInBase64")
}

// PrintHelpSelectStorage //////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpSelectStorage() {
	ShowInfoCyan("SELECT STORAGE")
	ShowInfoYellow(" Select storage:")
	fmt.Println("  qrclip storage")
}

// PrintHelpDecrypt ////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpDecrypt() {
	ShowInfoCyan("DECRYPT (OFFLINE)")
	fmt.Println(" qrclip d -f fileToDecrypt -k 32CharactersEncryptionKeyEncodedInBase64")
}

// PrintHelpGenerateKey ////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpGenerateKey() {
	ShowInfoCyan("GENERATE KEY ENCODED IN BASE64")
	ShowInfoYellow(" Generate a random key:")
	fmt.Println("  qrclip g")
	ShowInfoYellow(" Generate a key with phrase:")
	fmt.Println("  qrclip g -p TheMountainFlyingOverTheRedRiver")
	ShowInfoYellow("   < 32 characters, X's are appended")
	ShowInfoYellow("   > 32 characters, the text is shortened to 32 characters")
}
