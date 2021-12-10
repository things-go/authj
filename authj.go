package authj

import (
	"context"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer, so it fits in an interface{} without allocation.
type ctxAuthKey struct{}

// Config for Authorizer
type Config struct {
	errFallback        func(*gin.Context, error)
	forbiddenFallback  func(*gin.Context)
	skipAuthentication func(*gin.Context) bool
	subject            func(*gin.Context) string
}

// Option config option
type Option func(*Config)

// WithErrorFallback set the fallback handler when request are error happened.
// default: the 500 server error to the client
func WithErrorFallback(fn func(*gin.Context, error)) Option {
	return func(cfg *Config) {
		if fn != nil {
			cfg.errFallback = fn
		}
	}
}

// WithForbiddenFallback set the fallback handler when request are not allow.
// default: the 403 Forbidden to the client
func WithForbiddenFallback(fn func(*gin.Context)) Option {
	return func(cfg *Config) {
		if fn != nil {
			cfg.forbiddenFallback = fn
		}
	}
}

// WithSubject set the subject extractor of the requests.
// default: Subject
func WithSubject(fn func(*gin.Context) string) Option {
	return func(cfg *Config) {
		if fn != nil {
			cfg.subject = fn
		}
	}
}

// WithSkipAuthentication set the skip approve when it is return true.
// Default: always false
func WithSkipAuthentication(fn func(*gin.Context) bool) Option {
	return func(cfg *Config) {
		if fn != nil {
			cfg.skipAuthentication = fn
		}
	}
}

// Authorizer returns the authorizer
// uses a Casbin enforcer, and Subject as subject.
func Authorizer(e casbin.IEnforcer, opts ...Option) gin.HandlerFunc {
	cfg := Config{
		func(c *gin.Context, err error) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  "Permission validation errors occur!",
			})
		},
		func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "Permission denied!",
			})
		},
		func(c *gin.Context) bool { return false },
		Subject,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return func(c *gin.Context) {
		if !cfg.skipAuthentication(c) {
			// checks the subject,path,method permission combination from the request.
			allowed, err := e.Enforce(cfg.subject(c), c.Request.URL.Path, c.Request.Method)
			if err != nil {
				cfg.errFallback(c, err)
				return
			}
			if !allowed {
				cfg.forbiddenFallback(c)
				return
			}
		}
		c.Next()
	}
}

// Subject returns the value associated with this context for subjectCtxKey,
func Subject(c *gin.Context) string {
	val, _ := c.Request.Context().Value(ctxAuthKey{}).(string)
	return val
}

// ContextWithSubject return a copy of parent in which the value associated with
// subjectCtxKey is subject.
func ContextWithSubject(c *gin.Context, subject string) {
	ctx := context.WithValue(c.Request.Context(), ctxAuthKey{}, subject)
	c.Request = c.Request.WithContext(ctx)
}
