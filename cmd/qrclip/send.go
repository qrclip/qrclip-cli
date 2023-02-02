package main

import (
	"bytes"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// SendQRClip //////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SendQRClip(pFilePath string, pMessage string, pExpiration int, pMaxTransfers int, pAllowDelete bool) {
	ShowInfo("SENDING QRCLIP")

	// GENERATE THE ENCRYPTION KEY
	tKey := GenerateEncryptionKey()

	// GET UPDATE CLIP DTO
	tUpdateClipDto := getUpdateClipDtoObject(pFilePath, pExpiration, pMaxTransfers, pAllowDelete)

	// CHECK LIMITS FOR FILE SIZE AND TEXT SIZE ALSO
	CheckIfCanBeSent(&tUpdateClipDto) // EXITS THE PROGRAM IF NOT OK

	// CREATE QRCLIP
	tClipDto, tErr := CreateQRClip(false)
	if tErr != nil {
		ExitWithError("Send error: " + tErr.Error())
	}

	// GENERATE IV DATA
	tFilesChunkNumber := make([]int, 1)
	tFilesChunkNumber[0] = int(math.Ceil(float64(tUpdateClipDto.FileSize) / float64(gFileChunkSizeBytes)))

	tQRClipIVGenerator := CreateQRClipIVGenerator(tClipDto.SubId, 24, 2500)

	tIV, tErr := GetIV(&tQRClipIVGenerator, 0)
	if tErr != nil {
		ExitWithError(tErr.Error())
	}

	// ENCRYPT TEXT MESSAGE IF EXISTS
	if pMessage != "" {
		tUpdateClipDto.EncryptedText = EncryptText(pMessage, tKey, tIV)
	}

	// UPDATE QRCLIP
	tUpdateClipResponseDto, tErr := updateQRClip(tClipDto, tUpdateClipDto)
	if tErr != nil {
		ExitWithError("Error updating QRClip!")
	}

	// UPLOAD FILE
	if tUpdateClipDto.FileSize > 0 {
		// UPLOAD FILE
		tChunkCount := uploadFileChunkByChunk(tClipDto, tUpdateClipResponseDto, tUpdateClipDto.FileSize, pFilePath, tKey, &tQRClipIVGenerator)
		if tChunkCount == 0 {
			ExitWithError("Error uploading file!")
		}

		tIV, tErr := GetIV(&tQRClipIVGenerator, 1)
		if tErr != nil {
			ExitWithError(tErr.Error())
		}

		// SET FILE UPLOAD FINISHED
		setFileUploadFinished(pFilePath, tKey, tClipDto, tUpdateClipDto, tChunkCount, tIV)
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
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func setFileUploadFinished(pFilePath string, pKey string, pClipDto ClipDto, pUpdateClipDto UpdateClipDto,
	pChunkCount int, pIV string) {
	var tFileUploadFinishedFileDto FileUploadFinishedFileDto
	tFileUploadFinishedFileDto.Name = EncryptText(pFilePath, pKey, pIV)
	tFileUploadFinishedFileDto.Index = 0
	tFileUploadFinishedFileDto.Size = pUpdateClipDto.FileSize
	tFileUploadFinishedFileDto.ChunkCount = pChunkCount

	var tFileUploadFinishedDto FileUploadFinishedDto
	tFileUploadFinishedDto.Files = append(tFileUploadFinishedDto.Files, tFileUploadFinishedFileDto)

	// SET QRCLIP FILE UPLOAD HAS FINISHED
	tFileUploadFinishedResponseDto, tErr := fileUploadFinished(pClipDto, tFileUploadFinishedDto)
	if tErr != nil || !tFileUploadFinishedResponseDto.Ok {
		ExitWithError("Failed to upload file!")
	}
}

// getUpdateClipDtoObject //////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getUpdateClipDtoObject(pFilePath string, pExpiration int, pMaxTransfers int, pAllowDelete bool) UpdateClipDto {
	var tUpdateClipDto UpdateClipDto
	tUpdateClipDto.ExpiresInMinutes = pExpiration
	tUpdateClipDto.MaxTransfers = pMaxTransfers
	tUpdateClipDto.AllowDelete = pAllowDelete
	tUpdateClipDto.FileSize = 0
	tUpdateClipDto.Version = gClientVersion
	tUpdateClipDto.MacSize = gMACSize

	// SET STORAGE
	tConfig, tError := GetQRClipConfig()
	if tError == nil {
		if tConfig.Storage != "" {
			tUpdateClipDto.Storage = tConfig.Storage
		} else {
			tUpdateClipDto.Storage = "storj"
		}
	}

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

		tUpdateClipDto.FirstChunkSize = tUpdateClipDto.FirstChunkSize + int64(gMACSize)
	}
	return tUpdateClipDto
}

