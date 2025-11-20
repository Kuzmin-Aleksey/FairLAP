package entity

import "image"

type RectDetection struct {
	Id          int     `json:"id" db:"id"`
	DetectionId int     `json:"detection_id" db:"detection_id"`
	Width       int     `json:"width" db:"width"`
	Height      int     `json:"height" db:"height"`
	X0          int     `json:"x0" db:"x0"`
	Y0          int     `json:"y0" db:"y0"`
	X1          int     `json:"x1" db:"x1"`
	Y1          int     `json:"y1" db:"y1"`
	Confidence  float32 `json:"confidence" db:"confidence"`
}

func (r RectDetection) Rect() image.Rectangle {
	return image.Rect(r.X0, r.Y0, r.X1, r.Y1)
}

func (r RectDetection) ImgBounds() image.Rectangle {
	return image.Rect(0, 0, r.Width, r.Height)
}
