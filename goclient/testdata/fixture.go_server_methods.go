// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aarondl/chrono"
	"github.com/aarondl/json"
	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Authenticate post /auth
func (_c Client) Authenticate(ctx context.Context) (HTTPStatusOk, *http.Response, error) {
	var _resp HTTPStatusOk
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/auth`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = HTTPStatusOk{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestArrayRequest get /test/array/request
func (_c Client) TestArrayRequest(ctx context.Context, body TestArrayRequestInline) (HTTPStatusOk, *http.Response, error) {
	var _resp HTTPStatusOk
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/array/request`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodGet, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_bodyBytes, _err := json.Marshal(body)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_req.Body = io.NopCloser(bytes.NewReader(_bodyBytes))
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = HTTPStatusOk{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestMapsArrayInline get /test/arraymaps
func (_c Client) TestMapsArrayInline(ctx context.Context) (TestMapsArrayInline200Inline, *http.Response, error) {
	var _resp TestMapsArrayInline200Inline
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/arraymaps`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodGet, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		var _respObject TestMapsArrayInline200Inline
		_b, _err := io.ReadAll(_httpResp.Body)
		if _err != nil {
			return _resp, _httpResp, _err
		}
		if _err = json.Unmarshal(_b, &_respObject); _err != nil {
			return _resp, _httpResp, _err
		}
		_resp = _respObject
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestMapsArrayRef post /test/arraymaps
func (_c Client) TestMapsArrayRef(ctx context.Context) (TestMapsArrayRef200Inline, *http.Response, error) {
	var _resp TestMapsArrayRef200Inline
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/arraymaps`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		var _respObject TestMapsArrayRef200Inline
		_b, _err := io.ReadAll(_httpResp.Body)
		if _err != nil {
			return _resp, _httpResp, _err
		}
		if _err = json.Unmarshal(_b, &_respObject); _err != nil {
			return _resp, _httpResp, _err
		}
		_resp = _respObject
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestEnumQueryRequest get /test/enum/query/request
func (_c Client) TestEnumQueryRequest(ctx context.Context, body TestEnumQueryRequestInline, sort TestEnumQueryRequestGetSortParam) (HTTPStatusOk, *http.Response, error) {
	var _resp HTTPStatusOk
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/enum/query/request`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodGet, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_bodyBytes, _err := json.Marshal(body)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_req.Body = io.NopCloser(bytes.NewReader(_bodyBytes))
	var _query url.Values
	if _query == nil {
		_query = make(url.Values)
	}
	_query.Add(`sort`, fmt.Sprintf("%v", sort))
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = HTTPStatusOk{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestInlinePrimitiveBody get /test/inline
func (_c Client) TestInlinePrimitiveBody(ctx context.Context, body string) (HTTPStatusOk, *http.Response, error) {
	var _resp HTTPStatusOk
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/inline`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodGet, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_bodyBytes, _err := json.Marshal(body)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_req.Body = io.NopCloser(bytes.NewReader(_bodyBytes))
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = HTTPStatusOk{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestInline post /test/inline
func (_c Client) TestInline(ctx context.Context, body TestInlineInline) (TestInlineResponse, *http.Response, error) {
	var _resp TestInlineResponse
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/inline`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_bodyBytes, _err := json.Marshal(body)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_req.Body = io.NopCloser(bytes.NewReader(_bodyBytes))
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		var _respObject TestInline200Inline
		_b, _err := io.ReadAll(_httpResp.Body)
		if _err != nil {
			return _resp, _httpResp, _err
		}
		if _err = json.Unmarshal(_b, &_respObject); _err != nil {
			return _resp, _httpResp, _err
		}
		_resp = _respObject
	case 201:
		var _respObject TestInline201Inline
		_b, _err := io.ReadAll(_httpResp.Body)
		if _err != nil {
			return _resp, _httpResp, _err
		}
		if _err = json.Unmarshal(_b, &_respObject); _err != nil {
			return _resp, _httpResp, _err
		}
		_resp = _respObject
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestInlineResponseComponent post /test/inlineresponsecomponent
func (_c Client) TestInlineResponseComponent(ctx context.Context) (InlineResponseTestInline, *http.Response, error) {
	var _resp InlineResponseTestInline
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/inlineresponsecomponent`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		var _respObject InlineResponseTestInline
		_b, _err := io.ReadAll(_httpResp.Body)
		if _err != nil {
			return _resp, _httpResp, _err
		}
		if _err = json.Unmarshal(_b, &_respObject); _err != nil {
			return _resp, _httpResp, _err
		}
		_resp = _respObject
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestInlineResponseComponentMultiple post /test/inlineresponsecomponentmultiple
func (_c Client) TestInlineResponseComponentMultiple(ctx context.Context) (TestInlineResponseComponentMultipleResponse, *http.Response, error) {
	var _resp TestInlineResponseComponentMultipleResponse
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/inlineresponsecomponentmultiple`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		var _respObject InlineResponseTestInline
		_b, _err := io.ReadAll(_httpResp.Body)
		if _err != nil {
			return _resp, _httpResp, _err
		}
		if _err = json.Unmarshal(_b, &_respObject); _err != nil {
			return _resp, _httpResp, _err
		}
		_resp = _respObject
	case 201:
		_resp = HTTPStatusCreated{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestMapsInline get /test/maps
func (_c Client) TestMapsInline(ctx context.Context) (TestMapsInline200Inline, *http.Response, error) {
	var _resp TestMapsInline200Inline
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/maps`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodGet, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		var _respObject TestMapsInline200Inline
		_b, _err := io.ReadAll(_httpResp.Body)
		if _err != nil {
			return _resp, _httpResp, _err
		}
		if _err = json.Unmarshal(_b, &_respObject); _err != nil {
			return _resp, _httpResp, _err
		}
		_resp = _respObject
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestMapsRef post /test/maps
func (_c Client) TestMapsRef(ctx context.Context) (MapAny, *http.Response, error) {
	var _resp MapAny
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/maps`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		var _respObject MapAny
		_b, _err := io.ReadAll(_httpResp.Body)
		if _err != nil {
			return _resp, _httpResp, _err
		}
		if _err = json.Unmarshal(_b, &_respObject); _err != nil {
			return _resp, _httpResp, _err
		}
		_resp = _respObject
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestQueryIntArrayParam post /test/queryintarrayparam
func (_c Client) TestQueryIntArrayParam(ctx context.Context, intarray omit.Val[TestQueryIntArrayParamPostIntarrayParam], intarrayrequired TestQueryIntArrayParamPostIntarrayrequiredParam) (HTTPStatusOk, *http.Response, error) {
	var _resp HTTPStatusOk
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/queryintarrayparam`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if _query == nil {
		_query = make(url.Values)
	}
	if _val, _ok := intarray.Get(); _ok {
		for _, _v := range _val {
			_query.Add(`intarray`, fmt.Sprintf("%v", _v))
		}
	}
	for _, _v := range intarrayrequired {
		_query.Add(`intarrayrequired`, fmt.Sprintf("%v", _v))
	}
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = HTTPStatusOk{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestServerPathOverrideRequest get /test/servers
func (_c Client) TestServerPathOverrideRequest(ctx context.Context, baseURL URLBuilderTestservers) (HTTPStatusOk, *http.Response, error) {
	var _resp HTTPStatusOk
	var _httpResp *http.Response
	var _err error
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/servers`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodGet, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = HTTPStatusOk{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestServerOpOverrideRequest post /test/servers
func (_c Client) TestServerOpOverrideRequest(ctx context.Context, baseURL URLBuilderTestserversPost) (HTTPStatusOk, *http.Response, error) {
	var _resp HTTPStatusOk
	var _httpResp *http.Response
	var _err error
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/servers`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = HTTPStatusOk{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestSingleServerPathOverrideRequest get /test/single_servers
func (_c Client) TestSingleServerPathOverrideRequest(ctx context.Context) (HTTPStatusOk, *http.Response, error) {
	var _resp HTTPStatusOk
	var _httpResp *http.Response
	var _err error
	baseURL := Httppathdevlocal
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/single_servers`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodGet, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = HTTPStatusOk{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestSingleServerOpOverrideRequest post /test/single_servers
func (_c Client) TestSingleServerOpOverrideRequest(ctx context.Context) (HTTPStatusOk, *http.Response, error) {
	var _resp HTTPStatusOk
	var _httpResp *http.Response
	var _err error
	baseURL := Httpopdevlocal
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/single_servers`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = HTTPStatusOk{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestTypeOverrides get /test/type_overrides
func (_c Client) TestTypeOverrides(ctx context.Context, body *Primitives, number decimal.Decimal, date chrono.Date, numberNull null.Val[decimal.Decimal], dateNull null.Val[chrono.Date], numberNonReq omit.Val[decimal.Decimal], dateNonReq omit.Val[chrono.Date]) (HTTPStatusOk, *http.Response, error) {
	var _resp HTTPStatusOk
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/type_overrides`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodGet, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_bodyBytes, _err := json.Marshal(body)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_req.Body = io.NopCloser(bytes.NewReader(_bodyBytes))
	var _query url.Values
	if _query == nil {
		_query = make(url.Values)
	}
	_query.Add(`number`, fmt.Sprintf("%v", number))
	_query.Add(`date`, fmt.Sprintf("%v", date))
	if _val, _ok := numberNull.Get(); _ok {
		_query.Add(`number_null`, fmt.Sprintf("%v", _val))
	}
	if _val, _ok := dateNull.Get(); _ok {
		_query.Add(`date_null`, fmt.Sprintf("%v", _val))
	}
	if _val, _ok := numberNonReq.Get(); _ok {
		_query.Add(`number_non_req`, fmt.Sprintf("%v", _val))
	}
	if _val, _ok := dateNonReq.Get(); _ok {
		_query.Add(`date_non_req`, fmt.Sprintf("%v", _val))
	}
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = HTTPStatusOk{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// TestUnknownBodyType post /test/unknown/body/type
func (_c Client) TestUnknownBodyType(ctx context.Context, body io.ReadCloser) (TestUnknownBodyType200Inline, *http.Response, error) {
	var _resp TestUnknownBodyType200Inline
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/test/unknown/body/type`
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_req.Body = body
	var _query url.Values
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		_resp = _httpResp.Body
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// GetUser get /users/{id}
//
// Retrieves a user with a long description that spans multiple lines so
// that we can see that both wrapping and long-line support is not
// bleeding over the sacred 80 char limit.
func (_c Client) GetUser(ctx context.Context, id string, paramComponent string, validStr omitnull.Val[GetUserGetValidStrParam], reqValidStr null.Val[GetUserGetReqValidStrParam], validInt omit.Val[int], reqValidInt int, validNum omit.Val[float64], reqValidNum float64, validBool omit.Val[bool], reqValidBool bool, reqStrFormat uuid.UUID, dateTime chrono.DateTime, date chrono.Date, timeVal chrono.Time, durationVal time.Duration, arrayPrimExplode omit.Val[GetUserGetArrayPrimExplodeParam], arrayPrimFlat GetUserGetArrayPrimFlatParam, arrayPrimIntExplode omit.Val[GetUserGetArrayPrimIntExplodeParam], arrayPrimIntFlat GetUserGetArrayPrimIntFlatParam, arrayEnumExplode omit.Val[GetUserGetArrayEnumExplodeParam], arrayEnumFlat GetUserGetArrayEnumFlatParam) (HTTPStatusNotModified, *http.Response, error) {
	var _resp HTTPStatusNotModified
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/users/{id}`
	_urlStr = strings.Replace(_urlStr, `{id}`, fmt.Sprintf("%v", id), 1)
	_req, _err := http.NewRequestWithContext(ctx, http.MethodGet, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	if _val, _ok := validStr.Get(); _ok {
		_req.Header.Add(`valid_str`, fmt.Sprintf("%v", _val))
	}
	var _query url.Values
	if _query == nil {
		_query = make(url.Values)
	}
	_query.Add(`param_component`, fmt.Sprintf("%v", paramComponent))
	if _val, _ok := reqValidStr.Get(); _ok {
		_query.Add(`req_valid_str`, fmt.Sprintf("%v", _val))
	}
	if _val, _ok := validInt.Get(); _ok {
		_query.Add(`valid_int`, fmt.Sprintf("%v", _val))
	}
	_query.Add(`req_valid_int`, fmt.Sprintf("%v", reqValidInt))
	if _val, _ok := validNum.Get(); _ok {
		_query.Add(`valid_num`, fmt.Sprintf("%v", _val))
	}
	_query.Add(`req_valid_num`, fmt.Sprintf("%v", reqValidNum))
	if _val, _ok := validBool.Get(); _ok {
		_query.Add(`valid_bool`, fmt.Sprintf("%v", _val))
	}
	_query.Add(`req_valid_bool`, fmt.Sprintf("%v", reqValidBool))
	_query.Add(`req_str_format`, fmt.Sprintf("%v", reqStrFormat))
	_query.Add(`date_time`, fmt.Sprintf("%v", dateTime))
	_query.Add(`date`, fmt.Sprintf("%v", date))
	_query.Add(`time_val`, fmt.Sprintf("%v", timeVal))
	_query.Add(`duration_val`, fmt.Sprintf("%v", durationVal))
	if _val, _ok := arrayPrimExplode.Get(); _ok {
		for _, _v := range _val {
			_query.Add(`array_prim_explode`, fmt.Sprintf("%v", _v))
		}
	}
	var _arrayPrimFlatSlice []string
	for _, _v := range arrayPrimFlat {
		_arrayPrimFlatSlice = append(_arrayPrimFlatSlice, fmt.Sprintf("%v", _v))
	}
	_query.Set(`array_prim_flat`, strings.Join(_arrayPrimFlatSlice, ","))
	if _val, _ok := arrayPrimIntExplode.Get(); _ok {
		for _, _v := range _val {
			_query.Add(`array_prim_int_explode`, fmt.Sprintf("%v", _v))
		}
	}
	var _arrayPrimIntFlatSlice []string
	for _, _v := range arrayPrimIntFlat {
		_arrayPrimIntFlatSlice = append(_arrayPrimIntFlatSlice, fmt.Sprintf("%v", _v))
	}
	_query.Set(`array_prim_int_flat`, strings.Join(_arrayPrimIntFlatSlice, ","))
	if _val, _ok := arrayEnumExplode.Get(); _ok {
		for _, _v := range _val {
			_query.Add(`array_enum_explode`, fmt.Sprintf("%v", _v))
		}
	}
	var _arrayEnumFlatSlice []string
	for _, _v := range arrayEnumFlat {
		_arrayEnumFlatSlice = append(_arrayEnumFlatSlice, fmt.Sprintf("%v", _v))
	}
	_query.Set(`array_enum_flat`, strings.Join(_arrayEnumFlatSlice, ","))
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 304:
		_resp = HTTPStatusNotModified{}
	default:
		return _resp, _httpResp, fmt.Errorf("unknown response code %d", _httpResp.StatusCode)
	}

	return _resp, _httpResp, nil
}

// SetUser post /users/{id}
//
// Sets a user
func (_c Client) SetUser(ctx context.Context, body *Primitives, id string, paramComponent string) (SetUserResponse, *http.Response, error) {
	var _resp SetUserResponse
	var _httpResp *http.Response
	var _err error
	baseURL := _c.url
	_urlStr := strings.TrimSuffix(baseURL.ToURL(), "/") + `/users/{id}`
	_urlStr = strings.Replace(_urlStr, `{id}`, fmt.Sprintf("%v", id), 1)
	_req, _err := http.NewRequestWithContext(ctx, http.MethodPost, _urlStr, nil)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_bodyBytes, _err := json.Marshal(body)
	if _err != nil {
		return _resp, _httpResp, _err
	}
	_req.Body = io.NopCloser(bytes.NewReader(_bodyBytes))
	var _query url.Values
	if _query == nil {
		_query = make(url.Values)
	}
	_query.Add(`param_component`, fmt.Sprintf("%v", paramComponent))
	if len(_query) > 0 {
		_req.URL.RawQuery = _query.Encode()
	}

	_httpResp, _err = _c.doRequest(ctx, _req)
	if _err != nil {
		return _resp, _httpResp, _err
	}

	switch _httpResp.StatusCode {
	case 200:
		var _respObject SetUserWrappedResponse
		_b, _err := io.ReadAll(_httpResp.Body)
		if _err != nil {
			return _resp, _httpResp, _err
		}
		if _err = json.Unmarshal(_b, &_respObject.Body); _err != nil {
			return _resp, _httpResp, _err
		}
		if hdr := _httpResp.Header.Get(`X-Response-Header`); len(hdr) != 0 {
			_respObject.HeaderXResponseHeader.Set(hdr)
		}
		_resp = _respObject
	default:
		var _respObject Primitives
		_b, _err := io.ReadAll(_httpResp.Body)
		if _err != nil {
			return _resp, _httpResp, _err
		}
		if _err = json.Unmarshal(_b, &_respObject); _err != nil {
			return _resp, _httpResp, _err
		}
		_resp = _respObject
	}

	return _resp, _httpResp, nil
}
