/*
repositories/ is the data access layer.
These are the only files that run SQL queries (via database/sql or pgx).
They return models and hide database details from the services.
*/

package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/34codezak/makeworkflow/internal/models"
)

var ErrProjectNotFound = errors.New("project not found")

type ProjectRepository interface {
	Create(ctx context.Context, p *models.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	List(ctx context.Context) ([]*models.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, p *models.Project) error {
	p.ID = uuid.New()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO projects (id, name, description) VALUES ($1, $2, $3)`,
		p.ID, p.Name, p.Description,
	)
	return err
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	p := &models.Project{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, created_at, updated_at FROM projects WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProjectNotFound
	}
	return p, err
}

func (r *projectRepository) List(ctx context.Context) ([]*models.Project, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, description, created_at, updated_at FROM projects ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		p := &models.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM projects WHERE id = $1`, id,
	)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return ErrProjectNotFound
	}
	return nil
}   