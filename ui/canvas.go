package ui

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type canvasState struct {
	col color.RGBA

	translateX int
	translateY int

	clipX int
	clipY int
	clipW int
	clipH int
}

type Canvas struct {
	rgba  *image.RGBA
	state *canvasState
	stack []*canvasState
}

type VAlign int
type HAlign int

const VAlignTop VAlign = 0
const VAlignCenter VAlign = 1
const VAlignBottom VAlign = 2

const HAlignLeft HAlign = 0
const HAlignCenter HAlign = 1
const HAlignRight HAlign = 2

func NewCanvas(rgba *image.RGBA) *Canvas {
	var c Canvas
	c.rgba = rgba
	c.state = &canvasState{}
	c.stack = make([]*canvasState, 0)
	return &c
}

func (c *Canvas) Save() {
	oldState := *c.state
	c.stack = append(c.stack, c.state)
	c.state = &oldState
}

func (c *Canvas) Restore() {
	if len(c.stack) == 0 {
		return
	}
	c.state = c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
}

func (c *Canvas) SetColor(col color.RGBA) {
	c.state.col = col
}

func (c *Canvas) SetClip(x, y, w, h int) {
	c.state.clipX += x
	c.state.clipY += y
	c.state.clipW = w
	c.state.clipH = h
}

func (c *Canvas) Translate(x, y int) {
	c.state.translateX += x
	c.state.translateY += y
}

func (c *Canvas) SetPixel(x, y int) {
	x += c.state.translateX
	y += c.state.translateY
	if c.state.clipX != 0 || c.state.clipY != 0 || c.state.clipW != 0 || c.state.clipH != 0 {
		if x < c.state.clipX || y < c.state.clipY || x >= c.state.clipX+c.state.clipW || y >= c.state.clipY+c.state.clipH {
			return
		}
	}
	c.rgba.Set(x, y, c.state.col)
}

func (c *Canvas) MixPixel(x int, y int, rgba color.Color) {

	if x < c.state.clipX || x > c.state.clipX+c.state.clipW {
		return
	}
	if y < c.state.clipY || y > c.state.clipY+c.state.clipH {
		return
	}

	cOld := c.rgba.At(x, y)
	oR, oG, oB, _ := cOld.RGBA()
	cR, cG, cB, cA := rgba.RGBA()
	cR = cR >> 8
	cG = cG >> 8
	cB = cB >> 8
	cA = cA >> 8

	oR = oR >> 8
	oG = oG >> 8
	oB = oB >> 8

	alpha := uint32(cA)
	antialpha := 255 - alpha

	if cA > 0 && cA < 255 {
		cA += 1
	}

	//alpha, antialpha = antialpha, alpha
	nR := uint8(((uint32(oR) * antialpha) >> 8) + ((cR * alpha) >> 8))
	nG := uint8(((uint32(oG) * antialpha) >> 8) + ((cG * alpha) >> 8))
	nB := uint8(((uint32(oB) * antialpha) >> 8) + ((cB * alpha) >> 8))

	c.rgba.SetRGBA(x, y, color.RGBA{nR, nG, nB, 255})
}

func (c *Canvas) DrawLine(x1 int, y1 int, x2 int, y2 int, width int, color color.Color) {
	x1 = x1 + c.state.translateX
	y1 = y1 + c.state.translateY
	x2 = x2 + c.state.translateX
	y2 = y2 + c.state.translateY

	x1, y1, x2, y2, visible := CohenSutherland(x1, y1, x2, y2, c.state.clipX, c.state.clipY, c.state.clipX+c.state.clipW, c.state.clipY+c.state.clipH)
	if visible {
		if (x1 == x2 || y1 == y2) && width == 1 {
			if x1 == x2 {
				if y1 > y2 {
					y1, y2 = y2, y1
				}
				for y := y1; y < y2; y++ {
					c.MixPixel(x1, y, color)
				}
			}
			if y1 == y2 {
				if x1 > x2 {
					x1, x2 = x2, x1
				}
				for x := x1; x < x2; x++ {
					c.MixPixel(x, y1, color)
				}
			}
		} else {
			script := c.MakeScriptLine(float64(x1), float64(y1), float64(x2), float64(y2), float64(width))
			script.Bounds = image.Rectangle{Min: image.Point{X: c.state.clipX, Y: c.state.clipY}, Max: image.Point{X: c.state.clipX + c.state.clipW, Y: c.state.clipY + c.state.clipH}}
			script.DrawToRGBA(c.rgba, color)
		}
	}
}

