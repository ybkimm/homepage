package res

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	compilers2 "go.ybk.im/homepage/internal/app/skins/res/compilers"
)

type Handler struct {
	compiled    []byte
	contentType string
}

func NewHandler(fileName string, compiler compilers2.Compiler) *Handler {
	fileName = filepath.Clean(fileName)

	_, err := os.Stat(fileName)
	if err != nil {
		panic("File is not exists")
	}

	// Load file
	fp, err := os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		panic(fmt.Sprintf("Failed to open file: %s", err))
	}
	defer func() { _ = fp.Close() }()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file: %s", err))
	}

	compiled, err := compiler.Compile(data)
	if err != nil {
		panic(fmt.Sprintf("Compile failed: %s", err))
	}

	return &Handler{
		compiled:    compiled,
		contentType: compiler.ContentType(),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", h.contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.compiled)
}
