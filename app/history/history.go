package history

import (
	im "github.com/AllenDang/cimgui-go/imgui"
)

type change interface {
	Undo()
	Redo()
}

type historyData struct {
	history   []change
	idx       int
	undoLevel int
}

func Append(c change) {
	ctx.history = ctx.history[:len(ctx.history)-ctx.undoLevel]
	ctx.history = append(ctx.history, c)
	ctx.undoLevel = 0
}

func undo() {
	changeIdx := len(ctx.history) - 1 - ctx.undoLevel
	if changeIdx >= 0 {
		change := ctx.history[changeIdx]
		ctx.undoLevel += 1
		change.Undo()
	}
}

func redo() {
	if ctx.undoLevel > 0 {
		changesIdx := len(ctx.history) - int(ctx.undoLevel)
		if changesIdx >= 0 {
			change := ctx.history[changesIdx]
			ctx.undoLevel--
			change.Redo()
		}
	}
}

func Loop(id int32) {
	ctxManager.Check(id)

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
