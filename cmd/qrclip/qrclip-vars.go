package main

// API URL
var gApiUrl = "https://api.qrclip.io"
var gSpaUrl = "https://www.qrclip.io"

//var gApiUrl = "http://localhost:3000"
//var gSpaUrl = "http://localhost"

// PROGRESS BAR TEMPLATE
var gProgressBarTemplate = `{{ " " }} {{ bar . "|" "-" (cycle . "-" "|" "-" "|" ) "." "|"}} {{percent . "%06.2f%%" "?"}}`

// FILE CHUNK SIZE
var gFileChunkSizeBytes = 1000 * 1024 * 50

// QRCODE WITH HALF BLOCKS - SMALLER QRCODE (doesn't work on Windows)
var gHalfBlocks = false

// IV TRIES
var gIVTries = 10

// QRCLIP VERSION
var gClientVersion = 3
