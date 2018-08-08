package main

import (
	"github.com/fogleman/gg"
	"image/png"
	"bytes"
	"bufio"
	"github.com/labstack/echo"
	"net/http"
		"image"
		"strconv"
		"fmt"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/golang/freetype/truetype"
)

const (
	IMAGE_PNG_CONTENT_TYPE = "image/png"
)

type IconRequest struct {
	Style string
	Name string
	Color string
	Points float64
}

func main() {
	data, err := ioutil.ReadFile("fonts/fontawesome/icons.yml")

	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/:style/:name", func(c echo.Context) error {
		iconStyle := c.Param("style")
		iconName := c.Param("name")

		points := c.QueryParam("points")
		fmt.Println("points:",points)

		color := c.QueryParam("color")
		fmt.Println("color:", color)
		
		imageBytes, err := createImage(iconName, iconStyle, data)

		if err != nil {
			return c.HTML(http.StatusNotFound, err.Error())
		}

		return c.Blob(http.StatusOK, IMAGE_PNG_CONTENT_TYPE, imageBytes )
	})
	e.Logger.Fatal(e.Start(":1400"))
}

func fontAwesomeRuneMap() map[string]string {
	return map[string]string {
		"comment":"f075",
		"comment-alt":"f27a",
		"comment-dots":"f4ad",
	}
}

type FAIcon struct {
	Name struct {
		Changes []string `yaml:"changes"`
		Label string `yaml:"label"`
		Styles []string `yaml:"styles"`
		Unicode string `yaml:"unicode"`
	}
}

//https://stackoverflow.com/questions/32147325/how-to-parse-yaml-with-dyanmic-key-in-golang
func fontAwesomeIcon(name, style string, icons []byte) (string, error) {
	items := make(map[string]interface{})

	yaml.Unmarshal([]byte(icons), items)

	if items[name] == nil {
		return "", errors.New(fmt.Sprint("No rune found for", name))
	}


	// TODO: check style and throw error if not found for the icon.

	if foundRune, ok := items[name].(map[interface{}]interface{})["unicode"].(string); ok {
		n, _ := strconv.ParseUint(foundRune, 16, 32)

		return string(rune(n)), nil
	} else {
		return "", errors.New(fmt.Sprint("No rune found for", name))
	}
}


func createImage(runeName string, runeType string, icons []byte) ([]byte, error) {
	const S = 1024

	icon, err := fontAwesomeIcon(runeName, runeType, icons)
	if err != nil {
		return nil, err
	}

	dc := gg.NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, S, S)))
	dc.SetRGBA(0,0,0,0)
	dc.Clear()
	dc.SetRGB255(10,255,255)

	if runeType == "regular" {
		loadFontFaceFromCache("fonts/fontawesome/fa-regular-400.ttf", 1024, dc)
	} else if runeType == "solid" {
		loadFontFaceFromCache("fonts/fontawesome/fa-solid-900.ttf", 512, dc)
	} else if runeType == "brands" {
		loadFontFaceFromCache("fonts/fontawesome/fa-brands-400.ttf", 1024, dc)
	}

	dc.DrawStringAnchored(icon, S/2, S/2, 0.5, 0.5)

	imageBytes := getImageBytes(dc)
	return imageBytes, nil
}

func getImageBytes(dc *gg.Context) []byte {
	var b bytes.Buffer

	png.Encode(bufio.NewWriter(&b), dc.Image())

	return b.Bytes()
}

func loadFontFaceFromCache(fontFaceName string, points float64, dc *gg.Context) {
	asset, err := Asset(fontFaceName)
	if err != nil {
		panic(err)
	}

	f, err := truetype.Parse(asset)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})

	if err == nil {
		dc.SetFontFace(face)
		//dc.fontHeight = points * 72 / 96
	}
}