package benchmark

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BenchmarkJSONEncode(b *testing.B) {
	url := map[string]string{
		"original_url": "https://example.com/test/path",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(url)
	}
}

func BenchmarkJSONDecode(b *testing.B) {
	jsonData := []byte(`{"url":"https://example.com"}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]string
		json.Unmarshal(jsonData, &result)
	}
}

func BenchmarkHTTPRequest(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	req := httptest.NewRequest("GET", "/", nil)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkConcurrentRequests(b *testing.B) {
	done := make(chan bool)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		go func() {
			time.Sleep(1 * time.Millisecond)
			done <- true
		}()
	}
	
	for i := 0; i < b.N; i++ {
		<-done
	}
}