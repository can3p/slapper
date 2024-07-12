package runner

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

type targetTest struct {
	input    string
	base64   bool
	expected []request
}

var tests = []targetTest{
	{
		input: `POST http://127.0.0.1:5000/test`,
		expected: []request{
			{
				method: "POST",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte{},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test`,
		expected: []request{
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte{},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
GET http://127.0.0.1:5000/test`,
		expected: []request{
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte{},
			},
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte{},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
GET http://127.0.0.1:5000/test
$ {"foo": "bar"}`,
		expected: []request{
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte{},
			},
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte(`{"foo": "bar"}`),
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
$ {"foo": "bar"}
GET http://127.0.0.1:5000/test`,
		expected: []request{
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte(`{"foo": "bar"}`),
			},
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte{},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
{}
`,
		expected: []request{
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte{},
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
$ {"foo": "bar"}
`,
		expected: []request{
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte(`{"foo": "bar"}`),
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
$ {"foo": "bar"}

`,
		expected: []request{
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte(`{"foo": "bar"}`),
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
$ {"foo": "bar"}

GET http://www.example.com
$ {"spam": "eggs"}

`,
		expected: []request{
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte(`{"foo": "bar"}`),
			},
			{
				method: "GET",
				url:    "http://www.example.com",
				body:   []byte(`{"spam": "eggs"}`),
			},
		},
	},

	{
		input: `GET http://127.0.0.1:5000/test
$ Zm9v

`,
		base64: true,
		expected: []request{
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte(`foo`),
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
		expected: []request{
			{
				method: "GET",
				url:    "http://127.0.0.1:5000/test",
				body:   []byte(`{"foo": "bar"}`),
				header: map[string][]string{
					"Content-Type": {"application/json"},
					"X-Auth":       {"124"},
				},
			},
			{
				method: "GET",
				url:    "http://www.example.com",
				body:   []byte(`{"spam": "eggs"}`),
				header: map[string][]string{
					"X-Auth": {"124", "125"},
				},
			},
		},
	},
}

func TestNewTargeter(t *testing.T) {
	failed := 0

	for _, test := range tests {
		r := bufio.NewReader(strings.NewReader(test.input))

		trgt := targeter{}
		err := trgt.readTargets(r, test.base64)
		if err != nil {
			t.Error(err)
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
			if test.expected[req].method != trgt.requests[req].method {
				t.Errorf("Expected method '%s', got '%s'", test.expected[req].method, trgt.requests[req].method)
				failed++
				break
			}

			if test.expected[req].url != trgt.requests[req].url {
				t.Errorf("Expected URL '%s', got '%s'", test.expected[req].url, trgt.requests[req].url)
				failed++
				break
			}

			if !bytes.Equal(test.expected[req].body, trgt.requests[req].body) {
				t.Errorf(`Bad request body
Expected	%+v
Got		%+v"`, test.expected[req].body, trgt.requests[req].body)
				failed++
				break
			}

		}
	}

	if failed > 0 {
		t.Logf("Failed %d/%d tests\n", failed, len(tests))
	}
}
