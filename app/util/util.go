package util

type IdSystem struct {
	lastId int32
}

func NewIdSystem() (sys *IdSystem) {
	sys = new(IdSystem)
	sys.lastId = 0
	return
}

func (s *IdSystem) GetID() (id int32) {
	id = s.lastId
	s.lastId += 1
	return
}
