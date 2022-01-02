package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cheggaaa/pb/v3"
)

// ReceiveQRClip ///////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ReceiveQRClip(pId string, pSubId string, pKey string) {
	if pId != "" && pSubId != "" && pKey != "" {
		tClipDto := getQRClip(pId, pSubId)
		// IF EXISTS SHOW IT
		if tClipDto.Id != "" {
			// DECRYPT AND DISPLAY TEXT THEN DOWNLOAD FILES
			showQRClipInfo(tClipDto, pKey)

			// DOWNLOAD QRCLIP FILES
			downloadQRClipFiles(tClipDto, pKey)
		} else {
			ExitWithError("QRClip not found!")
		}
	} else {
		startReceiveMode()
	}
}

// startReceiveMode ////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func startReceiveMode() {
	ShowInfo("GENERATING RECEIVER")

	// GENERATE KEY
	tKey := GenerateEncryptionKey()

	// CREATE QRCLIP
	tClipDto := CreateQRClip(true)

	// DISPLAY QRCLIP
	DisplayQRClipQRCode(tClipDto, tKey)

	// WAIT FOR THE USER TO SEND
userInputGOTO: // FIRST GOTO I EVER USED :) AND THIS LOOKS A PERFECT PLACE TO USE IT
	ShowInfoYellow("Press any key after you sent in the other device or \"a\" to abort:")

	GetUserInput() // HERE IF USER WANTS ABORTS AND PROGRAM EXITS

	// GET CLIP
	tClipDtoFetched := getQRClip(tClipDto.Id, tClipDto.SubId)

	// IF EXISTS SHOW IT, IF NOT GOTO THE USER INPUT AGAIN
	if tClipDtoFetched.Id == "" {
		ShowError("QRClip not ready!")
		goto userInputGOTO // GOTO THE TOP AGAIN UNTIL ABORTS
	} else {
		tClipDto = tClipDtoFetched
		// DECRYPT TEXT AND DISPLAY
		showQRClipInfo(tClipDto, tKey)

		// DOWNLOAD QRCLIP FILES
		downloadQRClipFiles(tClipDto, tKey)
	}
}

// showQRClipInfo //////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func showQRClipInfo(pClipDto ClipDto, pKey string) {
	ShowSuccess("------------------------------------------------------------")
	ShowSuccess("Id         :" + pClipDto.Id)
	ShowSuccess("SubId      :" + pClipDto.SubId)
	ShowSuccess("Valid until:" + pClipDto.ValidUntil)
	ShowSuccess("Text       :")
	ShowSuccess(DecryptText(pClipDto.EncryptedText, pKey, pClipDto.SubId))
}

// downloadQRClipFiles /////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func downloadQRClipFiles(pClipDto ClipDto, pKey string) {
	// FOR EACH FILE
	for _, tFile := range pClipDto.Files {
		// DECRYPT THE FILE NAME
		tFile.Name = DecryptText(tFile.Name, pKey, pClipDto.SubId)
		// START DOWNLOADING
		downloadFileWithTicket(pClipDto, pKey, tFile)
	}
}

// downloadFileWithTicket //////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func downloadFileWithTicket(pClipDto ClipDto, pKey string, pClipFileDto ClipFileDto) {
	tDownloadTicket := getFileDownloadTicket(pClipDto, pClipFileDto.Index)
	if tDownloadTicket.Error != 0 {
		ExitWithError("Error getting download ticket:")
		if tDownloadTicket.Error != 1 {
			ExitWithError("Clip owner has no credits!")
		}
		if tDownloadTicket.Error != 2 {
			ExitWithError("You dont have credits!")
		}
		return
	}

	ShowInfo("DOWNLOADING ENCRYPTED FILE")

	// CREATE THE EMPTY ENCRYPTED FILE
	tEncryptedFileName := pClipFileDto.Name + ".enc"
	tOutFile, tErr := os.Create(tEncryptedFileName)
	if tErr != nil {
		return
	}
	defer tOutFile.Close()

	// DOWNLOAD THE FILE CHUNK BY CHUNK AND APPEND TO ENCRYPTED FILE
	for tChunk := 0; tChunk < pClipFileDto.ChunkCount; tChunk++ {
		if tChunk == 0 {
			// FIRST CHUNK USE THE URL FROM TICKET
			tErr := downloadFileAppend(tOutFile, tDownloadTicket.Url)
			if tErr != nil {
				ExitWithError("Error downloading first file chunk!")
			}
		} else {
			// REQUEST THE NEW CHUNK
			tChunkUrl := getFileDownloadChunk(pClipDto, tDownloadTicket)
			tErr := downloadFileAppend(tOutFile, tChunkUrl.Url)
			if tErr != nil {
				ExitWithError("Error downloading chunk!")
			}
		}
	}

	// CLOSE FILE
	tOutFile.Close()

	// DECRYPT THE FILE
	DecryptFile(tEncryptedFileName, pKey, pClipFileDto.Name, pClipFileDto.Size, pClipDto.SubId)

	// REMOVE ENCRYPTED FILE
	RemoveFile(tEncryptedFileName)
}

// downloadFileAppend //////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func downloadFileAppend(pFile *os.File, pUrl string) error {
	tResponse, tErr := http.Get(pUrl)
	if tErr != nil {
		return tErr
	}
	defer tResponse.Body.Close()

	tBar := pb.ProgressBarTemplate(gProgressBarTemplate).Start64(tResponse.ContentLength)
	defer tBar.Finish()

	tBarWriter := tBar.NewProxyWriter(pFile)

	_, tErr = io.Copy(tBarWriter, tResponse.Body)

	return tErr
}

// getFileDownloadChunk ////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getFileDownloadChunk(
	pClipDto ClipDto,
	pDownloadTicket GetFileDownloadTicketResponseDto) GetFileDownloadChunkResponseDto {

	tErrorPrefix := "Error getting Download Ticket, "
	tUrl := gApiUrl + "/clips/" + pClipDto.Id + "/" + pClipDto.SubId +
		"/file-download-chunk/" + pDownloadTicket.Id + "/" + pDownloadTicket.Key

	tResponse, tErr := http.Get(tUrl)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer tResponse.Body.Close()

	var tGetFileDownloadChunkResponseDto GetFileDownloadChunkResponseDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tGetFileDownloadChunkResponseDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	return tGetFileDownloadChunkResponseDto
}

// getFileDownloadTicket ///////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getFileDownloadTicket(pClipDto ClipDto, pFileIndex int) GetFileDownloadTicketResponseDto {
	tErrorPrefix := "Error getting Download Ticket, "
	tUrl := gApiUrl + "/clips/" + pClipDto.Id + "/" + pClipDto.SubId +
		"/file-download-ticket/" + fmt.Sprintf("%v", pFileIndex)

	tResponse, tErr := http.Get(tUrl)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer tResponse.Body.Close()

	var tGetFileDownloadTicketResponseDto GetFileDownloadTicketResponseDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tGetFileDownloadTicketResponseDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	return tGetFileDownloadTicketResponseDto
}

// getQRClip ///////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getQRClip(pId string, pSubId string) ClipDto {
	tErrorPrefix := "Error getting QRClip, "
	tUrl := gApiUrl + "/clips/" + pId + "/" + pSubId

	tResponse, tErr := http.Get(tUrl)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer tResponse.Body.Close()

	var tClipDto ClipDto
	tClipDto.Id = ""
	tErr = json.NewDecoder(tResponse.Body).Decode(&tClipDto)

	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	return tClipDto
}
