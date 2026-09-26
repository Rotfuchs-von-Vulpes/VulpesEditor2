package blocks

import (
	"VulpesEditor/app/file"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"slices"

	"github.com/AllenDang/cimgui-go/backend"
	im "github.com/AllenDang/cimgui-go/imgui"
)

type texture struct {
	name string
	data *backend.Texture
}

var (
	textures        []*texture
	textureNames    []string
	selectedTexture int32
)

type BlocksContext struct {
	blockNameInput   string
	appendedTextures []*texture

	top, bottom, left, right, front, back *texButton
}

func getAllTextures() {
	files := file.GetAllAssets("textures")
	for _, t := range files {
		if !slices.Contains(textureNames, t.Name) {
			if e := filepath.Ext(t.Name); e != ".png" {
				continue
			}
			f, err := os.Open(t.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer f.Close()
			img, err := png.Decode(f)
			if err != nil {
				fmt.Println(err)
				continue
			}
			tex := new(texture)
			tex.name = t.Name
			tex.data = backend.NewTextureFromRgba(backend.ImageToRgba(img))
			textures = append(textures, tex)
			textureNames = append(textureNames, t.Name)
		}
	}
}

var isTextureSelectionOpen bool

func openTextureSelection() {
	getAllTextures()
	if len(textures) == 0 {
		return
	}
	isTextureSelectionOpen = true
}

func closeTextureSelection() {
	isTextureSelectionOpen = false
	selectedTexture = 0
	im.CloseCurrentPopup()
}

func ShowSelectTexturePopUp() {
	if isTextureSelectionOpen {
		if !im.IsPopupOpenStr("Texture Selection") {
			im.OpenPopupStr("Texture Selection")
		}

		if im.BeginPopupModal("Texture Selection") {
			im.ComboStrarr("Texture", &selectedTexture, textureNames, int32(len(textureNames)))
			if im.Button("Append") {
				tex := textures[selectedTexture]
				if !slices.Contains(ctx.appendedTextures, tex) {
					ctx.appendedTextures = append(ctx.appendedTextures, tex)
					closeTextureSelection()
				}
			}
			im.SameLine()
			if im.Button("Close") {
				closeTextureSelection()
			}
			im.EndPopup()
		}
	}
}

func resetBlockCreation() {
	ctx.appendedTextures = nil
	ctx.top.ref = nil
	ctx.bottom.ref = nil
	ctx.front.ref = nil
	ctx.back.ref = nil
	ctx.left.ref = nil
	ctx.right.ref = nil
	ctx.blockNameInput = ""
}

type texButton struct {
	label string
	ref   *im.TextureRef
	idx   int
}

func (s *texButton) show() {
	if s.ref == nil {
		s.idx = 0
		s.ref = &ctx.appendedTextures[0].data.ID
	}
	if im.ImageButton(s.label, *s.ref, im.NewVec2(64, 64)) {
		s.idx++
		if s.idx >= len(ctx.appendedTextures) {
			s.idx = 0
		}
		s.ref = &ctx.appendedTextures[s.idx].data.ID
	}
}

func Show(id string) {
	ctxManager.Check(id)

	if im.Begin("Blocks") {

	}
	im.End()

	if im.Begin("New Block") {
		im.InputTextWithHint("Block Name", "Block...", &ctx.blockNameInput, im.InputTextFlagsNone, nil)
		if im.Button("Append Texture") {
			openTextureSelection()
		}
		if len(ctx.appendedTextures) > 0 {
			ctx.top.show()
			ctx.left.show()
			im.SameLine()
			ctx.front.show()
			im.SameLine()
			ctx.right.show()
			im.SameLine()
			ctx.back.show()
			ctx.bottom.show()
		}
		if im.Button("Add") {
			resetBlockCreation()
		}
		im.SameLine()
		if im.Button("Close") {
			resetBlockCreation()
		}
	}
	im.End()

	ShowSelectTexturePopUp()
}
