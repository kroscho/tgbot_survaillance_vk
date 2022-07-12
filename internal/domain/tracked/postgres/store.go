package postgres

import (
	"context"
	"fmt"
	"strings"

	trackedsvc "tgbot_surveillance/internal/domain/tracked"
	"tgbot_surveillance/internal/domain/user"
	"tgbot_surveillance/pkg/clock"
	"tgbot_surveillance/pkg/database/psql"
	govk "tgbot_surveillance/pkg/go-vk"
	vkmodels "tgbot_surveillance/pkg/go-vk/models"
	"tgbot_surveillance/pkg/stommer"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type store struct {
	db               psql.DB
	clock            clock.Clock
	sqlBuilder       sq.StatementBuilderType
	tableUserTracked string
	tableTracked     string
	tablePrevFriends string
}

// nolint:golint
func NewStore(db psql.DB, clock clock.Clock) *store {
	return &store{
		db:               db,
		clock:            clock,
		sqlBuilder:       sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		tableUserTracked: "usertracked",
		tableTracked:     "trackeds",
		tablePrevFriends: "prevfriends",
	}
}

func (s *store) DB(ctx context.Context) psql.DB {
	db := s.db

	return db
}

func (s store) Get(ctx context.Context, user *user.User) ([]*trackedsvc.TrackedInfo, error) {

	builder := s.sqlBuilder.Select("id_tracked", "vk_id").
		From(s.tableUserTracked + " as ut").
		InnerJoin(s.tableTracked + " as t on ut.tracked_id = t.id_tracked").
		Where(sq.Eq{"user_id": user.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var trackeds []*trackedInfo
	err = s.db.SelectContext(ctx, &trackeds, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select disputes")
	}

	if len(trackeds) == 0 {
		return nil, nil
	}
	result := make([]*trackedsvc.TrackedInfo, 0, len(trackeds))
	for _, tt := range trackeds {
		dd, err := tt.marshal()
		if err != nil {
			return nil, errors.Wrap(err, "marshal")
		}

		apiVk, _ := govk.NewApiClient(*user.Token)
		params := govk.UserGetParams{
			UserIDS: int64(tt.VkID),
			Fields:  "id, first_name, last_name",
		}
		res, err := apiVk.UserGet(params)
		if err != nil {
			return nil, errors.Wrap(err, "api vk")
		}
		dd.UserVK.FirstName = res.FirstName
		dd.UserVK.LastName = res.LastName
		result = append(result, dd)
	}
	return result, nil
}

func (s store) GetTrackedByVkID(ctx context.Context, user *user.User, vkId int64) (*trackedsvc.TrackedInfo, error) {

	builder := s.sqlBuilder.Select("id_tracked", "vk_id").
		From(s.tableTracked).
		Where(sq.Eq{"vk_id": vkId})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var tracked []*trackedInfo
	err = s.db.SelectContext(ctx, &tracked, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select disputes")
	}

	if tracked == nil {
		return nil, nil
	}

	tr := tracked[0]

	result, err := tr.marshal()
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	apiVk, _ := govk.NewApiClient(*user.Token)
	params := govk.UserGetParams{
		UserIDS: int64(tr.VkID),
		Fields:  "id, first_name, last_name",
	}
	res, err := apiVk.UserGet(params)
	if err != nil {
		return nil, errors.Wrap(err, "api vk")
	}
	result.UserVK.FirstName = res.FirstName
	result.UserVK.LastName = res.LastName

	return result, nil
}

func (s store) Create(ctx context.Context, user *user.User, trackedAdd *vkmodels.User) error {
	trackeds, err := s.Get(ctx, user)
	if err != nil {
		return errors.Wrap(err, "callback, api vk")
	}
	isExist := false
	for _, tt := range trackeds {
		if tt.UserVK.UID == trackedsvc.VkID(trackedAdd.UID) {
			isExist = true
			break
		}
	}
	if !isExist {
		fromTracked := trackedsvc.Tracked{
			VkID: trackedsvc.VkID(trackedAdd.UID),
		}

		var tr tracked
		tr.unmarshal(&fromTracked)

		st, err := stommer.New(tr, "id_tracked")
		if err != nil {
			return errors.Wrap(err, "create stommer")
		}

		returning := st.Columns
		returning = append(returning, "id_tracked")

		query, args, err := s.sqlBuilder.Insert(s.tableTracked).Suffix(fmt.Sprintf("RETURNING %s", strings.Join(returning, ", "))).
			Columns(st.Columns...).
			Values(st.Values...).ToSql()

		if err != nil {
			return errors.Wrap(err, "build query")
		}

		err = s.db.GetContext(ctx, &tr, query, args...)
		if err != nil {
			return errors.Wrap(err, "exec query")
		}

		fromUserTracked := trackedsvc.UserTracked{
			UserID:          user.ID,
			TrackedPersonID: trackedsvc.ID_TRACKED_PERSON(tr.ID),
		}

		var user_tracked userTracked
		user_tracked.unmarshal(&fromUserTracked)

		st, err = stommer.New(user_tracked, "id_user_tracked")
		if err != nil {
			return errors.Wrap(err, "create stommer")
		}

		query, args, err = s.sqlBuilder.Insert(s.tableUserTracked).
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

		trackedInfo := trackedsvc.TrackedInfo{
			ID: tr.ID,
			UserVK: trackedsvc.UserVK{
				UID:       trackedsvc.VkID(trackedAdd.UID),
				FirstName: trackedAdd.FirstName,
				LastName:  trackedAdd.LastName,
			},
		}

		apiVk, _ := govk.NewApiClient(*user.Token)
		params := govk.FriendsGetParams{
			UserID: trackedAdd.UID,
			Fields: "id, first_name, last_name",
		}
		res, err := apiVk.FriendsGet(params)
		if err != nil {
			return errors.Wrap(err, "api vk")
		}

		err = s.UpdatePrevFriends(ctx, &trackedInfo, res)
		if err != nil {
			return err
		}

		return nil
	} else {
		return trackedsvc.ErrTrackedAlreadyExist
	}
}

func (s store) GetPrevFriends(ctx context.Context, tracked *trackedsvc.TrackedInfo) ([]int64, error) {

	builder := s.sqlBuilder.Select("tracked_id", "vk_id").
		From(s.tablePrevFriends).
		Where(sq.Eq{"tracked_id": tracked.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var friends []*prevFriends
	err = s.db.SelectContext(ctx, &friends, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select disputes")
	}

	if len(friends) == 0 {
		return nil, nil
	}
	result := make([]int64, 0, len(friends))
	for _, tt := range friends {
		dd, err := tt.marshal()
		if err != nil {
			return nil, errors.Wrap(err, "marshal")
		}
		result = append(result, int64(dd.VkID))
	}
	return result, nil
}

func (s store) UpdatePrevFriends(ctx context.Context, tracked *trackedsvc.TrackedInfo, newFriends map[int64]vkmodels.User) error {

	for _, friend := range newFriends {
		s.AddUserInPrevFriends(ctx, &friend, tracked)
	}
	return nil
}

func (s store) DeleteUserFromPrevFriends(ctx context.Context, deleteUser *vkmodels.User, tracked *trackedsvc.TrackedInfo) error {

	/*query, args, err := sq.Delete(s.tablePrevFriends).
	Where(sq.Eq{"vk_id": deleteUser.UID}).
	Where(sq.Eq{"tracked_id": tracked.ID}).
	ToSql()*/

	//if err != nil {
	//	return errors.Wrap(err, "Internal Error")
	//}

	query := fmt.Sprintf("delete from %s where vk_id=%d and tracked_id=%d", s.tablePrevFriends, deleteUser.UID, tracked.ID)

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

func (s store) AddUserInPrevFriends(ctx context.Context, addUser *vkmodels.User, tracked *trackedsvc.TrackedInfo) error {

	fromPrevFriends := trackedsvc.PrevFriends{
		TrackedPersonID: trackedsvc.ID_TRACKED_PERSON(tracked.ID),
		VkID:            trackedsvc.VkID(addUser.UID),
	}

	var prev_friends prevFriends
	prev_friends.unmarshal(&fromPrevFriends)

	st, err := stommer.New(prev_friends, "id_prev_friends")
	if err != nil {
		return errors.Wrap(err, "create stommer")
	}

	query, args, err := s.sqlBuilder.Insert(s.tablePrevFriends).
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

	return nil
}

func (s store) CheckUsersAboutTracked(ctx context.Context, user *user.User, tracked *trackedsvc.TrackedInfo) (bool, error) {

	builder := s.sqlBuilder.Select("user_id").
		From(s.tableUserTracked).
		Where(sq.Eq{"tracked_id": tracked.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return true, errors.Wrap(err, "build query")
	}

	var userIds []int
	err = s.db.SelectContext(ctx, &userIds, query, args...)
	if err != nil {
		return true, errors.Wrap(err, "select disputes")
	}

	if len(userIds) == 0 {
		return false, nil
	}
	return true, nil
}

func (s store) DeleteUserFromTracked(ctx context.Context, user *user.User, tracked *trackedsvc.TrackedInfo) error {

	/*query, args, err := sq.Delete(s.tablePrevFriends).
	Where(sq.Eq{"vk_id": deleteUser.UID}).
	Where(sq.Eq{"tracked_id": tracked.ID}).
	ToSql()*/

	//if err != nil {
	//	return errors.Wrap(err, "Internal Error")
	//}

	query := fmt.Sprintf("delete from %s where user_id=%d and tracked_id=%d", s.tableUserTracked, user.ID, tracked.ID)

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

	isExistUsers, err := s.CheckUsersAboutTracked(ctx, user, tracked)
	if err != nil {
		return errors.Wrap(err, "check users")
	}

	if !isExistUsers {
		query = fmt.Sprintf("delete from %s where tracked_id=%d", s.tablePrevFriends, tracked.ID)

		result, err = s.db.ExecContext(ctx, query)
		if err != nil {
			return errors.Wrap(err, "exec query")
		}

		affected, err = result.RowsAffected()
		if affected == 0 {
			return errors.Wrap(err, "User not found")
		}
		if err != nil {
			return errors.Wrap(err, "Internal Error")
		}

		query = fmt.Sprintf("delete from %s where vk_id=%d", s.tableTracked, tracked.UserVK.UID)

		result, err = s.db.ExecContext(ctx, query)
		if err != nil {
			return errors.Wrap(err, "exec query")
		}

		affected, err = result.RowsAffected()
		if affected == 0 {
			return errors.Wrap(err, "User not found")
		}
		if err != nil {
			return errors.Wrap(err, "Internal Error")
		}
	}

	return nil
}