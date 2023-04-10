package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testJwt = "test_jwt"

func TestHttpDoGet(t *testing.T) {
	testAccessKey := "test_access_key"
	// Set up a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request header
		if r.Header.Get("Authorization") != "Bearer "+testJwt {
			t.Errorf("expected Authorization header value: Bearer %s, got: %s", testJwt, r.Header.Get("Authorization"))
		}
		if r.Header.Get("qrclip-access-key") != testAccessKey {
			t.Errorf("expected qrclip-access-key header value: %s, got: %s", testAccessKey, r.Header.Get("qrclip-access-key"))
		}
		fmt.Fprintln(w, "OK")
	}))
	defer ts.Close()

	tCopyApiUrl := gApiUrl
	gApiUrl = ts.URL
	response, err := HttpDoGet("/test", testJwt, &testAccessKey)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	if strings.TrimSpace(string(body)) != "OK" {
		t.Errorf("expected response body: OK, got: %s", string(body))
	}

	gApiUrl = tCopyApiUrl
}

func TestHttpDoPut(t *testing.T) {
	// Set up a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request header
		if r.Header.Get("Authorization") != "Bearer "+testJwt {
			t.Errorf("expected Authorization header value: Bearer %s, got: %s", testJwt, r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header value: application/json, got: %s", r.Header.Get("Content-Type"))
		}

		body, _ := ioutil.ReadAll(r.Body)
		var data map[string]string
		_ = json.Unmarshal(body, &data)

		if data["test_key"] != "test_value" {
			t.Errorf("expected request body value: test_value, got: %s", data["test_key"])
		}
		fmt.Fprintln(w, "OK")
	}))
	defer ts.Close()

	tCopyApiUrl := gApiUrl
	gApiUrl = ts.URL
	response, err := HttpDoPut("/test", testJwt, map[string]string{"test_key": "test_value"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	if strings.TrimSpace(string(body)) != "OK" {
		t.Errorf("expected response body: OK, got: %s", string(body))
	}

	gApiUrl = tCopyApiUrl
}

func TestHttpDoPost(t *testing.T) {
	// Set up a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request header
		if r.Header.Get("Authorization") != "Bearer "+testJwt {
			t.Errorf("expected Authorization header value: Bearer %s, got: %s", testJwt, r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type header value: application/json, got: %s", r.Header.Get("Content-Type"))
		}

		body, _ := ioutil.ReadAll(r.Body)
		var data map[string]string
		_ = json.Unmarshal(body, &data)

		if data["test_key"] != "test_value" {
			t.Errorf("expected request body value: test_value, got: %s", data["test_key"])
		}
		fmt.Fprintln(w, "OK")
	}))
	defer ts.Close()

	tCopyApiUrl := gApiUrl
	gApiUrl = ts.URL
	response, err := HttpDoPost("/test", testJwt, map[string]string{"test_key": "test_value"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	if strings.TrimSpace(string(body)) != "OK" {
		t.Errorf("expected response body: OK, got: %s", string(body))
	}

	gApiUrl = tCopyApiUrl
}