// getFileChunkUploadLink //////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getFileChunkUploadLink(pClipDto ClipDto,
	pGetFileChunkUploadLink GetFileChunkUploadLink) (FileChunkUploadLinkResponse, error) {

	tUrlPath := "/clips/" + pClipDto.Id + "/" + pClipDto.SubId + "/file-chunk"

	// REQUEST
	tResponse, tErr := HttpDoPut(tUrlPath, "", pGetFileChunkUploadLink)
	if tErr != nil {
		return FileChunkUploadLinkResponse{}, tErr
	}

	// PARSE  RESPONSE
	var tFileChunkUploadLinkResponse FileChunkUploadLinkResponse
	tErr = DecodeJSONResponse(tResponse, &tFileChunkUploadLinkResponse)
	if tErr != nil {
		return FileChunkUploadLinkResponse{}, tErr
	}

	return tFileChunkUploadLinkResponse, nil
}

// fileUploadFinished //////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func fileUploadFinished(pClipDto ClipDto, pFileUploadFinishedDto FileUploadFinishedDto) (FileUploadFinishedResponseDto, error) {
	tUrlPath := "/clips/" + pClipDto.Id + "/" + pClipDto.SubId + "/file-upload-finished"

	// REQUEST
	tResponse, tErr := HttpDoPut(tUrlPath, "", pFileUploadFinishedDto)
	if tErr != nil {
		return FileUploadFinishedResponseDto{}, tErr
	}

	// PARSE  RESPONSE
	var tFileUploadFinishedResponseDto FileUploadFinishedResponseDto
	tErr = DecodeJSONResponse(tResponse, &tFileUploadFinishedResponseDto)
	if tErr != nil {
		return FileUploadFinishedResponseDto{}, tErr
	}

	return tFileUploadFinishedResponseDto, nil
}

// addFormField ////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getFormDataForFileUpload(pBufferFile []byte, pS3PreSignedPost S3PreSignedPost) (*bytes.Buffer, string) {
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
	_, tErr = tFw.Write(pBufferFile)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	err := tWriter.Close()
	if err != nil {
		ExitWithError(err.Error())
	}

	return tBuffer, tWriter.FormDataContentType()
}

// uploadChunkPost /////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func uploadChunkPost(pS3PreSignedPost S3PreSignedPost, pBuffer []byte, pBar *pb.ProgressBar) {
	tErrorPrefix := "Failed to upload file, "

	tBody, tFormDataContentType := getFormDataForFileUpload(pBuffer, pS3PreSignedPost)

	tReader := io.Reader(tBody)
	tPr := &ProgressReader{tReader, pBar}

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
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(tResponse.Body)
}

// uploadChunkPost /////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func uploadChunkPut(pUrl string, pBuffer []byte, pBar *pb.ProgressBar) {
	tErrorPrefix := "Failed to upload file, "

	tReader := bytes.NewReader(pBuffer)
	tPr := &ProgressReader{tReader, pBar}

	// CREATE REQUEST
	tRequest, tErr := http.NewRequest("PUT", pUrl, tPr)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	tRequest.ContentLength = int64(len(pBuffer)) // IT'S NEEDED FOR S3

	// SEND REQUEST
	tClient := &http.Client{}
	tResponse, tErr := tClient.Do(tRequest)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(tResponse.Body)
}

