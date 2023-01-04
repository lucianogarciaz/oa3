// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/aarondl/opt/omit"
	"golang.org/x/time/rate"
)

type ctxKey string

const (
	ctxKeyDebug ctxKey = "debug"
)

// URLBuilder builds a base url. Implementations are likely simple fixed
// strings or slightly more complicated variable replacement strings with
// defaults.
//
// Implementors:
// - Httpdevlocal
// - Httpprodlocalonetwo
// - Httpvariableslocalvariable
type URLBuilder interface {
	ToURL() string
}

// URL is a simple base url builder that's just a static string
type URL string

func (b URL) ToURL() string { return string(b) }

// Local development
var Httpdevlocal = URL(`http://dev.local:3030`)

// Production
type Httpprodlocalonetwo struct {
	One string
	Two string
}

func (s Httpprodlocalonetwo) ToURL() string {
	uri := `http://prod.local:3030/{one}/{two}`
	if len(s.One) != 0 {
		uri = strings.ReplaceAll(uri, `{one}`, s.One)
	} else {
		uri = strings.ReplaceAll(uri, `{one}`, `one`)
	}
	if len(s.Two) != 0 {
		uri = strings.ReplaceAll(uri, `{two}`, s.Two)
	} else {
		uri = strings.ReplaceAll(uri, `{two}`, `two`)
	}
	return uri
}

// Variable path
type Httpvariableslocalvariable struct {
	Variable string
}

func (s Httpvariableslocalvariable) ToURL() string {
	uri := `http://variables.local:3030/{variable}`
	if len(s.Variable) != 0 {
		switch s.Variable {
		case `v1`, `v2`, `v3`:
		default:
			panic("unknown server variable enum value: " + s.Variable)
		}
		uri = strings.ReplaceAll(uri, `{variable}`, s.Variable)
	} else {
		uri = strings.ReplaceAll(uri, `{variable}`, `v1`)
	}
	return uri
}

// URLTestservers is a simple string url like URLSimple
type URLTestservers string

func (b URLTestservers) ToURL() string       { return string(b) }
func (b URLTestservers) TestserversSatisfy() {}

// URLBuilderTestservers builds a base url like URLBuilder
// but restricts the implementing types to a smaller subset.
//
// Implementors:
// - Httppathdevlocal
// - Httppathprodlocalonetwo
// - Httppathvariableslocalvariable
type URLBuilderTestservers interface {
	URLBuilder
	TestserversSatisfy()
}

// Local development
var Httppathdevlocal = URLTestservers(`http://path.dev.local:3030`)

// Production
type Httppathprodlocalonetwo struct {
	One string
	Two string
}

func (Httppathprodlocalonetwo) TestserversToURL() {}
func (s Httppathprodlocalonetwo) ToURL() string {
	uri := `http://path.prod.local:3030/{one}/{two}`
	if len(s.One) != 0 {
		uri = strings.ReplaceAll(uri, `{one}`, s.One)
	} else {
		uri = strings.ReplaceAll(uri, `{one}`, `one`)
	}
	if len(s.Two) != 0 {
		uri = strings.ReplaceAll(uri, `{two}`, s.Two)
	} else {
		uri = strings.ReplaceAll(uri, `{two}`, `two`)
	}
	return uri
}

// Variable path
type Httppathvariableslocalvariable struct {
	Variable string
}

func (Httppathvariableslocalvariable) TestserversToURL() {}
func (s Httppathvariableslocalvariable) ToURL() string {
	uri := `http://path.variables.local:3030/{variable}`
	if len(s.Variable) != 0 {
		switch s.Variable {
		case `v1`, `v2`, `v3`:
		default:
			panic("unknown server variable enum value: " + s.Variable)
		}
		uri = strings.ReplaceAll(uri, `{variable}`, s.Variable)
	} else {
		uri = strings.ReplaceAll(uri, `{variable}`, `v1`)
	}
	return uri
}

// URLTestserversPost is a simple url
type URLTestserversPost string

func (b URLTestserversPost) ToURL() string           { return string(b) }
func (b URLTestserversPost) TestserversPostSatisfy() {}

// URLBuilderTestserversPost builds a base url like URLBuilder
// but restricts the implementing types to a smaller subset.
//
// Implementors:
// - Httpopdevlocal
// - Httpopprodlocalonetwo
// - Httpopvariableslocalvariable
type URLBuilderTestserversPost interface {
	URLBuilder
	TestserversPostSatisfy()
}

// Local development
var Httpopdevlocal = URLTestserversPost(`http://op.dev.local:3030`)

// Production
type Httpopprodlocalonetwo struct {
	One string
	Two string
}

func (Httpopprodlocalonetwo) TestserversPostToURL() {}
func (s Httpopprodlocalonetwo) ToURL() string {
	uri := `http://op.prod.local:3030/{one}/{two}`
	if len(s.One) != 0 {
		uri = strings.ReplaceAll(uri, `{one}`, s.One)
	} else {
		uri = strings.ReplaceAll(uri, `{one}`, `one`)
	}
	if len(s.Two) != 0 {
		uri = strings.ReplaceAll(uri, `{two}`, s.Two)
	} else {
		uri = strings.ReplaceAll(uri, `{two}`, `two`)
	}
	return uri
}

