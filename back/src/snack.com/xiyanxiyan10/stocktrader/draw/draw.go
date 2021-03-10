package draw

// DrawHandler Draw draw interface
type DrawHandler interface {
	//  PlotKLine draw line of the diagram
	PlotKLine(time string, open, closed, low, high float32)
	//  PlotLine draw kline of the diagram
	PlotLine(name string, time string, v float32, shape string)
	//  reset diagram
	Reset()
	//  set path store diagram
	SetPath(path string)
	// get path store diagram
	GetPath() string
	// draw pic
	Display() error
}

// NewDrawHandler ...
func NewDrawHandler() DrawHandler {
	var draw LineService
	draw.kline = []KLineData{}
	draw.line = make(map[string][]LineData)
	return &draw
}