func (c *Canvas) MakeScriptLine(x1, y1, x2, y2, width float64) *DrawScript {
	x1 = math.Round(x1)
	y1 = math.Round(y1)
	x2 = math.Round(x2)
	y2 = math.Round(y2)
	width = math.Round(width)

	result := NewDrawScript()

	{
		minX := math.Min(x1, x2)
		maxX := math.Max(x1, x2)
		minY := math.Min(y1, y2)
		maxY := math.Max(y1, y2)

		if int(minX) > c.state.clipX+c.state.clipW {
			return result
		}
		if int(maxX) < c.state.clipX {
			return result
		}
		if int(minY) > c.state.clipY+c.state.clipH {
			return result
		}
		if int(maxY) < c.state.clipY {
			return result
		}
	}

	if width > 1.1 {

		width -= 1
		width /= 2

		deltaX := x2 - x1
		deltaY := y2 - y1

		tnA := deltaY / deltaX

		katY := width * math.Cos(math.Atan(tnA))
		katX := math.Sqrt(math.Abs(width*width - katY*katY))

		if deltaX <= 0 && deltaY <= 0 {
			katY = -katY
		}

		if deltaX <= 0 && deltaY > 0 {
			katX = -katX
			katY = -katY
		}

		if deltaX > 0 && deltaY > 0 {
			katX = -katX
		}

		a1_x := x1 - katX
		a1_y := y1 - katY
		a2_x := x1 + katX
		a2_y := y1 + katY

		b1_x := x2 - katX
		b1_y := y2 - katY
		b2_x := x2 + katX
		b2_y := y2 + katY

		var ps []image.Point
		ps = make([]image.Point, 0)
		ps = append(ps, image.Pt(int(a1_x), int(a1_y)))
		ps = append(ps, image.Pt(int(a2_x), int(a2_y)))
		ps = append(ps, image.Pt(int(b2_x), int(b2_y)))
		ps = append(ps, image.Pt(int(b1_x), int(b1_y)))
		result.append(c.fillEx(ps))

		result.append(c.MakeScriptLineWu(a1_x, a1_y, a2_x, a2_y))
		result.append(c.MakeScriptLineWu(b1_x, b1_y, b2_x, b2_y))
		result.append(c.MakeScriptLineWu(a1_x, a1_y, b1_x, b1_y))
		result.append(c.MakeScriptLineWu(a2_x, a2_y, b2_x, b2_y))
	} else {
		result.append(c.MakeScriptLineWu(x1, y1, x2, y2))
	}

	return result
}

func (c *Canvas) MakeScriptLineBresenham(x1 float64, y1 float64, x2 float64, y2 float64) *DrawScript {
	x1 = math.Round(x1)
	y1 = math.Round(y1)
	x2 = math.Round(x2)
	y2 = math.Round(y2)

	result := NewDrawScript()

	stepByX := true
	if math.Abs(x2-x1) < math.Abs(y2-y1) {
		stepByX = false
	}

	if stepByX {
		if x2 < x1 { // Swap
			xo := x1
			x1 = x2
			x2 = xo
			yo := y1
			y1 = y2
			y2 = yo
		}

		beginX := int(x1)
		endX := int(x2)

		deltaX := math.Abs(x2 - x1)
		deltaY := math.Abs(y2 - y1)

		y := int(y1)
		var err float64 = 0
		errDelta := deltaY / deltaX
		var yDir = 1
		if (y2 - y1) < 0 {
			yDir = -1
		}

		for x := beginX; x <= endX; x++ {
			result.plot(x, y, 1)
			err += errDelta
			if err > 0.5 {
				y += yDir
				err -= 1.0
			}
		}
	} else {
		if y2 < y1 {
			xo := x1
			x1 = x2
			x2 = xo
			yo := y1
			y1 = y2
			y2 = yo
		}

		beginY := int(y1)
		endY := int(y2)

		deltaX := math.Abs(x2 - x1)
		deltaY := math.Abs(y2 - y1)

		x := int(x1)
		err := 0.0
		errDelta := deltaX / deltaY
		xDir := 1
		if (x2 - x1) < 0 {
			xDir = -1
		}

		for y := beginY; y <= endY; y++ {
			result.plot(x, y, 1)
			err += errDelta
			if err > 0.5 {
				x += xDir
				err -= 1.0
			}
		}
	}

	return result
}

