package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const cHttpResponseTimeoutSeconds = 10

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
	}

	// CHECK STATUS CODE
	if strings.HasPrefix(tResponse.Status, "2") {
		return tResponse, nil
	} else {
		return nil, errors.New(tResponse.Status)
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
	}

	// CHECK STATUS CODE
	if strings.HasPrefix(tResponse.Status, "2") {
		return tResponse, nil
	} else {
		return nil, errors.New(tResponse.Status)
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
	}

	// CHECK STATUS CODE
	if strings.HasPrefix(tResponse.Status, "2") {
		return tResponse, nil
	} else {
		return nil, errors.New(tResponse.Status)
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
	tClient := &http.Client{Timeout: cHttpResponseTimeoutSeconds * time.Second}
	return tClient.Do(pRequest)
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
	} else {
		if pJwt != "" {
			tRequest.Header.Add("Authorization", "Bearer "+pJwt)
		}
		tRequest.Header.Set("Content-Type", "application/json")
		return tRequest, nil
	}
}
