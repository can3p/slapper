package curl

import (
	"net/http"
	"testing"

	"github.com/can3p/slapper/pkg/requests"
	"github.com/stretchr/testify/assert"
)

// curl 'https://www.livejournal.com/__api/' -X POST -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:129.0) Gecko/20100101 Firefox/129.0' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3' -H 'Accept-Encoding: gzip, deflate, br, zstd' -H 'Referer: https://can3p.livejournal.com/feed/' -H 'Content-Type: text/plain' -H 'Origin: https://can3p.livejournal.com' -H 'Connection: keep-alive' -H 'Cookie: luid=URNKAWV481hULnTiCPRWAgB=; ljuniq=BhtATj2nu6tnf1Y:1702425432:pgstats0; adtech_uid=3c9aa53c-df8f-43cd-b22f-5df835345dc0%3Alivejournal.com; top100_id=t1.1111412.1058155273.1702425432671; last_visit=1715547964603%3A%3A1715555164603; ab_d=1; bltsr=1; addruid=1m7x0p2M4E2u5l4k3iN79C6d4n; BMLschemepref=horizon; langpref=ru/1721674997; vpuid=1705190010.765-1725890650; prop_cookies_alert=1; t3_sid_1111412=s1.2007315197.1721781863357.1721782298278.3.14; splittest=none; prop_friendsfeed_tour=%7B%22three_posts_tour%22%3A0%2C%22friendsfeed%22%3A0%2C%22regionalrating%22%3A0%2C%22rss%22%3A0%2C%22video_update_tour%22%3A0%2C%22activity_engage_popup%22%3A%7B%22closeTimestamp%22%3A1721782303198%7D%7D; ljmastersession=v2:u6359713:s447:a1H9UVgoSJr:g84b1a9c39488ae25bd84a0d252ae06e44211b1b8//1; ljloggedin=v2:u6359713:s447:t1721674997:g237029c17a6c5d8d70b5fce14a284d545046caef; ljsession=v1:u6359713:s447:t1721674800:gb160ad94286d55ea147656aca427bfb055a7453c//1; crookie=hTuI6m6b8zwHuSSZE/uOQueNr1uBmauuNFH6iUyiQxGUXTIkAdw5TcMBAZOcFKuNv4gBnV0TlxNRHDVJXaVgGZIFnpw=; cmtchd=MTcyMzY1Nzk4NjYxNw==' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-site' --data-raw '[{"jsonrpc":"2.0","method":"notifications.get_events_counter","params":{"auth_token":"ajax:1724104800:6359713:447:/__api/::e230001f3e9ab8da2bac725e12095c603ffb5b91"},"id":199}]'

func TestParseCommand(t *testing.T) {
	var ex = []struct {
		cmd string
		out *requests.Request
	}{
		{
			cmd: `curl 'https://www.livejournal.com/__api/' -X POST -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:129.0) Gecko/20100101 Firefox/129.0' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3' -H 'Accept-Encoding: gzip, deflate, br, zstd' -H 'Referer: https://can3p.livejournal.com/feed/' -H 'Content-Type: text/plain' -H 'Origin: https://can3p.livejournal.com' -H 'Connection: keep-alive' -H 'Cookie: hide me' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-site' --data-raw '[{"jsonrpc":"2.0","method":"notifications.get_events_counter","params":{"auth_token":"ajax:1724104800:6359713:447:/__api/::e230001f3e9ab8da2bac725e12095c603ffb5b91"},"id":199}]'`,
			out: &requests.Request{
				Method: "POST",
				Url:    "https://www.livejournal.com/__api/",
				Body:   []byte(`[{"jsonrpc":"2.0","method":"notifications.get_events_counter","params":{"auth_token":"ajax:1724104800:6359713:447:/__api/::e230001f3e9ab8da2bac725e12095c603ffb5b91"},"id":199}]`),
				Header: http.Header{
					"User-Agent":      []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:129.0) Gecko/20100101 Firefox/129.0"},
					"Accept":          []string{"application/json, text/javascript, */*; q=0.01"},
					"Accept-Language": []string{"fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3"},
					"Accept-Encoding": []string{"gzip, deflate, br, zstd"},
					"Referer":         []string{"https://can3p.livejournal.com/feed/"},
					"Content-Type":    []string{"text/plain"},
					"Origin":          []string{"https://can3p.livejournal.com"},
					"Connection":      []string{"keep-alive"},
					"Cookie":          []string{"hide me"},
					"Sec-Fetch-Dest":  []string{"empty"},
					"Sec-Fetch-Mode":  []string{"cors"},
					"Sec-Fetch-Site":  []string{"same-site"},
				},
			},
		},
		{
			cmd: `curl http://www.example.com -X GET -H 'X-Extra: 124' -H 'X-Extra: 125' --data-raw='{"spam": "eggs"}'`,
			out: &requests.Request{
				Method: "GET",
				Url:    "http://www.example.com",
				Body:   []byte(`{"spam": "eggs"}`),
				Header: http.Header{
					"X-Extra": []string{"124", "125"},
				},
			},
		},
		{
			// unknown args should be simply ignored, order is not guaranteed
			cmd: `curl http://www.example.com -v --connect-timeout 20`,
			out: &requests.Request{
				Method: "GET",
				Url:    "http://www.example.com",
				Header: http.Header{},
			},
		},
	}

	for idx, cmd := range ex {
		req, err := ParseCommand(cmd.cmd)

		assert.NoError(t, err, "No errors expected in example %d", idx+1)
		assert.Equal(t, cmd.out, req, "request does not match in example %d", idx+1)
	}
}
