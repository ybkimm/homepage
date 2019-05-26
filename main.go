package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"os"

	"go.ybk.im/homepage/internal/app/skins"

	"go.ybk.im/homepage/internal/app/database"

	"go.ybk.im/homepage/internal/app/server"
)

var (
	dbFlag    = flag.String("db", "data.db", "데이터베이스 파일의 위치")
	skinFlag  = flag.String("skin", "skin/", "스킨 폴더의 위치")
	httpFlag  = flag.String("http", "", "HTTP 요청을 받을 주소")
	httpsFlag = flag.String("https", "", "HTTPS 요청을 받을 주소")
	certFlag  = flag.String("cert", "", "TLS 인증서 파일 위치")
	keyFlag   = flag.String("key", "", "TLS 인증서 키 파일 위치")
	acmeFlag  = flag.String("acme", "", "ACME 서버에서 받아온 인증서를 저장할 디렉토리")
)

func main() {
	fmt.Printf("Yongbin.Kim Server v1\n")

	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERRO] Fatal error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	if *httpFlag == "" && *httpsFlag == "" {
		return errors.New("-http, -https 둘 중 하나의 플래그가 필요합니다")
	}
	if *httpFlag != "" && *httpsFlag != "" {
		return errors.New("-http, -https 플래그를 동시에 쓸 수 없습니다")
	}
	if *httpsFlag == "" && *acmeFlag != "" {
		return errors.New("-https 없이 -acme 플래그를 쓸 수 없습니다")
	}
	if *httpsFlag != "" && *acmeFlag == "" && (*certFlag == "" || *keyFlag == "") {
		return errors.New("-https 플래그를 쓰고 -acme 플래그를 쓰지 않는 경우, -cert, -key 플래그를 같이 써야 합니다")
	}
	if *acmeFlag != "" && (*certFlag != "" || *keyFlag != "") {
		return errors.New("-acme 플래그를 -key 플래그나 -cert 플래그와 함께 쓸 수 없습니다")
	}

	srv := server.New()
	srv.ApplyRoutes()

	var err error

	// 데이터베이스 초기화
	err = database.Init(*dbFlag)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	// 스킨 초기화
	err = skins.Init(*skinFlag)
	if err != nil {
		panic(err)
	}

	// 서버 시작
	if *httpsFlag != "" {
		ch := make(chan error, 2)

		var tlsConfig *tls.Config

		// ACME 플래그가 세워졌으면 ACME 서버를 켜고,
		// 인증서를 받아올 수 있도록 설정함.
		if *acmeFlag != "" {
			// 인증서를 저장할 디렉토리
			if err := os.MkdirAll(*acmeFlag, 0700); err != nil {
				return err
			}

			// ACME 서버 시작 (포트 변경 불가, 80번 사용함)
			ach, conf := srv.StartACME(*acmeFlag)
			go func() {
				err := <-ach
				ch <- err
			}()
			tlsConfig = conf
		}

		// HTTPS 서버 시작
		sch := srv.StartTLS(*httpsFlag, *certFlag, *keyFlag, tlsConfig)
		go func() {
			err := <-sch
			ch <- err
		}()

		return <-ch
	} else {
		return <-srv.Start(*httpFlag)
	}
}
