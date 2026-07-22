package context

type Context interface {
	Use()
}

type Manager struct {
	data   map[int32]Context
	lastId int32
}

func (s *Manager) Add(id int32, value Context) {
	_, ok := s.data[id]
	if ok {
		panic("Alreade in use")
	}
	s.data[id] = value
}

func (s *Manager) Check(id int32) {
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
	ctxM.data = map[int32]Context{}
	ctxM.lastId = -1
	return
}
