package history

import (
	im "github.com/AllenDang/cimgui-go/imgui"
)

type change interface {
	Undo()
	Redo()
}

var history []change

func Append(c change) {
	history = history[:len(history)-undoLevel]
	history = append(history, c)
	undoLevel = 0
}

var idx = 0
var undoLevel = 0

func undo() {
	changeIdx := len(history) - 1 - undoLevel
	if changeIdx >= 0 {
		change := history[changeIdx]
		undoLevel += 1
		change.Undo()
	}
}

func redo() {
	if undoLevel > 0 {
		changesIdx := len(history) - int(undoLevel)
		if changesIdx >= 0 {
			change := history[changesIdx]
			undoLevel--
			change.Redo()
		}
	}
}

func Loop() {
	io := im.CurrentContext().IO()
	if io.KeyCtrl() && im.IsKeyPressedBoolV(im.KeyZ, true) {
		undo()
	}
	if io.KeyCtrl() && im.IsKeyPressedBoolV(im.KeyY, true) {
		redo()
	}
	buttons := io.MouseClicked()
	if buttons[3] {
		undo()
	}
	if buttons[4] {
		redo()
	}
}
