package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/lapotkin/file-storage/pkg/models"

	sq "github.com/Masterminds/squirrel"
	q "github.com/lapotkin/file-storage/internal/adapter/db/postgres/query"
	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type DocumentRepository struct {
	db sqlx.ExtContext
}

func NewDocumentRepository(db sqlx.ExtContext) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(ctx context.Context, doc *models.Document) error {
	query, args, err := q.NewInsert(ctx).
		Into("documents").
		Columns(
			"user_id",
			"folder_id",
			"name",
			"description",
			"file_path",
			"file_size",
			"mime_type",
			"is_public",
			"created_at",
		).
		Values(
			doc.UserID,
			doc.FolderID,
			doc.Name,
			doc.Description,
			doc.FilePath,
			doc.FileSize,
			doc.MimeType,
			doc.IsPublic,
			doc.CreatedAt,
		).
		ReturningID().
		ToSql()
	if err != nil {
		return err
	}

	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&doc.ID); err != nil {
		return err
	}

	return nil
}

func (r *DocumentRepository) GetByID(ctx context.Context, id uint64) (*models.Document, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			user_id,
			folder_id,
			name,
			description,
			file_path,
			file_size,
			mime_type,
			is_public,
			created_at,
			modified_at
		`).
		From("documents").
		Where(sq.Eq{"id": id}).
		IsDeleted(false).
		ToSql()
	if err != nil {
		return nil, err
	}

	doc := &models.Document{}
	if err = sqlx.GetContext(ctx, r.db, doc, query, args...); err != nil {
		return nil, err
	}

	return doc, nil
}

func (r *DocumentRepository) FetchDocuments(ctx context.Context, filters f.Filters) (models.Documents, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			user_id,
			folder_id,
			name,
			description,
			file_path,
			file_size,
			mime_type,
			is_public,
			created_at,
			modified_at
		`).
		From("documents").
		IsDeleted(false).
		ApplyQueryFilters(filters, f.Wheres{
			{
				Operator: f.And,
				Conditions: f.Conditions{
					{Name: models.UserIDParam, Column: "user_id", Operator: f.Equals},
					{Name: models.FolderIDParam, Column: "folder_id", Operator: f.Equals},
					{Name: models.IsPublicParam, Column: "is_public", Operator: f.Equals},
				},
			},
		}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	documents := make(models.Documents, 0)
	if err = sqlx.SelectContext(ctx, r.db, &documents, query, args...); err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *DocumentRepository) Update(ctx context.Context, doc *models.Document) error {

	query, args, err := q.NewUpdate(ctx).
		Table("documents").
		SetMap(sq.Eq{
			"name":        doc.Name,
			"description": doc.Description,
			"folder_id":   doc.FolderID,
			"is_public":   doc.IsPublic,
			"modified_at": doc.ModifiedAt,
		}).
		Where(sq.Eq{"id": doc.ID}).
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

func (r *DocumentRepository) Delete(ctx context.Context, id uint64) error {
	query, args, err := q.NewUpdate(ctx).
		Table("documents").
		Set(q.ColumnIsDeleted, true).
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
