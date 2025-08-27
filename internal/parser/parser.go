package parser

import (
	"log"
	"net/url"
	"regexp"
	"strings"
)

type Entry struct{
	Method string
	Path string //URL path without query
	Status int 
	Latency float64 //milli seconds
}

var (
	reReqTime = regexp.MustCompile(`"(?P<m>GET|POST|PUT|DELETE|PATCH|HEAD) (?P<u>[^" ]+) [^"]+" (?P<st>\d{3}).*?request_time=(?P<sec>\d+(?:\.\d+)?)`)
	reRTms    = regexp.MustCompile(`"(?P<m>GET|POST|PUT|DELETE|PATCH|HEAD) (?P<u>[^" ]+) [^"]+" (?P<st>\d{3}).*?\brt=(?P<ms>\d+)ms\b`)
)

func Parse(line string) (Entry, bool){
	if m:= reReqTime.FindStringSubmatch(line); m!= nil{
		method := m[1]
		rawURL := m[2]
		status := atoi(m[3])
		sec := atof(m[4])
		log.Printf("method %s", method)
		return Entry{Method: method, Path: stripQuery(rawURL), Status: status, Latency: sec*1000}, true
	}
	if m := reRTms.FindStringSubmatch(line); m != nil{
		method := m[1]
		rawURL := m[2]
		status := atoi(m[3])
		ms := atof(m[4])
		return Entry{Method: method, Path: stripQuery(rawURL), Status: status, Latency: ms}, true
	}
	return Entry{}, false
}
func stripQuery(u string) string{
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://"){
	if parsed, err:= url.Parse(u); err ==nil{
		return parsed.Path
	}
	}
	if i:= strings.IndexByte(u, '?'); i>=0{
		return u[:i]
	}
	return u
}

func atoi(s string) int {
	n := 0
	for i:= 0; i<len(s); i++{
		n = n*10 + int(s[i]-'0')
	}
	return n
}

func atof(s string) float64{
	var n float64
	var frac float64 = 1
	dot := false
	for i:= 0; i<len(s); i++{
		if s[i] == '.'{
			dot = true
			continue
		}
		if !dot {
			n = n*10 + float64(s[i] - '0')
		}else{
			frac *= 0.1
			n += float64(s[i] - '0')*frac
		}
	}
	return n
}