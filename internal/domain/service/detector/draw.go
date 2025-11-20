package detector

import (
	"FairLAP/pkg/yolo_model"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"slices"
)

var opentypeFont, _ = opentype.Parse(fontTTF)
var colors = []color.RGBA{
	{0, 255, 0, 0},   // зеленый
	{255, 0, 0, 0},   // синий
	{0, 0, 255, 0},   // красный
	{255, 255, 0, 0}, // голубой
	{255, 0, 255, 0}, // пурпурный
	{0, 255, 255, 0}, // желтый
	{128, 0, 128, 0}, // фиолетовый
	{255, 165, 0, 0}, // оранжевый
}

func drawDetected(img image.Image, detections []yolo_model.Detection) image.Image {
	res := image.NewRGBA(img.Bounds())
	draw.Draw(res, img.Bounds(), img, image.Point{}, draw.Src)

	fontSize := float64(img.Bounds().Dy()) / 40

	face, _ := opentype.NewFace(opentypeFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	fontHeight := face.Metrics().Height.Floor()

	drawer := &font.Drawer{
		Dst:  res,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	fontRects := make([]image.Rectangle, 0, len(detections))

	slices.SortFunc(detections, func(d1, d2 yolo_model.Detection) int {
		return d1.BBox.Min.Y - d2.BBox.Min.Y
	})

	thickness := max(int(float64(img.Bounds().Dy())/500), 1)

	for _, detection := range detections {
		c := colors[detection.ClassID%len(colors)]

		drawRect(res, c, detection.BBox, thickness)

		lbl := fmt.Sprintf("%s: %0.2f", detection.ClassName, detection.Confidence)

		width := font.MeasureString(face, lbl).Ceil()

		x, y := detection.BBox.Min.X+thickness+2, detection.BBox.Min.Y+fontHeight+thickness

		fontRect := image.Rect(x, y-fontHeight, x+width, y)

		for _, rect := range fontRects {
			if rect.Overlaps(fontRect) {
				fontRect.Min.Y += fontHeight
				fontRect.Max.Y += fontHeight
			}
		}
		drawer.Dot = fixed.P(fontRect.Min.X, fontRect.Max.Y)
		drawer.DrawString(lbl)

		fontRects = append(fontRects, fontRect)
	}

	return res
}

func drawRect(img draw.Image, c color.Color, r image.Rectangle, thickness int) {
	for j := range thickness {
		for i := r.Min.X; i <= r.Max.X; i++ {
			img.Set(i, r.Min.Y+j, c)
			img.Set(i, r.Max.Y+j, c)
		}
	}
	for j := range thickness {
		for i := r.Min.Y; i <= r.Max.Y; i++ {
			img.Set(r.Min.X+j, i, c)
			img.Set(r.Max.X+j, i, c)
		}
	}
}
