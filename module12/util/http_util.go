package util

import (
	"context"
	"fmt"
	"github.com/Kit-Hung/cncamp/module12/log"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	versionKey = "VERSION"
	headerLine = "=========================== Details of the http header: ==========================="
)

func addReqHeaderToResp(w *http.ResponseWriter, r *http.Request) string {
	// 将 request 中带的 header 写入 response header
	sb := strings.Builder{}
	sb.WriteString(headerLine)
	sb.WriteString("\n")

	writer := *w
	if len(r.Header) > 0 {
		for key, value := range r.Header {
			respValue := strings.Join(value, ";")
			writer.Header().Set(key, respValue)
			sb.WriteString(fmt.Sprintf("%s = %s \n", key, respValue))
		}
	}
	return sb.String()
}

func readEnvAndSetToHeader(w *http.ResponseWriter, envKey string) {
	// 读取环境变量, 并设置到 response header
	envValue := os.Getenv(envKey)
	writer := *w
	writer.Header().Set(envKey, envValue)
}

func recordClientIpAndHttpCode(w *http.ResponseWriter, r *http.Request, httpCode int) error {
	// 记录客户端 ip 和 http 返回码
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return err
	}

	if net.ParseIP(host) != nil {
		// 客户端 ip 输出到标准输出
		fmt.Printf("client ip: %v\n", host)
	}

	// http 返回码输出到标准输出
	fmt.Printf("http response code: %v\n", httpCode)
	writer := *w
	writer.WriteHeader(httpCode)
	return nil
}

func requestHandler(w *http.ResponseWriter, r *http.Request, httpCode int) (string, error) {
	/*
		接收客户端 request，并将 request 中带的 header 写入 response header
		读取当前系统的环境变量中的 VERSION 配置，并写入 response header
		Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	*/
	headerDetails := addReqHeaderToResp(w, r)
	readEnvAndSetToHeader(w, versionKey)
	return headerDetails, recordClientIpAndHttpCode(w, r, httpCode)
}

func getFinalOutput(headerDetails, value string) string {
	return fmt.Sprintf("%s \n \n %s", value, headerDetails)
}

func writeToResponse(ctx context.Context, funcName string, w *http.ResponseWriter, value string) {
	if writeString, err := io.WriteString(*w, value); err != nil {
		log.Logger.For(ctx).Error("write to response error: ", zap.Any(funcName, err))
	} else {
		log.Logger.For(ctx).Info("write to response: ", zap.Any(funcName, writeString))
	}
}

func WriteToResponseAndHandleError(funcName string, w *http.ResponseWriter, r *http.Request, value string, httpCode int) {
	headerDetails, err := requestHandler(w, r, httpCode)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		writeToResponse(r.Context(), funcName, w, err.Error())
		return
	}

	finalOutput := getFinalOutput(headerDetails, value)
	writeToResponse(r.Context(), funcName, w, finalOutput)
}
