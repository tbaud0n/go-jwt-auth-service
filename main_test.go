package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	testTokenKey      = `abcdefghijklmnopqrstuvwxyz0123456789`
	testTokenIssuer   = `test_issuer`
	testTokenDuration = 42
	testTokenAudience = `test_audience`
)

func getJWTService() *JWTService {
	return &JWTService{
		TokenKey:            testTokenKey,
		TokenIssuer:         testTokenIssuer,
		TokenExpirationTime: time.Duration(testTokenDuration) * time.Second,
		TokenAudience:       testTokenAudience,
	}
}

func TestTokenNewSuccess(t *testing.T) {
	data := map[string]interface{}{
		"lastname":  "LASTNAME",
		"firstname": "FIRSTNAME",
		"age":       42,
	}

	b, _ := json.Marshal(data)

	res := doTokenNewRequest(t, bytes.NewReader(b), http.StatusOK)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	t.Logf("Generated token : %s", buf.String())
}

func TestTokenNewNoDataError(t *testing.T) {
	_ = doTokenNewRequest(t, nil, http.StatusUnprocessableEntity)
}

func TestTokenNewInvalidJSONError(t *testing.T) {
	_ = doTokenNewRequest(t, bytes.NewReader([]byte(`this is an invalid JSON`)), http.StatusInternalServerError)
}

func doTokenNewRequest(t *testing.T, body io.Reader, expectedStatusCode int) *http.Response {
	req, err := http.NewRequest(http.MethodPost, `/tokens`, body)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	h := TokenNewHandler{
		JWTService: getJWTService(),
	}

	h.ServeHTTP(rec, req)

	res := rec.Result()

	if res.StatusCode != expectedStatusCode {
		t.Fatalf("Unexpected status code %d vs %d (%s)", res.StatusCode, expectedStatusCode, res.Status)
	}

	return res
}

func TestTokenCheckSuccess(t *testing.T) {
	data := map[string]interface{}{
		"lastname":  "LASTNAME",
		"firstname": "FIRSTNAME",
		"age":       42,
	}

	b, _ := json.Marshal(data)

	res := doTokenNewRequest(t, bytes.NewReader(b), http.StatusOK)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	token := ``
	if err := json.Unmarshal([]byte(buf.String()), &token); err != nil {
		t.Fatal(err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	res = doTokenCheckRequest(t, headers, http.StatusOK)
	buf.Reset()
	buf.ReadFrom(res.Body)

	d := map[string]bool{}
	if err := json.Unmarshal([]byte(buf.String()), &d); err != nil {
		t.Fatal(err)
	}

	if valid, ok := d["valid"]; !ok {
		t.Fatal("Response does not contain valid key")
	} else if !valid {
		t.Fatal("Generated token is not valid")
	}
}

func TestTokenCheckEmptyToken(t *testing.T) {
	_ = doTokenCheckRequest(t, nil, http.StatusUnauthorized)
}

func TestTokenCheckInvalidToken(t *testing.T) {
	h := map[string]string{
		"Authorization": "Bearer 1V@lidT0ken",
	}
	_ = doTokenCheckRequest(t, h, http.StatusUnauthorized)

}

func TestTokenCheckExpiredToken(t *testing.T) {
	h := map[string]string{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZ2UiOjQyLCJhdWQiOiJ0ZXN0X2F1ZGllbmNlIiwiZXhwIjoxNTQ3MTEzNTIxLCJmaXJzdG5hbWUiOiJGSVJTVE5BTUUiLCJpc3MiOiJ0ZXN0X2lzc3VlciIsImxhc3RuYW1lIjoiTEFTVE5BTUUifQ.-VrA-I9OSEPgjiwcvn_MP0wyzjvdaZCL5B__KXB_YXk",
	}

	res := doTokenCheckRequest(t, h, http.StatusOK)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	d := map[string]bool{}
	if err := json.Unmarshal([]byte(buf.String()), &d); err != nil {
		t.Fatal(err)
	}

	if valid, ok := d["valid"]; !ok {
		t.Fatal("Response does not contain valid key")
	} else if valid {
		t.Fatal("Token should not be valid")
	}
}

func doTokenCheckRequest(t *testing.T, headers map[string]string, expectedStatusCode int) *http.Response {

	req, err := http.NewRequest(http.MethodGet, `/tokens/check`, nil)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range headers {
		t.Logf("Add header : [%s : %s]", k, v)
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()

	h := TokenCheckHandler{
		JWTService: getJWTService(),
	}

	h.ServeHTTP(rec, req)

	res := rec.Result()

	if res.StatusCode != expectedStatusCode {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		t.Fatalf("Unexpected status code %d(%s) vs %d - Body is '%s'", res.StatusCode, res.Status, expectedStatusCode, buf.String())
	}

	return res
}

func TestTokenDecodeSuccess(t *testing.T) {
	data := map[string]interface{}{
		"lastname":  "LASTNAME",
		"firstname": "FIRSTNAME",
		"age":       42,
	}

	b, _ := json.Marshal(data)

	res := doTokenNewRequest(t, bytes.NewReader(b), http.StatusOK)

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	token := ``
	if err := json.Unmarshal([]byte(buf.String()), &token); err != nil {
		t.Fatal(err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	res = doTokenDecodeRequest(t, headers, http.StatusOK)
	buf.Reset()
	buf.ReadFrom(res.Body)

	resData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(buf.String()), &resData); err != nil {
		t.Fatal(err)
	}

	for k, v := range data {
		if value, ok := resData[k]; !ok {
			t.Fatalf("Missing field %s in token data", k)
		} else if fmt.Sprint(value) != fmt.Sprint(v) {
			t.Fatalf("Invalid decoded data encodedData.%s (%v) != decodedData.%s (%v)", k, v, k, value)
		} else {
			t.Logf("encodedData.%s (%v) = decodedData.%s (%v)", k, v, k, value)
		}
	}
}

func TestTokenDecodeEmptyToken(t *testing.T) {
	_ = doTokenDecodeRequest(t, nil, http.StatusUnauthorized)
}

func TestTokenDecodeInvalidToken(t *testing.T) {
	h := map[string]string{
		"Authorization": "Bearer 1V@lidT0ken",
	}
	_ = doTokenDecodeRequest(t, h, http.StatusUnauthorized)

}

func TestTokenDecodeExpiredToken(t *testing.T) {
	h := map[string]string{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZ2UiOjQyLCJhdWQiOiJ0ZXN0X2F1ZGllbmNlIiwiZXhwIjoxNTQ3MTEzNTIxLCJmaXJzdG5hbWUiOiJGSVJTVE5BTUUiLCJpc3MiOiJ0ZXN0X2lzc3VlciIsImxhc3RuYW1lIjoiTEFTVE5BTUUifQ.-VrA-I9OSEPgjiwcvn_MP0wyzjvdaZCL5B__KXB_YXk",
	}

	_ = doTokenDecodeRequest(t, h, http.StatusUnauthorized)
}

func doTokenDecodeRequest(t *testing.T, headers map[string]string, expectedStatusCode int) *http.Response {

	req, err := http.NewRequest(http.MethodGet, `/tokens/decode`, nil)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()

	h := TokenDecodeHandler{
		JWTService: getJWTService(),
	}

	h.ServeHTTP(rec, req)

	res := rec.Result()

	if res.StatusCode != expectedStatusCode {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		t.Fatalf("Unexpected status code %d(%s) vs %d - Body is '%s'", res.StatusCode, res.Status, expectedStatusCode, buf.String())
	}

	return res
}
