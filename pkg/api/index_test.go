package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv := NewMockServer()
	handler := http.HandlerFunc(srv.indexHandler)

	handler.ServeHTTP(rr, req)

	//expecting "MockHello"
	if response := rr.Body.String(); response != "MockHello" {
		t.Errorf("handler returned wrong response: got %v want %v",
			response, "MockHello")
	}
}
