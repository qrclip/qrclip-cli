package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

// DecryptText /////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DecryptText(pTextToDecrypt string, pKey string, pIV string) string {
	tKey, tErr := base64.StdEncoding.DecodeString(pKey)
	if tErr != nil {
		ExitWithError("Base64 key failed to decode.")
	}

	tErrorPrefix := "Decrypting text, "
	tText := DecodeBase64(pTextToDecrypt)

	tBlock, tErr := aes.NewCipher(tKey)
	if tErr != nil {
		ExitWithError(tErrorPrefix + "creating cipher.")
	}

	tStream := cipher.NewCFBDecrypter(tBlock, []byte(pIV[:aes.BlockSize]))
	tStream.XORKeyStream(tText, tText)

	return string(tText)
}

// DecryptFile /////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DecryptFile(pSrcFile string, pKey string, pDstFile string, pFileSize int64, pIV string) {
	tKey, tErr := base64.StdEncoding.DecodeString(pKey)
	if tErr != nil {
		ExitWithError("Base64 key failed to decode.")
	}

	ShowInfo("DECRYPTING FILE")
	tInFile, tErr := os.Open(pSrcFile)
	if tErr != nil {
		ExitWithError("Decrypting file " + pSrcFile + " , failed to open")
	}

	tBlock, tErr := aes.NewCipher(tKey)
	if tErr != nil {
		ExitWithError("Decrypting file, creating cipher")
	}

	tOutfile, tErr := os.OpenFile(pDstFile, os.O_RDWR|os.O_CREATE, 0777)
	if tErr != nil {
		ExitWithError("Decrypting file " + pDstFile + " , failed to create")
	}

	tBar := pb.ProgressBarTemplate(gProgressBarTemplate).Start64(pFileSize)

	tBuf := make([]byte, 1024)
	tStream := cipher.NewCFBDecrypter(tBlock, []byte(pIV[:aes.BlockSize]))
	for {
		tN, tErrR := tInFile.Read(tBuf)
		if tN > 0 {
			tStream.XORKeyStream(tBuf, tBuf[:tN])
			_, tErrW := tOutfile.Write(tBuf[:tN])
			if tErrW != nil {
				ExitWithError("Writing decrypted file, " + tErrW.Error())
			}
		}

		if tErrR == io.EOF {
			break
		}

		if tErrR != nil {
			ExitWithError("Reading encrypted file, " + tErrR.Error())
		}

		tBar.Add(tN)
	}
	tBar.Finish()
	tInFile.Close()
	tOutfile.Close()
}

// OfflineDecrypt //////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func OfflineDecrypt(pFile string, pKey string) {
	tFileStat, tErr := os.Stat(pFile)
	if tErr != nil {
		ExitWithError(pFile + " not found!")
	}

	tDstName := pFile + ".dec" // BY DEFAULT DECRYPTED FILE APPENDS .dec

	// BUT IF SOURCE FILE ENDS WITH .enc THAT IS REMOVED FROM THE NAME FOR THE FILE TO HAVE THE CORRECT EXTENSION
	tExtension := filepath.Ext(pFile)
	if tExtension == ".enc" {
		tDstName = strings.TrimSuffix(pFile, tExtension)
	}

	// CHECK IF FILE EXISTS
	_, tErr = os.Stat(tDstName)
	if tErr == nil {
		ShowWarning("A File already exists named: " + tDstName)
		ExitWithError("Please rename it or move it to decrypt the file.")
	}

	// FOR OFFLINE WE USE THE KEY FOR IV
	DecryptFile(pFile, pKey, tDstName, tFileStat.Size(), pKey)

	ShowSuccess("FILE WAS DECRYPTED TO: " + tDstName)
}
