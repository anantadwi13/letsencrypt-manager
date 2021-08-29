package main

import (
	"context"
	"fmt"
	"github.com/anantadwi13/letsencrypt-manager/internal"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	certMan := internal.NewCertbot()
	ctx := context.TODO()
	wg := sync.WaitGroup{}
	err := os.Mkdir("./public", 0777)
	if err != nil {
		log.Panicln(err)
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		s := internal.NewService()
		s.Start()
	}()

	time.Sleep(1 * time.Second)

	cert, err := certMan.Add(ctx, "anantadwi@aaa.com", "coba.anantadwi13.com", "www.coba.anantadwi13.com")
	if err != nil {
		log.Println("error", err)
	}
	fmt.Println(cert)
	certs, err := certMan.GetAll(ctx)
	if err != nil {
		log.Println(err)
	}
	if len(certs) <= 0 {
		fmt.Println("Not found")
	}
	for _, certificate := range certs {
		fmt.Println(certificate)
	}

	wg.Wait()
}
