package types

import "net/http"

type Server struct {
	HttpServer *http.Server
	Logger     interface{}
}
