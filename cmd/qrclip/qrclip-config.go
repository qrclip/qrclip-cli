package main

// API URL
var gApiUrl = "https://api.qrclip.io"

//var gApiUrl = "http://localhost:3000"

// PROGRESS BAR TEMPLATE
var gProgressBarTemplate = `{{ " " }} {{ bar . "|" "-" (cycle . "-" "|" "-" "|" ) "." "|"}} {{percent . "%06.2f%%" "?"}}`

//FILE CHUNK SIZE
var gFileChunkSizeBytes = 1000 * 1024 * 50

//QRCODE WITH HALF BLOCKS - SMALLER QRCODE (doesn't work on Windows)
var gHalfBlocks = false
