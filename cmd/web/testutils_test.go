package main

//
//import (
//	"bytes"
//	"io"
//	"log/slog"
//	"net/http"
//	"net/http/cookiejar"
//	"net/http/httptest"
//	"testing"
//)
//
//func newTestApplication(t *testing.T) *application {
//	return &application{
//		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
//	}
//}
//
//type testSever struct {
//	*httptest.Server
//}
//
//func newTestServer(t *testing.T, h http.Handler) *testSever {
//	ts := httptest.NewTLSServer(h)
//
//	// Initialize a new cookie jar.
//	jar, err := cookiejar.New(nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	ts.Client().Jar = jar
//
//	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
//		return http.ErrUseLastResponse
//	}
//	return &testSever{ts}
//}
//
//func (ts *testSever) get(t *testing.T, urlPath string) (int, http.Header, string) {
//	rs, err := ts.Client().Get(ts.URL + urlPath)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	defer rs.Body.Close()
//	body, err := io.ReadAll(rs.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//	body = bytes.TrimSpace(body)
//
//	return rs.StatusCode, rs.Header, string(body)
//}
