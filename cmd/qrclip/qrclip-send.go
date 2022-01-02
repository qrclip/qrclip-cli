package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

// SendQRClip //////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SendQRClip(pFilePath string, pMessage string, pExpiration int, pMaxTransfers int, pAllowDelete bool) {
	ShowInfo("SENDING QRCLIP")

	// GENERATE THE ENCRYPTION KEY
	tKey := GenerateEncryptionKey()

	// GET UPDATE CLIP DTO
	tUpdateClipDto := getUpdateClipDtoObject(pFilePath, pExpiration, pMaxTransfers, pAllowDelete)

	// CHECK LIMITS FOR FILE SIZE AND TEXT SIZE ALSO
	CheckIfCanBeSent(&tUpdateClipDto) // EXITS THE PROGRAM IF NOT OK

	// CREATE QRCLIP
	tClipDto := CreateQRClip(false)

	// ENCRYPT TEXT MESSAGE IF EXISTS
	if pMessage != "" {
		tUpdateClipDto.EncryptedText = EncryptText(pMessage, tKey, tClipDto.SubId)
	}

	// ENCRYPT THE FILE IF EXISTS
	if tUpdateClipDto.FileSize > 0 {
		EncryptFile(pFilePath, tKey, tClipDto.Id, tUpdateClipDto.FileSize, tClipDto.SubId)
	}

	// UPDATE QRCLIP
	tUpdateClipResponseDto := updateQRClip(tClipDto, tUpdateClipDto)

	// UPLOAD FILE
	if tUpdateClipDto.FileSize > 0 {
		// UPLOAD FILE
		tChunkCount := uploadFileChunkByChunk(tClipDto, tUpdateClipResponseDto.PreSignedPost, tUpdateClipDto.FileSize)
		if tChunkCount == 0 {
			ExitWithError("Error uploading file!")
		}

		// SET FILE UPLOAD FINISHED
		setFileUploadFinished(pFilePath, tKey, tClipDto, tUpdateClipDto, tChunkCount)

		// REMOVE FILE
		RemoveFile(tClipDto.Id)
	}

	// DISPLAY QRCLIP QR CODE
	if tUpdateClipResponseDto.Ok {
		ShowSuccess("QRCLIP READY - SCAN TO DOWNLOAD")
		DisplayQRClipQRCode(tClipDto, tKey)
	} else {
		ExitWithError("Error sending QRClip")
	}
}

// setFileUploadFinished ///////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func setFileUploadFinished(pFilePath string, pKey string, pClipDto ClipDto, pUpdateClipDto UpdateClipDto, pChunkCount int) {
	var tFileUploadFinishedFileDto FileUploadFinishedFileDto
	tFileUploadFinishedFileDto.Name = EncryptText(pFilePath, pKey, pClipDto.SubId)
	tFileUploadFinishedFileDto.Index = 0
	tFileUploadFinishedFileDto.Size = pUpdateClipDto.FileSize
	tFileUploadFinishedFileDto.ChunkCount = pChunkCount

	var tFileUploadFinishedDto FileUploadFinishedDto
	tFileUploadFinishedDto.Files = append(tFileUploadFinishedDto.Files, tFileUploadFinishedFileDto)

	// SET QRCLIP FILE UPLOAD HAS FINISHED
	tFileUploadFinishedResponseDto := fileUploadFinished(pClipDto, tFileUploadFinishedDto)
	if !tFileUploadFinishedResponseDto.Ok {
		ExitWithError("Failed to upload file!")
	}
}

