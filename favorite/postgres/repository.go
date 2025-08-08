package postgres

import (
	"context"
	"database/sql"

	"github.com/uesleicarvalhoo/aiqfome/favorite"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) favorite.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Find(ctx context.Context, clientID uuid.ID, productID int) (favorite.Favorite, error) {
	query := `
		SELECT
			client_id, product_id, registred_at
		FROM favorites
		WHERE
			client_id = $1
			AND product_id = $2
		`

	var f favorite.Favorite
	if err := r.db.QueryRowContext(ctx, query, clientID, productID).Scan(
		&f.ClientID,
		&f.ProductID,
		&f.RegistredAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return favorite.Favorite{}, &favorite.ErrFavoriteNotFound{
				ClientID:  clientID,
				ProductID: productID,
			}
		}
		return favorite.Favorite{}, err
	}

	return f, nil
}

func (r *repository) PaginateByClientID(ctx context.Context, clientID uuid.ID, page, pageSize int) ([]favorite.Favorite, int, error) {
	query := `
		SELECT
			client_id, product_id, registred_at
		FROM favorites
		WHERE
			client_id = $1
		ORDER BY product_id
		LIMIT $2 OFFSET $3
	`

	queryCount := `
		SELECT count(*) FROM favorites
		WHERE client_id = $1
	`

	var total int
	if err := r.db.QueryRow(queryCount, clientID).Scan(&total); err != nil {
		return []favorite.Favorite{}, 0, err
	}

	offset := page * pageSize

	rows, err := r.db.QueryContext(ctx, query, clientID, pageSize, offset)
	if err != nil {
		return []favorite.Favorite{}, 0, err
	}

	var ff []favorite.Favorite
	for rows.Next() {
		var f favorite.Favorite
		if err := rows.Scan(
			&f.ClientID,
			&f.ProductID,
			&f.RegistredAt,
		); err != nil {
			return []favorite.Favorite{}, 0, err
		}

		ff = append(ff, f)
	}

	return ff, total, nil
}

func (r *repository) Create(ctx context.Context, f favorite.Favorite) error {
	query := `
	INSERT INTO favorites(
		client_id, product_id, registred_at
	) VALUES (
	 $1, $2, $3
	 )
	`

	_, err := r.db.ExecContext(ctx, query, f.ClientID, f.ProductID, f.RegistredAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Remove(ctx context.Context, f favorite.Favorite) error {
	query := `
	DELETE FROM favorites WHERE client_id = $1 and product_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, f.ClientID, f.ProductID)
	if err != nil {
		return err
	}

	return nil
}
