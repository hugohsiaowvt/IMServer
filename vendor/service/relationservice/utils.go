package relationservice

const (
	TABLE_RELATIONSHIPS	string = "relationships"

	None =		0
	Pending =	1
	Acception =	2
	Rejection =	3
)

type RelationShips struct {
	UserOneId		string	`json:"user_one_id"`
	UserTwoId		string	`json:"user_two_id"`
	ActionUserId	string	`json:"action_user_id"`
	Time			int64	`json:"time"`
	Status			int		`json:"status"`
	RequestToken	string	`json:"request_token"`
}

type InviteFriendInput struct {
	ToId	string 	`json:"to_id"`
}

type AcceptFriendInput struct {
	RequestToken	string 	`json:"request_token"`
}