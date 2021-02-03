package draw

import (
	"github.com/go-echarts/go-echarts/charts"
	log "gopkg.in/logger.v1"
	"os"
	"sync"
)

// KLineData ...
type KLineData struct {
	Time string
	Data [4]float32
}

type ScatterData struct {
	Time string
	Data float32
	// 图的形状
	Shape string
}

// LineData ...
type LineData struct {
	Time string
	Data float32
	// 图的形状
	Shape string
}

// LineService ...
type LineService struct {
	BaseService
	mutex        sync.Mutex
	lineChart    *charts.Line
	klineChart   *charts.Kline
	scatterChart *charts.Scatter
	kline        []KLineData
	line         map[string][]LineData
	scatter      map[string][]ScatterData
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
	var file *os.File = nil
	p.lock()
	p.prevKLine()
	p.prevLine()
	//p.prevScatter()
	if len(p.kline) > 0 {
		p.klineChart.Overlap(p.lineChart)
	}

	p.unLock()

	DrawPath := p.GetPath()
	if _, err := os.Stat(DrawPath); err != nil {
		if !os.IsNotExist(err) {
			log.Error("State diagram fail", err)
			return err
		}
	} else {
		_ = os.Remove(DrawPath)
	}
	file, err := os.Create(DrawPath)
	if err != nil {
		log.Error("Create diagram fail", err)
		return err
	}
	if len(p.kline) > 0 {
		if err := p.klineChart.Render(file); err != nil {
			log.Error("Render diagram fail", err)
			return err
		}
	} else {
		if err := p.lineChart.Render(file); err != nil {
			log.Error("Render diagram fail", err)
			return err
		}

	}
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
		//charts.TitleOpts{Title: "line"},
		charts.XAxisOpts{SplitNumber: 20},
		charts.YAxisOpts{Scale: true},
		charts.TooltipOpts{Trigger: "axis"},
		charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
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

// PlotScatter Plot kline into pix
func (p *LineService) PlotScatter(name string, time string, a float32, shape string) {
	var data ScatterData
	data.Time = time
	data.Data = a
	p.lock()
	p.scatter[name] = append(p.scatter[name], data)
	p.unLock()
}

// PlotLine Plot line into pix
func (p *LineService) PlotLine(name string, time string, v float32, shape string) {
	var data LineData
	data.Time = time
	data.Data = v
	data.Shape = shape
	p.lock()
	p.line[name] = append(p.line[name], data)
	p.unLock()
}

// Reset Reset pix
func (p *LineService) Reset() {
	p.lock()
	p.kline = []KLineData{}
	p.line = make(map[string][]LineData)
	p.scatter = make(map[string][]ScatterData)
	p.unLock()
}

// prevScatter ...
func (p *LineService) prevScatter() {
	p.scatterChart = charts.NewScatter()
	for k, v := range p.scatter {
		x := make([]string, 0)
		y := make([]float32, 0)
		for i := 0; i < len(v); i++ {
			x = append(x, v[i].Time)
			y = append(y, v[i].Data)
		}
		p.scatterChart.AddXAxis(x).AddYAxis(k, y)
	}
}

// prevLine ...
func (p *LineService) prevLine() {
	p.lineChart = charts.NewLine()
	//p.lineChart.SetGlobalOptions(charts.TitleOpts{Title: "Line多线"}, charts.InitOpts{Theme: "shine"})
	var shape string
	for k, v := range p.line {
		x := make([]string, 0)
		y := make([]float32, 0)
		for i := 0; i < len(v); i++ {
			x = append(x, v[i].Time)
			y = append(y, v[i].Data)
			shape = v[i].Shape
		}
		if shape == StepLine {
			p.lineChart.AddXAxis(x).AddYAxis(k, y, charts.LineOpts{Step: true}, charts.LineOpts{ConnectNulls: false})
		} else if shape == SmoothLine {
			p.lineChart.AddXAxis(x).AddYAxis(k, y, charts.LineOpts{Smooth: true}, charts.LineOpts{ConnectNulls: false})
		} else {
			p.lineChart.AddXAxis(x).AddYAxis(k, y, charts.LineOpts{ConnectNulls: false})
		}
	}

	p.lineChart.SetGlobalOptions(
		charts.TooltipOpts{Trigger: "axis"},
	)
}
