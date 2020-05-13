// Package main ...
package main

import (
	"flag"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/IgorRybak2055/logistic-service/pkg/client"
)

var (
	httpAddr     = flag.String("http-addr", "http://127.0.0.1:8888", "http ragger address")
	httpTimeout  = flag.Int("http-timeout", 10, "http timeout for requests")
	userEmail    = flag.String("credential-email", "test@domain.com", "email of ragger user")
	userPassword = flag.String("credential-password", "password", "password of ragger user")
)

func main() {
	flag.Parse()

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	cl := client.NewClient(*httpAddr, time.Duration(*httpTimeout)*time.Second)

	cl.Run()
}
