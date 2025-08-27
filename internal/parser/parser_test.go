package parser

import "testing"

func TestParse_SupportedFormats(t *testing.T){
	cases :=[] struct{
		name string
		line string
		want struct{
			method string
			path string
			status int
			msMin float64
			msMax float64
		}
	}{
		{
			name: "nginx_request_time_secs",
			line: `127.0.0.1 - - [10/Oct/2024:13:55:36 +0000] "GET /v1/foo?id=1 HTTP/1.1" 200 123 "-" "curl/7.64" request_time=0.120`,
			want: struct {
				method string; path string; status int; msMin, msMax float64
			}{"GET", "/v1/foo", 200, 119, 121},
		},
		{
			name: "rt_ms",
			line: `127.0.0.1 - - [10/Oct/2024:13:55:37 +0000] "POST /v1/bar HTTP/1.1" 201 321 "-" "curl/7.64" rt=350ms`,
			want: struct {
				method string; path string; status int; msMin, msMax float64
			}{"POST", "/v1/bar", 201, 349, 351},
		},
		{
			name: "full_url_strip_query",
			line: `1.2.3.4 - - [10/Oct/2024:13:55:37 +0000] "GET https://api.example.com/v2/x?y=1 HTTP/1.1" 200 321 "-" "ua" request_time=1.000`,
			want: struct {
				method string; path string; status int; msMin, msMax float64
			}{"GET", "/v2/x", 200, 999, 1001},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, ok := Parse(tc.line)
			if !ok {
				t.Fatalf("Parse() failed")
			}
			if got.Method != tc.want.method {
				t.Fatalf("method got %q want %q", got.Method, tc.want.method)
			}
			if got.Path != tc.want.path {
				t.Fatalf("path got %q want %q", got.Path, tc.want.path)
			}
			if got.Status != tc.want.status {
				t.Fatalf("status got %d want %d", got.Status, tc.want.status)
			}
			if !(got.Latency >= tc.want.msMin && got.Latency <= tc.want.msMax) {
				t.Fatalf("latency got %.3fms want in [%.3f, %.3f]", got.Latency, tc.want.msMin, tc.want.msMax)
			}
		})
	}
}

func TestParse_Unsupported(t *testing.T) {
	lines := []string{
		`bad line with no quotes`,
		`"GET /just/method HTTP/1.1" ???`,
	}
	for _, ln := range lines {
		if _, ok := Parse(ln); ok {
			t.Fatalf("expected parse failure for %q", ln)
		}
	}
}