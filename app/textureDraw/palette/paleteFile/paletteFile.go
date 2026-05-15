package pallete_file

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func downloadFile(url, path string) (*bytes.Buffer, error) {
	fileName := filepath.Base(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad status code: %s", resp.Status)
	}

	out := bytes.NewBuffer(nil)

	io.Copy(out, resp.Body)
	if f, err := os.Create(path + "/" + fileName); err != nil {
		fmt.Println(err)
	} else {
		f.Write(out.Bytes())
		f.Close()
	}
	return out, nil
}

func hexToFloat(hex string) [3]float32 {
	showError := func(index int) {
		fmt.Println("Unrecognized character at index " + strconv.FormatInt(int64(index), 10) + ", in \"" + hex + "\"")
	}
	count := 0
	floatIndex := 0
	final := [3]float32{}
	buf := strings.Builder{}
	if len(hex) != 6 {
		fmt.Println("\"" + hex + "\" is not a hex color")
		return final
	}
	for i, r := range hex {
		buf.WriteRune(r)
		count += 1
		if count == 2 {
			num, err := strconv.ParseInt(buf.String(), 16, 64)
			if err != nil {
				fmt.Println(err)
				showError(i)
				return [3]float32{}
			}
			if floatIndex < 3 {
				final[floatIndex] = float32(num) / 255
				buf.Reset()
				floatIndex += 1
				count = 0
			}
		}
	}
	return final
}

type PaletteData struct {
	Name    string
	Creator string
	Colors  [][3]float32
}

func digestPalette(buff io.Reader) (palette PaletteData) {
	r := csv.NewReader(buff)
	d, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	var name string
	var creator string
	var colors [][3]float32
	for _, line := range d {
		for i, col := range line {
			switch i {
			case 0:
				name = col
			case 1:
				creator = col
			default:
				colors = append(colors, hexToFloat(col))
			}
		}
	}
	return PaletteData{name, creator, colors}
}

func GetPaletteFromLospec(paletteName string) (bool, PaletteData) {
	strings.ReplaceAll(paletteName, " ", "-")
	if buff, err := downloadFile("https://Lospec.com/palette-list/"+paletteName+".csv", "UserData/palettes"); err != nil {
		fmt.Println(err)
		return false, PaletteData{}
	} else {
		return true, digestPalette(buff)
	}
}

func GetPaletteFromLospecLink(palleteLink string) (bool, PaletteData) {
	if buff, err := downloadFile(palleteLink+".csv", "UserData/palettes"); err != nil {
		fmt.Println(err)
		return false, PaletteData{}
	} else {
		return true, digestPalette(buff)
	}
}

func GetAllPalettes() (final []PaletteData) {
	path := "UserData/palettes/"
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, file := range files {
			if !file.IsDir() {
				if f, err := os.Open(path + file.Name()); err != nil {
					fmt.Println(err)
				} else {
					p := digestPalette(f)
					if len(p.Colors) > 0 {
						final = append(final, p)
					}
					f.Close()
				}
			}
		}
	}
	return
}
