package tabs

import "slices"

type Tab interface {
	Name() string
	Show()
	Focus() bool
	Save()
}

var AllTabs []Tab

func Push(t Tab) {
	AllTabs = append(AllTabs, t)
}

func Close(idx int) {
	AllTabs = slices.Delete(AllTabs, idx, idx+1)
}
