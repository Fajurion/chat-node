package util

import (
	"crypto/rand"
	"math/big"
)

const NODE_TOKEN = "etMbbM3urDxkBRtphGWdIQ2XboZVgFcTbndiHB9PBgBhJsgFSetxYmT9Z9Ig5qK6mR8TjaFDrmq5Fi7VrrulXhKi3dnVy3gKESYmHPulv1yTCmGtuUkDiE3awAOO5y8Mxi9sTOfUFJZBncEYJcA0RPAqLrj3QSqfySBtEuMrq4DhcjtD9xzqylq4TpCWUHXIc6WpFmeiUTvgtAtp0mAsuNcfYPlpLKptO2mfFOmgbMx2hPMwX1jJa6FOB2vQg1lwMyMjGPexb1pHki26JPJJmCunIlWVMJsAObF2lrIXe4Py"
const NODE_ID = 1

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateToken(tkLength int32) string {

	s := make([]rune, tkLength)

	length := big.NewInt(int64(len(letters)))

	for i := range s {

		number, _ := rand.Int(rand.Reader, length)
		s[i] = letters[number.Int64()]
	}

	return string(s)
}