// Variable path
type Httpopvariableslocalvariable struct {
	Variable string
}

func (Httpopvariableslocalvariable) TestserversPostToURL() {}
func (s Httpopvariableslocalvariable) ToURL() string {
	uri := `http://op.variables.local:3030/{variable}`
	if len(s.Variable) != 0 {
		switch s.Variable {
		case `v1`, `v2`, `v3`:
		default:
			panic("unknown server variable enum value: " + s.Variable)
		}
		uri = strings.ReplaceAll(uri, `{variable}`, s.Variable)
	} else {
		uri = strings.ReplaceAll(uri, `{variable}`, `v1`)
	}
	return uri
}

var (
	apiHTTPClient = &http.Client{Timeout: time.Second * 5}
)

// Client is a generated package for consuming an openapi spec.
//
// A great api
type Client struct {
	httpClient  *http.Client
	httpHandler http.Handler
	limiter     *rate.Limiter

	url URLBuilder
}

// WithDebug creates a context that will emit debugging information to stdout
// for each request.
func WithDebug(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyDebug, "t")
}

func hasDebug(ctx context.Context) bool {
	v := ctx.Value(ctxKeyDebug)
	return v != nil && v.(string) == "t"
}

// NewClient constructs an api client, optionally using a supplied http.Client
// to be able to add instrumentation or customized timeouts.
//
// If nil is supplied then this package's generated apiHTTPClient is used which
// has reasonable defaults for timeouts.
//
// It also takes an optional rate limiter to implement rate limiting.
func NewClient(httpClient *http.Client, limiter *rate.Limiter, baseURL URLBuilder) Client {
	client := Client{httpClient: apiHTTPClient, limiter: limiter, url: baseURL}
	if httpClient != nil {
		client.httpClient = httpClient
	}
	return client
}

// NewLocalClient constructs an api client, but takes in a handler to call
// with the prepared requests instead of an http client that will touch the
// network. Useful for testing.
func NewLocalClient(httpHandler http.Handler) Client {
	return Client{httpHandler: httpHandler, url: URL("http://localhost")}
}

// WithURL sets the url for this client, the client is a shallow clone and
// therefore still shares the same http client, handler and rate limiter.
func (c Client) WithURL(url URLBuilder) Client {
	newClient := c
	newClient.url = url
	return newClient
}

func (c Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	if hasDebug(ctx) {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, fmt.Errorf("failed to emit debugging info: %w", err)
		}
		fmt.Printf("%s\n", reqDump)
	}

	var resp *http.Response
	if c.httpHandler != nil {
		w := httptest.NewRecorder()
		c.httpHandler.ServeHTTP(w, req)
		resp = w.Result()
	} else {
		var err error
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
	}

	if hasDebug(ctx) {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("failed to emit debugging info: %w", err)
		}
		fmt.Printf("%s\n", respDump)
	}

	return resp, nil
}

// TestInlineResponse one-of enforcer
//
// Implementors:
// - TestInline200Inline
// - TestInline201Inline
type TestInlineResponse interface {
	TestInlineImpl()
}

// TestInlineImpl implements TestInlineResponse(200) for TestInline200Inline
func (TestInline200Inline) TestInlineImpl() {}

// TestInlineImpl implements TestInlineResponse(201) for TestInline201Inline
func (TestInline201Inline) TestInlineImpl() {}

// TestInlineResponseComponentMultipleResponse one-of enforcer
//
// Implementors:
// - InlineResponseTestInline
// - HTTPStatusCreated
type TestInlineResponseComponentMultipleResponse interface {
	TestInlineResponseComponentMultipleImpl()
}

// TestInlineResponseComponentMultipleImpl implements TestInlineResponseComponentMultipleResponse(200) for InlineResponseTestInline
func (InlineResponseTestInline) TestInlineResponseComponentMultipleImpl() {}

// TestInlineResponseComponentMultipleImpl implements TestInlineResponseComponentMultipleResponse(201) for HTTPStatusCreated
func (HTTPStatusCreated) TestInlineResponseComponentMultipleImpl() {}

// SetUserResponse one-of enforcer
//
// Implementors:
// - SetUserWrappedResponse
type SetUserResponse interface {
	SetUserImpl()
}

// SetUserWrappedResponse wraps the normal body response with a
// struct to be able to additionally return headers or differentiate between
// multiple response codes with the same response body.
type SetUserWrappedResponse struct {
	Code                  int
	HeaderXResponseHeader omit.Val[string]
	Body                  Primitives
}

// SetUserImpl implements SetUserResponse(200) for SetUserWrappedResponse
func (SetUserWrappedResponse) SetUserImpl() {}

// HTTPStatusCreated is an empty response
type HTTPStatusCreated struct{}

// HTTPStatusNotModified is an empty response
type HTTPStatusNotModified struct{}

// HTTPStatusOk is an empty response
type HTTPStatusOk struct{}
type TestUnknownBodyType200Inline io.ReadCloser
