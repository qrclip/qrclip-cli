package main

import (
	"math"
	"strconv"
	"unicode"
)

type QrcIVData struct {
	Text      string
	FileNames []string
	Files     []string
}

// calculateSeed ///////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func calculateSeed(pSubID string) int64 {
	var tSeed int64 = 0
	for _, tChar := range pSubID {
		if unicode.IsNumber(tChar) {
			tCharStr := string(tChar)
			tCharInt, _ := strconv.Atoi(tCharStr)
			tSeed = tSeed + int64(tCharInt)
		}
	}
	return tSeed
}

// getNextIndex ////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getNextIndex(pSeed int64, pIVLength int) (int64, int) {
	pSeed = (pSeed*1664525 + 1013904223) % 4294967296
	tIndexF := float64(pSeed) / 4294967296.0 * float64(pIVLength-1)
	tIndex := int(math.Round(tIndexF))
	return pSeed, tIndex
}

// getNextIV ///////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getNextIV(pSubID string, pSeed int64, pIVLength int) (int64, string) {
	tIV := ""
	for i := 0; i < pIVLength; i++ {
		tIndex := 0
		pSeed, tIndex = getNextIndex(pSeed, pIVLength)
		tIV = tIV + string(pSubID[tIndex])
	}
	return pSeed, tIV
}

// GenerateIVData //////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateIVData(pSubID string, pFileNumber int) QrcIVData {
	tIVLength := 16

	// CALCULATE SEED
	tSeed := calculateSeed(pSubID)

	// CALCULATE DATA
	var tQrcIVData QrcIVData
	tIV := ""

	tSeed, tIV = getNextIV(pSubID, tSeed, tIVLength)
	tQrcIVData.Text = tIV

	for i := 0; i < pFileNumber; i++ {
		tSeed, tIV = getNextIV(pSubID, tSeed, tIVLength)
		tQrcIVData.FileNames = append(tQrcIVData.FileNames, tIV)

		tSeed, tIV = getNextIV(pSubID, tSeed, tIVLength)
		tQrcIVData.Files = append(tQrcIVData.Files, tIV)
	}

	return tQrcIVData
}
