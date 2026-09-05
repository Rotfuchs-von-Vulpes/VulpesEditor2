package context

type Context interface {
	Use()
}

type Manager struct {
	data   map[string]Context
	lastId string
}

func (s *Manager) Add(id string, value Context) {
	_, ok := s.data[id]
	if ok {
		panic("Alreade in use")
	}
	s.data[id] = value
}

func (s *Manager) Check(id string) {
	ctx, ok := s.data[id]
	if ok {
		if s.lastId != id {
			ctx.Use()
			s.lastId = id
		}
	} else {
		panic("Unknow Id")
	}
}

func New() (ctxM *Manager) {
	ctxM = new(Manager)
	ctxM.data = map[string]Context{}
	ctxM.lastId = ""
	return
}
