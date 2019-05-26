package handlers

import (
	"log"
	"net/http"

	"go.ybk.im/homepage/internal/app/skins"
)

type IndexPage struct {
	handlerBase
}

func (p *IndexPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	args := p.args()

	err := skins.Render("index.html", w, args)
	if err != nil {
		log.Printf("Failed to render index skin: %s\n", err)
	}
}
