package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"net/http"
	"os"
	"strings"
)

// ReceiveQRClip ///////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ReceiveQRClip(pId string, pSubId string, pKey string, pUrl string) {
	if pId != "" && pSubId != "" && pKey != "" {
		DisplayAndDownloadQRClip(pId, pSubId, pKey)
	} else {
		if pUrl != "" {
			handleQRClipUrl(pUrl)
		} else {
			startReceiveMode()
		}
	}
}

// DisplayAndDownloadQRClip ////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DisplayAndDownloadQRClip(pId string, pSubId string, pKey string) {
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
}

// handleQRClipUrl /////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func handleQRClipUrl(pUrl string) {
	tKey := ""
	tId := ""
	tSubId := ""
	tWorkingUrl := pUrl
	tSplitKey := strings.Split(tWorkingUrl, "#")
	if len(tSplitKey) == 2 {
		tKey = tSplitKey[1]
		tWorkingUrl = tSplitKey[0]
	}

	tSplitKey = strings.Split(tWorkingUrl, "&subId=")
	if len(tSplitKey) == 2 {
		tSubId = tSplitKey[1]
		tWorkingUrl = tSplitKey[0]
	}

	tSplitKey = strings.Split(tWorkingUrl, "receive/open?id=")
	if len(tSplitKey) == 2 {
		tId = tSplitKey[1]
	}

	DisplayAndDownloadQRClip(tId, tSubId, tKey)
}

// startReceiveMode ////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func showQRClipInfo(pClipDto ClipDto, pKey string) {
	ShowSuccess("------------------------------------------------------------")
	ShowSuccess("Id         :" + pClipDto.Id)
	ShowSuccess("SubId      :" + pClipDto.SubId)
	ShowSuccess("Valid until:" + pClipDto.ValidUntil)
	ShowSuccess("Text       :")
	ShowSuccess(DecryptText(pClipDto.EncryptedText, pKey, pClipDto.IVData.Text))
}

// downloadQRClipFiles /////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func downloadQRClipFiles(pClipDto ClipDto, pKey string) {
	// FOR EACH FILE
	for _, tFile := range pClipDto.Files {
		// DECRYPT THE FILE NAME
		tFile.Name = DecryptText(tFile.Name, pKey, pClipDto.IVData.FileNames[tFile.Index])
		// START DOWNLOADING
		downloadFileWithTicket(pClipDto, pKey, tFile)
	}
}

// downloadFileWithTicket //////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

	ShowInfo("DOWNLOADING FILE: " + pClipFileDto.Name)

	tBar := pb.ProgressBarTemplate(gProgressBarTemplate).Start64(pClipFileDto.Size)
	defer tBar.Finish()

	tDecryptedFile, tErr := os.OpenFile(pClipFileDto.Name, os.O_RDWR|os.O_CREATE, 0777)
	if tErr != nil {
		ExitWithError("Decrypted file " + pClipFileDto.Name + " , failed to create")
	}

	tEncryptedChunkBuffer := make([]byte, 1024)
	// DOWNLOAD THE FILE CHUNK BY CHUNK AND APPEND TO ENCRYPTED FILE
	for tChunk := 0; tChunk < pClipFileDto.ChunkCount; tChunk++ {

		if tChunk == 0 {
			// FIRST CHUNK USE THE URL FROM TICKET
			tEncryptedChunkBuffer, tErr = downloadFileChunk(tDownloadTicket.Url, tBar)
			if tErr != nil {
				ExitWithError("Error downloading first file chunk!")
			}

		} else {
			// REQUEST THE NEW CHUNK
			tChunkUrl := getFileDownloadChunk(pClipDto, tDownloadTicket)
			tEncryptedChunkBuffer, tErr = downloadFileChunk(tChunkUrl.Url, tBar)
			if tErr != nil {
				ExitWithError("Error downloading chunk!")
			}
		}

		tDecryptedChunk := DecryptBuffer(tEncryptedChunkBuffer, pKey, pClipDto.IVData.Files[pClipFileDto.Index])

		// COPY THE DECRYPTED CHUNK TO THE FILE
		_, tErr = tDecryptedFile.Write(tDecryptedChunk)
		if tErr != nil {
			ExitWithError("Failed to write chunk to file!")
		}

	}
	tBar.Finish()
}

// downloadFileChunk ///////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func downloadFileChunk(pUrl string, pBar *pb.ProgressBar) ([]byte, error) {

	tResponse, tErr := http.Get(pUrl)
	if tErr != nil {
		return nil, tErr
	}
	defer tResponse.Body.Close()

	// create proxy reader
	tReader := pBar.NewProxyReader(tResponse.Body)

	tBuf := new(bytes.Buffer)
	_, tErr = tBuf.ReadFrom(tReader)

	return tBuf.Bytes(), tErr
}

// getFileDownloadChunk ////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

	// IV
	tClipDto.IVData = getIVDataForClipDto(tClipDto.Version, tClipDto.SubId, len(tClipDto.Files))

	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	return tClipDto
}

// getIVDataForClipDto /////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getIVDataForClipDto(pVersion int, pSubId string, pFileNumber int) QrcIVData {
	var tQrcIVData QrcIVData

	// LEGACY - ALWAYS THE SAME (THIS CODE CAN BE REMOVED AFTER DEPLOY VERSION 2 PLUS 20 DAYS TO BE SAFE)
	if pVersion <= 1 {
		tQrcIVData.Text = pSubId[:16]
		for tFileCount := 0; tFileCount < pFileNumber; tFileCount++ {
			tQrcIVData.FileNames = append(tQrcIVData.FileNames, pSubId[:16])
			tQrcIVData.Files = append(tQrcIVData.FileNames, pSubId[:16])
		}
	}

	// NEW VERSION VARIABLE IV
	if pVersion == 2 {
		tQrcIVData = GenerateIVData(pSubId, pFileNumber)
	}

	return tQrcIVData
}
