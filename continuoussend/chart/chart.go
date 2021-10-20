package chart

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	outDir = "out"

	latencies = make([]int, 0)

	classMax int
	bar      = charts.NewBar()
)

func init() {
	flag.IntVar(&classMax, "max", 100, "max of latency class")
	flag.Parse()
}

func Record(latency int) {
	latencies = append(latencies, latency)
}

func Render() error {
	filename := time.Now().Format("2006-01-02_15:04:05") + ".html"
	f, err := makeFile(outDir, filename)
	if err != nil {
		return err
	}
	defer f.Close()

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Latency Histogram"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "latency [ms]"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "frequency"}),
	)

	bar.SetXAxis(makeRange(classMax)).
		AddSeries("HTTP/1.1", makeBars(latencies, classMax))

	if err = bar.Render(f); err != nil {
		return err
	}

	fmt.Printf("rendered to %s\n", f.Name())

	return nil
}

func makeRange(max int) []int {
	r := make([]int, max)
	for i := 1; i < max+1; i++ {
		r[i-1] = i
	}
	return r
}

func makeFile(dir, file string) (*os.File, error) {
	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		return nil, err
	}

	f, err := os.Create(dir + "/" + file)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func makeBars(data []int, max int) []opts.BarData {
	items := make([]int, max)

	for _, d := range data {
		prev := items[d]
		items[d] = prev + 1
	}

	bars := make([]opts.BarData, len(items))
	for i := 0; i < len(items); i++ {
		bars[i] = opts.BarData{Value: items[i]}
	}

	return bars
}
