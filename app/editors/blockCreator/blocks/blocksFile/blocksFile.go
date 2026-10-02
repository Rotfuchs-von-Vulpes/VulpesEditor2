package blocksfile

import (
	"VulpesEditor/app/editors/textureDraw/canvas/texture"
	"VulpesEditor/app/file"
	"VulpesEditor/app/util"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var blocksDir string
var AllBlocks []Block

func Init() {
	blocksDir = filepath.Join(util.AppDir, "resources", "blocks")
	if err := os.MkdirAll(blocksDir, os.ModePerm); err != nil {
		panic(err)
	}
	AllBlocks = GetAllBlocks()
}

type Block struct {
	Name          string
	Id            string
	Textures      [6]*texture.Texture
	TexturesNames []string
}

type BlockJSON struct {
	Name     string   `json:"name"`
	Id       string   `json:"id"`
	TexIdx   [6]int   `json:"textureIndex"`
	Textures []string `json:"texture"`
}

var allTextureFiles []file.Project

func getAllTextures() {
	allTextureFiles = nil
	files := file.GetAllAssets("textures")
	for _, t := range files {
		if e := filepath.Ext(t.Name); e != ".png" {
			continue
		}

		f, err := os.Open(t.Path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		allTextureFiles = append(allTextureFiles, t)
		f.Close()
	}
}

func getTexture(name string) (err error, tex *texture.Texture) {
	for _, t := range allTextureFiles {
		if t.Name == name {
			var f *os.File
			f, err = os.Open(t.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer f.Close()
			tex, err = texture.DecodePNG(f)
			return
		}
	}
	err = fmt.Errorf("Texture %s not fonded.", name)
	return
}

func digestBlock(in io.Reader) (err error, block Block) {
	buff, err := io.ReadAll(in)
	if err != nil {
		return
	}
	var b BlockJSON
	if err = json.Unmarshal(buff, &b); err != nil {
		return
	}
	err, block = b.ToBlock()
	return
}

func (b BlockJSON) ToBlock() (err error, block Block) {
	getAllTextures()
	block.Name = b.Name
	block.Id = b.Id
	var textures []*texture.Texture
	for _, texName := range b.Textures {
		var tex *texture.Texture
		err, tex = getTexture(texName)
		if err != nil {
			return
		}
		textures = append(textures, tex)
	}
	for i, texIdx := range b.TexIdx {
		if texIdx >= len(textures) {
			err = fmt.Errorf("Texture index overflow: index %d of %d", texIdx, len(textures))
		}
		block.Textures[i] = textures[texIdx]
	}
	return
}

func (b BlockJSON) SaveBlock() (err error) {
	path := filepath.Join(blocksDir, b.Name+".json")
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	encoder := jsontext.NewEncoder(f, jsontext.Multiline(true))
	err = json.MarshalEncode(encoder, b)
	return
}

func GetAllBlocks() (final []Block) {
	if files, err := os.ReadDir(blocksDir); err == nil {
		for _, file := range files {
			if !file.IsDir() {
				f, err := os.Open(filepath.Join(blocksDir, file.Name()))
				if err != nil {
					fmt.Println("error opening", file.Name())
					fmt.Println(err)
					continue
				}
				err, b := digestBlock(f)
				if err != nil {
					fmt.Println("error opening", file.Name())
					fmt.Println(err)
					continue
				}
				final = append(final, b)
				f.Close()
				fmt.Println(b.Name + " loaded!")
			}
		}
	}
	return
}
