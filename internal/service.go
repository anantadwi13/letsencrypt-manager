package internal

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Service interface {
	Start()
}

type service struct {
	config     Config
	certMan    CertificateManager
	static     *echo.Echo
	api        *echo.Echo
	shutdownWg sync.WaitGroup
}

func NewService(config Config) Service {
	return &service{config: config}
}

func (s *service) Start() {
	signalOS := make(chan os.Signal, 1)
	signal.Notify(signalOS, syscall.SIGINT, syscall.SIGTERM)

	s.registerDependencies()

	select {
	case <-signalOS:
		s.Shutdown()
	}
}

func (s *service) registerDependencies() {
	s.static = echo.New()
	s.api = echo.New()
	s.certMan = NewCertbot(s.config)

	s.static.Use(middleware.Logger())
	s.static.Use(middleware.Recover())
	s.api.Use(middleware.Logger())
	s.api.Use(middleware.Recover())

	err := os.Mkdir(s.config.PublicStaticPath(), 0777)
	if err != nil {
		log.Panicln(err)
	}
	s.static.Static("/", s.config.PublicStaticPath())
	RegisterHandlers(s.api, s)
	s.api.GET("/specs", func(c echo.Context) error {
		return c.File("./specification.yaml")
	})
	s.api.GET("/docs", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
			<!DOCTYPE html>
			<html>
			  <head>
				<title>ReDoc</title>
				<!-- needed for adaptive design -->
				<meta charset="utf-8"/>
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
			
				<!--
				ReDoc doesn't change outer page styles
				-->
				<style>
				  body {
					margin: 0;
					padding: 0;
				  }
				</style>
			  </head>
			  <body>
				<redoc spec-url='/specs'></redoc>
				<script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
			  </body>
			</html>
		`)
	})

	go func() {
		err := s.static.Start(":80")
		if err != nil {
			log.Fatalln(err)
		}
	}()
	go func() {
		err := s.api.Start(":" + strconv.Itoa(s.config.ApiPort()))
		if err != nil {
			log.Fatalln(err)
		}
	}()
}

func (s *service) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go func() {
		s.shutdownWg.Add(1)
		defer s.shutdownWg.Done()
		err := s.static.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
	}()
	go func() {
		s.shutdownWg.Add(1)
		defer s.shutdownWg.Done()
		err := s.api.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
	}()
	s.shutdownWg.Wait()
}

func (s *service) GetCertificates(c echo.Context) error {
	ctx := c.Request().Context()

	certs, err := s.certMan.GetAll(ctx)
	if err != nil {
		return responseServerErr(c, err)
	}

	var res []*CertificateRes
	for _, cert := range certs {
		res = append(res, &CertificateRes{
			Domains:      cert.Domains,
			ExpiryDate:   cert.ExpiryDate,
			KeyType:      cert.KeyType,
			Name:         cert.Name,
			PrivateCert:  cert.Private,
			PublicCert:   cert.Public,
			SerialNumber: cert.SerialNumber,
		})
	}
	if res == nil {
		res = []*CertificateRes{}
	}
	return c.JSON(200, res)
}

func (s *service) PostCertificates(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(CertificateReq)
	err := c.Bind(req)
	if err != nil {
		return responseServerErr(c, err)
	}

	if req.Email == "" || req.Domain == "" {
		return responseBadRequest(c, errors.New("make sure email and domain are filled"))
	}

	if req.AltDomains == nil {
		req.AltDomains = &[]string{}
	}

	cert, err := s.certMan.Add(ctx, req.Email, req.Domain, *req.AltDomains...)
	if err != nil {
		return responseServerErr(c, err)
	}

	certRes := &CertificateRes{
		Domains:      cert.Domains,
		ExpiryDate:   cert.ExpiryDate,
		KeyType:      cert.KeyType,
		Name:         cert.Name,
		PrivateCert:  cert.Private,
		PublicCert:   cert.Public,
		SerialNumber: cert.SerialNumber,
	}
	return c.JSON(http.StatusCreated, certRes)
}

func (s *service) PutCertificates(c echo.Context) error {
	ctx := c.Request().Context()
	err := s.certMan.RenewAll(ctx)
	if err != nil {
		return responseServerErr(c, err)
	}
	return c.JSON(http.StatusOK, GeneralRes{
		Code:    200,
		Message: "ok",
	})
}

func (s *service) DeleteCertificatesDomain(c echo.Context, domain string) error {
	ctx := c.Request().Context()
	cert, err := s.certMan.Get(ctx, domain)
	if err != nil {
		return responseServerErr(c, err)
	}
	if cert == nil {
		return responseNotFound(c, "Certificate is not found")
	}

	err = s.certMan.Delete(ctx, cert.Name)
	if err != nil {
		return responseServerErr(c, err)
	}
	return c.JSON(http.StatusOK, GeneralRes{
		Code:    200,
		Message: "ok",
	})
}

func (s *service) GetCertificatesDomain(c echo.Context, domain string) error {
	ctx := c.Request().Context()

	cert, err := s.certMan.Get(ctx, domain)
	if err != nil {
		return responseServerErr(c, err)
	}
	if cert == nil {
		return responseNotFound(c, "Certificate is not found")
	}
	certRes := &CertificateRes{
		Domains:      cert.Domains,
		ExpiryDate:   cert.ExpiryDate,
		KeyType:      cert.KeyType,
		Name:         cert.Name,
		PrivateCert:  cert.Private,
		PublicCert:   cert.Public,
		SerialNumber: cert.SerialNumber,
	}
	return c.JSON(http.StatusOK, certRes)
}

func (s *service) PutCertificatesDomain(c echo.Context, domain string) error {
	ctx := c.Request().Context()
	cert, err := s.certMan.Get(ctx, domain)
	if err != nil {
		return responseServerErr(c, err)
	}
	if cert == nil {
		return responseNotFound(c, "Certificate is not found")
	}

	err = s.certMan.Renew(ctx, domain)
	if err != nil {
		return responseServerErr(c, err)
	}
	return c.JSON(http.StatusOK, GeneralRes{
		Code:    200,
		Message: "ok",
	})
}

func responseNotFound(c echo.Context, message string) error {
	return c.JSON(http.StatusBadGateway, GeneralRes{
		Code:    404,
		Message: message,
	})
}

func responseServerErr(c echo.Context, err error) error {
	log.Println(err)
	return c.JSON(http.StatusBadGateway, GeneralRes{
		Code:    500,
		Message: err.Error(),
	})
}

func responseBadRequest(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, GeneralRes{
		Code:    400,
		Message: err.Error(),
	})
}
