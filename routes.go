package generator

type Routes struct {
	Fn     func(ctx Context)
	Method string
	Path   string
}

// Register add custom route
func (s *Server) Register(route interface{}) {
	s.Routes = append(s.Routes, route.(Routes))
}
