package draw

// DrawHandler Draw draw interface
type DrawHandler interface {
	//  PlotKLine draw line of the diagram
	PlotKLine(data KlineData)
	//  PlotLine draw kline of the diagram
	PlotLine(name string, data LineData)
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
	draw.kline = []KlineData{}
	draw.line = make(map[string][]LineData)
	return &draw
}
