package handlers

import (
	"fmt"
	"github.com/Kit-Hung/cncamp/module12/log"
	"github.com/Kit-Hung/cncamp/module12/util"
	"net/http"
)

func (s *Server) shutdown(w http.ResponseWriter, r *http.Request) {
	ClearResources()

	funcName := "shutdown"
	util.WriteToResponseAndHandleError(funcName, &w, r, "ok", http.StatusOK)
}

func ClearResources() {
	err := log.Logger.Sync()
	if err != nil {
		fmt.Printf("logger sync error: %v", err)
	}
}
