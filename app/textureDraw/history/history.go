package history

import (
	im "github.com/AllenDang/cimgui-go/imgui"
)

type change interface {
	Undo()
	Redo()
}

type data struct {
	history   []change
	idx       int
	undoLevel int
}

func Append(c change) {
	currentCtx.history = currentCtx.history[:len(currentCtx.history)-currentCtx.undoLevel]
	currentCtx.history = append(currentCtx.history, c)
	currentCtx.undoLevel = 0
}

func undo() {
	changeIdx := len(currentCtx.history) - 1 - currentCtx.undoLevel
	if changeIdx >= 0 {
		change := currentCtx.history[changeIdx]
		currentCtx.undoLevel += 1
		change.Undo()
	}
}

func redo() {
	if currentCtx.undoLevel > 0 {
		changesIdx := len(currentCtx.history) - int(currentCtx.undoLevel)
		if changesIdx >= 0 {
			change := currentCtx.history[changesIdx]
			currentCtx.undoLevel--
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