// getUpdateClipDtoObject ////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getUpdateClipDtoObject(pFilePath string, pExpiration int, pMaxTransfers int, pAllowDelete bool) UpdateClipDto {
	var tUpdateClipDto UpdateClipDto
	tUpdateClipDto.ExpiresInMinutes = pExpiration
	tUpdateClipDto.MaxTransfers = pMaxTransfers
	tUpdateClipDto.AllowDelete = pAllowDelete
	tUpdateClipDto.FileSize = 0

	// HANDLE FILE IF EXISTS
	if pFilePath != "" {
		tFileStat, tErr := os.Stat(pFilePath)
		if tErr != nil {
			ExitWithError("File not found " + pFilePath)
		}
		// SET THE FILE SIZE TO SEND
		tUpdateClipDto.FileSize = tFileStat.Size()
		// SEND FIRST CHUNK SIZE, IF FILE BIGGER THAN CHUNK SIZE SEND CHUNK SIZE IF NOT SEND THE FILE SIZE
		tUpdateClipDto.FirstChunkSize = tFileStat.Size()
		if tUpdateClipDto.FirstChunkSize > int64(gFileChunkSizeBytes) {
			tUpdateClipDto.FirstChunkSize = int64(gFileChunkSizeBytes)
		}
	}
	return tUpdateClipDto
}

// getFileChunkUploadLink //////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getFileChunkUploadLink(pClipDto ClipDto, pGetFileChunkUploadLink GetFileChunkUploadLink) FileChunkUploadLinkResponse {
	tErrorPrefix := "Get File Chunk Upload, "
	tJson, tErr := json.Marshal(pGetFileChunkUploadLink)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tUrl := gApiUrl + "/clips/" + pClipDto.Id + "/" + pClipDto.SubId + "/file-chunk"

	// CREATE REQUEST
	tRequest, tErr := http.NewRequest(http.MethodPut, tUrl, bytes.NewBuffer(tJson))
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	tRequest.Header.Set("Content-Type", "application/json")

	// SEND REQUEST
	tClient := &http.Client{}
	tResponse, tErr := tClient.Do(tRequest)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer tResponse.Body.Close()

	// PARSE  RESPONSE
	var tFileChunkUploadLinkResponse FileChunkUploadLinkResponse
	tErr = json.NewDecoder(tResponse.Body).Decode(&tFileChunkUploadLinkResponse)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	return tFileChunkUploadLinkResponse
}

// fileUploadFinished //////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func fileUploadFinished(pClipDto ClipDto, pFileUploadFinishedDto FileUploadFinishedDto) FileUploadFinishedResponseDto {
	tErrorPrefix := "Setting file upload finished, "
	tJson, tErr := json.Marshal(pFileUploadFinishedDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tUrl := gApiUrl + "/clips/" + pClipDto.Id + "/" + pClipDto.SubId + "/file-upload-finished"

	// CREATE REQUEST
	tRequest, tErr := http.NewRequest(http.MethodPut, tUrl, bytes.NewBuffer(tJson))
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	tRequest.Header.Set("Content-Type", "application/json")

	// SEND REQUEST
	tClient := &http.Client{}
	tResponse, tErr := tClient.Do(tRequest)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer tResponse.Body.Close()

	// PARSE  RESPONSE
	var tFileUploadFinishedResponseDto FileUploadFinishedResponseDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tFileUploadFinishedResponseDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	return tFileUploadFinishedResponseDto
}

// addFormField ////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func addFormField(pWriter *multipart.Writer, pFieldName string, pFieldValue string) error {
	tFw, tErr := pWriter.CreateFormField(pFieldName)
	if tErr != nil {
		return tErr
	}
	_, tErr = io.Copy(tFw, strings.NewReader(pFieldValue))
	if tErr != nil {
		return tErr
	}
	return nil
}

// getFormDataForFileUpload ////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getFormDataForFileUpload(tBufferFile []byte, pS3PreSignedPost S3PreSignedPost) (*bytes.Buffer, string) {
	tErrorPrefix := "Creating form data for file upload, "

	tBuffer := &bytes.Buffer{}
	tWriter := multipart.NewWriter(tBuffer)

	tErr := addFormField(tWriter, "key", pS3PreSignedPost.Fields.Key)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tErr = addFormField(tWriter, "bucket", pS3PreSignedPost.Fields.Bucket)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tErr = addFormField(tWriter, "X-Amz-Algorithm", pS3PreSignedPost.Fields.XAmzAlgorithm)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tErr = addFormField(tWriter, "X-Amz-Credential", pS3PreSignedPost.Fields.XAmzCredential)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tErr = addFormField(tWriter, "X-Amz-Date", pS3PreSignedPost.Fields.XAmzDate)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tErr = addFormField(tWriter, "Policy", pS3PreSignedPost.Fields.Policy)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tErr = addFormField(tWriter, "X-Amz-Signature", pS3PreSignedPost.Fields.XAmzSignature)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tFw, tErr := tWriter.CreateFormField("file")
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	_, tErr = tFw.Write(tBufferFile)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	tWriter.Close()

	return tBuffer, tWriter.FormDataContentType()
}

