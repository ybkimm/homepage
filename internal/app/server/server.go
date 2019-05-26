package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	router *mux.Router
}

func New() *Server {
	return &Server{
		router: mux.NewRouter(),
	}
}

func (s *Server) Start(addr string) <-chan error {
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Minute,
		Addr:         addr,
		Handler:      s.router,
	}

	ch := make(chan error, 1)
	go func() {
		fmt.Printf("[INFO] HTTP 서버가 %s에서 시작되었습니다.\n", addr)
		ch <- srv.ListenAndServe()
	}()

	return ch
}

func (s *Server) StartACME(acmeDir string) (<-chan error, *tls.Config) {
	m := autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       autocert.DirCache(acmeDir),
		RenewBefore: 24 * 30 * time.Hour,
		HostPolicy: autocert.HostWhitelist(
			"localhost",
			"yongbin.kim",
		),
		Email: "iam@yongbin.kim",
	}

	ch := make(chan error, 1)
	go func() {
		fmt.Printf("[INFO] ACME 서버가 :80에서 시작되었습니다.\n")
		ch <- http.ListenAndServe(":80", m.HTTPHandler(nil))
	}()

	return ch, &tls.Config{GetCertificate: m.GetCertificate}
}

func (s *Server) StartTLS(addr string, cert string, key string, conf *tls.Config) <-chan error {
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Minute,
		Addr:         addr,
		Handler:      s.router,
	}

	if conf != nil {
		srv.TLSConfig = conf
	}

	ch := make(chan error, 1)
	go func() {
		fmt.Printf("[INFO] HTTPS 서버가 %s에서 시작되었습니다.\n", addr)
		ch <- srv.ListenAndServeTLS(cert, key)
	}()

	return ch
}
