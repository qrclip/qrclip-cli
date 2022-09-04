package main

import (
	"flag"
	"os"
	"strconv"
)

// ///////////////////////////////////////////////
// MAIN
func main() {
	ShowInfoCyan("---------------------------------")
	ShowInfoCyan("--    QRCLIP - VERSION 0.9     --")
	ShowInfoCyan("--                             --")
	ShowInfoCyan("--    https://www.qrclip.io    --")
	ShowInfoCyan("---------------------------------")
	ShowInfo("")

	////
	// COMMANDS
	tLoginCommand := flag.NewFlagSet("l", flag.ExitOnError)
	tLogoutCommand := flag.NewFlagSet("logout", flag.ExitOnError)
	tSendCommand := flag.NewFlagSet("s", flag.ExitOnError)
	tReceiveCommand := flag.NewFlagSet("r", flag.ExitOnError)
	tCheckCommand := flag.NewFlagSet("c", flag.ExitOnError)
	tEncryptCommand := flag.NewFlagSet("e", flag.ExitOnError)
	tDecryptCommand := flag.NewFlagSet("d", flag.ExitOnError)
	tGenerateKeyCommand := flag.NewFlagSet("g", flag.ExitOnError)
	tHelpCommand := flag.NewFlagSet("h", flag.ExitOnError)
	tStorageCommand := flag.NewFlagSet("storage", flag.ExitOnError)

	////
	// SUB COMMANDS
	tLoginSubCmdUser := tLoginCommand.String("u", "", "Username")
	tLoginSubCmdPassword := tLoginCommand.String("p", "", "Password")

	tSendSubCmdPath := tSendCommand.String("f", "", "File to send")
	tSendSubCmdMessage := tSendCommand.String("m", "", "Message to send")
	tSendSubCmdExpiration := tSendCommand.String("e", "15", "Expiration in minutes")
	tSendSubCmdMaxTransfers := tSendCommand.String("mt", "2", "Max Transfers")
	tSendSubCmdAllowDelete := tSendCommand.String("ad", "true", "Allow Delete")

	tReceiveSubCmdId := tReceiveCommand.String("i", "", "QRClip ID")
	tReceiveSubCmdSubUrl := tReceiveCommand.String("u", "", "QRClip Url")
	tReceiveSubCmdSubId := tReceiveCommand.String("s", "", "QRClip Sub ID")
	tReceiveSubCmdSubKey := tReceiveCommand.String("k", "", "QRClip Key")

	tEncryptSubCmdFile := tEncryptCommand.String("f", "", "File to encrypt")
	tEncryptSubCmdKey := tEncryptCommand.String("k", "", "Encryption Key (if none, one is generated)")

	tDecryptSubCmdFile := tDecryptCommand.String("f", "", "File to decrypt")
	tDecryptSubCmdKey := tDecryptCommand.String("k", "", "Decryption Key")

	tGenerateKeySubCmdPhrase := tGenerateKeyCommand.String("p", "", "Phrase")

	////
	// IF NO COMMAND FOUND SHOW HELP
	if len(os.Args) < 2 {
		PrintHelp()
		ExitWithError("No command found!")
	}

	////
	// PARSE COMMANDS
	switch os.Args[1] {

	// LOGIN
	case "l":
		if tLoginCommand.Parse(os.Args[2:]) == nil {
			Login(*tLoginSubCmdUser, *tLoginSubCmdPassword)
		}

	// LOGOUT
	case "logout":
		if tLogoutCommand.Parse(os.Args[2:]) == nil {
			Logout()
		}

	// SEND
	case "s":
		if tSendCommand.Parse(os.Args[2:]) == nil {
			handleSendCommand(*tSendSubCmdPath, *tSendSubCmdMessage, *tSendSubCmdExpiration,
				*tSendSubCmdMaxTransfers, *tSendSubCmdAllowDelete)
		}

	// RECEIVE
	case "r":
		if tReceiveCommand.Parse(os.Args[2:]) == nil {
			ReceiveQRClip(*tReceiveSubCmdId, *tReceiveSubCmdSubId, *tReceiveSubCmdSubKey, *tReceiveSubCmdSubUrl)
		}

	// CHECK LIMITS
	case "c":
		if tCheckCommand.Parse(os.Args[2:]) == nil {
			CheckLimits()
		}

	// ENCRYPT OFFLINE
	case "e":
		if tEncryptCommand.Parse(os.Args[2:]) == nil {
			handleEncryptCommand(*tEncryptSubCmdFile, *tEncryptSubCmdKey)
		}

	// DECRYPT OFFLINE
	case "d":
		if tDecryptCommand.Parse(os.Args[2:]) == nil {
			handleDecryptCommand(*tDecryptSubCmdFile, *tDecryptSubCmdKey)
		}

	// GENERATE KEY
	case "g":
		if tGenerateKeyCommand.Parse(os.Args[2:]) == nil {
			handleGenerateKeyCommand(*tGenerateKeySubCmdPhrase)
		}

	// HELP
	case "h":
		if tHelpCommand.Parse(os.Args[2:]) == nil {
			PrintHelp()
		}

	// STORAGE
	case "storage":
		if tStorageCommand.Parse(os.Args[2:]) == nil {
			handleStorageCommand()
		}

	// COMMAND NOT FOUND
	default:
		{
			PrintHelp()
			ExitWithError("Command not found!")
		}
	}
}

