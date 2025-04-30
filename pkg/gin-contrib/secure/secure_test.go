package secure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	// Default
	Secure()(c)
	assert.Equal(t, "1; mode=block", rec.Header().Get(HeaderXXSSProtection))
	assert.Equal(t, "nosniff", rec.Header().Get(HeaderXContentTypeOptions))
	assert.Equal(t, "SAMEORIGIN", rec.Header().Get(HeaderXFrameOptions))
	assert.Equal(t, "", rec.Header().Get(HeaderStrictTransportSecurity))
	assert.Equal(t, "", rec.Header().Get(HeaderContentSecurityPolicy))
	assert.Equal(t, "", rec.Header().Get(HeaderReferrerPolicy))

	// Custom
	req.Header.Set(HeaderXForwardedProto, "https")
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = req
	SecureWithConfig(SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "origin",
	})(c)
	assert.Equal(t, "", rec.Header().Get(HeaderXXSSProtection))
	assert.Equal(t, "", rec.Header().Get(HeaderXContentTypeOptions))
	assert.Equal(t, "", rec.Header().Get(HeaderXFrameOptions))
	assert.Equal(t, "max-age=3600; includeSubdomains", rec.Header().Get(HeaderStrictTransportSecurity))
	assert.Equal(t, "default-src 'self'", rec.Header().Get(HeaderContentSecurityPolicy))
	assert.Equal(t, "", rec.Header().Get(HeaderContentSecurityPolicyReportOnly))
	assert.Equal(t, "origin", rec.Header().Get(HeaderReferrerPolicy))

	// Custom with CSPReportOnly flag
	req.Header.Set(HeaderXForwardedProto, "https")
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = req
	SecureWithConfig(SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
		CSPReportOnly:         true,
		ReferrerPolicy:        "origin",
	})(c)
	assert.Equal(t, "", rec.Header().Get(HeaderXXSSProtection))
	assert.Equal(t, "", rec.Header().Get(HeaderXContentTypeOptions))
	assert.Equal(t, "", rec.Header().Get(HeaderXFrameOptions))
	assert.Equal(t, "max-age=3600; includeSubdomains", rec.Header().Get(HeaderStrictTransportSecurity))
	assert.Equal(t, "default-src 'self'", rec.Header().Get(HeaderContentSecurityPolicyReportOnly))
	assert.Equal(t, "", rec.Header().Get(HeaderContentSecurityPolicy))
	assert.Equal(t, "origin", rec.Header().Get(HeaderReferrerPolicy))

	// Custom, with preload option enabled
	req.Header.Set(HeaderXForwardedProto, "https")
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = req
	SecureWithConfig(SecureConfig{
		HSTSMaxAge:         3600,
		HSTSPreloadEnabled: true,
	})(c)
	assert.Equal(t, "max-age=3600; includeSubdomains; preload", rec.Header().Get(HeaderStrictTransportSecurity))

	// Custom, with preload option enabled and subdomains excluded
	req.Header.Set(HeaderXForwardedProto, "https")
	rec = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(rec)
	c.Request = req
	SecureWithConfig(SecureConfig{
		HSTSMaxAge:            3600,
		HSTSPreloadEnabled:    true,
		HSTSExcludeSubdomains: true,
	})(c)
	assert.Equal(t, "max-age=3600; preload", rec.Header().Get(HeaderStrictTransportSecurity))
}