func (c *Canvas) fillEx(points []image.Point) *DrawScript {
	result := NewDrawScript()
	if len(points) < 3 {
		return result
	}

	// Borders
	lastPoint := points[0]
	for i := 1; i < len(points); i++ {
		result.append(c.MakeScriptLineBresenham(float64(lastPoint.X), float64(lastPoint.Y), float64(points[i].X), float64(points[i].Y)))
		lastPoint = points[i]
	}
	result.append(c.MakeScriptLineBresenham(float64(lastPoint.X), float64(lastPoint.Y), float64(points[0].X), float64(points[0].Y)))

	minX := int(math.MaxInt32)
	maxX := int(math.MinInt32)
	minY := int(math.MaxInt32)
	maxY := int(math.MinInt32)

	for _, pp := range points {
		if pp.X > maxX {
			maxX = pp.X
		}
		if pp.X < minX {
			minX = pp.X
		}
		if pp.Y > maxY {
			maxY = pp.Y
		}
		if pp.Y < minY {
			minY = pp.Y
		}
	}

	beginX := 0

	for y := minY + 1; y <= maxY-1; y++ {
		countOfTriggers := 0
		flag := false

		for x := minX; x <= maxX+1; x++ {
			if result.hasPixel(x, y) && !result.hasPixel(x+1, y) {
				countOfTriggers++
			}
		}

		if countOfTriggers > 1 {
			for x := minX; x <= maxX; x++ {
				if result.hasPixel(x, y) && !result.hasPixel(x+1, y) {
					if flag {
						var line DrawScriptHorLine
						line.X1 = beginX
						line.X2 = x
						line.Y = y
						line.C = 1
						result.horLines = append(result.horLines, line)
					} else {
						beginX = x
					}

					flag = !flag
				}
			}
		}
	}

	return result
}

func ipart(x float64) int {
	rn := math.Floor(x)
	return int(rn)
}

func round(x float64) int {
	iprt := ipart(x + 0.5)
	return iprt
}

func fpart(x float64) float64 {
	return x - float64(ipart(x))
}

func (c *Canvas) MakeScriptLineWu(x1, y1, x2, y2 float64) *DrawScript {

	result := NewDrawScript()

	x1 = math.Round(x1)
	y1 = math.Round(y1)
	x2 = math.Round(x2)
	y2 = math.Round(y2)

	if math.Abs(x1-x2) < 0.1 && math.Abs(y1-y2) < 0.1 {
		return result
	}

	stepByX := true
	if math.Abs(x2-x1) < math.Abs(y2-y1) {
		stepByX = false
	}

	if stepByX {
		if x2 < x1 {
			xo := x1
			x1 = x2
			x2 = xo
			yo := y1
			y1 = y2
			y2 = yo
		}

		dx := x2 - x1
		dy := y2 - y1
		gradient := dy / dx
		offset := y1 + gradient*(float64(round(x1))-x1)

		workFrom := int(x1)
		workTo := int(x2)

		xgap := 1 - fpart(x1+0.5)
		result.plot(workFrom, int(math.Floor(y1)), (1-fpart(offset))*xgap)
		result.plot(workFrom, int(math.Floor(y1))+1, fpart(offset)*xgap)

		xgap = fpart(x1 + 0.5)
		result.plot(workTo, int(math.Floor(y2)), (1-fpart(offset))*xgap)
		result.plot(workTo, int(math.Floor(y2))+1, fpart(offset)*xgap)

		offset += gradient

		for workIndex := workFrom + 1; workIndex <= workTo-1; workIndex++ {
			err := fpart(offset)
			if err > 0 {
				c1 := 1 - err
				c2 := err
				result.plot(workIndex, ipart(offset), c1)
				result.plot(workIndex, ipart(offset)+1, c2)
			} else {
				c1 := 1 - math.Abs(err)
				c2 := math.Abs(err)
				result.plot(workIndex, ipart(offset), c1)
				result.plot(workIndex, ipart(offset)-1, c2)
			}
			offset = offset + gradient
		}
	} else {
		if y2 < y1 {
			xo := x1
			x1 = x2
			x2 = xo
			yo := y1
			y1 = y2
			y2 = yo
		}

		dx := x2 - x1
		dy := y2 - y1

		gradient := dx / dy

		offset := x1 + gradient*(float64(round(y1))-y1)

		workFrom := int(y1)
		workTo := int(y2)

		xgap := 1 - fpart(y1+0.5)
		result.plot(int(math.Floor(x1)), workFrom, (1-fpart(offset))*xgap)
		result.plot(int(math.Floor(x1))+1, workFrom, fpart(offset)*xgap)

		xgap = fpart(y1 + 0.5)
		result.plot(int(math.Floor(x2)), workTo, (1-fpart(offset))*xgap)
		result.plot(int(math.Floor(x2))+1, workTo, fpart(offset)*xgap)

		offset += gradient

		for workIndex := workFrom + 1; workIndex <= workTo-1; workIndex++ {
			err := fpart(offset)
			if err > 0 {
				c1 := 1 - err
				c2 := err
				result.plot(ipart(offset), workIndex, c1)
				result.plot(ipart(offset)+1, workIndex, c2)
			} else {
				c1 := 1 - math.Abs(err)
				c2 := math.Abs(err)
				result.plot(ipart(offset), workIndex, c1)
				result.plot(ipart(offset)-1, workIndex, c2)
			}
			offset = offset + gradient
		}
	}

	return result
}

