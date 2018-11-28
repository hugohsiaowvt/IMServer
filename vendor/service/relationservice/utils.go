package relationservice

const (

	TABLE_RELATIONSHIPS		string = "relationships"
	REQUEST_TOKEN_PREFIX	string = "request_token:"

	None =		-1
	Pending =	1
	Acception =	2

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

type QueryFriendRecordInput struct {
	RequestToken	string 	`json:"request_token"`
}