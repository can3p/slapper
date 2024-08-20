package curl

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/can3p/slapper/pkg/requests"
	"github.com/jessevdk/go-flags"
	"github.com/kballard/go-shellquote"
)

type curlCmd struct {
	Method  string   `short:"X" description:"http method to use"`
	Header  []string `short:"H" description:"header value"`
	DataRaw string   `long:"data-raw" description:"raw data for request body"`
}

func ParseCommand(cmd string) (*requests.Request, error) {
	words, err := shellquote.Split(cmd)

	if err != nil {
		return nil, fmt.Errorf("Failed to split the command string, %w", err)
	}

	if len(words) == 0 {
		return nil, fmt.Errorf("empty command string")
	}

	if words[0] != "curl" {
		return nil, fmt.Errorf("not a curl command: %s", cmd)
	}

	var opts curlCmd

	args, err := flags.NewParser(&opts, flags.IgnoreUnknown).ParseArgs(words[1:])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse the arguments, %w", err)
	}

	// skip all unknown flags
	finalArgs := []string{}
	idx := 0

	for idx < len(args) {
		if strings.HasPrefix(args[idx], "--") {
			// skip unknown options like --arg val
			if idx+1 < len(args) && !strings.HasPrefix(args[idx+1], "-") {
				idx += 2
				continue
			}

			idx++
			continue
		}
		if strings.HasPrefix(args[idx], "-") {
			idx++
			continue
		}

		finalArgs = append(finalArgs, args[idx])
		idx++
	}

	if len(finalArgs) == 0 {
		return nil, fmt.Errorf("no url specified")
	}

	if len(finalArgs) > 1 {
		return nil, fmt.Errorf("to many positional arguments to curl: %v", finalArgs)
	}

	out := &requests.Request{
		Url:    args[0],
		Header: http.Header{},
	}

	for _, h := range opts.Header {
		parts := strings.SplitN(h, ":", 2)
		out.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}

	if opts.DataRaw != "" {
		out.Body = []byte(opts.DataRaw)
	}

	if opts.Method == "" {
		if len(opts.DataRaw) > 0 {
			out.Method = "POST"
		} else {
			out.Method = "GET"
		}
	} else {
		out.Method = opts.Method
	}

	return out, nil
}

func IsCurl(line string) bool {
	return strings.HasPrefix(line, "curl ")
}
