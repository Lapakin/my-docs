package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/lapotkin/file-storage/pkg/models"

	sq "github.com/Masterminds/squirrel"
	q "github.com/lapotkin/file-storage/internal/adapter/db/postgres/query"
	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type FolderRepository struct {
	db sqlx.ExtContext
}

func NewFolderRepository(db sqlx.ExtContext) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) Create(ctx context.Context, folder *models.Folder) error {
	query, args, err := q.NewInsert(ctx).
		Into("folders").
		Columns(`
			user_id,
			parent_id,
			name,
			path,
			is_public,
			color,
			icon,
			created_at
		`).
		Values(
			folder.UserID,
			folder.ParentID,
			folder.Name,
			folder.Path,
			folder.IsPublic,
			folder.Color,
			folder.Icon,
			folder.CreatedAt,
		).
		ReturningID().
		ToSql()
	if err != nil {
		return err
	}

	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&folder.ID); err != nil {
		return err
	}

	return nil
}

func (r *FolderRepository) GetByID(ctx context.Context, id uint64) (*models.Folder, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			user_id,
			parent_id,
			name,
			path,
			is_public,
			color,
			icon,
			created_at,
			modified_at
		`).
		From("folders").
		Where(sq.Eq{"id": id}).
		IsDeleted(false).
		ToSql()
	if err != nil {
		return nil, err
	}

	folder := &models.Folder{}
	if err = sqlx.GetContext(ctx, r.db, folder, query, args...); err != nil {
		return nil, err
	}

	return folder, nil
}

func (r *FolderRepository) FetchFolders(ctx context.Context, filters f.Filters) (models.Folders, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			user_id,
			parent_id,
			name,
			path,
			is_public,
			color,
			icon,
			created_at,
			modified_at
		`).
		From("folders").
		IsDeleted(false).
		ApplyQueryFilters(filters, f.Wheres{
			{
				Operator: f.And,
				Conditions: f.Conditions{
					{Name: models.UserIDParam, Column: "user_id", Operator: f.Equals},
					{Name: models.ParentIDParam, Column: "parent_id", Operator: f.Equals},
					{Name: models.IsPublicParam, Column: "is_public", Operator: f.Equals},
				},
			},
		}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	folders := make(models.Folders, 0)
	if err = sqlx.SelectContext(ctx, r.db, &folders, query, args...); err != nil {
		return nil, err
	}

	return folders, nil
}

func (r *FolderRepository) Update(ctx context.Context, folder *models.Folder) error {

	query, args, err := q.NewUpdate(ctx).
		Table("folders").
		SetMap(sq.Eq{
			"name":        folder.Name,
			"is_public":   folder.IsPublic,
			"color":       folder.Color,
			"icon":        folder.Icon,
			"modified_at": folder.ModifiedAt,
		}).
		Where(sq.Eq{"id": folder.ID}).
		IsDeleted(false).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *FolderRepository) Delete(ctx context.Context, id uint64) error {
	query, args, err := q.NewSoftDelete(ctx).
		Table("folders").
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

func (r *FolderRepository) GetPath(ctx context.Context, id uint64) (string, error) {
	query, args, err := q.NewSelect(ctx).
		Columns("path").
		From("folders").
		Where(sq.Eq{"id": id}).
		IsDeleted(false).
		ToSql()
	if err != nil {
		return "", err
	}

	var path string
	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&path); err != nil {
		return "", err
	}

	return path, nil
}
