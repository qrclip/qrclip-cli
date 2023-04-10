package main

import (
	"bytes"
	"testing"
)

// TestEncryptionAndDecryption /////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestTextEncryptionAndDecryption(t *testing.T) {
	tText := "Testing 123456789 10 11 finished"
	//goland:noinspection SpellCheckingInspection
	tKey := "Oryf-5fzB9xG3ZGhUEKBltmUx2Q0wFPKEgs442Bo6aM"
	//goland:noinspection SpellCheckingInspection
	tIV := "ABCDEFGIJKLMNOPQXCVFDSDF"

	tEncryptedText := EncryptText(tText, tKey, tIV)
	//goland:noinspection SpellCheckingInspection
	if tEncryptedText != "65lnyOSYZEEE-LbSZGTc1r-6iUiRjJomQ3WnEMPxxhSUoR9SRBlpAE86j9AOiMdM" {
		t.Errorf("Text encryption failed!")
	}

	tDecText := DecryptText(tEncryptedText, tKey, tIV)
	if tDecText != "Testing 123456789 10 11 finished" {
		t.Errorf("Text decryption failed!")
	}
}

// TestBufferEncryptionAndDecryption ///////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestBufferEncryptionAndDecryption(t *testing.T) {
	tBaseData := []byte{0, 100, 200, 250}
	tKey := GenerateEncryptionKey()
	//goland:noinspection SpellCheckingInspection
	tIV := "ABCDEFGIJKLMNOPQXCVFDSDF"

	tEncryptedData := EncryptBuffer(tBaseData, tKey, tIV)
	if len(tEncryptedData) != 20 {
		t.Errorf("Encrypted buffer wrong length")
	}

	tDecryptedData := DecryptBuffer(tEncryptedData, tKey, tIV)
	if len(tDecryptedData) != 4 {
		t.Errorf("Decrypted buffer wrong length")
	}

	if !bytes.Equal(tDecryptedData, tBaseData) {
		t.Errorf("Wrong decrypted data")
	}
}
