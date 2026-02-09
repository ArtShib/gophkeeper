package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/ArtShib/gophkeeper/internal/storage"
)

type UserStorage struct {
	store storage.GenericStorage[models.User]
}

func New(db *sql.DB, driver models.DriverType) *UserStorage {
	return &UserStorage{
		store: storage.GenericStorage[models.User]{
			DB:     db,
			Table:  "users",
			Driver: driver,
		},
	}
}

// AddUser регистрация пользователя
func (u *UserStorage) AddUser(ctx context.Context, user *models.User) (int64, error) {
	const op = "storage.user.AddUser"

	var id int64

	switch u.store.Driver {
	case models.DriverPostgres:
		query := `
				INSERT INTO users (login, pass_hash, created_at) 
				VALUES ($1, $2, $3) RETURNING id`
		if err := u.store.DB.QueryRowContext(ctx, query, user.Login, user.PasswordHash, time.Now().Unix()).Scan(&id); err != nil {
			return 0, u.store.HandleError(err, op)
		}
	case models.DriverSQLite:
		query := `
				INSERT INTO users (user_id, login, password_hash) 
				VALUES ($1, $2, $3)
				ON CONFLICT(user_id) DO NOTHING`

		if err := u.store.DB.QueryRowContext(ctx, query, user.ID, user.Login, user.PasswordHash).Scan(&id); err != nil {
			return 0, u.store.HandleError(err, op)
		}
	default:
		return 0, u.store.HandleError(errors.New(op), op)
	}

	return id, nil
}

// GetUser select пользователя по логину
func (u *UserStorage) GetUser(ctx context.Context, login string) (*models.User, error) {
	const op = "storage.user.GetUser"

	query := `SELECT id, login, pass_hash FROM users WHERE login = $1`
	//query = u.store.DB.Rebind(query) // на случай sqlx
	row := u.store.DB.QueryRowContext(ctx, query, login)

	var user models.User
	if err := row.Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, u.store.HandleError(err, op)
			//return &models.User{}, fmt.Errorf("%s: %w", op, models.ErrNotFound)
		}
		return &models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
