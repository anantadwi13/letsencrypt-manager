package internal

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
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

func (s *service) GetCertificates(ctx echo.Context) error {
	panic("implement me")
}

func (s *service) PostCertificates(ctx echo.Context) error {
	panic("implement me")
}

func (s *service) PutCertificates(ctx echo.Context) error {
	panic("implement me")
}

func (s *service) DeleteCertificatesDomain(ctx echo.Context, domain string) error {
	panic("implement me")
}

func (s *service) GetCertificatesDomain(ctx echo.Context, domain string) error {
	panic("implement me")
}

func (s *service) PutCertificatesDomain(ctx echo.Context, domain string) error {
	panic("implement me")
}