func CohenSutherland(x1, y1, x2, y2, left, top, right, bottom int) (int, int, int, int, bool) {

	invalid := false

	for {
		// Left-Right-Top-Bottom
		code1 := 0
		if x1 < left {
			code1 |= 8
		}
		if x1 > right {
			code1 |= 4
		}
		if y1 < top {
			code1 |= 2
		}
		if y1 > bottom {
			code1 |= 1
		}

		code2 := 0
		if x2 < left {
			code2 |= 8
		}
		if x2 > right {
			code2 |= 4
		}
		if y2 < top {
			code2 |= 2
		}
		if y2 > bottom {
			code2 |= 1
		}

		if code1 == 0 && code2 == 0 {
			break
		}

		if (code1 & code2) != 0 {
			invalid = true
			break // Outside
		}

		code := code1
		if code == 0 {
			code = code2
		}

		x := 0
		y := 0

		if (code & 2) != 0 {
			x = x1 + (x2-x1)*(top-y1)/(y2-y1)
			y = top
		} else {
			if (code & 1) != 0 {
				x = x1 + (x2-x1)*(bottom-y1)/(y2-y1)
				y = bottom
			} else {
				if (code & 4) != 0 {
					y = y1 + (y2-y1)*(right-x1)/(x2-x1)
					x = right
				} else {
					if (code & 8) != 0 {
						y = y1 + (y2-y1)*(left-x1)/(x2-x1)
						x = left
					}
				}
			}
		}

		if code1 != 0 {
			x1 = x
			y1 = y
		} else {
			x2 = x
			y2 = y
		}
	}

	return x1, y1, x2, y2, !invalid
}

func (c *Canvas) FillRect(x int, y int, width int, height int, colr color.Color) {

	x = x + c.state.translateX
	y = y + c.state.translateY

	if x < 0 {
		width += x
		x = 0
	}

	if y < 0 {
		height += y
		y = 0
	}

	if x+width > c.rgba.Rect.Max.X {
		width = c.rgba.Rect.Max.X - x
	}

	if y+height > c.rgba.Rect.Max.Y {
		height = c.rgba.Rect.Max.Y - y
	}

	if width < 0 {
		return
	}

	if height < 0 {
		return
	}

	if x < c.state.clipX {
		width -= c.state.clipX - x
		x = c.state.clipX
	}

	if y < c.state.clipY {
		height -= c.state.clipY - y
		y = c.state.clipY
	}

	if x+width > c.state.clipX+c.state.clipW {
		width = c.state.clipX + c.state.clipW - x
	}

	if y+height > c.state.clipY+c.state.clipH {
		height = c.state.clipY + c.state.clipH - y
	}

	if width < 1 {
		return
	}

	if height < 1 {
		return
	}

	if _, _, _, a := colr.RGBA(); a == 65535 {
		line := make([]uint8, width*4)
		cc := color.RGBAModel.Convert(colr).(color.RGBA)
		for xx := 0; xx < width*4; xx += 4 {
			line[xx] = cc.R
			line[xx+1] = cc.G
			line[xx+2] = cc.B
			line[xx+3] = cc.A
		}

		for yy := y; yy < y+height; yy++ {
			offset := yy*c.rgba.Stride + x*4
			copy(c.rgba.Pix[offset:offset+width*4], line[:])
		}
	} else {
		for yy := y; yy < y+height; yy++ {
			for xx := x; xx < x+width; xx++ {
				c.MixPixel(xx, yy, colr)
			}
		}
	}
}

