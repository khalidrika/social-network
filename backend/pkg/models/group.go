package models

type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AdminID     string `json:"adminId"`
	Privacy     string `json:"privacy"`     // "public" or "private"
	Image       string `json:"image"`       // profile/banner image path
	MemberCount int    `json:"memberCount"` // total number of members including admin

	Invitations []string `json:"invitations"`

	Member         bool `json:"member"`         // true if current user is a member
	Administrator  bool `json:"admin"`          // true if current user is admin
	RequestPending bool `json:"requestPending"` // true if request to join is pending
}

type GroupRepository interface {
	GetAllAndRelations(userId string) ([]Group, error)
	GetUserGroups(userId string) ([]Group, error)
	New(Group) error                               //create new group
	GetData(groupId string) (Group, error)         //get info- name and desc
	GetMembers(groupId string) ([]User, error)     // get all group members and admin
	GetAdmin(groupId string) (string, error)       //get admin id
	IsMember(groupId, userId string) (bool, error) //checks if user is a member
	IsAdmin(groupId, userId string) (bool, error)  //checks if user is admin

	SaveMember(userId, groupId string) error
	RemoveMember(userId, groupId string) error
}
