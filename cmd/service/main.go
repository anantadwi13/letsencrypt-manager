package main

import (
	"context"
	"fmt"
	"github.com/anantadwi13/letsencrypt-manager/internal"
)

func main() {
	certMan := internal.NewCertbot()
	certs, err := certMan.GetAll(context.TODO())
	if err != nil {
		panic(err)
	}
	if len(certs) <= 0 {
		fmt.Println("Not found")
	}
	for _, certificate := range certs {
		fmt.Println(certificate)
	}
	//s := internal.NewService()
	//s.Start()
}
