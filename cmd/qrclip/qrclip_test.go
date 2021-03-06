package main

import (
	"os"
	"testing"
)

// TestTextEncryption //////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestTextEncryption(t *testing.T) {
	tText := "Testing 123456789 10 11 finished"
	//goland:noinspection SpellCheckingInspection
	tKey := "YuIiUrhYAdWqdhXT/kl1t/ETxq/R7E+UXciBQt7ZJRM="
	//goland:noinspection SpellCheckingInspection
	tIV := "ABCDEFGIJKLMNOPQ"

	tEncryptedText := EncryptText(tText, tKey, tIV)
	//goland:noinspection SpellCheckingInspection
	if tEncryptedText != "qCBBayULihr2cEaJ3P9dDDqj1WSMeeV14Z4vbODkwbI=" {
		t.Errorf("Text encryption failed!")
	}
}

// TestTextDecryption //////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestTextDecryption(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tEncryptedText := "qCBBayULihr2cEaJ3P9dDDqj1WSMeeV14Z4vbODkwbI="
	//goland:noinspection SpellCheckingInspection
	tKey := "YuIiUrhYAdWqdhXT/kl1t/ETxq/R7E+UXciBQt7ZJRM="
	//goland:noinspection SpellCheckingInspection
	tIV := "ABCDEFGIJKLMNOPQ"

	tText := DecryptText(tEncryptedText, tKey, tIV)
	if tText != "Testing 123456789 10 11 finished" {
		t.Errorf("Text decryption failed!")
	}
}

// TestFileEncryption //////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestFileEncryption(t *testing.T) {
	tDecryptedFileName := "test-decrypted-a.file"
	tEncryptedFileName := "test-encrypted-a.file"
	//goland:noinspection SpellCheckingInspection
	tKey := "YuIiUrhYAdWqdhXT/kl1t/ETxq/R7E+UXciBQt7ZJRM="
	//goland:noinspection SpellCheckingInspection
	tIV := "ABCDEFGIJKLMNOPQ"

	// CREATE FILE
	tWriteBytes := []byte("Testing 123456789 10 11 finished")
	tErr := os.WriteFile(tDecryptedFileName, tWriteBytes, 0644)
	if tErr != nil {
		t.Errorf("Failed to create " + tDecryptedFileName)
	}

	// ENCRYPT FILE
	EncryptFile(tDecryptedFileName, tKey, tEncryptedFileName, 0, tIV)

	// READ FILE
	tReadBytes, tErr := os.ReadFile(tEncryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to read " + tEncryptedFileName)
	}
	tResult := EncodeBase64(tReadBytes)
	//goland:noinspection SpellCheckingInspection
	if tResult != "qCBBayULihr2cEaJ3P9dDDqj1WSMeeV14Z4vbODkwbI=" {
		t.Errorf("File decryption failed!")
	}

	// REMOVE FILES
	tErr = os.Remove(tDecryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to delete file " + tDecryptedFileName)
	}
	tErr = os.Remove(tEncryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to delete file " + tEncryptedFileName)
	}
}

// TestFileDecryption //////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestFileDecryption(t *testing.T) {
	tEncryptedFileName := "test-encrypted-b.file"
	tDecryptedFileName := "test-decrypted-b.file"
	//goland:noinspection SpellCheckingInspection
	tKey := "YuIiUrhYAdWqdhXT/kl1t/ETxq/R7E+UXciBQt7ZJRM="
	//goland:noinspection SpellCheckingInspection
	tIV := "ABCDEFGIJKLMNOPQ"

	// CREATE FILE
	//goland:noinspection SpellCheckingInspection
	tWriteBytes := DecodeBase64("qCBBayULihr2cEaJ3P9dDDqj1WSMeeV14Z4vbODkwbI=")
	tErr := os.WriteFile(tEncryptedFileName, tWriteBytes, 0644)
	if tErr != nil {
		t.Errorf("Failed to create " + tEncryptedFileName)
	}

	// ENCRYPT FILE
	DecryptFile(tEncryptedFileName, tKey, tDecryptedFileName, 0, tIV)

	// READ FILE
	tReadBytes, tErr := os.ReadFile(tDecryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to read " + tDecryptedFileName)
	}

	tResult := string(tReadBytes)
	if tResult != "Testing 123456789 10 11 finished" {
		t.Errorf("File decryption failed!")
	}

	// REMOVE FILES
	tErr = os.Remove(tDecryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to delete file " + tDecryptedFileName)
	}
	tErr = os.Remove(tEncryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to delete file " + tEncryptedFileName)
	}
}

// TestOfflineEncrypt //////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestOfflineEncrypt(t *testing.T) {
	tDecryptedFileName := "test-decrypted-offline-a.file"
	tEncryptedFileName := "test-decrypted-offline-a.file.enc"
	//goland:noinspection SpellCheckingInspection
	tKey := "K+bW7JuRetdfqN6Y+yuPgzrdbem0BIDPsIbpwb5heCI="

	// CREATE FILE
	tWriteBytes := []byte("Testing 123456789 10 11 finished")
	tErr := os.WriteFile(tDecryptedFileName, tWriteBytes, 0644)
	if tErr != nil {
		t.Errorf("Failed to create " + tDecryptedFileName)
	}

	OfflineEncrypt(tDecryptedFileName, tKey)

	// READ FILE
	tReadBytes, tErr := os.ReadFile(tEncryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to read " + tEncryptedFileName)
	}
	tResult := EncodeBase64(tReadBytes)
	//goland:noinspection SpellCheckingInspection
	if tResult != "gfR+7M5kAspDjYndprmf0hJKgHoUkeePZ2BK7vmcZ8k=" {
		t.Errorf("File decryption failed!")
	}

	// REMOVE FILES
	tErr = os.Remove(tDecryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to delete file " + tDecryptedFileName)
	}
	tErr = os.Remove(tEncryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to delete file " + tEncryptedFileName)
	}
}

// TestOfflineDecrypt //////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestOfflineDecrypt(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tKey := "K+bW7JuRetdfqN6Y+yuPgzrdbem0BIDPsIbpwb5heCI="
	tEncryptedFileName := "test-encrypted-offline.file.enc"
	tDecryptedFileName := "test-encrypted-offline.file"
	// CREATE FILE
	//goland:noinspection SpellCheckingInspection
	tWriteBytes := DecodeBase64("gfR+7M5kAspDjYndprmf0hJKgHoUkeePZ2BK7vmcZ8k=")
	tErr := os.WriteFile(tEncryptedFileName, tWriteBytes, 0644)
	if tErr != nil {
		t.Errorf("Failed to create " + tEncryptedFileName)
	}

	OfflineDecrypt(tEncryptedFileName, tKey)

	// READ FILE
	tReadBytes, tErr := os.ReadFile(tDecryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to read " + tDecryptedFileName)
	}

	tResult := string(tReadBytes)
	if tResult != "Testing 123456789 10 11 finished" {
		t.Errorf("Offline decryption failed!")
	}

	// REMOVE FILES
	tErr = os.Remove(tDecryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to delete file " + tDecryptedFileName)
	}
	tErr = os.Remove(tEncryptedFileName)
	if tErr != nil {
		t.Errorf("Failed to delete file " + tEncryptedFileName)
	}
}
