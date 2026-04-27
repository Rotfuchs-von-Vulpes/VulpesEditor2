package util

type IdSystem struct {
	lastId uint32
}

func NewIdSystem() IdSystem {
	return IdSystem{0}
}

func (s *IdSystem) GetID() (id uint32) {
	id = s.lastId
	s.lastId += 1
	return
}
