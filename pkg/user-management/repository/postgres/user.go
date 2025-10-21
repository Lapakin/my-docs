package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/lapotkin/file-storage/pkg/models"

	sq "github.com/Masterminds/squirrel"
	q "github.com/lapotkin/file-storage/internal/adapter/db/postgres/query"
	f "github.com/lapotkin/file-storage/internal/app/filter"
)

type UserRepository struct {
	db sqlx.ExtContext
}

func NewUserRepository(db sqlx.ExtContext) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query, args, err := q.NewInsert(ctx).
		Into("users").
		Columns(`
			email,
			username,
			first_name,
			last_name,
			role,
			is_active,
			is_deleted,
			created_at
		`).
		Values(
			user.Email,
			user.Username,
			user.FirstName,
			user.LastName,
			user.Role,
			user.IsActive,
			user.IsDeleted,
			user.CreatedAt,
		).
		ReturningID().
		ToSql()
	if err != nil {
		return err
	}

	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			email,
			username,
			first_name,
			last_name,
			role,
			is_active,
			is_deleted,
			created_at,
			modified_at
		`).
		From("users").
		WhereID(id).
		IsDeleted(false).
		ToSql()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	if err = sqlx.GetContext(ctx, r.db, user, query, args...); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			email,
			username,
			first_name,
			last_name,
			role,
			is_active,
			is_deleted,
			created_at,
			modified_at
		`).
		From("users").
		Where(sq.Eq{"email": email}).
		IsDeleted(false).
		ToSql()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	if err = sqlx.GetContext(ctx, r.db, user, query, args...); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			email,
			username,
			first_name,
			last_name,
			role,
			is_active,
			is_deleted,
			created_at,
			modified_at
		`).
		From("users").
		Where(sq.Or{sq.Eq{"email": login}, sq.Eq{"username": login}}).
		IsDeleted(false).
		ToSql()
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	if err = sqlx.GetContext(ctx, r.db, user, query, args...); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FetchUsers(ctx context.Context, filters f.Filters) (models.Users, error) {
	query, args, err := q.NewSelect(ctx).
		Columns(`
			id,
			email,
			username,
			first_name,
			last_name,
			role,
			is_active,
			is_deleted,
			created_at,
			modified_at
		`).
		From("users").
		IsDeleted(false).
		ApplyQueryFilters(filters, f.Wheres{
			{
				Operator: f.And,
				Conditions: f.Conditions{
					{Name: models.IDsParam, Column: "id", Operator: f.In},
					{Name: models.EmailParam, Column: "email", Operator: f.Equals},
				},
			},
		}).
		ToSql()

	users := make(models.Users, 0)
	if err = sqlx.SelectContext(ctx, r.db, &users, query, args...); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query, args, err := q.NewUpdate(ctx).
		Table("users").
		SetMap(sq.Eq{
			"email":       user.Email,
			"username":    user.Username,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"role":        user.Role,
			"is_active":   user.IsActive,
			"modified_at": user.ModifiedAt,
		}).
		WhereID(user.ID).
		Returning(q.ColumnCreatedAt).
		ToSql()
	if err != nil {
		return err
	}

	if err = r.db.QueryRowxContext(ctx, query, args...).Scan(&user.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uint64) error {
	query, args, err := q.NewSoftDelete(ctx).
		Table("users").
		WhereID(id).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) HardDeleteUser(ctx context.Context, id uint64) error {
	query, args, err := q.NewHardDelete(ctx).
		From("users").
		WhereID(id).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeactivateUser(ctx context.Context, id uint64) error {
	query, args, err := q.NewUpdate(ctx).
		Table("users").
		Set("is_active", false).
		WhereID(id).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) ActivateUser(ctx context.Context, id uint64) error {
	query, args, err := q.NewUpdate(ctx).
		Table("users").
		Set("is_active", true).
		WhereID(id).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
