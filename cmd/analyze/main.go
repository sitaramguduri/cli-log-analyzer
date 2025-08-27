package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/sitaram/cli-log-analyzer/internal/parser"
	"github.com/sitaram/cli-log-analyzer/internal/stats"
)

func main() {
	file := flag.String("file", "", "path to log file(or - for stdin)")
	topN := flag.Int("top", 5, "show top-N endpoints by count")
	quantiles := flag.String("q", "0.5,0.95", "comma-seperated quantles(e.g., 0.5,0.95)")
	flag.Parse()
	log.Printf("file: %d", *file)
	if *file == "" {
		log.Fatal("missing -file")
	}

	var in *os.File
	var err error
	if *file == "_" {
		in = os.Stdin
	} else {
		in, err = os.Open(*file)
		if err != nil {
			log.Fatalf("open %s: %v", *file, err)
		}
		defer in.Close()
	}

	qVals := []float64{0.5, 0.95}
	if *quantiles != "" {
		qVals = nil
		for _, s := range strings.Split(*quantiles, ",") {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			var f float64
			_, err := fmt.Sscanf(s, "%f", &f)
			if err != nil || f <= 0 || f >= 1 {
				log.Fatalf("bad quaantile %q", s)
			}
			qVals = append(qVals, f)
		}
	}
	sc := bufio.NewScanner(in)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 2<<20)

	acc := stats.NewAccumulator()
	for sc.Scan() {
		line := sc.Text()
		ent, ok := parser.Parse(line)
		if !ok {
			continue // skip unparseable lines
		}
		acc.Add(ent.Path, ent.Latency)
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan: %v", err)
	}
	total := acc.Total()
	fmt.Printf("lines=%d endpoints=%d\n", total, acc.NumEndpoints())

	// Global quantiles
	all := acc.AllLatencies()
	sort.Float64s(all)
	for _, q := range qVals {
		fmt.Printf("p%02.0f_ms=%.0f\n", q*100, stats.QuantileSorted(all, q))
	}

	top := acc.TopEndpoints(*topN)
	fmt.Println("top_endpoints:")
	for i, t := range top {
		qs := make([]string, 0, len(qVals))
		lats := acc.LatenciesFor(t.Path)
		sort.Float64s(lats)
		for _, q := range qVals {
			qs = append(qs, fmt.Sprintf("p%02.0f=%.0fms", q*100, stats.QuantileSorted(lats, q)))
		}
		fmt.Printf("%d. %s count=%d %s\n", i+1, t.Path, t.Count, strings.Join(qs, " "))
	}

	_ = time.Now() // keep time import for future extension

}
