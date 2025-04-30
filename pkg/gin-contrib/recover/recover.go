package recover

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type (
	Skipper func(c *gin.Context) bool
	// LogErrorFunc defines a function for custom logging in the middleware.
	LogErrorFunc func(c *gin.Context, err error, stack []byte)

	// RecoverConfig defines the config for Recover middleware.
	RecoverConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Size of the stack to be printed.
		// Optional. Default value 4KB.
		StackSize int `yaml:"stack_size"`

		// DisableStackAll disables formatting stack traces of all other goroutines
		// into buffer after the trace for the current goroutine.
		// Optional. Default value false.
		DisableStackAll bool `yaml:"disable_stack_all"`

		// DisablePrintStack disables printing stack trace.
		// Optional. Default value as false.
		DisablePrintStack bool `yaml:"disable_print_stack"`

		// LogErrorFunc defines a function for custom logging in the middleware.
		// If it's set you don't need to provide LogLevel for config.
		LogErrorFunc LogErrorFunc
	}
)

var (
	// DefaultRecoverConfig is the default Recover middleware config.
	DefaultRecoverConfig = RecoverConfig{
		Skipper:           DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogErrorFunc:      nil,
	}
)

func DefaultSkipper(c *gin.Context) bool {
	return false
}

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() gin.HandlerFunc {
	return RecoverWithConfig(DefaultRecoverConfig)
}

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func RecoverWithConfig(config RecoverConfig) gin.HandlerFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig.StackSize
	}

	return func(c *gin.Context) {
		if config.Skipper(c) {
			c.Next()
		}

		defer func() {
			if r := recover(); r != nil {
				if r == http.ErrAbortHandler {
					panic(r)
				}
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				var stack []byte
				var length int

				if !config.DisablePrintStack {
					stack = make([]byte, config.StackSize)
					length = runtime.Stack(stack, !config.DisableStackAll)
					stack = stack[:length]
				}

				if config.LogErrorFunc != nil {
					config.LogErrorFunc(c, err, stack)
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}

}
