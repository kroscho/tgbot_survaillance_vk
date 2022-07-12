package postgres

import (
	"context"
	userVksvc "tgbot_surveillance/internal/domain/userVk"
	"tgbot_surveillance/pkg/clock"
	"tgbot_surveillance/pkg/database/psql"
	"tgbot_surveillance/pkg/stommer"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type store struct {
	db           psql.DB
	clock        clock.Clock
	sqlBuilder   sq.StatementBuilderType
	tableUsersVk string
}

// nolint:golint
func NewStore(db psql.DB, clock clock.Clock) *store {
	return &store{
		db:           db,
		clock:        clock,
		sqlBuilder:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		tableUsersVk: "usersvk",
	}
}

func (s *store) DB(ctx context.Context) psql.DB {
	db := s.db

	return db
}

func (s store) Create(ctx context.Context, from *userVksvc.UserVk) error {
	var in userVk
	in.unmarshal(from)

	filters := userVksvc.Filters{
		VkIDs: []userVksvc.VkID{from.VkID},
	}

	getOpts := userVksvc.WithFilters(filters)

	res, err := s.Get(ctx, getOpts)
	if err != nil {
		return errors.Wrap(err, "exec query")
	}

	if len(res) == 0 {
		st, err := stommer.New(in, "id_user_vk")
		if err != nil {
			return errors.Wrap(err, "create stommer")
		}

		query, args, err := s.sqlBuilder.Insert(s.tableUsersVk).
			Columns(st.Columns...).
			Values(st.Values...).ToSql()

		if err != nil {
			return errors.Wrap(err, "build query")
		}

		result, err := s.db.ExecContext(ctx, query, args...)
		if err != nil {
			return errors.Wrap(err, "exec query")
		}
		affected, err := result.RowsAffected()
		if affected == 0 {
			return errors.Wrap(err, "Request failed.")
		}
		if err != nil {
			return errors.Wrap(err, "Internal Error")
		}
	}

	return nil
}

func (s store) Get(ctx context.Context, opts ...userVksvc.GetOptFunc) ([]*userVksvc.UserVk, error) {
	getOpts := userVksvc.NewGetOptions(opts...)

	st, err := stommer.New(&userVk{})
	if err != nil {
		return nil, errors.Wrap(err, "create stommer")
	}

	builder := s.sqlBuilder.Select(st.Columns...).From(s.tableUsersVk)
	builder = applyFilters(builder, getOpts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var usersVk []*userVk
	err = s.DB(ctx).SelectContext(ctx, &usersVk, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select disputes")
	}

	result := make([]*userVksvc.UserVk, 0, len(usersVk))
	for _, mm := range usersVk {
		dd, err := mm.marshal()
		if err != nil {
			return nil, errors.Wrap(err, "marshal")
		}
		result = append(result, dd)
	}

	return result, nil

}

func applyFilters(builder sq.SelectBuilder, opts *userVksvc.GetOptions) sq.SelectBuilder {
	if opts == nil {
		return builder
	}

	if len(opts.Filters.VkIDs) > 0 {
		builder = builder.Where(sq.Eq{"vk_id": opts.Filters.VkIDs})
	}

	return builder
}
