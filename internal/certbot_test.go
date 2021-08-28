package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	certman CertificateManager
	ctx     context.Context
)

func init() {
	certman = NewCertbot()
	ctx = context.TODO()
}

func TestCertbot_GetAll(t *testing.T) {
	all, err := certman.GetAll(ctx)
	assert.Empty(t, err)
	assert.NotEmpty(t, all)
}
