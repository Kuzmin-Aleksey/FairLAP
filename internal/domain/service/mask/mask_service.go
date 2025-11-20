package mask

import (
	"FairLAP/internal/domain/entity"
	"context"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
)

type Repo interface {
	GetRect(ctx context.Context, detectionId int) (*entity.RectDetection, string, error)
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

var opentypeFont, _ = opentype.Parse(fontTTF)
var colors = []color.RGBA{
	{0, 255, 0, 255},   // зеленый
	{255, 0, 0, 255},   // синий
	{0, 0, 255, 255},   // красный
	{255, 255, 0, 255}, // голубой
	{255, 0, 255, 255}, // пурпурный
	{0, 255, 255, 255}, // желтый
	{128, 0, 128, 255}, // фиолетовый
	{255, 165, 0, 255}, // оранжевый
}

func (s *Service) GetMask(ctx context.Context, detectionId int) (image.Image, error) {
	const op = "service.GetMask"

	rectDetection, class, err := s.repo.GetRect(ctx, detectionId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rectBounds := rectDetection.Rect()

	mask := image.NewRGBA(rectDetection.ImgBounds())

	fontSize := float64(rectDetection.ImgBounds().Dy()) / 40

	face, _ := opentype.NewFace(opentypeFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	thickness := max(int(float64(rectDetection.ImgBounds().Dy())/500), 1)

	c := colors[rectDetection.Id%len(colors)]

	drawer := &font.Drawer{
		Dst:  mask,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	drawRect(mask, c, rectBounds, thickness)

	lbl := fmt.Sprintf("%s: %0.2f", class, rectDetection.Confidence)
	x, y := rectBounds.Min.X+thickness+2, rectBounds.Min.Y+thickness

	drawer.Dot = fixed.P(x, y)
	drawer.DrawString(lbl)

	return mask, nil

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
