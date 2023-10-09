package main

import "time"

var gVersion = "1.10"

// API URL
var gApiUrl = "https://api.qrclip.io"

const gSpaUrl = "https://app.qrclip.io"

//var gApiUrl = "http://localhost:3000"
//var gSpaUrl = "http://localhost"

// PROGRESS BAR TEMPLATE
const gProgressBarTemplate = `{{ " " }} {{ bar . "|" "-" (cycle . "-" "|" "-" "|" ) "." "|"}} {{percent . "%06.2f%%" "?"}}`

// FILE CHUNK SIZE
const gFileChunkSizeBytes = 1000 * 1024 * 50

// QRCODE WITH HALF BLOCKS - SMALLER QRCODE (doesn't work on Windows)
const gHalfBlocks = false

// QRCLIP VERSION
const gClientVersion = 6

// XChaCha20-Poly1305 MAC SIZE
const gMACSize = 16

const (
	gHttpResponseTimeoutSeconds = 10
	gHttpMaxRetries             = 5
	gHttpBaseDelay              = 1 * time.Second
)
