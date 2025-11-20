package metrics

import (
	"FairLAP/internal/domain/aggregate"
	"FairLAP/internal/domain/entity"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type DetectionsRepo interface {
	GetByGroup(ctx context.Context, group int) ([]entity.Detection, error)
	IsExistProblem(ctx context.Context, lapId int) (bool, error)
}

type GroupsRepo interface {
	GetLaps(ctx context.Context) ([]aggregate.LapLastDetect, error)
}

type Service struct {
	groups     GroupsRepo
	detections DetectionsRepo
}

func NewService(groups GroupsRepo, detections DetectionsRepo) *Service {
	return &Service{
		groups:     groups,
		detections: detections,
	}
}

type LapItem struct {
	HaveProblems bool      `json:"have_problems"`
	LastDetect   time.Time `json:"last_detect"`
}

func (s *Service) GetLaps(ctx context.Context) (map[int]*LapItem, error) {
	const op = "metrics_service.GetLaps"

	laps, err := s.groups.GetLaps(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lapMap := make(map[int]*LapItem)

	for _, lap := range laps {
		lapMap[lap.LapId] = &LapItem{
			LastDetect: lap.LastDetect,
		}
	}

	for lap := range lapMap {
		isHaveProblems, err := s.detections.IsExistProblem(ctx, lap)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		lapMap[lap].HaveProblems = isHaveProblems
	}

	return lapMap, nil
}

type GroupMetric struct {
	ImageCount      int `json:"image_count"`
	DetectionsCount int `json:"detections_count"`
	ProblemCount    int `json:"problem_count"`

	Classes        []string                      `json:"classes"`
	ProblemClasses []string                      `json:"problem_classes"`
	Detections     map[string][]entity.Detection `json:"detections"`
}

func (s *Service) GetGroupMetric(ctx context.Context, groupId int) (*GroupMetric, error) {
	const op = "metrics_service.GetGroupMetric"

	detections, err := s.detections.GetByGroup(ctx, groupId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	imagesMap := make(map[uuid.UUID]struct{})
	classesMap := make(map[string]struct{})
	problemClassesMap := make(map[string]struct{})

	metric := &GroupMetric{
		Detections: make(map[string][]entity.Detection),
	}

	log.Println(len(classesMap))

	for _, detection := range detections {
		imagesMap[detection.ImageUid] = struct{}{}

		if detection.IsProblem {
			metric.ProblemCount++
			problemClassesMap[detection.Class] = struct{}{}
		} else {
			classesMap[detection.Class] = struct{}{}
			metric.DetectionsCount++
		}

		metric.Detections[detection.Class] = append(metric.Detections[detection.Class], detection)
	}

	metric.ImageCount = len(imagesMap)

	metric.Classes = make([]string, 0, len(classesMap))
	metric.ProblemClasses = make([]string, 0, len(problemClassesMap))

	for class := range classesMap {
		metric.Classes = append(metric.Classes, class)
	}
	for class := range problemClassesMap {
		metric.ProblemClasses = append(metric.ProblemClasses, class)
	}

	return metric, nil
}
