package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"

	"project/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) 	ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
	);
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) Add(ctx context.Context, user models.User) error {
	const query = `INSERT INTO users (name,email, password_hash,role) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	const query = `
			SELECT
				id,
				name,
				email,
				password_hash,
				role,
				created_at
			FROM users
			WHERE email = $1 AND deleted_at IS NULL
`
	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (models.User, error) {
	const query = `
			SELECT
				id,
				name,
				email,
				password_hash,
				role,
				created_at
			FROM users
			WHERE id = $1 AND deleted_at IS NULL
`
	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user models.User) error {
	const query = `UPDATE users SET name = $1, phone = $2 WHERE id = $3 AND deleted_at IS NULL`
	rows, err := r.db.Exec(ctx, query, user.Name, user.Phone, user.ID)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const query2 = `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL` 

	rows, err := r.db.Exec(ctx, query2, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserRepository) CreateOrder(ctx context.Context, userID int, description string, amount float64) error {
	const query = `INSERT INTO orders (user_id, description, amount, status) VALUES ($1, $2, $3, 'pending')`
	_, err := r.db.Exec(ctx, query, userID, description, amount)
	if err != nil {
		return err
	}
	return nil
}


func (r *UserRepository) GetOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error) {
	const query = `
		SELECT 
			id,
			user_id,
			description,
			amount,
			status
		FROM orders
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var ord models.Order
		err := rows.Scan(
		&ord.ID,
		&ord.UserID,
		&ord.Description,
		&ord.Amount,
		&ord.Status)

		if err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	return orders, nil
}

func (r *UserRepository) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	const query = `
		SELECT
			id,
			user_id,
			description,
			amount,
			status
		FROM orders
		WHERE id = $1 AND deleted_at IS NULL
	`
		var ord models.Order
		err := r.db.QueryRow(ctx, query, id).Scan(
		&ord.ID,
		&ord.UserID,
		&ord.Description,
		&ord.Amount,
		&ord.Status)

	if err != nil {
		return models.Order{}, err
	}
	return ord, nil
}

func (r *UserRepository) UpdateOrder(ctx context.Context, o models.Order) error {
	const query = `
		UPDATE orders 
		SET description = $1, amount = $2, status = $3 
		WHERE id = $4 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, o.Description, o.Amount, o.Status, o.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteOrder(ctx context.Context, id int) error {
	const query = `UPDATE orders SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) AdminGetAllUsers(ctx context.Context) ([]models.User, error) {
	const query = `SELECT 
		id,
		name,
		email,
		role 
		FROM users WHERE deleted_at IS NULL`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) AdminGetAllOrders(ctx context.Context) ([]models.Order, error) {
	const query = `SELECT 
		id,
		user_id,
		description,
		amount,
		status
	FROM orders WHERE deleted_at IS NULL`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.UserID, &o.Description, &o.Amount, &o.Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *UserRepository) AdminUpdateRole(ctx context.Context, id int, role string) error {
	const query = `UPDATE users SET role = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, role, id)
	if err != nil {
		return err
	}
	return nil
}
