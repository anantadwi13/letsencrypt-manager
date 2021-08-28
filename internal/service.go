package internal

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Service interface {
	Start()
}

type service struct {
	shutdownWg sync.WaitGroup
}

func NewService() Service {
	return &service{}
}

func (s *service) Start() {
	signalOS := make(chan os.Signal, 1)
	signal.Notify(signalOS, syscall.SIGINT, syscall.SIGTERM)

	static := echo.New()
	api := echo.New()

	static.Use(middleware.Logger())
	static.Use(middleware.Recover())
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())

	static.Static("/", "./public")
	RegisterHandlers(api, s)

	go func() {
		err := static.Start(":80")
		if err != nil {
			log.Fatalln(err)
		}
	}()
	go func() {
		err := api.Start(":5555")
		if err != nil {
			log.Fatalln(err)
		}
	}()

	select {
	case <-signalOS:
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		go func() {
			s.shutdownWg.Add(1)
			defer s.shutdownWg.Done()
			err := static.Shutdown(ctx)
			if err != nil {
				log.Println(err)
			}
		}()
		go func() {
			s.shutdownWg.Add(1)
			defer s.shutdownWg.Done()
			err := api.Shutdown(ctx)
			if err != nil {
				log.Println(err)
			}
		}()

		s.shutdownWg.Wait()
	}
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
