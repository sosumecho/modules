package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sosumecho/modules/i18n"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/middlewares"
	"github.com/sosumecho/modules/utils"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"os"
)

const (
	Debug   = "debug"
	Release = "release"
	Test    = "test"
)

// HTTPServer http 服务器
type HTTPServer struct {
	name         string
	service      func(engine *gin.Engine)
	dependencies []func()
	logger       *logger.Logger
	locale       *i18n.I18N
	conf         *Conf
	middlewares  []gin.HandlerFunc
}

type Conf struct {
	Addr string `mapstructure:"addr"`
	PID  string `mapstructure:"pid-path"`
	Env  string `mapstructure:"env"`
}

func (h *HTTPServer) GetPID() string {
	return fmt.Sprintf("%s/%s.pid", h.conf.PID, h.name)
}

func (h *HTTPServer) SavePID() error {
	pid := cast.ToString(os.Getpid())
	_, err := utils.Write(h.GetPID(), pid)
	h.logger.Debug("pid", zap.String("pid", pid))
	return err
}

func (h *HTTPServer) AddDependencies(dependencies ...func()) *HTTPServer {
	h.dependencies = append(h.dependencies, dependencies...)
	return h
}

func (h *HTTPServer) AddMiddlewares(middlewares ...gin.HandlerFunc) *HTTPServer {
	h.middlewares = append(h.middlewares, middlewares...)
	return h
}

func (h *HTTPServer) newHandler() *gin.Engine {
	gin.SetMode(h.conf.Env)
	router := gin.New()
	router.Use(
		h.middlewares...,
	)
	return router
}

// NewHTTPService 新建http服务
func NewHTTPService(name string, conf *Conf, service func(handler *gin.Engine), logf *logger.Logger, locale *i18n.I18N) *HTTPServer {
	s := &HTTPServer{
		name:    name,
		conf:    conf,
		service: service,
		logger:  logf,
		locale:  locale,
		middlewares: []gin.HandlerFunc{
			middlewares.Recovery(logf, locale),
			middlewares.ForceUA(locale, logf),
		},
	}
	if s.conf.PID == "" {
		s.conf.PID = utils.RuntimeDir()
	}
	return s
}
