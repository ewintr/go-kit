package doc_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	httpmock "git.ewintr.nl/go-kit/doc"
	"git.ewintr.nl/go-kit/test"
)

func TestFooClientDoStuff(t *testing.T) {
	path := "/path"
	username := "username"
	password := "password"

	for _, tc := range []struct {
		name      string
		param     string
		respCode  int
		respBody  string
		expErr    error
		expResult string
	}{
		{
			name:     "upstream failure",
			respCode: http.StatusInternalServerError,
			expErr:   httpmock.ErrUpstreamFailure,
		},
		{
			name:     "incorrect response body",
			respCode: http.StatusOK,
			respBody: `{"what?`,
			expErr:   httpmock.ErrUpstreamFailure,
		},
		{
			name:      "valid response to bar",
			param:     "bar",
			respCode:  http.StatusOK,
			respBody:  `{"result":"ok"}`,
			expResult: "ok",
		},
		{
			name:      "valid response to baz",
			param:     "baz",
			respCode:  http.StatusOK,
			respBody:  `{"result":"also ok"}`,
			expResult: "also ok",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var record test.MockAssertion
			mockServer := test.NewMockServer(&record, test.MockServerProcedure{
				URI:        path,
				HTTPMethod: http.MethodPost,
				Response: test.MockResponse{
					StatusCode: tc.respCode,
					Body:       []byte(tc.respBody),
				},
			})

			client := httpmock.NewFooClient(mockServer.URL, username, password)

			actResult, actErr := client.DoStuff(tc.param)

			// ceck result
			test.Equals(t, true, errors.Is(actErr, tc.expErr))
			test.Equals(t, tc.expResult, actResult)

			// check request was done
			test.Equals(t, 1, record.Hits(path, http.MethodPost))

			// check request body
			expBody := fmt.Sprintf(`{"param":%q}`, tc.param)
			actBody := string(record.Body(path, http.MethodPost)[0])
			test.Equals(t, expBody, actBody)

			// check request headers
			expHeaders := []http.Header{{
				"Authorization": []string{"Basic dXNlcm5hbWU6cGFzc3dvcmQ="},
				"Content-Type":  []string{"application/json;charset=utf-8"},
			}}
			test.Equals(t, expHeaders, record.Headers(path, http.MethodPost))
		})
	}
}
