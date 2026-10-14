package tabs

import "slices"

type Tab interface {
	ID() string
	Name() string
	Type() string
	Show()
	Focus() bool
	Save()
	Export()
	Saved() bool
	Preview(string)

	ParentID() string
	AddPreview(preview func(string))
}

var AllTabs []Tab

func Push(t Tab) {
	AllTabs = append(AllTabs, t)
}

func Close(idx int) {
	AllTabs = slices.Delete(AllTabs, idx, idx+1)
}

func appendParent(child, parent Tab) {
	child.AddPreview(parent.Preview)
}

func Loop() {
	for _, t := range AllTabs {
		if id := t.ParentID(); id != "" {
			idx := slices.IndexFunc(AllTabs, func(tab Tab) bool {
				if id == tab.ID() {
					return true
				}
				return false
			})
			if idx >= 0 {
				appendParent(t, AllTabs[idx])
			}
		}
	}
}