// ///////////////////////////////////////////////
// HANDLE SEND COMMAND
func handleSendCommand(
	pSendSubCmdPath string,
	pSendSubCmdMessage string,
	pSendSubCmdExpiration string,
	pSendSubCmdMaxTransfers string,
	pSendSubCmdAllowDelete string,
) {
	if pSendSubCmdPath == "" && pSendSubCmdMessage == "" {
		PrintHelpSend()
		ExitWithError("Command needs more parameters")
	}
	tExpirationInMinutes, tErr := strconv.Atoi(pSendSubCmdExpiration)
	if tErr != nil {
		ShowWarning("Expiration invalid using default 15 minutes")
		tExpirationInMinutes = 15
	}
	tMaxTransfers, tErr := strconv.Atoi(pSendSubCmdMaxTransfers)
	if tErr != nil {
		ShowWarning("Max transfers invalid using default 2")
		tMaxTransfers = 2
	}
	tAllowDelete, tErr := strconv.ParseBool(pSendSubCmdAllowDelete)
	if tErr != nil {
		ShowWarning("Anyone can delete invalid using default true")
		tAllowDelete = true
	}
	SendQRClip(pSendSubCmdPath, pSendSubCmdMessage, tExpirationInMinutes, tMaxTransfers, tAllowDelete)
}

// ///////////////////////////////////////////////
// HANDLE ENCRYPT COMMAND
func handleEncryptCommand(pEncryptSubCmdFile string, pEncryptSubCmdKey string) {
	if pEncryptSubCmdFile == "" {
		PrintHelpEncrypt()
		ExitWithError("File is needed")
	}
	OfflineEncrypt(pEncryptSubCmdFile, pEncryptSubCmdKey)
}

// ///////////////////////////////////////////////
// HANDLE DECRYPT COMMAND
func handleDecryptCommand(pDecryptSubCmdFile string, pDecryptSubCmdKey string) {
	if pDecryptSubCmdFile == "" || pDecryptSubCmdKey == "" {
		PrintHelpDecrypt()
		if pDecryptSubCmdFile == "" && pDecryptSubCmdKey == "" {
			ExitWithError("File and key are needed!")
		}
		if pDecryptSubCmdFile == "" {
			ExitWithError("File is needed!")
		}
		if pDecryptSubCmdKey == "" {
			ExitWithError("Key is needed!")
		}
	}
	OfflineDecrypt(pDecryptSubCmdFile, pDecryptSubCmdKey)
}

// ///////////////////////////////////////////////
// HANDLE GENERATE KEY COMMAND
func handleGenerateKeyCommand(pPhrase string) {
	if pPhrase == "" {
		ShowSuccess(GenerateEncryptionKey())
	} else {
		ShowSuccess(GenerateEncryptionKeyWithPhrase(pPhrase))
	}
}

// ///////////////////////////////////////////////
// HANDLE STORAGE COMMAND
func handleStorageCommand() {
	SelectStorage()
}
