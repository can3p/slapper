package runner

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/can3p/slapper/pkg/requests"
)

type targetTest struct {
	input    string
	base64   bool
	expected []requests.Request
}

var tests = []targetTest{
	{
		input: `POST http://127.0.0.1:5000/test`,
		expected: []requests.Request{
			{
				Method: "POST",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte{},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte{},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
GET http://127.0.0.1:5000/test`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte{},
			},
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte{},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
GET http://127.0.0.1:5000/test
$ {"foo": "bar"}`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte{},
			},
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte(`{"foo": "bar"}`),
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
$ {"foo": "bar"}
GET http://127.0.0.1:5000/test`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte(`{"foo": "bar"}`),
			},
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte{},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
{}
`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte{},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
$ {"foo": "bar"}
`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte(`{"foo": "bar"}`),
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
$ {"foo": "bar"}

`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte(`{"foo": "bar"}`),
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
$ {"foo": "bar"}

GET http://www.example.com
$ {"spam": "eggs"}

`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte(`{"foo": "bar"}`),
			},
			{
				Method: "GET",
				Url:    "http://www.example.com",
				Body:   []byte(`{"spam": "eggs"}`),
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
$ Zm9v

`,
		base64: true,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte(`foo`),
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
H Content-Type: application/json
H X-Auth: 124
$ {"foo": "bar"}

GET http://www.example.com
H X-Extra: 124
H X-Extra: 125
$ {"spam": "eggs"}

`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte(`{"foo": "bar"}`),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
					"X-Auth":       {"124"},
				},
			},
			{
				Method: "GET",
				Url:    "http://www.example.com",
				Body:   []byte(`{"spam": "eggs"}`),
				Header: map[string][]string{
					"X-Auth": {"124", "125"},
				},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
H Content-Type: application/json
H X-Auth: 124
$ {"foo": "bar"}

curl http://www.example.com -X GET -H 'X-Extra: 124' -H 'X-Extra: 125' --data-raw='{"spam": "eggs"}'

`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte(`{"foo": "bar"}`),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
					"X-Auth":       {"124"},
				},
			},
			{
				Method: "GET",
				Url:    "http://www.example.com",
				Body:   []byte(`{"spam": "eggs"}`),
				Header: map[string][]string{
					"X-Auth": {"124", "125"},
				},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
H Content-Type: application/json
H X-Auth: 124
$ {"foo": "bar"}

curl http://www.example.com \
    -X GET \
    -H 'X-Extra: 124' \
    -H 'X-Extra: 125' \
    --data-raw='{"spam": "eggs"}'

`,
		expected: []requests.Request{
			{
				Method: "GET",
				Url:    "http://127.0.0.1:5000/test",
				Body:   []byte(`{"foo": "bar"}`),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
					"X-Auth":       {"124"},
				},
			},
			{
				Method: "GET",
				Url:    "http://www.example.com",
				Body:   []byte(`{"spam": "eggs"}`),
				Header: map[string][]string{
					"X-Auth": {"124", "125"},
				},
			},
		},
	},
}

func TestNewTargeter(t *testing.T) {
	failed := 0

	for idx, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.input))

		trgt := targeter{}
		err := trgt.readTargets(r, test.base64)
		if err != nil {
			t.Errorf("[%d] Unexpected error: %v", idx+1, err)
			failed++
			continue
		}

		if len(test.expected) != len(trgt.requests) {
			t.Errorf("Input: %+v\n", test)
			t.Errorf("Expected %d requests, got %d requests", len(test.expected), len(trgt.requests))
			failed++
			continue
		}

		for req := 0; req < len(trgt.requests); req++ {
			if test.expected[req].Method != trgt.requests[req].Method {
				t.Errorf("[%d/%d] Expected method '%s', got '%s'", idx+1, req+1, test.expected[req].Method, trgt.requests[req].Method)
				failed++
				break
			}

			if test.expected[req].Url != trgt.requests[req].Url {
				t.Errorf("[%d/%d] Expected URL '%s', got '%s'", idx+1, req+1, test.expected[req].Url, trgt.requests[req].Url)
				failed++
				break
			}

			if !bytes.Equal(test.expected[req].Body, trgt.requests[req].Body) {
				t.Errorf(`Bad request body
Expected	%+v
Got		%+v"`, test.expected[req].Body, trgt.requests[req].Body)
				failed++
				break
			}

		}
	}

	if failed > 0 {
		t.Logf("Failed %d/%d tests\n", failed, len(tests))
	}
}
