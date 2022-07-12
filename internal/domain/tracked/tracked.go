package tracked

import (
	"tgbot_surveillance/internal/domain/user"

	"github.com/pkg/errors"
)

type ID int64

type ID_TRACKED_PERSON int64

type VkID int64

type Tracked struct {
	ID   ID
	VkID VkID
}

type UserTracked struct {
	ID              ID
	UserID          user.ID
	TrackedPersonID ID_TRACKED_PERSON
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

type UserVK struct {
	UID       VkID   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Sex       int    `json:"sex"`
}

var (
	ErrTrackedAlreadyExist = errors.New("tracked already exist")
)
