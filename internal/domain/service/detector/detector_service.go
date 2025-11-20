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
	Save(ctx context.Context, detections []entity.Detection) error
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

	renderedImg := drawDetected(img, modelsDetections)

	imgUid, err := s.images.Save(groupId, renderedImg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	detections := make([]entity.Detection, len(modelsDetections))

	for i, detection := range modelsDetections {
		detections[i] = entity.Detection{
			GroupId:  groupId,
			ImageUid: imgUid,
			Class:    detection.ClassName,
		}
	}

	if err := s.repo.Save(ctx, detections); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
