package main

import "golang.org/x/crypto/chacha20poly1305"

// EncryptText /////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func EncryptText(pTextToEncrypt string, pKey string, pIV string) string {
	return EncodeBase64UrlNoPad(EncryptBuffer([]byte(pTextToEncrypt), pKey, pIV))
}

// DecryptText /////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DecryptText(pTextToDecrypt string, pKey string, pIV string) string {
	if pTextToDecrypt == "" {
		return ""
	}

	tText2Decrypt, tErr := DecodeBase64Url(pTextToDecrypt)
	if tErr != nil {
		ExitWithError("Decoding Base64")
	}

	return string(DecryptBuffer(tText2Decrypt, pKey, pIV))
}

// EncryptBuffer   /////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func EncryptBuffer(pDecryptedBuffer []byte, pKey string, pIV string) []byte {
	tKey, tErr := DecodeBase64Url(pKey)
	if tErr != nil {
		ExitWithError("Decoding Base64")
	}

	tAEAD, tErr := chacha20poly1305.NewX(tKey)
	if tErr != nil {
		ExitWithError("Encrypting, using key")
	}

	return tAEAD.Seal(nil, []byte(pIV), pDecryptedBuffer, nil)
}

// DecryptBuffer ///////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DecryptBuffer(pEncryptedBuffer []byte, pKey string, pIV string) []byte {
	tKey, tErr := DecodeBase64Url(pKey)
	if tErr != nil {
		ExitWithError("Decoding Base64")
	}

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
