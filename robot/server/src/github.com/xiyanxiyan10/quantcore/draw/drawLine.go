package draw

import (
	"github.com/go-echarts/go-echarts/charts"
	//"github.com/go-openapi/errors"
	"os"
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

// LineService ...
type LineService struct {
	BaseService
	lineChart  *charts.Line
	klineChart *charts.Kline
	kline      []klineData
	line       map[string][]LineData
}

// Draw draw
func (p *LineService) Draw() error {
	var file *os.File = nil
	//if p.lineChart == nil || p.klineChart == nil {
	//		return errors.New(400, "Fail with no lineChat or klineChart")
	//}
	p.prevDrawKline()
	p.prevDrawLine()
	//p.klineChart.Overlap(p.lineChart)
	DrawPath := p.GetPath()
	if _, err := os.Stat(DrawPath); err != nil {
		if os.IsNotExist(err) {
			if os.IsNotExist(err) {
				file, err = os.Create(DrawPath)
				if err != nil {
					return err
				}
			}
		}
	} else {
		if file, err = os.Open(p.GetPath()); err != nil {
			return err
		}
	}
	if err := p.klineChart.Render(file); err != nil {
		return err
	}
	return nil
}

// prevDrawKline ...
func (p *LineService) prevDrawKline() {
	p.klineChart = charts.NewKLine()
	x := make([]string, 0)
	y := make([][4]float32, 0)
	for i := 0; i < len(p.kline); i++ {
		x = append(x, p.kline[i].date)
		y = append(y, p.kline[i].data)
	}
	p.klineChart.AddXAxis(x).AddYAxis("kline", y)
	p.klineChart.SetGlobalOptions(
		charts.TitleOpts{Title: "KLine"},
		charts.XAxisOpts{SplitNumber: 20},
		charts.YAxisOpts{Scale: true},
		charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
	)
}

// prevDrawLine ...
func (p *LineService) prevDrawLine() {
	p.lineChart = charts.NewLine()
	x := make([]string, 0)
	y := make([]float32, 0)
	//p.lineChart.SetGlobalOptions(charts.TitleOpts{Title: "Line多线"}, charts.InitOpts{Theme: "shine"})
	for k, v := range p.line {
		for i := 0; i < len(v); i++ {
			x = append(x, v[i].date)
			y = append(y, v[i].data)
		}
		p.lineChart.AddXAxis(x).AddYAxis(k, y)
	}
}

// PlotKLine Plot kline into pix
func (p *LineService) PlotKLine(data klineData) {
	p.kline = append(p.kline, data)
}

// PlotLine Plot line into pix
func (p *LineService) PlotLine(name string, data LineData) {
	p.line[name] = append(p.line[name], data)
}

// Reset Reset pix
func (p *LineService) Reset() {
	p.kline = []klineData{}
	p.line = make(map[string][]LineData)
	if len(p.GetPath()) > 0 {
		_ = os.Remove(p.GetPath())
	}
}
