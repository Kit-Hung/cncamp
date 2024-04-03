package handlers

import (
	"github.com/Kit-Hung/cncamp/module12/util"
	"net/http"
)

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	// 访问 localhost/healthz 时，返回 200
	funcName := "healthz"
	util.WriteToResponseAndHandleError(funcName, &w, r, "200", http.StatusOK)
}
