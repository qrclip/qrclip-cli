package main

import (
	"github.com/cheggaaa/pb/v3"
	"io"
)

type ProgressReader struct {
	io.Reader
	mBar *pb.ProgressBar
}

func (pPr *ProgressReader) Read(tBuffer []byte) (tN int, tErr error) {
	tN, tErr = pPr.Reader.Read(tBuffer)
	pPr.mBar.Add(tN)
	return
}
