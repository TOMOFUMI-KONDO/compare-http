package chart

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	outDir = "out"

	latencies = make([]int, 0)

	classWidth = 10
	classRange = 10
	bar        = charts.NewBar()
)

func Record(latency int) {
	latencies = append(latencies, latency)
}

func Render() error {
	now := time.Now().Format("2006-01-02_15:04:05")
	f, err := makeFile(now + ".html")
	if err != nil {
		return err
	}
	defer f.Close()

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Latency Histogram"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "latency [ms]"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "frequency"}),
	)

	bar.SetXAxis(makeRange(classRange)).
		AddSeries("HTTP/1.1", makeBars(latencies))

	if err = bar.Render(f); err != nil {
		return err
	}

	fmt.Printf("rendered to %s\n", f.Name())

	return nil
}

func makeRange(n int) []int {
	r := make([]int, n)
	for i := 0; i < n; i++ {
		r[i] = (i + 1) * 10
	}
	return r
}

func makeFile(file string) (*os.File, error) {
	if err := os.MkdirAll(outDir, os.FileMode(0755)); err != nil {
		return nil, err
	}

	f, err := os.Create(outDir + "/" + file)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func makeBars(data []int) []opts.BarData {
	items := make([]int, classRange)

	for i := 0; i < len(data); i++ {
		// round nearest 10th
		class := int(math.Round(float64(data[i]/10)) * 10)

		prev := items[class/classWidth]
		items[class/classWidth] = prev + 1
	}

	bars := make([]opts.BarData, len(items))
	for i := 0; i < len(items); i++ {
		bars[i] = opts.BarData{Value: items[i]}
	}

	return bars
}
