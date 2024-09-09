//go:build go1.22

// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Send an Email with custom text
	// (POST /email/basic/send)
	SendEmail(w http.ResponseWriter, r *http.Request)
	// Get Template by Name
	// (GET /email/templates/{templateName})
	GetTemplateByName(w http.ResponseWriter, r *http.Request, templateName string)
	// Add new template
	// (POST /email/templates/{templateName})
	AddNewTemplate(w http.ResponseWriter, r *http.Request, templateName string)
	// Fill placeholders of template
	// (POST /email/templates/{templateName}/fill)
	FillTemplate(w http.ResponseWriter, r *http.Request, templateName string)
	// Get Template Placeholders
	// (GET /email/templates/{templateName}/placeholders)
	GetTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request, templateName string)
	// Send an templated Email
	// (POST /email/templates/{templateName}/send)
	SendTemplatedEmail(w http.ResponseWriter, r *http.Request, templateName string)
	// Get describing html of openapi spec
	// (GET /info/openapi.html)
	GetOpenAPIHTML(w http.ResponseWriter, r *http.Request)
	// Get openapi spec as json
	// (GET /info/openapi.json)
	GetOpenAPIJSON(w http.ResponseWriter, r *http.Request)
	// Get status of the service
	// (GET /info/status)
	GetStatus(w http.ResponseWriter, r *http.Request)
	// Get version info of the service
	// (GET /info/version)
	GetVersion(w http.ResponseWriter, r *http.Request)
	// Save a static PDF
	// (POST /pdf/basic)
	SavePDF(w http.ResponseWriter, r *http.Request)
	// Get PDF Template by Name
	// (GET /pdf/templates/{templateName})
	GetPDFTemplateByName(w http.ResponseWriter, r *http.Request, templateName string)
	// Add new PDF template
	// (POST /pdf/templates/{templateName})
	AddNewPDFTemplate(w http.ResponseWriter, r *http.Request, templateName string)
	// Fill placeholders of PDF template
	// (POST /pdf/templates/{templateName}/fill)
	FillPDFTemplate(w http.ResponseWriter, r *http.Request, templateName string)
	// Get PDF Template Placeholders
	// (GET /pdf/templates/{templateName}/placeholders)
	GetPDFTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request, templateName string)
	// Send a SMS with custom text
	// (POST /sms/basic/send)
	SendBasicSMS(w http.ResponseWriter, r *http.Request)
	// Send a templated SMS with custom text
	// (POST /sms/template/{templateName}/send)
	SendTemplatedSMS(w http.ResponseWriter, r *http.Request, templateName string)
	// Get SMS Template by Name
	// (GET /sms/templates/{templateName})
	GetSMSTemplateByName(w http.ResponseWriter, r *http.Request, templateName string)
	// Add new SMS template
	// (POST /sms/templates/{templateName})
	AddNewSMSTemplate(w http.ResponseWriter, r *http.Request, templateName string)
	// Fill placeholders of SMS template
	// (POST /sms/templates/{templateName}/fill)
	FillSMSTemplate(w http.ResponseWriter, r *http.Request, templateName string)
	// Get SMS Template Placeholders
	// (GET /sms/templates/{templateName}/placeholders)
	GetSMSTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request, templateName string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// SendEmail operation middleware
