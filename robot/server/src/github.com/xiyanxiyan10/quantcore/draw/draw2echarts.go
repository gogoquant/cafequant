package draw

import (
	"github.com/go-echarts/go-echarts/charts"
	"github.com/go-openapi/errors"
	"os"
	"path"
)

// klineData ...
type klineData struct {
	date string
	data [4]float32
}

// LineData ...
type LineData struct {
	date string
	data float32
}

// KlineService ...
type KlineService struct {
	lineChart  *charts.Line
	klineChart *charts.Kline
	kline      []klineData
	line       map[string][]LineData
}

// Draw draw
func (p *KlineService) Draw() error {
	if p.lineChart == nil || p.klineChart == nil {
		return errors.New(1, "line or kline charts are not init")
	}
	p.prevDrawKline()
	p.prevDrawLine()
	p.klineChart.Overlap(p.lineChart)
	page := charts.NewPage(orderRouters("kline")...)
	page.Add(
		p.klineChart,
	)
	f, err := os.Create(getRenderPath("kline.html"))
	if err != nil {
		return err
	}
	page.Render(w, f)
	return nil
}

// prevDrawKline ...
func (p *KlineService) prevDrawKline() {
	p.klineChart = charts.NewKLine()
	x := make([]string, 0)
	y := make([][4]float32, 0)
	for i := 0; i < len(p.kline); i++ {
		x = append(x, p.kline[i].date)
		y = append(y, p.kline[i].data)
	}
	p.klineChart.AddXAxis(x).AddYAxis("kline", y)
	p.klineChart.SetGlobalOptions(
		charts.TitleOpts{Title: "Kline"},
		charts.XAxisOpts{SplitNumber: 20},
		charts.YAxisOpts{Scale: true},
		charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
	)
}

// prevDrawLine ...
func (p *KlineService) prevDrawLine() {
	line := charts.NewLine()
	x := make([]string, 0)
	y := make([]float32, 0)
	line.SetGlobalOptions(charts.TitleOpts{Title: "Line多线"}, charts.InitOpts{Theme: "shine"})
	for k, v := range p.line {
		for i := 0; i < len(v); i++ {
			x = append(x, v[i].date)
			y = append(y, v[i].data)
		}
		line.AddXAxis(x).AddYAxis(k, y)
	}
}

// PlotKLine Plot kline into pix
func (p *KlineService) PlotKLine(data klineData) {
	p.kline = append(p.kline, data)
}

// PlotLine Plot line into pix
func (p *KlineService) PlotLine(name string, data LineData) {
	p.line[name] = append(p.line[name], data)
}

// Reset Reset pix
func (p *KlineService) Reset() {
	p.kline = []klineData{}
	p.line = make(map[string][]LineData)
}

type router struct {
	name string
	charts.RouterOpts
}

var (
	routers = []router{
		{"kline", charts.RouterOpts{URL: "127.0.0.1" + "/kline", Text: "Kline-K 线图"}},
	}
)

func orderRouters(chartType string) []charts.RouterOpts {
	for i := 0; i < len(routers); i++ {
		if routers[i].name == chartType {
			routers[i], routers[0] = routers[0], routers[i]
			break
		}
	}

	rs := make([]charts.RouterOpts, 0)
	for i := 0; i < len(routers); i++ {
		rs = append(rs, routers[i].RouterOpts)
	}
	return rs
}

func getRenderPath(f string) string {
	return path.Join("html", f)
}
