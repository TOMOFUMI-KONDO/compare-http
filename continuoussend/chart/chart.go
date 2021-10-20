package chart

import (
	"fmt"
	"time"

	"gonum.org/v1/plot/vg"

	"gonum.org/v1/plot/plotter"

	"gonum.org/v1/plot"
)

var (
	outDir    = "out"
	latencies = make([]int, 0)
)

func Record(latency int) {
	latencies = append(latencies, latency)
}

func Render() error {
	p := plot.New()
	p.Title.Text = "Latency Histogram"

	values := make(plotter.Values, len(latencies))
	for i, v := range latencies {
		values[i] = float64(v)
	}

	h, err := plotter.NewHist(values, 100)
	if err != nil {
		return err
	}
	p.Add(h)

	filename := time.Now().Format("2006-01-02_15:04:05") + ".png"
	file := outDir + "/" + filename

	if err = p.Save(4*vg.Inch, 4*vg.Inch, file); err != nil {
		return err
	}

	fmt.Printf("rendered to %s\n", file)

	return nil
}
