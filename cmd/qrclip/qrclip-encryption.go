package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

// EncryptText /////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func EncryptText(pTextToEncrypt string, pKey string, pIV string) string {
	tKey, tErr := base64.StdEncoding.DecodeString(pKey)
	if tErr != nil {
		ExitWithError("Base64 key failed to decode.")
	}

	tBlock, tErr := aes.NewCipher(tKey)
	if tErr != nil {
		ExitWithError("Encrypting text, creating cipher.")
	}

	tText := []byte(pTextToEncrypt)

	tEncryptedText := make([]byte, len(tText))

	tStream := cipher.NewCFBEncrypter(tBlock, []byte(pIV[:aes.BlockSize]))
	tStream.XORKeyStream(tEncryptedText, tText)

	return EncodeBase64(tEncryptedText)
}

// EncryptFile   ///////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func EncryptFile(pFile string, pKey string, pFileName string, pFileSize int64, pIV string) {
	tKey, tErr := base64.StdEncoding.DecodeString(pKey)
	if tErr != nil {
		ExitWithError("Base64 key failed to decode.")
	}

	ShowInfo("ENCRYPTING FILE")
	tInFile, tErr := os.Open(pFile)
	if tErr != nil {
		ExitWithError("Encrypting file " + pFile + " , failed to open")
	}

	tBlock, tErr := aes.NewCipher(tKey)
	if tErr != nil {
		ExitWithError("Encrypting file, creating cipher")
	}

	tOutFile, tErr := os.OpenFile(pFileName, os.O_RDWR|os.O_CREATE, 0777)
	if tErr != nil {
		ExitWithError("Encrypting file, failed to create file " + pFileName)
	}

	tBar := pb.ProgressBarTemplate(gProgressBarTemplate).Start64(pFileSize)
	defer tBar.Finish()

	tBuf := make([]byte, 1024)
	var tProgressCount int64
	tProgressCount = 0
	tStream := cipher.NewCFBEncrypter(tBlock, []byte(pIV[:aes.BlockSize]))
	for {
		tN, tErrR := tInFile.Read(tBuf)
		if tN > 0 {
			tStream.XORKeyStream(tBuf, tBuf[:tN])
			_, tErrW := tOutFile.Write(tBuf[:tN])
			if tErrW != nil {
				ExitWithError("Encrypting file, failed to write file")
			}
		}
		tProgressCount = tProgressCount + int64(tN)

		if tErrR == io.EOF {
			break
		}

		if tErrR != nil {
			ExitWithError("Encrypting file, failed to read file")
		}

		tBar.Add(tN)
	}
	tInFile.Close()
	tOutFile.Close()
}

// OfflineEncrypt //////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func OfflineEncrypt(pFile string, pKey string) {
	if pKey == "" {
		pKey = GenerateEncryptionKey()
	}
	tFileStat, tErr := os.Stat(pFile)
	if tErr != nil {
		ExitWithError(pFile + " not found!")
	}
	tIV := GenerateRandomIV()

	tHexIV := hex.EncodeToString([]byte(tIV))
	tEncryptedFile := pFile + "_iv_" + tHexIV + ".enc" // ENCRYPTED FILE NAME

	// FOR OFFLINE WE USE THE KEY FOR IV
	EncryptFile(pFile, pKey, tEncryptedFile, tFileStat.Size(), tIV)

	ShowSuccess("ENCRYPTED FILE: " + tEncryptedFile)
	ShowSuccess("ENCRYPTION KEY: " + pKey)
}

// EncryptBuffer   /////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func EncryptBuffer(pDecryptedBuffer []byte, pKey string, pIV string) []byte {
	tKey, tErr := base64.StdEncoding.DecodeString(pKey)
	if tErr != nil {
		ExitWithError("Base64 key failed to decode.")
	}

	tBlock, tErr := aes.NewCipher(tKey)
	if tErr != nil {
		ExitWithError("Encrypting file, creating cipher")
	}

	tEncryptedBuf := make([]byte, len(pDecryptedBuffer))

	tStream := cipher.NewCFBEncrypter(tBlock, []byte(pIV[:aes.BlockSize]))

	tStream.XORKeyStream(tEncryptedBuf, pDecryptedBuffer)

	return tEncryptedBuf
}
