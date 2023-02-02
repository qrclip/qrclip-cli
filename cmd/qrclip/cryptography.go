package main

import (
	"golang.org/x/crypto/chacha20poly1305"
)

// EncryptText /////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func EncryptText(pTextToEncrypt string, pKey string, pIV string) string {
	return EncodeBase64UrlNoPad(EncryptBuffer([]byte(pTextToEncrypt), pKey, pIV))
}

// EncryptBuffer   /////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func EncryptBuffer(pDecryptedBuffer []byte, pKey string, pIV string) []byte {
	tKey := DecodeBase64Url(pKey)

	tAEAD, tErr := chacha20poly1305.NewX(tKey)
	if tErr != nil {
		ExitWithError("Encrypting, using key")
	}

	return tAEAD.Seal(nil, []byte(pIV), pDecryptedBuffer, nil)
}

// DecryptText /////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DecryptText(pTextToDecrypt string, pKey string, pIV string) string {
	return string(DecryptBuffer(DecodeBase64Url(pTextToDecrypt), pKey, pIV))
}

// DecryptBuffer ///////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DecryptBuffer(pEncryptedBuffer []byte, pKey string, pIV string) []byte {
	tKey := DecodeBase64Url(pKey)

	tAEAD, tErr := chacha20poly1305.NewX(tKey)
	if tErr != nil {
		ExitWithError("Decrypting, using key")
	}

	tDecData, tDecErr := tAEAD.Open(nil, []byte(pIV), pEncryptedBuffer, nil)
	if tDecErr != nil {
		ExitWithError(tDecErr.Error())
	}

	return tDecData
}
