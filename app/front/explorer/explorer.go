package explorer

import (
	"VulpesEditor/app/file"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AllenDang/cimgui-go/backend"
	im "github.com/AllenDang/cimgui-go/imgui"
)

type assetIcon struct {
	name   string
	ref    *im.TextureRef
	width  float32
	height float32
}

func (a assetIcon) show(width float32) {
	im.Image(*a.ref, im.NewVec2(width*a.width, width*a.height))
}

var allAssets []assetIcon

func getAll() {
	assets := file.GetAllAssets("textures")
	for i, asset := range assets {
		ext := filepath.Ext(asset.Name)
		if ext != ".png" {
			continue
		}
		f, err := os.Open(asset.Path)
		if err != nil {
			fmt.Println(asset.Name, asset.Path, err)
			continue
		}
		defer f.Close()
		img, err := file.LoadPixelImage(f, 10)
		if err != nil {
			fmt.Println(asset.Name, err)
			continue
		}
		ref := backend.NewTextureFromRgba(img).ID
		allAssets = append(allAssets, assetIcon{"Fox" + strconv.FormatInt(int64(i), 10), &ref, 1, float32(img.Rect.Dy()) / float32(img.Rect.Dx())})
	}
}

func AfterCreateContext() {
	getAll()
}

var table [][]assetIcon
var tableWidth float32
var tableSize []float32

func Resize(width int32) {
	table = nil
	tableSize = nil
	count := width / 200
	tableWidth = float32(width) / float32(count)
	for range count {
		table = append(table, []assetIcon{})
		tableSize = append(tableSize, 0)
	}
	var minIdx int
	for _, asset := range allAssets {
		min := tableSize[minIdx]
		for idx, size := range tableSize {
			if min > size {
				minIdx = idx
				min = size
			}
		}
		table[minIdx] = append(table[minIdx], asset)
		tableSize[minIdx] += asset.height
	}
}

var viwerSize [2]float32

func Show() {
	if im.Begin("Artworks") {
		wSize := im.ContentRegionAvail()

		if viwerSize[0] != wSize.X || viwerSize[1] != wSize.Y {
			viwerSize[0] = wSize.X
			viwerSize[1] = wSize.Y

			Resize(int32(wSize.X) - 50)
		}

		for i, collum := range table {
			if im.BeginChildStrV("collum "+strconv.FormatInt(int64(i), 10), im.NewVec2(tableWidth, 0), im.ChildFlagsAutoResizeY|im.ChildFlagsAlwaysAutoResize, im.WindowFlagsNone) {
				for _, asset := range collum {
					asset.show(tableWidth)
				}
			}
			im.EndChild()
			im.SameLine()
		}
	}
	im.End()
}
