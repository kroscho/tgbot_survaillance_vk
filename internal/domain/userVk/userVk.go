package userVk

type ID int64

type VkID int64

type UserVk struct {
	ID        ID
	VkID      VkID
	FirstName string
	LastName  string
}

func New(vkID VkID, firstName string, lastName string) *UserVk {
	return &UserVk{VkID: vkID, FirstName: firstName, LastName: lastName}
}
