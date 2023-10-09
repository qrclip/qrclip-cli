package main

import (
	"fmt"
)

// PrintHelp ///////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelp() {

	ShowInfo("Secure, Fast, End-to-End Encrypted File & Text Sharing Right from Your Terminal.")
	fmt.Println("")
	fmt.Println("COMMANDS:")
	fmt.Println("")
	ShowInfoBlue("AUTHENTICATION:")
	fmt.Println("")

	// LOGIN
	PrintHelpLogin()

	fmt.Println("")

	// LOGOUT
	PrintHelpLogout()

	fmt.Println("")
	ShowInfoBlue("DATA TRANSFER:")
	fmt.Println("")

	// SEND
	PrintHelpSend()

	fmt.Println("")

	// RECEIVE
	PrintHelpReceive()

	fmt.Println("")
	ShowInfoBlue("UTILITY:")
	fmt.Println("")

	// CHECK LIMITS
	PrintHelpCheckLimits()

	fmt.Println("")

	// SELECT STORAGE
	PrintHelpSelectStorage()

	fmt.Println("")
	ShowInfoBlue("HELP & INFO:")
	fmt.Println("")

	// PRINT HELP
	PrintHelpHelp() // :)

	fmt.Println("")
}

// PrintHelpHelp //////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpHelp() {
	ShowInfoCyan("  HELP (help,h,--help)")
	ShowInfoYellow("   * View this guide:")
	fmt.Println("      qrclip help")
	fmt.Println("")
	ShowInfoCyan("  VERSION (version,v,--version)")
	ShowInfoYellow("   * View tool version:")
	fmt.Println("      qrclip --version")
}

// PrintHelpLogin //////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpLogin() {
	ShowInfoCyan("  LOGIN (login,l)")
	ShowInfoYellow("   * To login using qr code")
	fmt.Println("      qrclip login")
	ShowInfoYellow("   * With username and password")
	fmt.Println("      qrclip login -u myemail@email.com -p \"MySecretPassword\"")
	ShowInfoYellow("   * With username (password will be asked)")
	fmt.Println("      qrclip login -u myemail@email.com")
}

// PrintHelpLogout /////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpLogout() {
	ShowInfoCyan("  LOGOUT (logout,q)")
	ShowInfoYellow("   * To clear the credentials")
	fmt.Println("      qrclip logout")
}

// PrintHelpCheckLimits ////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpCheckLimits() {
	ShowInfoCyan("  CHECK LIMITS (check,c,limits,limit)")
	ShowInfoYellow("   * View user limits")
	fmt.Println("      qrclip check")
}

// PrintHelpSend ///////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpSend() {
	ShowInfoCyan("  SEND (send,s)")
	ShowInfoYellow("   * Send a message")
	fmt.Println("      qrclip send -m \"Message to Send\"")
	ShowInfoYellow("   * Send a file (provide path)")
	fmt.Println("      qrclip send -f  /path/to/fileToSend")
	ShowInfoYellow("   * Send both")
	fmt.Println("      qrclip send -m \"Message to Send\" -f /path/to/fileToSend")
	fmt.Println("")
	ShowInfoYellow("   Other Options:")
	ShowInfoYellow("     * Expiry (default 2880 mins):")
	fmt.Println("        -e 2880")
	ShowInfoYellow("     * Max transfers (default 5):")
	fmt.Println("        -mt 5")
	ShowInfoYellow("     * Allow deletion (default true):")
	fmt.Println("        -ad true")
}

// PrintHelpReceive ////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpReceive() {
	ShowInfoCyan("  RECEIVE (receive,r)")
	ShowInfoYellow("   * Generate Receiver:")
	fmt.Println("      qrclip receive")
	ShowInfoYellow("   * Fetch specific QRClip:")
	fmt.Println("      qrclip receive -i QRClipID -s QRClipSubID -k 32CharactersEncryptionKeyEncodedInBase64Url")
	ShowInfoYellow("   * Fetch by URL:")
	fmt.Println("      qrclip receive -u \"QRClipURL\"")
}

// PrintHelpSelectStorage //////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func PrintHelpSelectStorage() {
	ShowInfoCyan("  SELECT STORAGE")
	ShowInfoYellow("   * Choose storage option:")
	fmt.Println("      qrclip storage")
}
