package test_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
	"testing"

	"git.sr.ht/~ewintr/go-kit/test"
)

func TestHTTPMock(t *testing.T) {

	procs := []test.MockServerProcedure{
		test.MockServerProcedure{
			URL:        "/",
			HTTPMethod: "GET",
			Response: test.MockResponse{
				Body: []byte("getRoot"),
			},
		},
		test.MockServerProcedure{
			URL:        "/",
			HTTPMethod: "POST",
			Response: test.MockResponse{
				Body: []byte("postRoot"),
			},
		},
		test.MockServerProcedure{
			URL:        "/get/header",
			HTTPMethod: "GET",
			Response: test.MockResponse{
				StatusCode: http.StatusAccepted,
				Headers: http.Header{
					"some-key": []string{"some-value"},
				},
				Body: []byte("getResponseHeader"),
			},
		},
		test.MockServerProcedure{
			URL:        "/get/auth",
			HTTPMethod: "GET",
			Response: test.MockResponse{
				Body: []byte("getRootAuth"),
			},
		},
		test.MockServerProcedure{
			URL:        "/my_account",
			HTTPMethod: "GET",
			Response: test.MockResponse{
				Body: []byte("getAccount"),
			},
		},
		test.MockServerProcedure{
			URL:        "/my_account.json",
			HTTPMethod: "GET",
			Response: test.MockResponse{
				Body: []byte("getAccountJSON"),
			},
		},
	}

	var record test.MockAssertion
	testMockServer := test.NewMockServer(&record, procs...)

	type mockRequest struct {
		uri            string
		method         string
		user, password string
		header         http.Header
		body           []byte
		hits           int
	}

	canonical := textproto.CanonicalMIMEHeaderKey

	for _, tc := range []struct {
		m        string
		request  mockRequest
		response test.MockResponse
	}{
		{
			m: "method get root path",
			request: mockRequest{
				uri:    "/",
				method: http.MethodGet,
				hits:   2,
			},
			response: test.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("getRoot"),
			},
		},
		{
			m: "method get root path with headers",
			request: mockRequest{
				uri:    "/",
				method: http.MethodGet,
				header: http.Header{
					canonical("input-header-key"): []string{"Just the Value"},
				},
				hits: 2,
			},
			response: test.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("getRoot"),
			},
		},
		{
			m: "method get root path with body",
			request: mockRequest{
				uri:    "/",
				method: http.MethodGet,
				body:   []byte("input"),
				hits:   2,
			},
			response: test.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("getRoot"),
			},
		},
		{
			m: "method get root path with headers and body",
			request: mockRequest{
				uri:    "/",
				method: http.MethodGet,
				header: http.Header{
					canonical("input-header-key"): []string{"Just the Value"},
				},
				body: []byte("input"),
				hits: 2,
			},
			response: test.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("getRoot"),
			},
		},
		{
			m: "method post root path",
			request: mockRequest{
				uri:    "/",
				method: http.MethodPost,
				hits:   2,
			},
			response: test.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("postRoot"),
			},
		},
		{
			m: "method post root path with basic authentication",
			request: mockRequest{
				uri:      "/",
				method:   http.MethodPost,
				user:     "my-user",
				password: "my-password",
				hits:     1,
			},
			response: test.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("postRoot"),
			},
		},
		{
			m: "unmatched uri path",
			request: mockRequest{
				uri:    "/unmatched",
				method: http.MethodGet,
				hits:   0,
			},
			response: test.MockResponse{
				StatusCode: http.StatusNotFound,
				Body:       []byte{},
			},
		},
	} {
		t.Run(tc.m, func(t *testing.T) {
			test.OK(t, record.Reset())

			for _ = range make([]int, tc.request.hits) {
				url, errU := url.Parse(testMockServer.URL + tc.request.uri)
				test.OK(t, errU)

				req, errReq := http.NewRequest(
					tc.request.method,
					url.String(),
					bytes.NewReader(tc.request.body),
				)
				test.OK(t, errReq)

				for k, v := range tc.request.header {
					req.Header[k] = v
				}

				// testing authentication in the request
				if len(tc.request.user) > 0 || len(tc.request.password) > 0 {
					req.SetBasicAuth(tc.request.user, tc.request.password)

					if tc.request.header == nil {
						tc.request.header = make(http.Header)
					}

					auth := tc.request.user + ":" + tc.request.password
					tc.request.header["Authorization"] = []string{
						fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))}
				}

				client := new(http.Client)
				resp, errResp := client.Do(req)
				test.OK(t, errResp)

				actualBody, err := ioutil.ReadAll(resp.Body)
				test.OK(t, err)
				defer resp.Body.Close()

				test.Equals(t, tc.response.StatusCode, resp.StatusCode)
				test.Equals(t, tc.response.Body, actualBody)
			}
			test.Equals(t, tc.request.hits, record.Hits(tc.request.uri, tc.request.method))

			// assert if all request had the correct header
			for _, h := range record.Headers(tc.request.uri, tc.request.method) {
				test.Equals(t, tc.request.header, h)
			}

			// assert if all request had the correct body
			for _, b := range record.Body(tc.request.uri, tc.request.method) {
				test.Equals(t, tc.request.body, b)
			}
		})
	}
}

func ExampleMockAssertion_Hits() {
	var record test.MockAssertion
	uri := "/"

	server := test.NewMockServer(&record, test.MockServerProcedure{
		URL:        uri,
		HTTPMethod: http.MethodGet,
	})

	http.Get(server.URL)

	fmt.Println(record.Hits(uri, http.MethodGet))
	// Output: 1
}

func ExampleMockAssertion_Headers() {
	var record test.MockAssertion
	uri := "/"

	server := test.NewMockServer(&record, test.MockServerProcedure{
		URL:        uri,
		HTTPMethod: http.MethodPost,
	})

	http.Post(server.URL, "application/json", nil)

	fmt.Println(record.Headers(uri, http.MethodPost))
	// Output: [map[Content-Type:[application/json]]]
}

func ExampleMockAssertion_Body() {
	var record test.MockAssertion
	uri := "/"

	server := test.NewMockServer(&record, test.MockServerProcedure{
		URL:        uri,
		HTTPMethod: http.MethodPost,
	})

	http.Post(server.URL, "text/plain", bytes.NewBufferString("hi there"))

	for _, b := range record.Body(uri, http.MethodPost) {
		fmt.Println(string(b))
	}
	// Output: hi there
}
