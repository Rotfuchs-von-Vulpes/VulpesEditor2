package context

type Context interface {
	Use()
	Reset()
}

type Manager struct {
	data   map[int32]Context
	inside bool
	ctx    Context
}

func (s *Manager) Add(key int32, value Context) {
	_, ok := s.data[key]
	if ok {
		panic("Alreade in use")
	}
	s.data[key] = value
}

func (s *Manager) Begin(id int32) {
	if s.inside {
		panic("Too much begin call")
	}
	s.inside = true
	var ok bool
	s.ctx, ok = s.data[id]
	if !ok {
		panic("Unknow Id")
	}
	s.ctx.Use()
}

func (s *Manager) End() {
	if !s.inside {
		panic("Too much end call")
	}
	s.ctx.Reset()
	s.inside = false
}

func New() (ctxM *Manager) {
	ctxM = new(Manager)
	ctxM.data = map[int32]Context{}
	return
}
