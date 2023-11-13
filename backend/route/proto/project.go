package proto

type GetMembersData struct {
	Page     int `form:"page" binding:"omitempty,gte=1"`
	PageSize int `form:"page_size" binding:"omitempty,gte=1,lte=100"`
}

type AddMemberData struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password" binding:"required,gte=6,lte=255"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

type UserIDData struct {
	UserID uint `uri:"user-id" binding:"required,gte=1"`
}

type SetMemberData struct {
	Email     string `json:"email" binding:"omitempty,email,lte=255"`
	Password  string `json:"password" binding:"omitempty,gte=6,lte=255"`
	Role      string `json:"role" binding:"omitempty,oneof=admin user"`
	IsEnabled int    `json:"is_enabled" binding:"oneof=0 1"`
}

type ProjectGroupCreateData struct {
	Name string `json:"name" binding:"required,lte=255"`
}

type ProjectGroupRenameUriData struct {
	ID uint `uri:"group_id" binding:"required,gt=0"`
}

type ProjectGroupRenameData struct {
	Name string `json:"name" binding:"required,lte=255"`
}

type ProjectGroupDeleteUriData struct {
	ID uint `uri:"group_id" binding:"required,gt=0"`
}

type ProjectGroupOrderData struct {
	IDs []uint `json:"ids" binding:"required,dive,gte=0"`
}

type ProjectMembersListData struct {
	Page     int `form:"page" binding:"omitempty,gte=1"`
	PageSize int `form:"page_size" binding:"omitempty,gte=1,lte=100"`
}

type GetPathUserID struct {
	UserID uint `uri:"user-id" binding:"required"`
}

type CreateProjectMemberData struct {
	UserIDs   []uint `json:"user_ids" binding:"required,gt=0,dive,required"`
	Authority string `json:"authority" binding:"required,oneof=manage write read"`
}

type UpdateProjectMemberAuthData struct {
	Authority string `json:"authority" binding:"required,oneof=manage write read"`
}

type ProjectMemberData struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Authority string `json:"authority"`
	CreatedAt string `json:"created_at"`
}

type ProjectSharingSwitchData struct {
	Share string `json:"share" binding:"required,oneof=open close"`
}

type ProjectShareSecretkeyCheckData struct {
	SecretKey string `json:"secret_key" binding:"required,lte=255"`
}

type CreateProject struct {
	Title      string `json:"title" binding:"required,lte=255"`
	Data       string `json:"data"`
	Cover      string `json:"cover" binding:"lte=255"`
	Visibility string `json:"visibility" binding:"required,oneof=private public"`
	DataType   string `json:"data_type" binding:"omitempty,oneof=apicat swagger openapi postman"`
	GroupID    uint   `json:"group_id" binding:"omitempty"`
}

type UpdateProject struct {
	Title       string `json:"title" binding:"required,lte=255"`
	Description string `json:"description" binding:"lte=255"`
	Cover       string `json:"cover" binding:"lte=255"`
	Visibility  string `json:"visibility" binding:"required,oneof=private public"`
}

type ProjectID struct {
	ID string `uri:"project-id" binding:"required"`
}

type ExportProject struct {
	Type     string `form:"type" binding:"required,oneof=apicat swagger openapi3.0.0 openapi3.0.1 openapi3.0.2 openapi3.1.0 HTML md"`
	Download string `form:"download" binding:"omitempty,oneof=true false"`
}

type TranslateProject struct {
	MemberID uint `json:"member_id" binding:"required,lte=255"`
}

type ProjectsListData struct {
	Auth       []string `form:"auth" binding:"omitempty,dive,oneof=manage write read"`
	GroupID    uint     `form:"group_id"`
	IsFollowed bool     `form:"is_followed"`
}

type ProjectChangeGroupData struct {
	TargetGroupID uint `json:"target_group_id" binding:"lte=255"`
}
