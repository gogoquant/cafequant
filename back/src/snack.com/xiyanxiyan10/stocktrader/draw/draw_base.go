package draw

// BaseService ...
type BaseService struct {
	// where to store the pic
	path string
}

// setPath ...
func (d *BaseService) SetPath(path string) {
	d.path = path
}

// getPath ...
func (d *BaseService) GetPath() string {
	return d.path
}
