package server

import (
	"fmt"
	"github.com/fvbock/endless"
	"net/http"
	"time"
)

// Start 启动
func (h *HTTPServer) Start() {
	for _, dependency := range h.dependencies {
		dependency()
	}
	handler := h.newHandler()
	h.service(handler)
	endlessServer := endless.NewServer(h.conf.Addr, handler)
	endlessServer.BeforeBegin = func(add string) {
		if err := h.SavePID(); err != nil {
			panic(err)
		}
	}
	err := endlessServer.ListenAndServe()
	if err == http.ErrServerClosed {
		h.logger.Debug(fmt.Sprintf("%s server closed at %d", h.name, time.Now().Unix()))
	} else if err != nil {
		h.logger.Debug(fmt.Sprintf("start %s server failed, err is %s", h.name, err.Error()))
	} else {
		h.logger.Debug("server closed")
	}
}
