package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lapotkin/file-storage/pkg/models"

	sq "github.com/Masterminds/squirrel"
	q "github.com/lapotkin/file-storage/internal/adapter/db/postgres/query"
)

type AuthRepository struct {
	db sqlx.ExtContext
}

func NewAuthRepository(db sqlx.ExtContext) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, token string, currentTime time.Time) (uint64, error) {
	query, args, err := q.NewSelect(ctx).
		Columns("user_id").
		From("refresh_tokens").
		Where(sq.Eq{"token": token}).
		Where(sq.Eq{"revoked": false}).
		Where(sq.Gt{"expires_at": currentTime.Unix()}).
		ToSql()
	if err != nil {
		return 0, err
	}

	var userID uint64
	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, data *models.RefreshTokenData) error {
	query, args, err := q.NewInsert(ctx).
		Into("refresh_tokens").
		Columns(`
			user_id,
			token,
			expires_at,
			created_at
		`).
		Values(
			data.UserID,
			data.Token,
			data.ExpiresAt,
			data.CreatedAt,
		).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	query, args, err := q.NewUpdate(ctx).
		Table("refresh_tokens").
		Set("revoked", true).
		Where(sq.Eq{"token": token}).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetPasswordHash(ctx context.Context, userID uint64) (string, error) {
	query, args, err := q.NewSelect(ctx).
		Columns("password_hash").
		From("user_passwords").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return "", err
	}

	var hash string
	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&hash); err != nil {
		return "", err
	}

	return hash, nil
}

func (r *AuthRepository) CreatePassword(ctx context.Context, data *models.PasswordData) error {
	query, args, err := q.NewInsert(ctx).
		Into("user_passwords").
		Columns(`
			user_id,
			password_hash,
			modified_at
		`).
		Values(
			data.UserID,
			data.Hash,
			data.ModifiedAt).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) UpdatePassword(ctx context.Context, data *models.PasswordData) error {
	query, args, err := q.NewUpdate(ctx).
		Table("user_passwords").
		SetMap(sq.Eq{
			"password_hash": data.Hash,
			"modified_at":   data.ModifiedAt,
		}).
		Where(sq.Eq{"user_id": data.UserID}).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
