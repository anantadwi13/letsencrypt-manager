package internal

import (
	"context"
	"time"
)

type CertificateManager interface {
	GetAll(ctx context.Context) ([]*Certificate, error)
	Get(ctx context.Context, domain string) (*Certificate, error)
	Add(ctx context.Context, domain, email string) (*Certificate, error)
	Delete(ctx context.Context, domain string) error
	RenewAll(ctx context.Context) error
	Renew(ctx context.Context, domain string) error
}

type Certificate struct {
	Name         string
	SerialNumber string
	Domains      []string
	KeyType      string
	ExpiryDate   time.Time
	Public       string
	Private      string
}
