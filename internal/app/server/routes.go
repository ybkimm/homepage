package server

import (
	"net/http"

	"go.ybk.im/homepage/internal/app/server/handlers"
	"go.ybk.im/homepage/internal/app/skins/res"
	"go.ybk.im/homepage/internal/app/skins/res/compilers/raw"
)

const CssPath = "skin/style.css"

func (s *Server) ApplyRoutes() {
	router := s.router

	router.
		Methods("GET").
		Path("/res/style.css").
		Handler(res.NewHandler(CssPath, raw.NewCompiler("text/css;charset=UTF-8")))

	router.
		Methods(http.MethodGet).
		Path("/").
		Handler(&handlers.IndexPage{})
}
