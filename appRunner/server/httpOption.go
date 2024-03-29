package zenServer

import (
	"context"
	"net/http"
	"time"

	zenValidator "github.com/hansvintsugata/go-utilities/validator"

	"github.com/labstack/echo/v4"
)

type (
	HealthCheckHook  func(ctx context.Context) (int, map[string]map[string]interface{})
	HttpHook         func(ec *echo.Echo)
	HttpErrorHandler func(err error, eCtx echo.Context)

	httpOptions struct {
		// Server
		enable         bool
		port           int
		gracefulPeriod time.Duration
		validator      zenValidator.Validator
		// logger               sunerLogger.Logger
		healthCheck          HealthCheckHook
		healthCheckAccessKey string

		// Hooks
		register     HttpHook
		beforeRun    HttpHook
		afterRun     HttpHook
		beforeExit   HttpHook
		afterExit    HttpHook
		errorHandler HttpErrorHandler
	}
	HttpOption interface {
		Apply(o *httpOptions)
	}
)

func defaultHttpOption() httpOptions {
	validator, _ := zenValidator.New()
	return httpOptions{
		enable:         true,
		port:           8090,
		gracefulPeriod: 12 * time.Second,
		validator:      validator,
		// logger:         sunerLogrus.DefaultLog(),
		register:   nil,
		beforeRun:  nil,
		afterRun:   nil,
		beforeExit: nil,
		afterExit:  nil,
		healthCheck: func(ctx context.Context) (int, map[string]map[string]interface{}) {
			return http.StatusOK, make(map[string]map[string]interface{})
		},
		healthCheckAccessKey: "",
	}
}

type withHttpEnable bool

func (w withHttpEnable) Apply(o *httpOptions) {
	o.enable = bool(w)
}

func WithHttpEnable(enable bool) HttpOption {
	return withHttpEnable(enable)
}

type withHttpPort int

func (w withHttpPort) Apply(o *httpOptions) {
	o.port = int(w)
}

func WithHttpPort(port int) HttpOption {
	return withHttpPort(port)
}

type withHttpHealthAccessKey string

func (w withHttpHealthAccessKey) Apply(o *httpOptions) {
	o.healthCheckAccessKey = string(w)
}

func WithHttpHealthAccessKey(accessKey string) HttpOption {
	return withHttpHealthAccessKey(accessKey)
}

type withHttpGraceFulPeriod time.Duration

func (w withHttpGraceFulPeriod) Apply(o *httpOptions) {
	o.gracefulPeriod = time.Duration(w)
}

func WithHttpGraceFulPeriod(duration time.Duration) HttpOption {
	return withHttpGraceFulPeriod(duration)
}

type withHttpValidator struct{ zenValidator.Validator }

func (w withHttpValidator) Apply(o *httpOptions) {
	o.validator = w.Validator
}

func WithHttpValidator(validator zenValidator.Validator) HttpOption {
	return withHttpValidator{validator}
}

// type withHttpLogger struct{ sunerLogger.Logger }

// func (w withHttpLogger) Apply(o *httpOptions) {
// 	o.logger = w.Logger
// }

// func WithHttpLogger(logger sunerLogger.Logger) HttpOption {
// 	return withHttpLogger{logger}
// }

type withHttpHealthCheck HealthCheckHook

func (w withHttpHealthCheck) Apply(o *httpOptions) {
	o.healthCheck = HealthCheckHook(w)
}

func WithHttpHealthCheck(hook HealthCheckHook) HttpOption {
	return withHttpHealthCheck(hook)
}

type withHttpRegister HttpHook

func (w withHttpRegister) Apply(o *httpOptions) {
	o.register = HttpHook(w)
}

func WithHttpRegister(hook HttpHook) HttpOption {
	return withHttpRegister(hook)
}

type withHttpBeforeRun HttpHook

func (w withHttpBeforeRun) Apply(o *httpOptions) {
	o.beforeRun = HttpHook(w)
}

func WithHttpBeforeRun(hook HttpHook) HttpOption {
	return withHttpBeforeRun(hook)
}

type withHttpAfterRun HttpHook

func (w withHttpAfterRun) Apply(o *httpOptions) {
	o.afterRun = HttpHook(w)
}

func WithHttpAfterRun(hook HttpHook) HttpOption {
	return withHttpAfterRun(hook)
}

type withHttpBeforeExit HttpHook

func (w withHttpBeforeExit) Apply(o *httpOptions) {
	o.beforeExit = HttpHook(w)
}

func WithHttpBeforeExit(hook HttpHook) HttpOption {
	return withHttpBeforeExit(hook)
}

type withHttpAfterExit HttpHook

func (w withHttpAfterExit) Apply(o *httpOptions) {
	o.afterExit = HttpHook(w)
}

func WithHttpAfterExit(hook HttpHook) HttpOption {
	return withHttpAfterExit(hook)
}

type withHttpErrorHandler HttpErrorHandler

func (w withHttpErrorHandler) Apply(o *httpOptions) {
	o.errorHandler = HttpErrorHandler(w)
}

func WithHttpErrorHandler(hook HttpErrorHandler) HttpOption {
	return withHttpErrorHandler(hook)
}