func (siw *ServerInterfaceWrapper) SendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SendEmail(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetTemplateByName operation middleware
func (siw *ServerInterfaceWrapper) GetTemplateByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTemplateByName(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// AddNewTemplate operation middleware
func (siw *ServerInterfaceWrapper) AddNewTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddNewTemplate(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// FillTemplate operation middleware
func (siw *ServerInterfaceWrapper) FillTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FillTemplate(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetTemplatePlaceholdersByName operation middleware
func (siw *ServerInterfaceWrapper) GetTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTemplatePlaceholdersByName(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// SendTemplatedEmail operation middleware
func (siw *ServerInterfaceWrapper) SendTemplatedEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SendTemplatedEmail(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetOpenAPIHTML operation middleware
func (siw *ServerInterfaceWrapper) GetOpenAPIHTML(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOpenAPIHTML(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetOpenAPIJSON operation middleware
func (siw *ServerInterfaceWrapper) GetOpenAPIJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOpenAPIJSON(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetStatus operation middleware
func (siw *ServerInterfaceWrapper) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetStatus(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetVersion operation middleware
func (siw *ServerInterfaceWrapper) GetVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetVersion(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// SavePDF operation middleware
func (siw *ServerInterfaceWrapper) SavePDF(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SavePDF(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetPDFTemplateByName operation middleware
func (siw *ServerInterfaceWrapper) GetPDFTemplateByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPDFTemplateByName(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// AddNewPDFTemplate operation middleware
func (siw *ServerInterfaceWrapper) AddNewPDFTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddNewPDFTemplate(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// FillPDFTemplate operation middleware
func (siw *ServerInterfaceWrapper) FillPDFTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FillPDFTemplate(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetPDFTemplatePlaceholdersByName operation middleware
func (siw *ServerInterfaceWrapper) GetPDFTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetPDFTemplatePlaceholdersByName(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// SendBasicSMS operation middleware
func (siw *ServerInterfaceWrapper) SendBasicSMS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SendBasicSMS(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// SendTemplatedSMS operation middleware
func (siw *ServerInterfaceWrapper) SendTemplatedSMS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SendTemplatedSMS(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetSMSTemplateByName operation middleware
func (siw *ServerInterfaceWrapper) GetSMSTemplateByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSMSTemplateByName(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// AddNewSMSTemplate operation middleware
func (siw *ServerInterfaceWrapper) AddNewSMSTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddNewSMSTemplate(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// FillSMSTemplate operation middleware
func (siw *ServerInterfaceWrapper) FillSMSTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FillSMSTemplate(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetSMSTemplatePlaceholdersByName operation middleware
func (siw *ServerInterfaceWrapper) GetSMSTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "templateName" -------------
	var templateName string

	err = runtime.BindStyledParameterWithOptions("simple", "templateName", r.PathValue("templateName"), &templateName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "templateName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSMSTemplatePlaceholdersByName(w, r, templateName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       *http.ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m *http.ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m *http.ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("POST "+options.BaseURL+"/email/basic/send", wrapper.SendEmail)
	m.HandleFunc("GET "+options.BaseURL+"/email/templates/{templateName}", wrapper.GetTemplateByName)
	m.HandleFunc("POST "+options.BaseURL+"/email/templates/{templateName}", wrapper.AddNewTemplate)
	m.HandleFunc("POST "+options.BaseURL+"/email/templates/{templateName}/fill", wrapper.FillTemplate)
	m.HandleFunc("GET "+options.BaseURL+"/email/templates/{templateName}/placeholders", wrapper.GetTemplatePlaceholdersByName)
	m.HandleFunc("POST "+options.BaseURL+"/email/templates/{templateName}/send", wrapper.SendTemplatedEmail)
	m.HandleFunc("GET "+options.BaseURL+"/info/openapi.html", wrapper.GetOpenAPIHTML)
	m.HandleFunc("GET "+options.BaseURL+"/info/openapi.json", wrapper.GetOpenAPIJSON)
	m.HandleFunc("GET "+options.BaseURL+"/info/status", wrapper.GetStatus)
	m.HandleFunc("GET "+options.BaseURL+"/info/version", wrapper.GetVersion)
	m.HandleFunc("POST "+options.BaseURL+"/pdf/basic", wrapper.SavePDF)
	m.HandleFunc("GET "+options.BaseURL+"/pdf/templates/{templateName}", wrapper.GetPDFTemplateByName)
	m.HandleFunc("POST "+options.BaseURL+"/pdf/templates/{templateName}", wrapper.AddNewPDFTemplate)
	m.HandleFunc("POST "+options.BaseURL+"/pdf/templates/{templateName}/fill", wrapper.FillPDFTemplate)
	m.HandleFunc("GET "+options.BaseURL+"/pdf/templates/{templateName}/placeholders", wrapper.GetPDFTemplatePlaceholdersByName)
	m.HandleFunc("POST "+options.BaseURL+"/sms/basic/send", wrapper.SendBasicSMS)
	m.HandleFunc("POST "+options.BaseURL+"/sms/template/{templateName}/send", wrapper.SendTemplatedSMS)
	m.HandleFunc("GET "+options.BaseURL+"/sms/templates/{templateName}", wrapper.GetSMSTemplateByName)
	m.HandleFunc("POST "+options.BaseURL+"/sms/templates/{templateName}", wrapper.AddNewSMSTemplate)
	m.HandleFunc("POST "+options.BaseURL+"/sms/templates/{templateName}/fill", wrapper.FillSMSTemplate)
	m.HandleFunc("GET "+options.BaseURL+"/sms/templates/{templateName}/placeholders", wrapper.GetSMSTemplatePlaceholdersByName)

	return m
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
