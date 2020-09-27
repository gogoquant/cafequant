package draw

import (
	log "gopkg.in/logger.v1"
	"os"
	"sync"

	"github.com/go-echarts/go-echarts/charts"
)

// KLineData ...
type KLineData struct {
	Time string
	Data [4]float32
}

// LineData ...
type LineData struct {
	Time string
	Data float32
}

// LineService ...
type LineService struct {
	BaseService
	mutex      sync.Mutex
	lineChart  *charts.Line
	klineChart *charts.Kline
	kline      []KLineData
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

// Display draw
func (p *LineService) Display() error {
	go func() {
		var file *os.File = nil
		p.lock()
		p.prevKLine()
		p.prevLine()
		p.klineChart.Overlap(p.lineChart)
		p.unLock()

		DrawPath := p.GetPath()
		if _, err := os.Stat(DrawPath); err != nil {
			if !os.IsNotExist(err) {
				log.Error("State diagram fail", err)
				return
			}
		} else {
			_ = os.Remove(DrawPath)
		}
		file, err := os.Create(DrawPath)
		if err != nil {
			log.Error("Create diagram fail", err)
			return
		}
		if err := p.klineChart.Render(file); err != nil {
			log.Error("Render diagram fail", err)
			return
		}
	}()
	return nil
}

// prevKLine ...
func (p *LineService) prevKLine() {
	p.klineChart = charts.NewKLine()
	x := make([]string, 0)
	y := make([][4]float32, 0)
	for i := 0; i < len(p.kline); i++ {
		x = append(x, p.kline[i].Time)
		y = append(y, p.kline[i].Data)
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

// prevLine ...
func (p *LineService) prevLine() {
	p.lineChart = charts.NewLine()
	//p.lineChart.SetGlobalOptions(charts.TitleOpts{Title: "Line多线"}, charts.InitOpts{Theme: "shine"})
	for k, v := range p.line {
		x := make([]string, 0)
		y := make([]float32, 0)
		for i := 0; i < len(v); i++ {
			x = append(x, v[i].Time)
			y = append(y, v[i].Data)
		}
		p.lineChart.AddXAxis(x).AddYAxis(k, y)
	}
	p.lineChart.SetGlobalOptions(
		charts.TooltipOpts{Trigger: "axis"},
	)
}

// PlotKLine Plot kline into pix
func (p *LineService) PlotKLine(time string, a, b, c, d float32) {
	var data KLineData
	data.Time = time
	data.Data[0], data.Data[1], data.Data[2], data.Data[3] = a, b, c, d
	p.lock()
	p.kline = append(p.kline, data)
	p.unLock()
}

// PlotLine Plot line into pix
func (p *LineService) PlotLine(name string, time string, v float32) {
	var data LineData
	data.Time = time
	data.Data = v
	p.lock()
	p.line[name] = append(p.line[name], data)
	p.unLock()
}

// Reset Reset pix
func (p *LineService) Reset() {
	p.lock()
	p.kline = []KLineData{}
	p.line = make(map[string][]LineData)
	p.unLock()
}
