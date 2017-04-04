package main

import (
	"testing"
	"net/http/httptest"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize()
	code := m.Run()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestInvalidRequest(t *testing.T) {
	req, _ := http.NewRequest("POST", "/parse", strings.NewReader(""))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid Email Format" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid Email Format'. Got '%s'", m["error"])
	}
}

func TestValidRequest(t *testing.T) {
	dat, _ := os.Open("test/santa-test-email.txt")
	req, _ := http.NewRequest("POST", "/parse", dat)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["To"] != "santa@northpole.com" {
		t.Errorf("Expected the 'to' key of the response to be set to 'santa@northpole.com'. Got '%s'", m["To"])
	}
	if m["From"] != "\"L.L.Bean\" <llbean@e1.llbean.com>" {
		t.Errorf("Expected the 'to' key of the response to be set to '\"L.L.Bean\" <llbean@e1.llbean.com>'. Got '%s'", m["To"])
	}
	if m["Subject"] != "Climbing mountains. Breaking barriers." {
		t.Errorf("Expected the 'to' key of the response to be set to 'Climbing mountains. Breaking barriers.'. Got '%s'", m["To"])
	}
	if m["Message-ID"] != "<0.0.4D.537.1D29E4787E56956.0@omp.e1.llbean.com>" {
		t.Errorf("Expected the 'to' key of the response to be set to '<0.0.4D.537.1D29E4787E56956.0@omp.e1.llbean.com>'. Got '%s'", m["To"])
	}
	if m["Date"] != "Thu, 16 Mar 2017 04:22:00 -0700" {
		t.Errorf("Expected the 'to' key of the response to be set to 'Thu, 16 Mar 2017 04:22:00 -0700'. Got '%s'", m["To"])
	}
}