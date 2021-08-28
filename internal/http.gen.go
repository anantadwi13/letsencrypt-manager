// Package internal provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package internal

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// CertificateRes defines model for certificate-res.
type CertificateRes struct {
	Domain      *string `json:"domain,omitempty"`
	PrivateCert *string `json:"private_cert,omitempty"`
	PublicCert  *string `json:"public_cert,omitempty"`
}

// ErrorRes defines model for error-res.
type ErrorRes struct {
	Code    *int    `json:"code"`
	Message *string `json:"message"`
}

// BadRequest defines model for bad-request.
type BadRequest ErrorRes

// DefaultError defines model for default-error.
type DefaultError ErrorRes

// NotFound defines model for not-found.
type NotFound ErrorRes

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all certificates
	// (GET /certificates)
	GetCertificates(ctx echo.Context) error
	// Create certificate
	// (POST /certificates)
	PostCertificates(ctx echo.Context) error
	// Renew all certificates
	// (PUT /certificates)
	PutCertificates(ctx echo.Context) error
	// Delete certificate
	// (DELETE /certificates/{domain})
	DeleteCertificatesDomain(ctx echo.Context, domain string) error
	// Get certificate by domain name
	// (GET /certificates/{domain})
	GetCertificatesDomain(ctx echo.Context, domain string) error
	// Renew certificate for selected domain
	// (PUT /certificates/{domain})
	PutCertificatesDomain(ctx echo.Context, domain string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetCertificates converts echo context to params.
func (w *ServerInterfaceWrapper) GetCertificates(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCertificates(ctx)
	return err
}

// PostCertificates converts echo context to params.
func (w *ServerInterfaceWrapper) PostCertificates(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostCertificates(ctx)
	return err
}

// PutCertificates converts echo context to params.
func (w *ServerInterfaceWrapper) PutCertificates(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutCertificates(ctx)
	return err
}

// DeleteCertificatesDomain converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteCertificatesDomain(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "domain" -------------
	var domain string

	err = runtime.BindStyledParameterWithLocation("simple", false, "domain", runtime.ParamLocationPath, ctx.Param("domain"), &domain)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter domain: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteCertificatesDomain(ctx, domain)
	return err
}

// GetCertificatesDomain converts echo context to params.
func (w *ServerInterfaceWrapper) GetCertificatesDomain(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "domain" -------------
	var domain string

	err = runtime.BindStyledParameterWithLocation("simple", false, "domain", runtime.ParamLocationPath, ctx.Param("domain"), &domain)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter domain: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCertificatesDomain(ctx, domain)
	return err
}

// PutCertificatesDomain converts echo context to params.
func (w *ServerInterfaceWrapper) PutCertificatesDomain(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "domain" -------------
	var domain string

	err = runtime.BindStyledParameterWithLocation("simple", false, "domain", runtime.ParamLocationPath, ctx.Param("domain"), &domain)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter domain: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutCertificatesDomain(ctx, domain)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/certificates", wrapper.GetCertificates)
	router.POST(baseURL+"/certificates", wrapper.PostCertificates)
	router.PUT(baseURL+"/certificates", wrapper.PutCertificates)
	router.DELETE(baseURL+"/certificates/:domain", wrapper.DeleteCertificatesDomain)
	router.GET(baseURL+"/certificates/:domain", wrapper.GetCertificatesDomain)
	router.PUT(baseURL+"/certificates/:domain", wrapper.PutCertificatesDomain)

}