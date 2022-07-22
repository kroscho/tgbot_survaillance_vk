package tracked

import (
	"tgbot_surveillance/internal/domain/user"
	"time"

	"github.com/pkg/errors"
)

type ID int64

type ID_TRACKED_PERSON int64

type VkID int64

type Tracked struct {
	ID        ID
	VkID      VkID
	FirstName string
	LastName  string
}

type UserTracked struct {
	ID              ID
	UserID          user.ID
	TrackedPersonID ID_TRACKED_PERSON
}

type UserTrackedInfo struct {
	ID        ID
	TgID      ID
	VkID      VkID
	FirstName string
	LastName  string
}

type TrackedInfo struct {
	ID     ID
	UserVK UserVK
}

type PrevFriends struct {
	ID              ID
	TrackedPersonID ID_TRACKED_PERSON
	VkID            VkID
}

type HistoryFriends struct {
	ID              ID
	TrackedPersonID ID_TRACKED_PERSON
	CreatedAt       time.Time
}

type HistoryVk struct {
	VkID      VkID
	FirstName string
	LastName  string
	CreatedAt time.Time
}

type UserVK struct {
	UID       VkID   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Sex       int    `json:"sex"`
}

var (
	ErrTrackedAlreadyExist = errors.New("tracked already exist")
	ErrTrackedToMach       = errors.New("maximum number of tracked 6")
)
