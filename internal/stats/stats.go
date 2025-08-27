package stats

import "sort"

type Accumulator struct {
	all []float64
	byPath map[string][]float64
	pathCount map[string]int
}

func NewAccumulator()*Accumulator{
	return &Accumulator{
		byPath: make(map[string][]float64),
		pathCount: make(map[string]int),
	}
}

func (a *Accumulator) Add(path string, latencyMs float64){
	a.all = append(a.all, latencyMs)
	a.byPath[path] = append(a.byPath[path], latencyMs)
	a.pathCount[path]++
}
func (a *Accumulator) Total() int{
	return len(a.all)
}

func (a *Accumulator) NumEndpoints() int{
	return len(a.byPath)
}

func (a *Accumulator) AllLatencies() []float64{
	out:= make([]float64,len(a.all))
	copy(out, a.all)
	return out
}

func (a *Accumulator) LatenciesFor(path string) []float64{
	src:=a.byPath[path]
	out := make([]float64, len(src))
	copy(out, src)
	return out
}
type Top struct{
	Path string
	Count int
}

func (a *Accumulator) TopEndpoints(n int)[]Top{
	tmp := make([]Top,0,len(a.pathCount))
	for p,c := range a.pathCount{
		tmp = append(tmp,Top{Path:p, Count:c})
	}
	sort.Slice(tmp, func(i,j int)bool{
		if tmp[i].Count == tmp[j].Count{
			return tmp[i].Path < tmp[j].Path
		}
		return tmp[i].Count > tmp[j].Count
	})

	if n>len(tmp){
		n = len(tmp)
	}
	return tmp[:n]
}

func QuantileSorted(sorted []float64, q float64) float64{
	n := len(sorted)
	if n==0{
		return 0
	}
	pos := q*float64(n-1)
	i := int(pos)
	if i == n - 1{
		return sorted[i]
	}
	frac := pos - float64(i)
	return sorted[i]*(1-frac) + sorted[i+1]*frac
}