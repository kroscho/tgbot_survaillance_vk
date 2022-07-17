package postgres

import (
	"context"
	"fmt"
	"strings"
	"tgbot_surveillance/config"
	usersvc "tgbot_surveillance/internal/domain/user"
	"tgbot_surveillance/pkg/clock"
	"tgbot_surveillance/pkg/database/psql"
	Encrypt "tgbot_surveillance/pkg/encrypt"
	"tgbot_surveillance/pkg/stommer"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type store struct {
	db                 psql.DB
	clock              clock.Clock
	sqlBuilder         sq.StatementBuilderType
	cfg                *config.Config
	tableUsers         string
	tableSubscribes    string
	tableUserSubscribe string
	tableUserTracked   string
}

// nolint:golint
func NewStore(db psql.DB, clock clock.Clock, cfg *config.Config) *store {
	return &store{
		db:                 db,
		clock:              clock,
		sqlBuilder:         sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		cfg:                cfg,
		tableUsers:         "users",
		tableSubscribes:    "subscribes",
		tableUserSubscribe: "usersubscribe",
		tableUserTracked:   "usertracked",
	}
}

func (s *store) DB(ctx context.Context) psql.DB {
	db := s.db

	return db
}

func (s store) Create(ctx context.Context, from *usersvc.User) (*usersvc.User, error) {
	var in user
	in.unmarshal(from)

	//st, err := stommer.New(in, "id_user")
	//if err != nil {
	//	return nil, errors.Wrap(err, "create stommer")
	//}

	//[created_at enabled last_payment plan price tg_id user_name user_token vk_id]
	//returning := st.Columns
	//returning = append(returning, "id_user")
	returning := []string{"tg_id", "user_name", "created_at"}
	returning1 := append(returning, "id_user")

	query, args, err := s.sqlBuilder.Insert(s.tableUsers).Suffix(fmt.Sprintf("RETURNING %s", strings.Join(returning1, ", "))).
		Columns(returning...).
		Values(from.TgID, from.Username, from.CreatedAt).ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	err = s.db.GetContext(ctx, &in, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "exec query")
	}

	return in.marshal()
}

func (s store) Get(ctx context.Context, opts ...usersvc.GetOptFunc) ([]*usersvc.User, error) {
	//getOpts := usersvc.NewGetOptions(opts...)

	st, err := stommer.New(&user{})
	if err != nil {
		return nil, errors.Wrap(err, "create stommer")
	}
	builder := s.sqlBuilder.Select(st.Columns...).From(s.tableUsers)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var users []*user
	err = s.DB(ctx).SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select disputes")
	}

	result := make([]*usersvc.User, 0, len(users))
	for _, mm := range users {
		dd, err := mm.marshal()
		if err != nil {
			return nil, errors.Wrap(err, "marshal")
		}
		result = append(result, dd)
	}

	return result, nil
}

func (s store) Update(ctx context.Context, from *usersvc.User) error {
	var in user
	in.unmarshal(from)

	st, err := stommer.New(in, "id_user")
	if err != nil {
		return errors.Wrap(err, "create stommer")
	}
	setMap := make(map[string]interface{}, 0)

	for i := 0; i < len(st.Columns); i++ {
		setMap[st.Columns[i]] = st.Values[i]
	}

	builder := s.sqlBuilder.Update(s.tableUsers).
		SetMap(setMap).
		Where(sq.Eq{"tg_id": from.TgID})

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	result, err := s.DB(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "update disputes")
	}
	affected, err := result.RowsAffected()
	if affected == 0 {
		return errors.Wrap(err, "Request failed. User may not exist, or request had no updates.")
	}
	if err != nil {
		return errors.Wrap(err, "Internal Error")
	}

	return nil
}

func (s store) GetUserByTgID(ctx context.Context, tgID usersvc.TelegramID) (*usersvc.User, error) {
	st, err := stommer.New(&user{})
	if err != nil {
		return nil, errors.Wrap(err, "create stommer")
	}

	//builder := s.sqlBuilder.Select(st.Columns...).
	//	From(s.tableUsers).
	//	Where(sq.Eq{"tg_id": tgID})

	builder := s.sqlBuilder.Select(st.Columns...).
		From(s.tableUsers + " as u").
		LeftJoin(s.tableUserSubscribe + " as u_sub on u.id_user = u_sub.user_id").
		LeftJoin(s.tableSubscribes + " as s on u_sub.subscribe_id = s.id_subscribe").
		Where(sq.Eq{"tg_id": tgID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var users []*user
	err = s.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select disputes")
	}

	if len(users) == 0 {
		return nil, nil
	}
	result := make([]*usersvc.User, 0, len(users))
	for _, mm := range users {
		dd, err := mm.marshal()
		if err != nil {
			return nil, errors.Wrap(err, "marshal")
		}

		if dd.Token != nil {
			decToken, err := Encrypt.Decrypt(*dd.Token, s.cfg.Secret)
			if err != nil {
				fmt.Println("error decrypting your encrypted text: ", err)
			}
			dd.Token = &decToken
		}
		result = append(result, dd)
	}

	return result[0], nil
}

func (s store) Delete(ctx context.Context, user *usersvc.User) error {

	query := fmt.Sprintf("delete from %s where user_id=%d", s.tableUserTracked, user.ID)

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "exec query")
	}

	query = fmt.Sprintf("delete from %s where id_user=%d", s.tableUsers, user.ID)

	result, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "exec query")
	}

	affected, err := result.RowsAffected()
	if affected == 0 {
		return errors.Wrap(err, "User not found")
	}
	if err != nil {
		return errors.Wrap(err, "Internal Error")
	}

	return nil
}
