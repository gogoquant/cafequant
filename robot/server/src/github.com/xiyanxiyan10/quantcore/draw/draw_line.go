package draw

import (
	"github.com/go-echarts/go-echarts/charts"
	log "gopkg.in/logger.v1"

	//"github.com/go-openapi/errors"
	"os"
	"sync"
)

// klineData ...
type KlineData struct {
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
	mutex      sync.Mutex
	lineChart  *charts.Line
	klineChart *charts.Kline
	kline      []KlineData
	line       map[string][]LineData
}

// lock draw
func (p *LineService) lock() {
	p.mutex.Lock()
}

// unLock draw
func (p *LineService) unLock() {
	p.mutex.Unlock()
}

// Draw draw
func (p *LineService) Draw() error {
	go func() {
		var file *os.File = nil
		p.lock()
		p.prevDrawKline()
		p.prevDrawLine()
		p.klineChart.Overlap(p.lineChart)
		p.unLock()

		DrawPath := p.GetPath()
		if _, err := os.Stat(DrawPath); err != nil {
			if !os.IsNotExist(err) {
				log.Error("State pic fail ", err)
				return
			}
		} else {
			os.Remove(DrawPath)
		}
		file, err := os.Create(DrawPath)
		if err != nil {
			log.Error("Create pic fail ", err)
			return
		}
		if err := p.klineChart.Render(file); err != nil {
			log.Error("Render pic fail", err)
			return
		}
	}()
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
		charts.TooltipOpts{Trigger: "axis"},
		charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
	)
}

// prevDrawLine ...
func (p *LineService) prevDrawLine() {
	p.lineChart = charts.NewLine()
	//p.lineChart.SetGlobalOptions(charts.TitleOpts{Title: "Line多线"}, charts.InitOpts{Theme: "shine"})
	for k, v := range p.line {
		x := make([]string, 0)
		y := make([]float32, 0)
		for i := 0; i < len(v); i++ {
			x = append(x, v[i].date)
			y = append(y, v[i].data)
		}
		p.lineChart.AddXAxis(x).AddYAxis(k, y)
	}
	p.lineChart.SetGlobalOptions(
		charts.TooltipOpts{Trigger: "axis"},
	)
}

// PlotKLine Plot kline into pix
func (p *LineService) PlotKLine(data KlineData) {
	p.lock()
	p.kline = append(p.kline, data)
	p.unLock()
}

// PlotLine Plot line into pix
func (p *LineService) PlotLine(name string, data LineData) {
	p.lock()
	p.line[name] = append(p.line[name], data)
	p.unLock()
}

// Reset Reset pix
func (p *LineService) Reset() {
	p.lock()
	p.kline = []KlineData{}
	p.line = make(map[string][]LineData)
	p.unLock()
}
