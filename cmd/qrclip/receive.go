package main

import (
	"bytes"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
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
	tBeforeOpenInfoDto, tErr := getQRClipBeforeOpenInfo(pId, pSubId)
	if tErr != nil {
		ExitWithError(tErr.Error())
	}

	// CHECK VERSION
	if tBeforeOpenInfoDto.Version > gClientVersion {
		ExitWithError("Update CLI to latest version first!")
	}

	// CHECK PASSWORD AND ASK FOR IT
	var tAccessKey string
	if tBeforeOpenInfoDto.PasswordProtected {
		pKey, tAccessKey = GetUserPassword(pKey, pSubId)
	}

	tClipDto, tErr := getQRClip(pId, pSubId, &tAccessKey)
	if tErr != nil {
		ExitWithError(tErr.Error())
	}

	// IF EXISTS SHOW IT
	if tClipDto.Id != "" {
		// DECRYPT AND DISPLAY TEXT THEN DOWNLOAD FILES
		showQRClipInfo(tClipDto, pKey)

		// DOWNLOAD QRCLIP FILES
		downloadQRClipFiles(tClipDto, pKey, &tAccessKey)
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
	tClipDto, tErr := CreateQRClip(true)
	if tErr != nil {
		ExitWithError(tErr.Error())
	}

	// DISPLAY QRCLIP
	DisplayQRClipQRCode(tClipDto, tKey)

	// WAIT FOR THE USER TO SEND
userInputGOTO: // FIRST GOTO I EVER USED :) AND THIS LOOKS A PERFECT PLACE TO USE IT
	ShowInfoYellow("Press any key after you sent in the other device or \"a\" to abort:")

	GetUserInput() // HERE IF USER WANTS ABORTS AND PROGRAM EXITS

	tBeforeOpenInfoDto, tErr := getQRClipBeforeOpenInfo(tClipDto.Id, tClipDto.SubId)
	if tErr != nil {
		ExitWithError(tErr.Error())
	}

	if tBeforeOpenInfoDto.Version == 0 {
		ShowError("QRClip not ready!")
		goto userInputGOTO // GOTO THE TOP AGAIN UNTIL ABORTS
	}

	DisplayAndDownloadQRClip(tClipDto.Id, tClipDto.SubId, tKey)
}

// showQRClipInfo //////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func showQRClipInfo(pClipDto ClipDto, pKey string) {
	ShowSuccess("------------------------------------------------------------")
	ShowSuccess("Id         :" + pClipDto.Id)
	ShowSuccess("SubId      :" + pClipDto.SubId)
	ShowSuccess("Valid until:" + pClipDto.ValidUntil)
	ShowSuccess("Text       :")
	tIV, tErr := GetIV(&pClipDto.IVGen, 0)
	if tErr != nil {
		ExitWithError(tErr.Error())
	}
	ShowSuccess(DecryptText(pClipDto.EncryptedText, pKey, tIV))
}

// downloadQRClipFiles /////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func downloadQRClipFiles(pClipDto ClipDto, pKey string, pAccessKey *string) {
	tFileChunkIndex := len(pClipDto.Files) + 1 // CHUNK INDEX FOR IV INDEX CALCULATION
	// FOR EACH FILE
	for _, tFile := range pClipDto.Files {
		fmt.Println(GetIV(&pClipDto.IVGen, int64(tFile.Index+1)))

		tIV, tErr := GetIV(&pClipDto.IVGen, int64(tFile.Index+1))
		if tErr != nil {
			ExitWithError(tErr.Error())
		}

		// DECRYPT THE FILE NAME
		tFile.Name = DecryptText(tFile.Name, pKey, tIV)

		// START DOWNLOADING
		downloadFileWithTicket(pClipDto, pKey, tFile, tFileChunkIndex, pAccessKey)

		tFileChunkIndex = tFileChunkIndex + tFile.ChunkCount
	}
}