func (c *Canvas) DrawText(x int, y int, text string, fontFamily string, fontSize float64, colr color.Color, underline bool) {
	if math.IsNaN(fontSize) {
		return
	}

	x += c.state.translateX
	y += c.state.translateY

	fontContext, err := FontContext(fontFamily, fontSize, false, false, colr)
	if err != nil {
		return
	}

	rectBounds := image.Rectangle{}
	rectBounds.Min.X = c.state.clipX
	rectBounds.Min.Y = c.state.clipY
	rectBounds.Max.X = c.state.clipX + c.state.clipW
	rectBounds.Max.Y = c.state.clipY + c.state.clipH

	fontContext.SetClip(rectBounds) // cnv.image.Bounds()
	fontContext.SetDst(c.rgba)

	_, face, err := Font(fontFamily, fontSize, false, false)
	if err != nil {
		return
	}

	//face := FontFace(f, fontSize)
	pt := freetype.Pt(x, y+face.Metrics().Ascent.Ceil())
	fontContext.DrawString(text, pt)
}

func (c *Canvas) TranslatedX() int {
	return c.state.translateX
}

func (c *Canvas) TranslatedY() int {
	return c.state.translateY
}

func (c *Canvas) ClipIn(x, y, width, height int) {
	nClipX := x
	nClipY := y
	nClipW := width
	nClipH := height

	if nClipX < c.state.clipX {
		nClipW -= c.state.clipX - nClipX
		nClipX = c.state.clipX
	}

	if nClipX+nClipW > c.state.clipX+c.state.clipW {
		nClipW -= (nClipX + nClipW) - (c.state.clipX + c.state.clipW)
	}

	if nClipY < c.state.clipY {
		nClipH -= c.state.clipY - nClipY
		nClipY = c.state.clipY
	}

	if nClipY+nClipH > c.state.clipY+c.state.clipH {
		nClipH -= (nClipY + nClipH) - (c.state.clipY + c.state.clipH)
	}

	if nClipW < 0 {
		nClipW = 0
	}

	if nClipH < 0 {
		nClipH = 0
	}

	c.state.clipX = nClipX
	c.state.clipY = nClipY
	c.state.clipW = nClipW
	c.state.clipH = nClipH
}

func (c *Canvas) DrawTextMultiline(x int, y int, width int, height int, hAlign HAlign, vAlign VAlign, text string, colr color.Color, fontFamily string, fontSize float64, underline bool) {
	lines := strings.Split(text, "\r\n")

	yOffset := 0

	_, textHeight, err := MeasureText(fontFamily, fontSize, false, false, "Йg", true)
	if err != nil {
		return
	}

	//textHeight := 20

	fulltextHeight := textHeight * len(lines)

	switch vAlign {
	case VAlignTop:
		yOffset = 0
	case VAlignCenter:
		yOffset = height/2 - fulltextHeight/2
	case VAlignBottom:
		yOffset = height - fulltextHeight
	}

	c.Save()
	c.ClipIn(x+c.TranslatedX(), y+c.TranslatedY(), width, height)

	for _, str := range lines {
		xx := x
		textWidth, _, err := MeasureText(fontFamily, fontSize, false, false, str, false)
		if err != nil {
			return
		}

		if hAlign != HAlignLeft {

			switch hAlign {
			case HAlignLeft:
				xx = x
			case HAlignCenter:
				xx = (width / 2) - (textWidth / 2) + x
			case HAlignRight:
				xx = width - textWidth + x
			}
		}

		c.DrawText(xx, yOffset+y, str, fontFamily, fontSize, colr, underline)

		if underline {
			underLineWidth := fontSize / 20
			if underLineWidth < 1 {
				underLineWidth = 1
			}
			c.DrawLine(xx, yOffset+y+textHeight-1, xx+textWidth, yOffset+y+textHeight-1, int(underLineWidth), colr)
		}

		yOffset += textHeight
	}

	c.Restore()
}

func MeasureText(family string, size float64, bold bool, italic bool, text string, multiline bool) (int, int, error) {
	measureData := GetFontMeasureData(family, size, bold, italic, text, multiline)
	return measureData.Width, measureData.Height, nil
}

type FontMeasureData struct {
	Width  int
	Height int
}

var fontMeasureData map[string]*FontMeasureData

