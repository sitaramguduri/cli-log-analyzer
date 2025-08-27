package stats

import (
	"sort"
	"testing"
)

func TestAccumulator_AddAndCounts(t *testing.T) {
	a := NewAccumulator()
	a.Add("/a", 100)
	a.Add("/a", 200)
	a.Add("/b", 300)

	if got := a.Total(); got != 3 {
		t.Fatalf("Total got %d want 3", got)
	}
	if got := a.NumEndpoints(); got != 2 {
		t.Fatalf("NumEndpoints got %d want 2", got)
	}
	al := a.AllLatencies()
	sort.Float64s(al)
	want := []float64{100, 200, 300}
	for i := range want {
		if al[i] != want[i] {
			t.Fatalf("AllLatencies[%d]=%.0f want %.0f", i, al[i], want[i])
		}
	}
	bl := a.LatenciesFor("/a")
	sort.Float64s(bl)
	if len(bl) != 2 || bl[0] != 100 || bl[1] != 200 {
		t.Fatalf("LatenciesFor(/a) got %v", bl)
	}
}

func TestAccumulator_TopEndpoints(t *testing.T) {
	a := NewAccumulator()
	a.Add("/x", 10)
	a.Add("/x", 20)
	a.Add("/y", 30)
	a.Add("/z", 40)
	a.Add("/z", 50)
	top := a.TopEndpoints(2)
	if len(top) != 2 {
		t.Fatalf("len(top)=%d", len(top))
	}
	if top[0].Path != "/x" || top[0].Count != 2 {
		t.Fatalf("top[0]=%+v", top[0])
	}
	if top[1].Path != "/z" || top[1].Count != 2 {
		t.Fatalf("top[1]=%+v", top[1])
	}
}

func TestQuantileSorted(t *testing.T) {
	if got := QuantileSorted(nil, 0.5); got != 0 {
		t.Fatalf("empty got %.0f want 0", got)
	}
	d := []float64{10, 20, 30, 40, 50}
	if got := QuantileSorted(d, 0.5); got != 30 {
		t.Fatalf("p50 got %.0f want 30", got)
	}
	if got := QuantileSorted(d, 0.25); got != 20 {
		t.Fatalf("p25 got %.0f want 20", got)
	}
	if got := QuantileSorted(d, 0.9); int(got) <= 40 {
		t.Fatalf("p90 too small: %.0f", got)
	}
}
