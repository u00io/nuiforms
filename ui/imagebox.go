package ui

import (
	"image"

	"github.com/nfnt/resize"
)

type ImageBox struct {
	Widget
	img     image.Image
	scaling ImageBoxScale
}

type ImageBoxScale int

const (
	ImageBoxScaleNoScaleAdjustBox           ImageBoxScale = 0
	ImageBoxScaleNoScaleImageInLeftTop      ImageBoxScale = 1
	ImageBoxScaleNoScaleImageInCenter       ImageBoxScale = 2
	ImageBoxScaleStretchImage               ImageBoxScale = 3
	ImageBoxScaleAdjustImageKeepAspectRatio ImageBoxScale = 4
)

func NewImageBox() *ImageBox {
	var c ImageBox
	c.InitWidget()
	c.SetTypeName("ImageBox")
	c.SetOnPaint(c.draw)
	c.scaling = ImageBoxScaleNoScaleAdjustBox

	return &c
}

func (c *ImageBox) Image() image.Image {
	return c.img
}

func (c *ImageBox) SetImage(img image.Image) {
	c.img = img
	UpdateMainForm()
}

func (c *ImageBox) Scaling() ImageBoxScale {
	return c.scaling
}

func (c *ImageBox) SetScaling(scaling ImageBoxScale) {
	if c.scaling == scaling {
		return
	}
	c.scaling = scaling
	UpdateMainForm()
}

func (c *ImageBox) draw(cnv *Canvas) {
	if c.img == nil {
		return
	}

	if c.scaling == ImageBoxScaleNoScaleImageInLeftTop || c.scaling == ImageBoxScaleNoScaleAdjustBox {
		//b := c.img.Bounds()
		cnv.DrawImage(0, 0, c.img)
	}

	if c.scaling == ImageBoxScaleNoScaleImageInCenter {
		b := c.img.Bounds()
		offsetX := (c.Width() - b.Max.X) / 2
		offsetY := (c.Height() - b.Max.Y) / 2
		cnv.DrawImage(offsetX, offsetY, c.img)
	}

	if c.scaling == ImageBoxScaleStretchImage {
		img := resize.Resize(uint(c.Width()), uint(c.Height()), c.img, resize.Bicubic)
		cnv.DrawImage(0, 0, img)
	}

	if c.scaling == ImageBoxScaleAdjustImageKeepAspectRatio {
		b := c.img.Bounds()
		aspRatioImg := float64(b.Max.X) / float64(b.Max.Y)
		aspRationWidget := float64(c.Width()) / float64(c.Height())
		if aspRatioImg > aspRationWidget {
			img := resize.Resize(uint(c.Width()), 0, c.img, resize.Bicubic)
			b := img.Bounds()
			offsetX := (c.Width() - b.Max.X) / 2
			offsetY := (c.Height() - b.Max.Y) / 2
			cnv.DrawImage(offsetX, offsetY, img)
		} else {
			img := resize.Resize(0, uint(c.Height()), c.img, resize.Bicubic)
			b := img.Bounds()
			offsetX := (c.Width() - b.Max.X) / 2
			offsetY := (c.Height() - b.Max.Y) / 2
			cnv.DrawImage(offsetX, offsetY, img)
		}
	}
}

func (c *ImageBox) updateInnerSize() {
	minWidth := 0
	minHeight := 0
	maxWidth := 10000
	maxHeight := 10000
	if c.scaling == ImageBoxScaleNoScaleAdjustBox {
		if c.img != nil {
			minWidth = c.img.Bounds().Dx()
			minHeight = c.img.Bounds().Dy()
			maxWidth = minWidth
			maxHeight = minHeight
		}
	}
	c.SetMinSize(minWidth, minHeight)
	c.SetMaxSize(maxWidth, maxHeight)
}
