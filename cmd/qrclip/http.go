package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HttpDoGet ///////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HttpDoGet(pUrlPath string, pJwt string, pAccessKey *string) (*http.Response, error) {
	tRequest, tErr := newHttpGetRequest(pUrlPath, pJwt)
	if tErr != nil {
		return nil, tErr
	}

	if pAccessKey != nil {
		tRequest.Header.Set("qrclip-access-key", *pAccessKey)
	}

	tResponse, tErr := httpDoRequest(tRequest)
	if tErr != nil {
		return nil, tErr
	} else {
		return tResponse, nil
	}
}

// HttpDoGetWithTempTicket /////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HttpDoGetWithTempTicket(pUrlPath string, pJwt string, pAccessKey *string, pTempToken ClipTempTokenDto) (*http.Response, error) {
	tRequest, tErr := newHttpGetRequest(pUrlPath, pJwt)
	if tErr != nil {
		return nil, tErr
	}

	if pAccessKey != nil {
		tRequest.Header.Set("qrclip-access-key", *pAccessKey)
	}

	tRequest.Header.Set("qrclip-temp-token-id", pTempToken.Id)
	tRequest.Header.Set("qrclip-temp-token-key", pTempToken.Key)

	tResponse, tErr := httpDoRequest(tRequest)
	if tErr != nil {
		return nil, tErr
	} else {
		return tResponse, nil
	}
}

// HttpDoPut ///////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HttpDoPut(pUrlPath string, pJwt string, pBody interface{}) (*http.Response, error) {
	tJson, tErr := json.Marshal(pBody)
	if tErr != nil {
		return nil, tErr
	}

	tRequest, tErr := newHttpPutRequest(pUrlPath, pJwt, tJson)
	if tErr != nil {
		return nil, tErr
	}

	tResponse, tErr := httpDoRequest(tRequest)
	if tErr != nil {
		return nil, tErr
	} else {
		return tResponse, nil
	}
}

// HttpDoPost //////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HttpDoPost(pUrlPath string, pJwt string, pBody interface{}) (*http.Response, error) {
	tJson, tErr := json.Marshal(pBody)
	if tErr != nil {
		return nil, tErr
	}

	tRequest, tErr := newHttpPostRequest(pUrlPath, pJwt, tJson)
	if tErr != nil {
		return nil, tErr
	}

	tResponse, tErr := httpDoRequest(tRequest)
	if tErr != nil {
		return nil, tErr
	} else {
		return tResponse, nil
	}
}

// DecodeJSONResponse //////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DecodeJSONResponse(tResponse *http.Response, v interface{}) error {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(tResponse.Body)
	return json.NewDecoder(tResponse.Body).Decode(v)
}

// httpDoRequest ///////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func httpDoRequest(pRequest *http.Request) (*http.Response, error) {
	var tLastError error
	for i := 0; i < gHttpMaxRetries; i++ {
		tResponse, tErr := gHttpClient.Do(pRequest)
		if tErr == nil {
			var tResponseOK = false
			// IF BETWEEN 200 and 300 it's OK
			if tResponse.StatusCode >= 200 && tResponse.StatusCode < 300 {
				tResponseOK = true
			}
			// IF UNAUTHORIZED DO NOT RETRY
			if tResponse.StatusCode == 400 || tResponse.StatusCode == 401 {
				tResponseOK = true
			}

			if tResponseOK {
				if i > 0 {
					ShowSuccess("Request retry ok ...")
				}
				return tResponse, nil
			}
		}

		if tErr != nil {
			tLastError = tErr
		}

		// If we're here, it means we need to retry. Compute the delay.
		tDelay := computeRetryDelay(i + 1)
		ShowInfoYellow(fmt.Sprintf("Request failed, retrying in %.2f seconds...", tDelay.Seconds()))
		fmt.Println(tErr)
		time.Sleep(tDelay)
	}

	return nil, tLastError // if we're here, all retries have failed
}

// computeRetryDelay ///////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func computeRetryDelay(tRetry int) time.Duration {
	var tMulti = tRetry * tRetry
	if tMulti > 16 {
		tMulti = 16
	}
	return gHttpBaseDelay * time.Duration(tMulti)
}

// newHttpGetRequest ///////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func newHttpGetRequest(pUrlPath string, pJwt string) (*http.Request, error) {
	return newHttpRequest(http.MethodGet, pUrlPath, pJwt, nil)
}

// newHttpPutRequest ///////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func newHttpPutRequest(pUrlPath string, pJwt string, pBody []byte) (*http.Request, error) {
	return newHttpRequest(http.MethodPut, pUrlPath, pJwt, bytes.NewBuffer(pBody))
}

// newHttpPostRequest //////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func newHttpPostRequest(pUrlPath string, pJwt string, pBody []byte) (*http.Request, error) {
	return newHttpRequest(http.MethodPost, pUrlPath, pJwt, bytes.NewBuffer(pBody))
}

// newHttpRequest //////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func newHttpRequest(pMethod string, pUrlPath string, pJwt string, tBody io.Reader) (*http.Request, error) {
	tRequest, tErr := http.NewRequest(pMethod, gApiUrl+pUrlPath, tBody)
	if tErr != nil {
		return tRequest, tErr
	}

	if pJwt != "" {
		tRequest.Header.Add("Authorization", "Bearer "+pJwt)
	}
	tRequest.Header.Set("Content-Type", "application/json")
	return tRequest, nil
}
