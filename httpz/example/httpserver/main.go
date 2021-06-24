package main

import (
	"fmt"
	"github.com/seepre/go-lamp/httpz"
	"github.com/seepre/go-lamp/httpz/router"
	"net/http"
	"os"
)

func main() {
	c := httpz.HttpConf{
		Addr: ":443",
		CertFile: "./xxxx.pem",
		KeyFile: "./xxxx.key",
	}

	p, _ := os.UserHomeDir()
	fmt.Println("PATH:", p)

	srv := httpz.MustNewServe(c)

	srv.AddRoutes([]router.Route{
		{
			Method: http.MethodPost,
			Path: "/tls/test",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("TLS"))
			},
		},
	})

	fmt.Println("Listen on: 443")

	srv.Run()
}