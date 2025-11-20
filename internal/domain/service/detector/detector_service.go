package detector

import (
	"FairLAP/internal/domain/entity"
	"FairLAP/pkg/yolo_model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"image"
)

type Repo interface {
	Save(ctx context.Context, detections *entity.Detection) error
	SaveRects(ctx context.Context, rects []entity.RectDetection) error
}

type ImageRepo interface {
	Save(groupId int, img image.Image) (uuid.UUID, error)
}

type Service struct {
	model  *yolo_model.Model
	repo   Repo
	images ImageRepo
}

func NewService(model *yolo_model.Model, repo Repo, images ImageRepo) *Service {
	return &Service{
		model:  model,
		repo:   repo,
		images: images,
	}
}

func (s *Service) Detect(ctx context.Context, groupId int, img image.Image) error {
	const op = "detector_service.Detect"

	modelsDetections, err := s.model.Detect(img)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	imgUid, err := s.images.Save(groupId, img)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rects := make([]entity.RectDetection, len(modelsDetections))

	for i, detection := range modelsDetections {
		d := &entity.Detection{
			GroupId:  groupId,
			ImageUid: imgUid,
			Class:    detection.ClassName,
		}
		if err := s.repo.Save(ctx, d); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		rects[i] = entity.RectDetection{
			DetectionId: d.Id,
			Width:       img.Bounds().Dx(),
			Height:      img.Bounds().Dy(),
			X0:          detection.BBox.Min.X,
			Y0:          detection.BBox.Min.Y,
			X1:          detection.BBox.Max.X,
			Y1:          detection.BBox.Max.Y,
			Confidence:  detection.Confidence,
		}
	}

	if err := s.repo.SaveRects(ctx, rects); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
