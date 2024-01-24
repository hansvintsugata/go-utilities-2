package zenServer

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	zenMiddleware "github.com/hansvintsugata/go-utilities/appRunner/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type (
	HttpServer struct {
		ec                                    *echo.Echo
		serviceName, serviceVersion, basePath string
		option                                httpOptions
		serverTemplate
	}
)

func NewHttp(serviceName, serviceVersion, basePath string, opts ...HttpOption) *HttpServer {
	o := defaultHttpOption()
	for _, opt := range opts {
		opt.Apply(&o)
	}
	return &HttpServer{
		ec:             echo.New(),
		basePath:       basePath,
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
		option:         o,
	}
}

func (s *HttpServer) Serve(sig chan os.Signal) {
	if s.option.validator != nil {
		s.ec.Validator = s.option.validator
	}

	s.ec.HideBanner = true

	s.ec.Use(
		middleware.Recover(),
		middleware.Gzip(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding},
		}))

	s.ec.GET(path.Join(s.basePath, "/healthz"), func(ctx echo.Context) error {
		httpStatus, dependencies := s.option.healthCheck(ctx.Request().Context())
		return ctx.JSON(httpStatus, map[string]interface{}{
			"service_name":    s.serviceName,
			"service_version": s.serviceVersion,
			"status":          "UP",
			"dependencies":    dependencies,
		})
	}, zenMiddleware.AccessKey(s.option.healthCheckAccessKey))

	if s.option.errorHandler != nil {
		s.ec.HTTPErrorHandler = func(err error, c echo.Context) {
			s.option.errorHandler(err, c)
		}
	}

	s.serve(sig, serveParam{
		serve: func(sig chan os.Signal) {
			log.Info("[HTTP-SERVER] starting server")
			go func() {
				if err := s.ec.Start(fmt.Sprintf(":%d", s.option.port)); err != nil {
					log.Errorf("[HTTP-SERVER] server interrupted %s", err.Error())
					sig <- os.Interrupt
				}
			}()
			time.Sleep(time.Second)
		},
		register: func() {
			if !s.option.enable {
				return
			}

			if s.option.register != nil {
				log.Debug("[HTTP-SERVER] starting register hooks")
				s.option.register(s.ec)
			}
		},
		beforeRun: func() {
			if s.option.beforeRun != nil {
				log.Debug("[HTTP-SERVER] starting before run hooks")
				s.option.beforeRun(s.ec)
			}
		},
		afterRun: func() {
			if s.option.afterRun != nil {
				log.Debug("[HTTP-SERVER] starting after run hooks")
				s.option.afterRun(s.ec)
			}
		},
	})
}

func (s *HttpServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.option.gracefulPeriod)
	defer cancel()

	s.shutdown(shutdownParam{
		shutdown: func() {
			log.Info("[HTTP-SERVER] shutting down server")
			if err := s.ec.Shutdown(ctx); err != nil {
				log.Errorf("[HTTP-SERVER] server can not be shutdown %s", err.Error())
			}
		},
		beforeExit: func() {
			if s.option.beforeExit != nil {
				log.Debug("[HTTP-SERVER] starting before exit hooks")
				s.option.beforeExit(s.ec)
			}
		},
		afterExit: func() {
			if s.option.afterExit != nil {
				log.Debug("[HTTP-SERVER] starting after exit hooks")
				s.option.afterExit(s.ec)
			}
		},
	})
}
