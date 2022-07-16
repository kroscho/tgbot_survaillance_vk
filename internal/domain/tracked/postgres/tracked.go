package postgres

import (
	trackedsvc "tgbot_surveillance/internal/domain/tracked"
	usersvc "tgbot_surveillance/internal/domain/user"
	"time"
)

type userTracked struct {
	ID              trackedsvc.ID                `db:"id_user_tracked"`
	UserID          usersvc.ID                   `db:"user_id"`
	TrackedPersonID trackedsvc.ID_TRACKED_PERSON `db:"tracked_id"`
}

type tracked struct {
	ID        trackedsvc.ID   `db:"id_tracked"`
	VkID      trackedsvc.VkID `db:"vk_id"`
	FirstName string          `db:"first_name"`
	LastName  string          `db:"last_name"`
}

type prevFriends struct {
	ID              trackedsvc.ID                `db:"id_prev_friends"`
	TrackedPersonID trackedsvc.ID_TRACKED_PERSON `db:"tracked_id"`
	VkID            trackedsvc.VkID              `db:"vk_id"`
}

type trackedInfo struct {
	ID        trackedsvc.ID   `db:"id_tracked"`
	VkID      trackedsvc.VkID `db:"vk_id"`
	FirstName string          `db:"first_name"`
	LastName  string          `db:"last_name"`
}

type userTrackedInfo struct {
	ID        trackedsvc.ID   `db:"id_tracked"`
	TgID      trackedsvc.ID   `db:"tg_id"`
	VkID      trackedsvc.VkID `db:"vk_id"`
	FirstName string          `db:"first_name"`
	LastName  string          `db:"last_name"`
}

type historyFriends struct {
	ID              trackedsvc.ID                `db:"id_history"`
	TrackedPersonID trackedsvc.ID_TRACKED_PERSON `db:"tracked_id"`
	CreatedAt       time.Time                    `db:"date_of_change"`
}

type historyVk struct {
	VkID      trackedsvc.VkID `db:"vk_id"`
	FirstName string          `db:"first_name"`
	LastName  string          `db:"last_name"`
	CreatedAt time.Time       `db:"date_of_change"`
}

func (t userTracked) marshal() (*trackedsvc.UserTracked, error) {
	return &trackedsvc.UserTracked{
		ID:              trackedsvc.ID(t.ID),
		UserID:          usersvc.ID(t.UserID),
		TrackedPersonID: trackedsvc.ID_TRACKED_PERSON(t.TrackedPersonID),
	}, nil
}

func (t *userTracked) unmarshal(from *trackedsvc.UserTracked) {
	*t = userTracked{
		ID:              trackedsvc.ID(from.ID),
		UserID:          usersvc.ID(from.UserID),
		TrackedPersonID: trackedsvc.ID_TRACKED_PERSON(from.TrackedPersonID),
	}
}

func (t historyFriends) marshal() (*trackedsvc.HistoryFriends, error) {
	return &trackedsvc.HistoryFriends{
		ID:              trackedsvc.ID(t.ID),
		TrackedPersonID: trackedsvc.ID_TRACKED_PERSON(t.TrackedPersonID),
		CreatedAt:       t.CreatedAt,
	}, nil
}

func (t *historyFriends) unmarshal(from *trackedsvc.HistoryFriends) {
	*t = historyFriends{
		ID:              trackedsvc.ID(from.ID),
		TrackedPersonID: trackedsvc.ID_TRACKED_PERSON(from.TrackedPersonID),
		CreatedAt:       from.CreatedAt,
	}
}

func (t historyVk) marshal() (*trackedsvc.HistoryVk, error) {
	return &trackedsvc.HistoryVk{
		VkID:      trackedsvc.VkID(t.VkID),
		FirstName: t.FirstName,
		LastName:  t.LastName,
		CreatedAt: t.CreatedAt,
	}, nil
}

func (t *historyVk) unmarshal(from *trackedsvc.HistoryVk) {
	*t = historyVk{
		VkID:      trackedsvc.VkID(from.VkID),
		FirstName: from.FirstName,
		LastName:  from.LastName,
		CreatedAt: from.CreatedAt,
	}
}

func (t tracked) marshal() (*trackedsvc.Tracked, error) {
	return &trackedsvc.Tracked{
		ID:        trackedsvc.ID(t.ID),
		VkID:      trackedsvc.VkID(t.VkID),
		FirstName: t.FirstName,
		LastName:  t.LastName,
	}, nil
}

func (t *tracked) unmarshal(from *trackedsvc.Tracked) {
	*t = tracked{
		ID:        trackedsvc.ID(from.ID),
		VkID:      trackedsvc.VkID(from.VkID),
		FirstName: from.FirstName,
		LastName:  from.LastName,
	}
}

func (t trackedInfo) marshal() (*trackedsvc.TrackedInfo, error) {
	userVk := trackedsvc.UserVK{
		UID:       trackedsvc.VkID(t.VkID),
		FirstName: string(t.FirstName),
		LastName:  string(t.LastName),
	}
	return &trackedsvc.TrackedInfo{
		ID:     trackedsvc.ID(t.ID),
		UserVK: userVk,
	}, nil
}

func (t *trackedInfo) unmarshal(from *trackedsvc.TrackedInfo) {
	*t = trackedInfo{
		ID:        trackedsvc.ID(from.ID),
		VkID:      trackedsvc.VkID(from.UserVK.UID),
		FirstName: string(from.UserVK.FirstName),
		LastName:  string(from.UserVK.LastName),
	}
}

func (t userTrackedInfo) marshal() (*trackedsvc.UserTrackedInfo, error) {
	return &trackedsvc.UserTrackedInfo{
		ID:        t.ID,
		TgID:      trackedsvc.ID(t.TgID),
		VkID:      trackedsvc.VkID(t.VkID),
		FirstName: t.FirstName,
		LastName:  t.LastName,
	}, nil
}

func (t *userTrackedInfo) unmarshal(from *trackedsvc.UserTrackedInfo) {
	*t = userTrackedInfo{
		ID:        from.ID,
		TgID:      trackedsvc.ID(from.TgID),
		VkID:      trackedsvc.VkID(from.VkID),
		FirstName: string(from.FirstName),
		LastName:  string(from.LastName),
	}
}

func (t prevFriends) marshal() (*trackedsvc.PrevFriends, error) {
	return &trackedsvc.PrevFriends{
		ID:              trackedsvc.ID(t.ID),
		TrackedPersonID: trackedsvc.ID_TRACKED_PERSON(t.TrackedPersonID),
		VkID:            trackedsvc.VkID(t.VkID),
	}, nil
}

func (t *prevFriends) unmarshal(from *trackedsvc.PrevFriends) {
	*t = prevFriends{
		ID:              trackedsvc.ID(from.ID),
		TrackedPersonID: trackedsvc.ID_TRACKED_PERSON(from.TrackedPersonID),
		VkID:            trackedsvc.VkID(from.VkID),
	}
}
