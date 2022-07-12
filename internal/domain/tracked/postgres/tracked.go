package postgres

import (
	trackedsvc "tgbot_surveillance/internal/domain/tracked"
	usersvc "tgbot_surveillance/internal/domain/user"
)

type userTracked struct {
	ID              trackedsvc.ID                `db:"id_user_tracked"`
	UserID          usersvc.ID                   `db:"user_id"`
	TrackedPersonID trackedsvc.ID_TRACKED_PERSON `db:"tracked_id"`
}

type tracked struct {
	ID   trackedsvc.ID   `db:"id_tracked"`
	VkID trackedsvc.VkID `db:"vk_id"`
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

func (t tracked) marshal() (*trackedsvc.Tracked, error) {
	return &trackedsvc.Tracked{
		ID:   trackedsvc.ID(t.ID),
		VkID: trackedsvc.VkID(t.VkID),
	}, nil
}

func (t *tracked) unmarshal(from *trackedsvc.Tracked) {
	*t = tracked{
		ID:   trackedsvc.ID(from.ID),
		VkID: trackedsvc.VkID(from.VkID),
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
