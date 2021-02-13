package db

const (
	UserAlreadyExists = "user already exists"
	UserNotFound      = "user not found"
	UserDeleted       = "user deleted"
	UserBlocked       = "user blocked"

	GroupAlreadyExists = "group already exists"
	GroupNotFound      = "group not found"
	GroupDeleted       = "group deleted"

	GroupChatAlreadyExists = "groupchat already exists"
	GroupChatNotFound      = "groupchat not found"

	MeetingroomAlreadyExists = "meetingroom already exists"
	MeetingroomNotFound      = "meetingroom not found"
	MeetingroomDeleted       = "meetingroom deleted"
	MeetingroomBlocked       = "meetingroom blocked"

	Deleted = "deleted"
	Blocked = "blocked"
	New     = "new"
	Member  = "member"
	Admin   = "admin"
	Owner   = "owner"
)