func GetFontMeasureData(family string, size float64, bold bool, italic bool, text string, multiline bool) *FontMeasureData {

	if fontMeasureData == nil {
		fontMeasureData = make(map[string]*FontMeasureData)
	}

	/*lines := strings.FieldsFunc(text, func(r rune) bool {
		return r == '\n'
	})
	linesCount := len(lines)
	if linesCount < 1 {
		linesCount = 1
	}*/

	fontHash := "font-" + family + "-size-" + strconv.FormatFloat(size, 'f', 10, 64) + "-bold-" + strconv.FormatBool(bold) + "-italic-" + strconv.FormatBool(italic) + text
	if fmd, ok := fontMeasureData[fontHash]; ok {
		return fmd
	}

	var fmdNew FontMeasureData
	fmdNew.Width, fmdNew.Height, _ = MeasureTextFreeType(family, size, bold, italic, text, multiline)
	fontMeasureData[fontHash] = &fmdNew

	return &fmdNew
}

func MeasureTextFreeType(family string, size float64, bold bool, italic bool, text string, multiline bool) (int, int, error) {

	if math.IsNaN(size) {
		return 0, 0, errors.New("font size is NaN")
	}

	_, face, err := Font(family, size, bold, italic)

	if err != nil {
		return 0, 0, err
	}

	//face := FontFace(f, size)

	d := &font.Drawer{
		Dst:  nil,
		Src:  nil,
		Face: face,
	}

	var textWidth int
	var textHeight int

	var lines []string

	if multiline {
		lines = strings.Split(text, "\r\n")
	} else {
		lines = make([]string, 0)
		lines = append(lines, text)
	}

	for _, str := range lines {
		//pt := freetype.Pt(x, y+face.Metrics().Ascent.Ceil())

		w := d.MeasureString(str) / 64

		h := face.Metrics().Ascent + face.Metrics().Descent
		textHeight += int(h)
		if int(w) > textWidth {
			textWidth = int(w)
		}
	}

	return int(textWidth), int(textHeight / 64), nil
}

type FontInfo struct {
	family         string
	fontRegular    *truetype.Font
	fontBold       *truetype.Font
	fontItalic     *truetype.Font
	fontBoldItalic *truetype.Font
	faces          map[float64]font.Face
}

var globalFonts map[string]*FontInfo

//go:embed "fonts/Roboto_Regular.ttf"
var fontRoboto []byte

//go:embed "fonts/RobotoMono_Regular.ttf"
var fontRobotoMono []byte

func init() {
	_ = embed.FS.Open

	globalFonts = make(map[string]*FontInfo)
	addFont("roboto", fontRoboto)
	addFont("robotomono", fontRobotoMono)
}

func addFont(name string, data []byte) {

	if _, ok := globalFonts[name]; !ok {
		var fontInfo FontInfo
		fontInfo.family = strings.ToLower(name)
		fontInfo.faces = make(map[float64]font.Face)
		globalFonts[name] = &fontInfo
	}

	val := globalFonts[name]
	f, err := freetype.ParseFont(data)
	if err != nil {
		return
	}
	val.fontRegular = f
	return
}

func Font(family string, size float64, bold bool, italic bool) (result *truetype.Font, face font.Face, err error) {
	var val *FontInfo
	ok := false

	if val, ok = globalFonts[strings.ToLower(family)]; ok {
		if !bold && !italic {
			result = val.fontRegular
		}
		if bold && !italic {
			result = val.fontBold
		}
		if !bold && italic {
			result = val.fontItalic
		}
		if bold && italic {
			result = val.fontBoldItalic
		}
	}

	if result == nil || val == nil {
		err = fmt.Errorf("no font found")
		return
	}

	if _, ok := val.faces[size]; !ok {
		val.faces[size] = MakeFontFace(result, size)
	}

	return result, val.faces[size], err
}

var dpi = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")

func MakeFontFace(f *truetype.Font, size float64) font.Face {
	h := font.HintingNone
	face := truetype.NewFace(f, &truetype.Options{
		Size:              size,
		DPI:               *dpi,
		Hinting:           h,
		GlyphCacheEntries: 64,
	})
	return face
}

func FontContext(family string, size float64, bold bool, italic bool, color color.Color) (*Context, error) {
	fg := image.NewUniform(color)
	f, _, err := Font(family, size, bold, italic)
	if err != nil {
		return nil, err
	}

	c := NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	//c.SetClip(cnv.image.Bounds())
	//c.SetDst(cnv.image)
	c.SetSrc(fg)
	return c, nil
}
