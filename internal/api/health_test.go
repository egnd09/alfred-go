package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.GET("/", HealthCheck)
	router.GET("/health", HealthCheck)
	
	tests := []struct {
		name     string
		endpoint string
		wantCode int
		wantBody string
	}{
		{"root endpoint", "/", http.StatusOK, `{"status":"ok"}`},
		{"health endpoint", "/health", http.StatusOK, `{"status":"ok"}`},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.endpoint, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.wantCode {
				t.Errorf("HealthCheck() status = %v, want %v", w.Code, tt.wantCode)
			}
			if w.Body.String() != tt.wantBody {
				t.Errorf("HealthCheck() body = %v, want %v", w.Body.String(), tt.wantBody)
			}
		})
	}
}
