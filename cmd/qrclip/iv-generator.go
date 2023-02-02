package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

type QRClipIVGenerator struct {
	SubId          string
	UniqueKey      string
	Factorials     []int64
	Max            int64
	Multiplier     int64
	Size           int
	TotalIVs       int64
	NumberOfDigits int
}

const cMaxSize = 20

// CreateQRClipIVGenerator /////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func CreateQRClipIVGenerator(pSubID string, pSize int, pTotalIVs int) QRClipIVGenerator {
	var tQRClipIVGenerator QRClipIVGenerator
	tQRClipIVGenerator.NumberOfDigits = len(strconv.Itoa(pTotalIVs))
	tQRClipIVGenerator.Size = pSize - tQRClipIVGenerator.NumberOfDigits
	if tQRClipIVGenerator.Size > cMaxSize {
		tQRClipIVGenerator.NumberOfDigits += tQRClipIVGenerator.Size - cMaxSize
		tQRClipIVGenerator.Size = cMaxSize
	}

	tQRClipIVGenerator.SubId = pSubID
	tQRClipIVGenerator.TotalIVs = int64(pTotalIVs)

	tQRClipIVGenerator.buildUniqueKey()

	tQRClipIVGenerator.buildFactorials()

	return tQRClipIVGenerator
}

// buildUniqueKey //////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (pQRClipIVGenerator *QRClipIVGenerator) buildUniqueKey() {
	pQRClipIVGenerator.UniqueKey = ""
	for i := 0; i < len(pQRClipIVGenerator.SubId); i++ {
		if !strings.Contains(pQRClipIVGenerator.UniqueKey, string(pQRClipIVGenerator.SubId[i])) {
			pQRClipIVGenerator.UniqueKey = pQRClipIVGenerator.UniqueKey + string(pQRClipIVGenerator.SubId[i])
			if len(pQRClipIVGenerator.UniqueKey) == pQRClipIVGenerator.Size {
				return
			}
		}
	}

	//goland:noinspection SpellCheckingInspection
	tAllChars := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < len(tAllChars); i++ {
		if !strings.Contains(pQRClipIVGenerator.UniqueKey, string(tAllChars[i])) {
			pQRClipIVGenerator.UniqueKey = pQRClipIVGenerator.UniqueKey + string(tAllChars[i])
			if len(pQRClipIVGenerator.UniqueKey) > pQRClipIVGenerator.Size {
				return
			}
		}
	}

	pQRClipIVGenerator.UniqueKey = ""
}

// buildFactorials /////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (pQRClipIVGenerator *QRClipIVGenerator) buildFactorials() {
	pQRClipIVGenerator.Factorials = make([]int64, 0)
	pQRClipIVGenerator.Factorials = append(pQRClipIVGenerator.Factorials, 1)
	pQRClipIVGenerator.Max = 0
	for i := 1; i <= pQRClipIVGenerator.Size; i++ {
		pQRClipIVGenerator.Factorials = append(pQRClipIVGenerator.Factorials, pQRClipIVGenerator.Factorials[i-1]*int64(i))
		pQRClipIVGenerator.Max = pQRClipIVGenerator.Factorials[i]
	}

	if pQRClipIVGenerator.TotalIVs > pQRClipIVGenerator.Max {
		pQRClipIVGenerator.TotalIVs = pQRClipIVGenerator.Max
	}
	pQRClipIVGenerator.Multiplier = int64(math.Floor(float64(pQRClipIVGenerator.Max) / float64(pQRClipIVGenerator.TotalIVs)))
}

// GetIV ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GetIV(pQRClipIVGenerator *QRClipIVGenerator, pIndex int64) (string, error) {
	if pIndex > pQRClipIVGenerator.Max {
		return "", errors.New("get IV Index to big")
	}
	if pQRClipIVGenerator.UniqueKey == "" {
		return "", errors.New("iv generator unique key is empty")
	}

	return getPermutation(pQRClipIVGenerator, pIndex)
}

// getPermutation //////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getPermutation(pQRClipIVGenerator *QRClipIVGenerator, pPermutationNo int64) (string, error) {
	// PART OF THIS ALGORITHM WAS TAKEN FROM HERE (Filip Nguyen)
	//https://stackoverflow.com/questions/2799078/permutation-algorithm-without-recursion-java/11471673#11471673

	tPermutationNoCopy := int(pPermutationNo) // MAKE A COPY FOR LATER
	pPermutationNo = pPermutationNo * pQRClipIVGenerator.Multiplier

	if pPermutationNo > pQRClipIVGenerator.Max {
		return "", errors.New("get IV Index to big")
	}

	// THE SIZE OF THE PERMUTED STRING PLUS THE NUMBER OF DIGITS
	tPermutation := make([]byte, pQRClipIVGenerator.Size+pQRClipIVGenerator.NumberOfDigits)
	tBase := pQRClipIVGenerator.UniqueKey

	tTemp := tBase
	tOnlyChars := make([]byte, 0)
	tSwitch := true

	// BUILD THE FIRST PART OF THE IV CONTAINING PERMUTED CHARACTERS BY THE INDEX
	for tPosition := len(tBase); tPosition > 0; tPosition-- {
		tSelected := pPermutationNo / pQRClipIVGenerator.Factorials[tPosition-1]
		tPermutation[tPosition-1] = tTemp[tSelected]
		pPermutationNo = pPermutationNo % pQRClipIVGenerator.Factorials[tPosition-1]
		tTemp = tTemp[0:tSelected] + tTemp[tSelected+1:]

		// THIS IS JUST TO MAKE A STRING CONTAINING ONLY CHARACTERS FOR LATER USE
		if !(tPermutation[tPosition-1] >= 48 && tPermutation[tPosition-1] <= 57) {
			if tSwitch {
				tOnlyChars = append([]byte{tPermutation[tPosition-1]}, tOnlyChars...)
			} else {
				tOnlyChars = append(tOnlyChars, tPermutation[tPosition-1])
			}
		} else {
			tSwitch = !tSwitch
		}
	}

	// NO WE ADD THE INDEX TO THE END OF THE IV AND REPLACE THE PADDED ZEROS WITH THE LETTERS FROM tOnlyChars
	if pQRClipIVGenerator.NumberOfDigits > 0 {
		tNumberText := strconv.Itoa(tPermutationNoCopy)
		tFillCharCount := pQRClipIVGenerator.NumberOfDigits - len(tNumberText)
		// GET THE LETTERS FROM ONLY CHARS TO PAD THE INDEX TO THE NEEDED SIZE
		for tP := 0; tP < tFillCharCount; tP++ {
			tPermutation[pQRClipIVGenerator.Size+tP] = tOnlyChars[(tPermutationNoCopy+tP)%len(tOnlyChars)]
		}
		// ADD JUST THE NUMBER PART
		for tP := 0; tP < len(tNumberText); tP++ {
			tPermutation[pQRClipIVGenerator.Size+tFillCharCount+tP] = tNumberText[tP]
		}
	}

	// EXAMPLE RESULT "ABCDEFGHIJKLM1000"

	return string(tPermutation), nil
}
