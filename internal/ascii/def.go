package ascii

// returns

// OK - OK string
const OK = `OK`

// ER - Error string
const ER = `ERROR`

// ascii

// CR - Carriage return
const CR = 0x0D

// LF - Line feed
const LF = 0x0A

// SP - Space
const SP = 0x20

// DEL - Delete
const DEL = 0x7F

// BS - Backspace
const BS = 0x08

// CRLF - CR + LF
const CRLF string = string(CR) + string(LF)
