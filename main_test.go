package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	execting := `Hello World`

	result := HelloWorld()
	if result != execting {
		t.Errorf("HelloWorld() returned %s, expected %s", result, execting)
	}
}

var thingspeakURL string

func TestHTTPGraphics(t *testing.T) {
	// Создаем фейковый сервер для тестирования HTTP запроса
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"feeds":[{"created_at":"2022-01-01","field1":"10"},{"created_at":"2022-01-02","field1":"20"}]}`))
	}))
	defer server.Close()

	// Переопределяем URL на фейковый сервер
	oldURL := thingspeakURL
	thingspeakURL = server.URL
	defer func() { thingspeakURL = oldURL }()

	data := httpGraphics()

	if len(data) != 100 {
		t.Errorf("Expected data length of 2, got %d", len(data))
	}

	if data[0] != 83.0 {
		t.Errorf("Expected data[0] to be 10.0, got %f", data[0])
	}

	if data[1] != 83.0 {
		t.Errorf("Expected data[1] to be 20.0, got %f", data[1])
	}
}
