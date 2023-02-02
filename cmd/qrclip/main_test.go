package main

import (
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

// TestIvGeneration /////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestIvGeneration(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tSubID := "IjjQwUqHbqqOliBz9EXcKOjCf5iGanW4"
	tQRClipIVGenerator := CreateQRClipIVGenerator(tSubID, 24, 2500)

	tResult0, _ := GetIV(&tQRClipIVGenerator, 0)
	//goland:noinspection SpellCheckingInspection
	if tResult0 != "fCKcXE9zBilObHqUwQjIzBi0" {
		t.Errorf("IV INDEX 0")
	}

	tResult1, _ := GetIV(&tQRClipIVGenerator, 1)
	//goland:noinspection SpellCheckingInspection
	if tResult1 != "9zlQKEicbHfwOCqBXUjIlQK1" {
		t.Errorf("IV INDEX 1")
	}

	tResult2, _ := GetIV(&tQRClipIVGenerator, 2)
	//goland:noinspection SpellCheckingInspection
	if tResult2 != "fOUwXzK9iqCHEclQBbjICHE2" {
		t.Errorf("IV INDEX 2")
	}

	tResult50, _ := GetIV(&tQRClipIVGenerator, 50)
	//goland:noinspection SpellCheckingInspection
	if tResult50 != "fCKXEilbHqUwQj9BczOIli50" {
		t.Errorf("IV INDEX 50")
	}

	tResult1003, _ := GetIV(&tQRClipIVGenerator, 1003)
	//goland:noinspection SpellCheckingInspection
	if tResult1003 != "CKwjXlQHBqcbfE9zUiIO1003" {
		t.Errorf("IV INDEX 50")
	}

}
