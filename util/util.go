package util

import (
	"crypto/rand"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
)

var Testing = false
var FilePath = "C:/Users/thisi/OneDrive/Dokumente/We(e)b Stuff/FJ/chat/env/node"

var NODE_TOKEN = "etMbbM3urDxkBRtphGWdIQ2XboZVgFcTbndiHB9PBgBhJsgFSetxYmT9Z9Ig5qK6mR8TjaFDrmq5Fi7VrrulXhKi3dnVy3gKESYmHPulv1yTCmGtuUkDiE3awAOO5y8Mxi9sTOfUFJZBncEYJcA0RPAqLrj3QSqfySBtEuMrq4DhcjtD9xzqylq4TpCWUHXIc6WpFmeiUTvgtAtp0mAsuNcfYPlpLKptO2mfFOmgbMx2hPMwX1jJa6FOB2vQg1lwMyMjGPexb1pHki26JPJJmCunIlWVMJsAObF2lrIXe4Py"
var NODE_ID = 1

const StatusOnline = 0
const StatusOffline = 1
const StatusError = 2

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

var Protocol = "http://"
var BasePath = "http://localhost:3000"

func PostRequest(url string, body map[string]interface{}) (map[string]interface{}, error) {
	return PostRaw(BasePath+url, body)
}

func PostRaw(url string, body map[string]interface{}) (map[string]interface{}, error) {

	req, _ := sonic.Marshal(body)

	reader := strings.NewReader(string(req))

	res, err := http.Post(url, "application/json", reader)
	if err != nil {
		return nil, err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)

	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = sonic.Unmarshal([]byte(buf.String()), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func First(a interface{}, _ interface{}) interface{} {
	return a
}
