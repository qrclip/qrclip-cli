package main

import "testing"

// TestIvGeneration ////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
