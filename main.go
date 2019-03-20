package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	srv := &http.Server{
		Addr:    ":7280",
		Handler: router,
	}

	errChan := make(chan error)
	endChan := make(chan *struct{})

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				fmt.Printf("[오류] 서버에서 오류가 발생했습니다.")
			}

		case <-endChan:
			fmt.Printf("[경고] 서버를 종료합니다.")
			os.Exit(0)
		}
	}
}
