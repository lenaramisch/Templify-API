package router

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	server "templify/pkg/server/generated"
	"templify/pkg/server/handler/apihandler"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	oapiMW "github.com/oapi-codegen/nethttp-middleware"
)

// New returns a new *http.ServeMux to be used for HTTP request routing.
func New(apiHandle *apihandler.APIHandler, cfg *Config) *http.ServeMux {
	// Create a new ServeMux
	mux := http.NewServeMux()
	handler := addMiddlewares(server.HandlerFromMux(apiHandle, nil), cfg)
	mux.Handle("/", handler)

	return mux
}

// addMiddleware chains together all the middleware functions.
func addMiddlewares(handler http.Handler, cfg *Config) http.Handler {
	handler = oapiMiddleware(handler)
	handler = corsMiddleware(handler, cfg.CORS)
	handler = timeoutMiddleware(handler, cfg.Timeout)
	handler = loggingMiddleware(handler)
	return handler
}

func oapiMiddleware(handler http.Handler) http.Handler {
	swagger, err := server.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Validate requests against OpenAPI spec
	validatorOptions := &oapiMW.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
				return nil
			},
		},
	}

	return oapiMW.OapiRequestValidatorWithOptions(swagger, validatorOptions)(handler)
}

func loggingMiddleware(next http.Handler) http.Handler {
	QuietDownRoutes := []string{
		"/info/version",
		"/info/status",
		"/info/openapi.json",
		"/info/openapi.html",
	}

	HideHeaders := []string{
		"Authorization",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request path is in the excludedPaths list
		shouldExclude := false
		for _, path := range QuietDownRoutes {
			if r.URL.Path == path {
				shouldExclude = true
				break
			}
		}

		// Log the request if it's not in the excluded list
		if !shouldExclude {
			// create duplicate of header map
			headers := make(http.Header)
			for k, v := range r.Header {
				headers[k] = v
			}
			// hide sensitive headers
			for _, header := range HideHeaders {
				// read the header value length
				length := len(strings.Join(headers[header], ""))
				redactedText := fmt.Sprintf("[REDACTED - %d bytes]", length)
				headers[header] = []string{redactedText}
			}

			slog.With(
				"Path", r.URL.Path,
				"Method", r.Method,
				"Header", headers,
				"Body", r.Body,
			).Debug("Request")
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})

}

// corsMiddleware adds CORS headers based on the provided configuration.
func corsMiddleware(next http.Handler, cfg CORSConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// Check if the origin is allowed
		allowedOrigin := false
		for _, o := range cfg.Origins {
			if o == "*" || o == origin {
				allowedOrigin = true
				break
			}
		}

		if allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		// Handle other CORS headers
		if r.Method == http.MethodOptions {
			// Preflight request handling
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.Methods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.Headers, ","))
			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// timeoutMiddleware adds timeout handling to requests.
func timeoutMiddleware(next http.Handler, timeout time.Duration) http.Handler {
	return http.TimeoutHandler(next, timeout, "Timeout")
}
