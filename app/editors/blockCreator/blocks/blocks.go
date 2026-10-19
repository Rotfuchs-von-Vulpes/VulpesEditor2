package blocks

import (
	blocksfile "VulpesEditor/app/editors/blockCreator/blocks/blocksFile"
	"VulpesEditor/app/editors/blockCreator/chunk"
	"VulpesEditor/app/file"
	"VulpesEditor/app/front/renderer"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"uuid"

	"github.com/AllenDang/cimgui-go/backend"
	im "github.com/AllenDang/cimgui-go/imgui"
)

func Init() {
	blocksfile.Init()
}

func AfterCreateContext() {
	blocksfile.AfterCreateContext()
	for i, b := range blocksfile.AllBlocks {
		blocks = append(blocks, newBlockButton(b.Name, uint32(i+1)))
	}
}

type texture struct {
	name string
	data *backend.Texture
}

var (
	blocks []*blockButton

	textures        []*texture
	textureNames    []string
	selectedTexture int32
)

type BlocksContext struct {
	selectedBlock uint32

	blockNameInput   string
	appendedTextures []*texture

	textures [6]*texButton
}

func GetBlockSelected() uint32 {
	return ctx.selectedBlock
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
			img, err := file.LoadImage(f)
			if err != nil {
				fmt.Println(err)
				continue
			}
			tex := new(texture)
			tex.name = t.Name
			tex.data = backend.NewTextureFromRgba(img)
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
	for _, texButton := range ctx.textures {
		texButton.idx = -1
		texButton.ref = nil
	}
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

func getTextures() (err error, textures [6]string) {
	for i, texButton := range ctx.textures {
		textures[i] = ctx.appendedTextures[texButton.idx].name
	}
	return
}

func addBlock() (err error) {
	if ctx.blockNameInput == "" {
		err = fmt.Errorf("Blank name")
		return
	}
	var b blocksfile.BlockJSON
	b.Id = uuid.NewV7().String()
	b.Name = ctx.blockNameInput
	var textures [6]string
	err, textures = getTextures()
	if err != nil {
		return
	}
	for i, tex := range textures {
		idx := slices.Index(b.Textures, tex)
		if idx == -1 {
			b.TexIdx[i] = len(b.Textures)
			b.Textures = append(b.Textures, tex)
		} else {
			b.TexIdx[i] = idx
		}
	}
	b.SaveBlock()
	err, _ = b.ToBlock()
	return
}

type blockButton struct {
	label   string
	blockID uint32
	ref     *im.TextureRef

	viewer *renderer.FrameBuffer
	mesh   *renderer.BlocksMesh
}

func newBlockButton(blockName string, blockId uint32) (b *blockButton) {
	b = new(blockButton)
	b.label = blockName
	b.blockID = blockId
	b.viewer = renderer.CreateFramebuffer(48, 48, true)
	b.mesh = renderer.NewBlocksMesh()
	b.mesh.SetVertices(chunk.Block(b.blockID).ToBuffer())
	b.viewer.RenderBlock(b.mesh)
	b.ref = im.NewTextureRefTextureID(im.TextureID(b.viewer.Image()))
	return
}

func (s *blockButton) show() {
	if im.ImageButtonV(s.label, *s.ref,
		im.NewVec2(48, 48),
		im.NewVec2(0, 1), im.NewVec2(1, 0),
		im.NewVec4(0, 0, 0, 0), im.NewVec4(1, 1, 1, 1)) {
		ctx.selectedBlock = s.blockID
	}
}

const width = 64

func Show(id string) {
	ctxManager.Check(id)

	if im.Begin("Blocks") {
		for i, b := range blocks {
			availableSpace := im.ContentRegionAvail().X
			b.show()
			if i != len(blocks)-1 && availableSpace-2*width > 0 {
				im.SameLine()
			}
		}
	}
	im.End()

	if im.Begin("New Block") {
		im.InputTextWithHint("Block Name", "Block...", &ctx.blockNameInput, im.InputTextFlagsNone, nil)
		if im.Button("Append Texture") {
			openTextureSelection()
		}
		if len(ctx.appendedTextures) > 0 {
			ctx.textures[0].show()
			ctx.textures[2].show()
			im.SameLine()
			ctx.textures[3].show()
			im.SameLine()
			ctx.textures[4].show()
			im.SameLine()
			ctx.textures[5].show()
			ctx.textures[1].show()
		}
		if im.Button("Add") {
			if err := addBlock(); err != nil {
				fmt.Println(err)
			} else {
				resetBlockCreation()
				blocks = nil
				for i, b := range blocksfile.AllBlocks {
					blocks = append(blocks, newBlockButton(b.Name, uint32(i+1)))
				}
			}
		}
		im.SameLine()
		if im.Button("Close") {
			resetBlockCreation()
		}
	}
	im.End()

	ShowSelectTexturePopUp()
}
