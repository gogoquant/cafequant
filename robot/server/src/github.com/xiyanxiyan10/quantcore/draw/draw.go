package draw

// Draw draw interface
type LineDrawer interface {
	//  PlotKLine draw line of the pic
	PlotKLine(data klineData)
	//  PlotLine draw kline of the pic
	PlotLine(name string, data LineData)
	//  reset pic
	Reset()
	//  set path store pic
	SetPath(path string)
	// get path store pc
	GetPath() string
	// draw pic
	Draw() error
}

// GetLineDrawer ...
func GetLineDrawer(path string) LineDrawer {
	var draw LineService
	draw.kline = []klineData{}
	draw.line = make(map[string][]LineData)
	draw.SetPath(path)
	return &draw
}
