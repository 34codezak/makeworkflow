package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/34codezak/makeworkflow/internal/models"
	"github.com/google/uuid"
)

var ErrTaskNotFound = errors.New("Task not found")

type TaskRepository interface {
	Create(ctx context.Context, t *models.Task) error

	GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error)

	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.Task, error)

	UpdateStatus(ctx context.Context, id uuid.UUID, status models.TaskStatus) error

	Delete(ctx context.Context, id uuid.UUID) error
}

type taskRepository struct {
	db *sql.DB
}

func newTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, t *models.Task) error {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (id, title, description, status, project_id)
			VALUES ($1, $2, $3, $4, $5)`,
		t.ID, t.Title, t.Description, t.Status, t.ProjectID,
	)
	return err
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	t := &models.Task{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, status, project_id, created_at, updated_at
			FROM tasks WHERE id = $1`, id,
	).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.ProjectID, &t.CreatedAt, &t.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *taskRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, description, status, project_id, created_at, updated_at
			FROM tasks WHERE project_id = $1 ORDER BY created_at DESC`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		t := &models.Task{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.ProjectID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *taskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TaskStatus) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET status = $1, updated_at = NOW() WHERE id = $2`, status, id,
	)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM tasks WHERE id = $1`, id,
	)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrTaskNotFound
	}
	return nil
}