// uploadFileChunkByChunk //////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func uploadFileChunkByChunk(pClipDto ClipDto, pUpdateClipResponseDto UpdateClipResponseDto, pFileSize int64, pFilePath string,
	pKey string, pQRClipIVGenerator *QRClipIVGenerator) int {
	tChunkCount := int(math.Ceil(float64(pFileSize) / float64(gFileChunkSizeBytes)))

	ShowInfo("UPLOADING ENCRYPTED FILE")

	tFile, err := os.Open(pFilePath)
	if err != nil {
		fmt.Println("FAILED TO OPEN FILE")
		return 0
	}

	tBar := pb.ProgressBarTemplate(gProgressBarTemplate).Start64(pFileSize)
	defer tBar.Finish()

	tBufferSize := int64(gFileChunkSizeBytes)
	if pFileSize < tBufferSize {
		tBufferSize = pFileSize
	}
	tBuffer := make([]byte, tBufferSize)
	for tChunkIndex := 0; tChunkIndex < tChunkCount; tChunkIndex++ {
		tIV, tErr := GetIV(pQRClipIVGenerator, int64(tChunkIndex+2))
		if tErr != nil {
			ExitWithError(tErr.Error())
		}

		tN, _ := tFile.Read(tBuffer)
		if tChunkIndex == 0 {

			// ENCRYPT BUFFER
			tEncryptedBuffer := EncryptBuffer(tBuffer, pKey, tIV)

			// UPLOAD - FIRST CHUNK USE THE PRE SIGNED POST RECEIVED WHEN UPDATING
			if pUpdateClipResponseDto.PreSignedPut != "" {
				uploadChunkPut(pUpdateClipResponseDto.PreSignedPut, tEncryptedBuffer, tBar)
			} else {
				uploadChunkPost(pUpdateClipResponseDto.PreSignedPost, tEncryptedBuffer, tBar)
			}

		} else {
			// ENCRYPT BUFFER
			tEncryptedBuffer := EncryptBuffer(tBuffer[0:tN], pKey, tIV)

			// UPLOAD - AFTER FIRST CHUNK ASK A NEW URL FOR EACH ONE
			var tGetFileChunkUploadLink GetFileChunkUploadLink
			tGetFileChunkUploadLink.ChunkIndex = tChunkIndex
			tGetFileChunkUploadLink.FileIndex = 0
			tGetFileChunkUploadLink.Size = int64(tN + gMACSize)
			tGetFileChunkUploadLink.MacSize = gMACSize
			tNewLink, tErr := getFileChunkUploadLink(pClipDto, tGetFileChunkUploadLink)
			if tErr != nil {
				ExitWithError(err.Error())
			}
			if tNewLink.PreSignedPut != "" {
				uploadChunkPut(tNewLink.PreSignedPut, tEncryptedBuffer, tBar)
			} else {
				uploadChunkPost(tNewLink.PreSignedPost, tEncryptedBuffer, tBar)
			}
		}
	}
	tBar.Finish()
	err = tFile.Close()
	if err != nil {
		ExitWithError(err.Error())
	}

	return tChunkCount
}

// updateQRClip ////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func updateQRClip(pClipDto ClipDto, pUpdateClipDto UpdateClipDto) (UpdateClipResponseDto, error) {
	tJwt := CheckJwtToken()

	tUrlPath := "/clips/" + pClipDto.Id + "/" + pClipDto.SubId
	if tJwt != "" {
		tUrlPath = tUrlPath + "/user"
	}

	// REQUEST
	tResponse, tErr := HttpDoPut(tUrlPath, tJwt, pUpdateClipDto)
	if tErr != nil {
		return UpdateClipResponseDto{}, tErr
	}

	// PARSE RESPONSE
	var tUpdateClipResponseDto UpdateClipResponseDto
	tErr = DecodeJSONResponse(tResponse, &tUpdateClipResponseDto)
	if tErr != nil {
		return UpdateClipResponseDto{}, tErr
	}

	return tUpdateClipResponseDto, nil
}
