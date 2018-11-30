package groupservice

const (

	TABLE_GROUPS				string	= "groups"
	TABLE_GROUP_MEMBERS			string	= "group_members"

	GROUP_ID_PREFIX				string	= "group:"

	ROLE_MEMBER					int		= 1
	ROLE_GROUP_MANAGER			int		= 2
	ROLE_MANAGER				int		= 3

	GROUP_MEMBER_STATUS_OK		int		= 1
	GROUP_MEMBER_STATUS_LOCKED	int		= -1

)

type Group struct {
	Id				int		`json:"id"`
	GroupId			string	`json:"group_id"`
	GroupName		string	`json:"group_name"`
	MaxMember		int		`json:"max_member"`
	Introduce		string	`json:"introduce"`
	IsOpen			int		`json:"is_open"`
	Status			int		`json:"status"`
}

type GroupMember struct {
	GroupId			string	`json:"group_id"`
	OpenId			string	`json:"open_id"`
	NickName		string	`json:"nick_name"`
	Role			int		`json:"role"`
	Status			int		`json:"status"`
}

type CreateGroupInput struct {
	GroupName		string				`json:"group_name"`
	MaxMember		int					`json:"max_member"`
	Introduce		string				`json:"introduce"`
	IsOpen			int					`json:"is_open"`
	Members   		[]GroupMemberInput	`json:"members"`
}

type GroupMemberInput struct {
	OpenId			string `json:"open_id"`
}

type CreateGroupResponse struct {
	GroupId			string	`json:"group_id"`
}

func (self CreateGroupInput) ToModel() *Group {
	model := &Group{}
	model.GroupName = self.GroupName
	model.MaxMember = self.MaxMember
	model.Introduce = self.Introduce
	model.IsOpen = self.IsOpen
	return model
}