// uploadChunk /////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func uploadChunk(pS3PreSignedPost S3PreSignedPost, tBuffer []byte, tBar *pb.ProgressBar) {
	tErrorPrefix := "Failed to upload file, "

	tBody, tFormDataContentType := getFormDataForFileUpload(tBuffer, pS3PreSignedPost)

	tReader := io.Reader(tBody)
	tPr := &ProgressReader{tReader, tBar}

	// CREATE REQUEST
	tRequest, tErr := http.NewRequest("POST", pS3PreSignedPost.Url, tPr)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	tRequest.Header.Set("Content-Type", tFormDataContentType)
	tRequest.ContentLength = int64(tBody.Len()) // IT'S NEEDED FOR S3

	// SEND REQUEST
	tClient := &http.Client{}
	tResponse, tErr := tClient.Do(tRequest)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer tResponse.Body.Close()
}

// uploadFileChunkByChunk //////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func uploadFileChunkByChunk(tClipDto ClipDto, pS3PreSignedPost S3PreSignedPost, pFileSize int64) int {
	tChunkCount := int(math.Ceil(float64(pFileSize) / float64(gFileChunkSizeBytes)))

	ShowInfo("UPLOADING ENCRYPTED FILE")

	tFile, err := os.Open(tClipDto.Id)
	if err != nil {
		fmt.Println("FAILED TO OPEN FILE")
		return 0
	}
	defer tFile.Close()

	tBar := pb.ProgressBarTemplate(gProgressBarTemplate).Start64(pFileSize)
	defer tBar.Finish()

	tBufferSize := int64(gFileChunkSizeBytes)
	if pFileSize < tBufferSize {
		tBufferSize = pFileSize
	}
	tBuffer := make([]byte, tBufferSize)
	for tChunkIndex := 0; tChunkIndex < tChunkCount; tChunkIndex++ {
		tN, _ := tFile.Read(tBuffer)
		if tChunkIndex == 0 {
			// FIRST CHUNK USE THE PRE SIGNED POST RECEIVED WHEN UPDATING
			uploadChunk(pS3PreSignedPost, tBuffer, tBar)
		} else {
			// AFTER FIRST CHUNK ASK A NEW URL FOR EACH ONE
			var tGetFileChunkUploadLink GetFileChunkUploadLink
			tGetFileChunkUploadLink.ChunkIndex = tChunkIndex
			tGetFileChunkUploadLink.FileIndex = 0
			tGetFileChunkUploadLink.Size = int64(tN)
			tNewLink := getFileChunkUploadLink(tClipDto, tGetFileChunkUploadLink)

			// UPLOAD USING THE NEW LINK
			uploadChunk(tNewLink.PreSignedPost, tBuffer[0:tN], tBar)
		}
	}
	tBar.Finish()

	tFile.Close()
	RemoveFile(tClipDto.Id)

	return tChunkCount
}

// updateQRClip ////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func updateQRClip(pClipDto ClipDto, pUpdateClipDto UpdateClipDto) UpdateClipResponseDto {
	tErrorPrefix := "Updating QRClip, "
	tJwt := CheckJwtToken()

	tJson, tErr := json.Marshal(pUpdateClipDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tUrl := gApiUrl + "/clips/" + pClipDto.Id + "/" + pClipDto.SubId
	if tJwt != "" {
		tUrl = tUrl + "/user"
	}

	// CREATE REQUEST
	tRequest, tErr := http.NewRequest(http.MethodPut, tUrl, bytes.NewBuffer(tJson))
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	tRequest.Header.Set("Content-Type", "application/json")

	// IN CASE THERE'S A JWT TOKEN USE IT
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
	var tUpdateClipResponseDto UpdateClipResponseDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tUpdateClipResponseDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	return tUpdateClipResponseDto
}
