// main_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewArticle(t *testing.T) {

	req, err := http.NewRequest("POST", "http://localhost:3000/articles/add", strings.NewReader("message"))
	if err != nil {
		t.Fatal(err)
	}

	//теперь попробуем получить ответ
	resp := httptest.NewRecorder()

	//resp.Header().Set("x-req-id", time.Now().Format("15:04:05"))

	ctr := new(Counter)
	handler := http.HandlerFunc(ctr.NewArticle)

	handler.ServeHTTP(resp, req)

	//проверяем статус
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("wrong status code %v", status)
	}

	//проверяем ответ
	reply := `{"text":"message1"}`
	if resp.Body.String() != reply {
		t.Errorf("body unexpected got %v want %v", resp.Body.String(), reply)
	}

	// проверить хэдер?
	// header := resp.Header().Get("x-req-id")
	// if header != time.Now().Format("15:04:05") {
	// 	t.Errorf("header doesnt match, got %v, want %v", header, time.Now().Format("15:04:05"))
	// }

}
