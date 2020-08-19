package oauthsocialite

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: time.Duration(10 * time.Second),
}

type Provider interface {
	UserFromToken(string) (user User, status int, error error)
}

func RequestGet(url string, headers map[string]string) (map[string]interface{}, error) {
	res, err := requestRun(url, "GET", headers)
	if err != nil {
		return nil, err
	}
	var parse map[string]interface{}
	json.NewDecoder(res.Body).Decode(&parse)
	parse["status"] = res.StatusCode
	return parse, nil
}

func requestRun(url, method string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	return response, err
}
func GenerateHashMac(secret string, data string) string {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(data))

	// Get result and encode as hexadecimal string
	return hex.EncodeToString(h.Sum(nil))
}