// downloadFileWithTicket //////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func downloadFileWithTicket(pClipDto ClipDto, pKey string, pClipFileDto ClipFileDto, pFileChunkIndexStart int, pAccessKey *string) {
	tDownloadTicket, tErr := getFileDownloadTicket(pClipDto, pClipFileDto.Index, pAccessKey)
	if tErr != nil {
		ExitWithError("Error getting download ticket:" + tErr.Error())
	}
	if tDownloadTicket.Error != 0 {
		if tDownloadTicket.Error == 1 {
			ExitWithError("Clip owner has no credits!")
		}
		if tDownloadTicket.Error == 2 {
			ExitWithError("You dont have credits!")
		}
		if tDownloadTicket.Error == 3 {
			ExitWithError("Not found!")
		}
		ExitWithError("Error getting download ticket!")
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
			tChunkUrl, tErr := getFileDownloadChunk(pClipDto, tDownloadTicket)
			if tErr != nil {
				ExitWithError("Error getting download url!")
			}
			tEncryptedChunkBuffer, tErr = downloadFileChunk(tChunkUrl.Url, tBar)
			if tErr != nil {
				ExitWithError("Error downloading chunk!")
			}
		}

		tIVChunkIndex := pFileChunkIndexStart + tChunk

		tIV, tErr := GetIV(&pClipDto.IVGen, int64(tIVChunkIndex))
		if tErr != nil {
			ExitWithError(tErr.Error())
		}

		tDecryptedChunk := DecryptBuffer(tEncryptedChunkBuffer, pKey, tIV)

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
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(tResponse.Body)

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
	pDownloadTicket GetFileDownloadTicketResponseDto) (GetFileDownloadChunkResponseDto, error) {

	tUrlPath := "/clips/" + pClipDto.Id + "/" + pClipDto.SubId +
		"/file-download-chunk/" + pDownloadTicket.Id + "/" + pDownloadTicket.Key

	tResponse, tErr := HttpDoGet(tUrlPath, "", nil)
	if tErr != nil {
		return GetFileDownloadChunkResponseDto{}, tErr
	}

	var tGetFileDownloadChunkResponseDto GetFileDownloadChunkResponseDto
	tErr = DecodeJSONResponse(tResponse, &tGetFileDownloadChunkResponseDto)
	if tErr != nil {
		return GetFileDownloadChunkResponseDto{}, tErr
	}

	return tGetFileDownloadChunkResponseDto, nil
}

// getFileDownloadTicket ///////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getFileDownloadTicket(pClipDto ClipDto, pFileIndex int, pAccessKey *string) (GetFileDownloadTicketResponseDto, error) {
	tUrlPath := "/clips/" + pClipDto.Id + "/" + pClipDto.SubId +
		"/file-download-ticket/" + fmt.Sprintf("%v", pFileIndex)

	tResponse, tErr := HttpDoGet(tUrlPath, "", pAccessKey)
	if tErr != nil {
		return GetFileDownloadTicketResponseDto{}, tErr
	}

	var tGetFileDownloadTicketResponseDto GetFileDownloadTicketResponseDto
	tErr = DecodeJSONResponse(tResponse, &tGetFileDownloadTicketResponseDto)
	if tErr != nil {
		return GetFileDownloadTicketResponseDto{}, tErr
	}

	return tGetFileDownloadTicketResponseDto, nil
}

// getQRClip ///////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getQRClip(pId string, pSubId string, pAccessKey *string) (ClipDto, error) {
	tUrlPath := "/clips/" + pId + "/" + pSubId

	// REQUEST
	tResponse, tErr := HttpDoGet(tUrlPath, "", pAccessKey)
	if tErr != nil {
		return ClipDto{}, tErr
	}

	var tClipDto ClipDto
	tClipDto.Id = ""
	tErr = DecodeJSONResponse(tResponse, &tClipDto)
	if tErr != nil {
		return ClipDto{}, tErr
	}

	// COUNT FILE CHUNKS
	tFilesChunkNumber := make([]int, len(tClipDto.Files))
	for _, tFile := range tClipDto.Files {
		tFilesChunkNumber[tFile.Index] = tFile.ChunkCount
	}

	tClipDto.IVGen = CreateQRClipIVGenerator(tClipDto.SubId, 24, 2500)

	return tClipDto, nil
}

// getQRClipBeforeOpenInfo /////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getQRClipBeforeOpenInfo(pId string, pSubId string) (BeforeOpenInfoDto, error) {
	tUrlPath := "/clips/" + pId + "/" + pSubId + "/before-open-info"

	// REQUEST
	tResponse, tErr := HttpDoGet(tUrlPath, "", nil)
	if tErr != nil {
		return BeforeOpenInfoDto{}, tErr
	}

	var tBeforeOpenInfoDto BeforeOpenInfoDto
	tErr = DecodeJSONResponse(tResponse, &tBeforeOpenInfoDto)
	if tErr != nil {
		return BeforeOpenInfoDto{}, tErr
	}

	return tBeforeOpenInfoDto, nil
}
