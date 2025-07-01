package ui

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var loadedFonts = make(map[string][]byte)

//go:embed "fonts/Roboto_Regular.ttf"
var fontRoboto []byte

//go:embed "fonts/RobotoMono_Regular.ttf"
var fontRobotoMono []byte

type renderedText struct {
	key          string
	lastAccessDT time.Time
	textImage    *image.RGBA
}

var renderedTexts = make(map[string]*renderedText)
var renderedTextLastClearDT time.Time

func init() {
	loadedFonts["roboto"] = fontRoboto
	loadedFonts["robotomono"] = fontRobotoMono
}

func clearRenderedTexts() {
	if renderedTextLastClearDT.IsZero() {
		renderedTextLastClearDT = time.Now()
		return
	}
	if time.Since(renderedTextLastClearDT) < 1*time.Minute {
		return
	}
	renderedTextLastClearDT = time.Now()
	for k, rt := range renderedTexts {
		if time.Since(rt.lastAccessDT) > 30*time.Second {
			delete(renderedTexts, k)
		}
	}
}

func DrawText(rgba *image.RGBA, text string, textColor color.Color, fontFamily string, fontSize float64, x int, y int, clipX int, clipY int, clipWidth int, clipHeight int) {
	var textImage *image.RGBA
	key := fontFamily + "_" + fmt.Sprint(fontSize) + "_" + fmt.Sprintf("%v", textColor) + "_" + text

	if img, ok := renderedTexts[key]; ok {
		textImage = img.textImage
		img.lastAccessDT = time.Now()
	} else {
		textImage = renderText(text, textColor, fontFamily, fontSize)
		if textImage == nil {
			return
		}
		var img renderedText
		img.key = key
		img.lastAccessDT = time.Now()
		img.textImage = textImage
		renderedTexts[key] = &img
		clearRenderedTexts()
	}

	textBounds := textImage.Bounds()
	if textBounds.Dx() <= 0 || textBounds.Dy() <= 0 {
		return
	}

	dstRect := image.Rect(x, y, x+textBounds.Dx(), y+textBounds.Dy())
	clipRect := image.Rect(clipX, clipY, clipX+clipWidth, clipY+clipHeight)
	dstRect = dstRect.Intersect(clipRect)

	if dstRect.Empty() {
		return
	}

	srcStart := textBounds.Min.Add(dstRect.Min.Sub(image.Pt(x, y)))
	draw.Draw(rgba, dstRect, textImage, srcStart, draw.Over)
}

func GetCharPositions(fontFamily string, fontSize float64, text string) ([]int, error) {
	positions := make([]int, len(text)+1) // на 1 больше, чтобы последняя позиция = ширине всей строки
	var advance fixed.Int26_6

	face, err := getFace(fontFamily, fontSize)
	if err != nil {
		return positions, err
	}
	defer face.Close()

	for i, r := range text {
		positions[i] = advance.Round()

		a, ok := face.GlyphAdvance(r)
		if !ok {
			continue
		}
		advance += a
	}
	positions[len(text)] = advance.Round()
	return positions, nil
}

func MeasureText(fontFamily string, fontSize float64, text string) (int, int, error) {
	face, err := getFace(fontFamily, fontSize)
	if err != nil {
		return 0, 0, err
	}
	defer face.Close()
	metrics := face.Metrics()
	textWidth := font.MeasureString(face, text).Ceil()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	return textWidth, textHeight, nil
}

func getFace(fontFamily string, fontSize float64) (font.Face, error) {
	fontFamily = strings.ToLower(fontFamily)
	var fontBytes []byte
	if bs, ok := loadedFonts[fontFamily]; ok {
		fontBytes = bs
	}

	ft, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	return face, err
}

func renderText(text string, textColor color.Color, fontFamily string, fontSize float64) *image.RGBA {
	face, err := getFace(fontFamily, fontSize)
	if err != nil {
		return nil
	}
	defer face.Close()

	textWidth := font.MeasureString(face, text).Ceil()
	_ = textWidth
	metrics := face.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()

	rgba := image.NewRGBA(image.Rect(0, 0, textWidth, textHeight))

	//posX := x + (w-textWidth)/2
	posX := 0
	//posY := y + (h-textHeight)/2 + metrics.Ascent.Ceil()
	posY := textHeight - metrics.Descent.Ceil()

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(textColor),
		Face: face,
		Dot:  fixed.P(posX, posY),
	}
	d.DrawString(text)
	return rgba
}
