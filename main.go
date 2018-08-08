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
)

const (
	IMAGE_PNG_CONTENT_TYPE = "image/png"
)

func main() {
	e := echo.New()
	e.GET("/:rune_name", func(c echo.Context) error {
		rune_name := c.Param("rune_name")


		return c.Blob(http.StatusOK, IMAGE_PNG_CONTENT_TYPE,  createImage(rune_name))
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

func fontAwesomeIcon(iconName string) (string, error) {

	foundRune := fontAwesomeRuneMap()[iconName]
	if foundRune == "" {
		return "", errors.New(fmt.Sprint("No rune found for", iconName))
	}

	n, _ := strconv.ParseUint(foundRune, 16, 32)

	return string(rune(n)), nil
}


func createImage(rune_name string) []byte {
	const S = 1024

	icon, _ := fontAwesomeIcon(rune_name)

	dc := gg.NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, S, S)))
	dc.SetRGBA(0,0,0,0)
	dc.Clear()
	dc.SetRGB255(10,255,255)
	//dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("./fonts/fontawesome/fa-regular-400.ttf", 1024); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(icon, S/2, S/2, 0.5, 0.5)

	imageBytes := getImageBytes(dc)
	return imageBytes
}

func getImageBytes(dc *gg.Context) []byte {
	var b bytes.Buffer

	png.Encode(bufio.NewWriter(&b), dc.Image())

	return b.Bytes()
}
