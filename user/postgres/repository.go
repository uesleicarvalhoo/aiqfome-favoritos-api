package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) user.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Find(ctx context.Context, id uuid.ID) (user.User, error) {
	query := `
		SELECT 
			id, name, email, password_hash, role, active, created_at
		FROM users
		WHERE
			id = $1
		`

	var u user.User
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Active,
		&u.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, user.ErrNotFound
		}

		return user.User{}, err
	}

	return u, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (user.User, error) {
	query := `
		SELECT 
			id, name, email, password_hash, role, active, created_at
		FROM users
		WHERE
			email = $1
		`

	var u user.User
	if err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Active,
		&u.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, user.ErrNotFound
		}

		return user.User{}, err
	}

	return u, nil
}

func (r *repository) Paginate(ctx context.Context, page, pageSize int) ([]user.User, int, error) {
	query := `
		SELECT 
			id, name, email, password_hash, role, active, created_at
		FROM users LIMIT $1 OFFSET $2
		`
	countQuery := `SELECT COUNT(*) FROM users`

	offset := page * pageSize

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return []user.User{}, 0, err
	}

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return []user.User{}, 0, err
	}

	var uu []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.PasswordHash,
			&u.Role,
			&u.Active,
			&u.CreatedAt,
		); err != nil {
			return []user.User{}, 0, err
		}
	}

	if err := rows.Err(); err != nil {
		return []user.User{}, 0, err
	}

	return uu, total, nil
}

func (r *repository) Create(ctx context.Context, c user.User) error {
	query := `
		INSERT INTO users (
			id, name, email, password_hash, role, active, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	_, err := r.db.ExecContext(ctx, query, c.ID, c.Name, c.Email, c.PasswordHash, c.Role, c.Active, c.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, u user.User) error {
	query := `DELETE FROM clientes WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(ctx context.Context, u user.User) error {
	query := `UPDATE users
		SET name = $2, email = $3, password_hash = $4, role = $5, active = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, u.ID, u.Name, u.Email, u.PasswordHash, u.Role, u.Active, time.Now())
	if err != nil {
		return err
	}

	return nil
}
