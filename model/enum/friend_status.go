package enum

type FriendRequestStatus string

const (
	Pending = FriendRequestStatus("pending")
	Accept  = FriendRequestStatus("accept")
	Reject  = FriendRequestStatus("reject")
)
