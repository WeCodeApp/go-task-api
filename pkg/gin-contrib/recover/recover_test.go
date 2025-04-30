package recover

import (
	"bufio"
	"bytes"
	//"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRecover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := new(bytes.Buffer)
	logger := bufio.NewWriter(buf)

	router := gin.New()

	router.Use(RecoverWithConfig(RecoverConfig{
		LogErrorFunc: func(c *gin.Context, err error, stack []byte) {
			logger.WriteString(fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack))
			logger.Flush()
		},
	}))
	router.GET("/", func(c *gin.Context) {
		panic("test")
	})

	w := PerformRequest(router, "GET", "/")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buf.String(), "PANIC RECOVER")
}

func TestRecoverErrAbortHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := new(bytes.Buffer)
	logger := bufio.NewWriter(buf)

	router := gin.New()

	router.Use(RecoverWithConfig(RecoverConfig{
		LogErrorFunc: func(c *gin.Context, err error, stack []byte) {
			logger.WriteString(fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack))
			logger.Flush()
		},
	}))

	router.GET("/", func(c *gin.Context) {
		panic(http.ErrAbortHandler)
	})

	defer func() {
		r := recover()
		if r == nil {
			assert.Fail(t, "expecting `http.ErrAbortHandler`, got `nil`")
		} else {
			if err, ok := r.(error); ok {
				assert.ErrorIs(t, err, http.ErrAbortHandler)
			} else {
				assert.Fail(t, "not of error type")
			}
		}
	}()
	w := PerformRequest(router, "GET", "/")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NotContains(t, buf.String(), "PANIC RECOVER")
}
