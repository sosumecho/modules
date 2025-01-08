package server

import (
	"fmt"
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

	err := http.ListenAndServe(h.conf.Addr, handler)
	if err == http.ErrServerClosed {
		h.logger.Debug(fmt.Sprintf("%s server closed at %d", h.name, time.Now().Unix()))
	} else if err != nil {
		h.logger.Debug(fmt.Sprintf("start %s server failed, err is %s", h.name, err.Error()))
	} else {
		h.logger.Debug("servere closed")
	}
}
