package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(LoggingMiddleware("../middleware/test.log"))

	defer func() {
		if err := os.Remove("test.log"); err != nil {
			if !os.IsNotExist(err) {
				t.Errorf("error deleting test.log file: %v", err)
			}
		}
	}()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	req.Header.Set("Referer", "http://localhost:8080")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	logFile, err := os.Open("../middleware/test.log")
	assert.NoError(t, err)
	defer logFile.Close()

	logBytes, err := ioutil.ReadAll(logFile)
	assert.NoError(t, err)

	logMsg := fmt.Sprintf(`192.168.1.1 [%s] "GET / HTTP/1.1" 200 2 "http://localhost:8080" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36" 0s`, time.Now().Format("02/Jan/2006:15:04:05 -0700"))
	assert.Contains(t, string(logBytes), logMsg)
}
