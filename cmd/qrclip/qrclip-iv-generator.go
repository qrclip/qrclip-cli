package main

import (
	"fmt"
	"math"
	"strconv"
	"unicode"
)

type QrcIVData struct {
	Text      string
	FileNames []string
	Chunks    []string
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
func getNextIV(pSubID string, pSeed int64, pIVLength int, pTries int, pExisting map[string]bool, pCounter int) (int64, string) {
	tIV := ""
	pCounter = pCounter + 1
	for i := 0; i < pIVLength; i++ {
		tIndex := 0
		pSeed, tIndex = getNextIndex(pSeed, pIVLength)
		tIV = tIV + string(pSubID[tIndex])
	}

	// CHECK IF EXISTS ALREADY
	if pTries > 0 {
		if pExisting[tIV] {
			if pCounter == pTries {
				return pSeed, tIV
			}
			return getNextIV(pSubID, pSeed, pIVLength, pTries, pExisting, pCounter)
		} else {
			pExisting[tIV] = true
		}
	}

	return pSeed, tIV
}

// GenerateIVData //////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateIVData(pVersion int, pSubID string, pFilesChunkNumber []int) QrcIVData {
	switch pVersion {
	case 1:
		return GenerateIVDataV1(pSubID, pFilesChunkNumber, 16)
	case 2:
		return GenerateIVDataV2(pSubID, pFilesChunkNumber, 16)
	case 3:
		return GenerateIVDataV3(pSubID, pFilesChunkNumber, 16)
	}

	ExitWithError("Please update the client!")
	var tQrcIVData QrcIVData // NEVER GETS HERE
	return tQrcIVData
}

// GenerateIVDataV1 //////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateIVDataV1(pSubID string, pFilesChunkNumber []int, pIVLength int) QrcIVData {
	// CALCULATE V1 IVS, IT WAS THE SAME FOR EVERYTHING
	var tQrcIVData QrcIVData
	tIV := pSubID[:pIVLength]

	tQrcIVData.Text = tIV
	for i := 0; i < len(pFilesChunkNumber); i++ {
		tQrcIVData.FileNames = append(tQrcIVData.FileNames, tIV)
		for k := 0; k < pFilesChunkNumber[i]; k++ {
			tQrcIVData.Chunks = append(tQrcIVData.Chunks, tIV)
		}
	}
	return tQrcIVData
}

// GenerateIVDataV2 //////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateIVDataV2(pSubID string, pFilesChunkNumber []int, pIVLength int) QrcIVData {
	// CALCULATE V2 IVS, WERE DIFFERENT FOR TEXT, FILENAMES AND FILES
	// ,BUT ALL CHUNKS FROM THE SAME FILE USED THE SAME IV

	tDummyIVs := map[string]bool{}
	tSeed := calculateSeed(pSubID)

	// CALCULATE DATA
	var tQrcIVData QrcIVData
	tIV := ""

	tSeed, tIV = getNextIV(pSubID, tSeed, pIVLength, 0, tDummyIVs, 0)
	tQrcIVData.Text = tIV

	for i := 0; i < len(pFilesChunkNumber); i++ {
		tSeed, tIV = getNextIV(pSubID, tSeed, pIVLength, 0, tDummyIVs, 0)
		tQrcIVData.FileNames = append(tQrcIVData.FileNames, tIV)

		tSeed, tIV = getNextIV(pSubID, tSeed, pIVLength, 0, tDummyIVs, 0)
		for k := 0; k < pFilesChunkNumber[i]; k++ {
			tQrcIVData.Chunks = append(tQrcIVData.Chunks, tIV)
		}
	}

	return tQrcIVData
}

// GenerateIVDataV3 //////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateIVDataV3(pSubID string, pFilesChunkNumber []int, pIVLength int) QrcIVData {
	// CALCULATE V3 IVS ARE DIFFERENT FOR EVERYTHING INCLUDING FOR EACH FILE CHUNK
	tCalculatedIVs := map[string]bool{}

	// CALCULATE SEED
	tSeed := calculateSeed(pSubID)

	// CALCULATE DATA
	var tQrcIVData QrcIVData
	tIV := ""

	tSeed, tIV = getNextIV(pSubID, tSeed, pIVLength, gIVTries, tCalculatedIVs, 0)
	tQrcIVData.Text = tIV

	// FIRST CALCULATE ALL THE FILE NAME IVS
	tTotalChunks := 0
	for i := 0; i < len(pFilesChunkNumber); i++ {
		tSeed, tIV = getNextIV(pSubID, tSeed, pIVLength, gIVTries, tCalculatedIVs, 0)
		tQrcIVData.FileNames = append(tQrcIVData.FileNames, tIV)
		tTotalChunks = tTotalChunks + pFilesChunkNumber[i]
	}

	// CALCULATE ALL THE CHUNK IVS
	for i := 0; i < tTotalChunks; i++ {
		tSeed, tIV = getNextIV(pSubID, tSeed, pIVLength, gIVTries, tCalculatedIVs, 0)
		tQrcIVData.Chunks = append(tQrcIVData.Chunks, tIV)
	}
	fmt.Println(tQrcIVData.Chunks)
	return tQrcIVData
}
