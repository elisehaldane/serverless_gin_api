package lambda

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusRoute(t *testing.T) {
	router := createRouter()

	response := httptest.NewRecorder()

	request, _ := http.NewRequest("GET", "/v1/status", nil)

	router.ServeHTTP(response, request)

	// Assert that the API is returning its status
	assert.Equal(t, "{\"message\":\"okay\"}", response.Body.String())

	log.Println("All tests concluded")

	signalTermination(router)
}
