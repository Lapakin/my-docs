package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lapotkin/file-storage/pkg/models"

	sq "github.com/Masterminds/squirrel"
	q "github.com/lapotkin/file-storage/internal/adapter/db/postgres/query"
	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type ShareRepository struct {
	db sqlx.ExtContext
}

func NewShareRepository(db sqlx.ExtContext) *ShareRepository {
	return &ShareRepository{db: db}
}

func (r *ShareRepository) Create(ctx context.Context, share *models.Share) error {
	query, args, err := q.NewInsert(ctx).
		Into("shares").
		Columns(`
			document_id,
			owner_id,
			shared_with,
			share_link,
			permission,
			expires_at,
			access_count,
			max_access,
			password,
			created_at
		`).
		Values(
			share.DocumentID,
			share.OwnerID,
			share.SharedWith,
			share.ShareLink,
			share.Permission,
			share.ExpiresAt,
			share.AccessCount,
			share.MaxAccess,
			share.Password,
			share.CreatedAt,
		).
		ReturningID().
		ToSql()
	if err != nil {
		return err
	}

	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&share.ID); err != nil {
		return err
	}

	return nil
}

func (r *ShareRepository) GetByID(ctx context.Context, id uint64) (*models.Share, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			document_id,
			owner_id,
			shared_with,
			share_link,
			permission,
			expires_at,
			access_count,
			max_access,
			password,
			created_at,
			modified_at
		`).
		From("shares").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	share := &models.Share{}
	if err = sqlx.GetContext(ctx, r.db, share, query, args...); err != nil {
		return nil, err
	}

	return share, nil
}

func (r *ShareRepository) GetByLink(ctx context.Context, link string) (*models.Share, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			document_id,
			owner_id,
			shared_with,
			share_link,
			permission,
			expires_at,
			access_count,
			max_access,
			password,
			created_at,
			modified_at
		`).
		From("shares").
		Where(sq.Eq{"share_link": link}).
		ToSql()
	if err != nil {
		return nil, err
	}

	share := &models.Share{}
	if err = sqlx.GetContext(ctx, r.db, share, query, args...); err != nil {
		return nil, err
	}

	return share, nil
}

func (r *ShareRepository) FetchShares(ctx context.Context, filters f.Filters) (models.Shares, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			document_id,
			owner_id,
			shared_with,
			share_link,
			permission,
			expires_at,
			access_count,
			max_access,
			password,
			created_at,
			modified_at
		`).
		From("shares").
		ApplyQueryFilters(filters, f.Wheres{
			{
				Operator: f.And,
				Conditions: f.Conditions{
					{Name: models.OwnerIDParam, Column: "owner_id", Operator: f.Equals},
					{Name: models.DocumentIDParam, Column: "document_id", Operator: f.Equals},
				},
			},
		}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	shares := make(models.Shares, 0)
	if err = sqlx.SelectContext(ctx, r.db, &shares, query, args...); err != nil {
		return nil, err
	}

	return shares, nil
}

func (r *ShareRepository) Update(ctx context.Context, share *models.Share) error {
	query, args, err := q.NewUpdate(ctx).
		Table("shares").
		SetMap(sq.Eq{
			"permission":  share.Permission,
			"expires_at":  share.ExpiresAt,
			"max_access":  share.MaxAccess,
			"password":    share.Password,
			"modified_at": share.ModifiedAt,
		}).
		Where(sq.Eq{"id": share.ID}).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *ShareRepository) Delete(ctx context.Context, id uint64) error {
	query, args, err := q.NewHardDelete(ctx).
		From("shares").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *ShareRepository) IncrementAccessCount(ctx context.Context, id uint64) error {
	query, args, err := q.NewUpdate(ctx).
		Table("shares").
		Set("access_count", sq.Expr("access_count + 1")).
		Set("modified_at", time.Now()).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
