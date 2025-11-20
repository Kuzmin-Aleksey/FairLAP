package mysql

import (
	"FairLAP/internal/domain/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type DetectionsRepo struct {
	db *sqlx.DB
}

func NewDetectionsRepo(db *sqlx.DB) *DetectionsRepo {
	return &DetectionsRepo{
		db: db,
	}
}

func (r *DetectionsRepo) Save(ctx context.Context, detections []entity.Detection) error {
	const op = "DetectionsRepo.Save"
	query := "INSERT INTO detections (group_id, image_uid, class, is_problem) VALUES"
	args := make([]any, len(detections)*4)

	for i, detection := range detections {
		query += " (?, ?, ?, ?),"
		args[i*4] = detection.GroupId
		args[i*4+1] = detection.ImageUid
		args[i*4+2] = detection.Class
		args[i*4+3] = detection.IsProblem
	}

	query = strings.TrimSuffix(query, ",")

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%w: %s", err, op)
	}

	return nil
}

func (r *DetectionsRepo) GetByGroup(ctx context.Context, group int) ([]entity.Detection, error) {
	const op = "DetectionsRepo.GetByGroup"
	var detections []entity.Detection

	if err := r.db.SelectContext(ctx, &detections, "SELECT * FROM detections WHERE group_id=?", group); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", err, op)
		}
	}

	return detections, nil
}

func (r *DetectionsRepo) IsExistProblem(ctx context.Context, lapId int) (bool, error) {
	const op = "DetectionsRepo.IsExistProblem"

	query := "SELECT EXISTS(SELECT * FROM detections INNER JOIN `groups` ON detections.group_id = `groups`.id WHERE detections.is_problem AND `groups`.lap_id=?)"
	var exist bool
	if err := r.db.QueryRowContext(ctx, query, lapId).Scan(&exist); err != nil {
		return false, fmt.Errorf("%w: %s", err, op)
	}
	return exist, nil
}
