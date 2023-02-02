package main

import (
	"flag"
	"os"
	"strconv"
)

// main ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	ShowInfoCyan("---------------------------------")
	ShowInfoCyan("--    QRCLIP - VERSION 0.98    --")
	ShowInfoCyan("--                             --")
	ShowInfoCyan("--    https://app.qrclip.io    --")
	ShowInfoCyan("---------------------------------")
	ShowInfo("")

	////
	// COMMANDS
	tLoginCommand := flag.NewFlagSet("login", flag.ExitOnError)
	tLogoutCommand := flag.NewFlagSet("logout", flag.ExitOnError)
	tSendCommand := flag.NewFlagSet("send", flag.ExitOnError)
	tReceiveCommand := flag.NewFlagSet("receive", flag.ExitOnError)
	tCheckCommand := flag.NewFlagSet("check", flag.ExitOnError)
	tHelpCommand := flag.NewFlagSet("help", flag.ExitOnError)
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
	case "login", "l":
		if tLoginCommand.Parse(os.Args[2:]) == nil {
			Login(*tLoginSubCmdUser, *tLoginSubCmdPassword)
		}

	// LOGOUT
	case "logout", "q":
		if tLogoutCommand.Parse(os.Args[2:]) == nil {
			Logout()
		}

	// SEND
	case "send", "s":
		if tSendCommand.Parse(os.Args[2:]) == nil {
			handleSendCommand(*tSendSubCmdPath, *tSendSubCmdMessage, *tSendSubCmdExpiration,
				*tSendSubCmdMaxTransfers, *tSendSubCmdAllowDelete)
		}

	// RECEIVE
	case "receive", "r":
		if tReceiveCommand.Parse(os.Args[2:]) == nil {
			ReceiveQRClip(*tReceiveSubCmdId, *tReceiveSubCmdSubId, *tReceiveSubCmdSubKey, *tReceiveSubCmdSubUrl)
		}

	// CHECK LIMITS
	case "check", "c":
		if tCheckCommand.Parse(os.Args[2:]) == nil {
			CheckLimits()
		}

	// HELP
	case "help", "h", "-h", "--help":
		if tHelpCommand.Parse(os.Args[2:]) == nil {
			PrintHelp()
		}

	// STORAGE
	case "storage":
		if tStorageCommand.Parse(os.Args[2:]) == nil {
			SelectStorage()
		}

	// COMMAND NOT FOUND
	default:
		{
			PrintHelp()
			ExitWithError("Command not found!")
		}
	}
}

// handleSendCommand ///////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
