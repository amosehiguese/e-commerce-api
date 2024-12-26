package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/amosehiguese/ecommerce-api/pkg/logger"
	"github.com/amosehiguese/ecommerce-api/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id" validate:"required,uuid4"`
	FirstName    string    `json:"first_name" validate:"required,min=2,max=100"`
	LastName     *string   `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
	Email        string    `json:"email" validate:"required,email,max=255"`
	PasswordHash string    `json:"-" validate:"required,min=8,max=255"`
	Role         string    `json:"role" validate:"required,oneof=user admin"`
	CreatedAt    time.Time `json:"created_at" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at" validate:"required"`
}

func (u *User) ComparePasswordHash(inputPwd string) bool {
	userPassword := utils.NormalizePassword(u.PasswordHash)
	inputPassword := utils.NormalizePassword(inputPwd)

	if err := bcrypt.CompareHashAndPassword(userPassword, inputPassword); err != nil {
		return false
	}
	return true

}

func (q *Query) CreateUser(ctx context.Context, user *User) (*User, error) {
	log := logger.Get()

	query := `
		INSERT INTO "user" (id, first_name, last_name, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := q.DB.ExecContext(ctx, query, user.ID, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Error("Failed to create user in database",
			zap.Error(err),
			zap.String("user_id", user.ID.String()),
		)
		return nil, err
	}

	log.Info("User created successfully",
		zap.String("user_id", user.ID.String()),
	)
	return user, nil
}

func (q *Query) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	log := logger.Get()
	query := `
		SELECT id, first_name, last_name, email, password_hash, role, created_at, updated_at
		FROM "user"
		WHERE email = $1
	`

	// Log the query execution
	log.Info("Executing query to fetch user by email", zap.String("email", email))

	var user User

	row := q.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("No user found for the provided email", zap.String("email", email))
			return nil, nil
		}

		log.Error("Failed to fetch user by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	log.Info("Successfully fetched user by email", zap.String("email", email), zap.String("user_id", user.ID.String()))
	return &user, nil
}
