package request

type AgeUpdate struct {
	Age int `json:"age"`
}

type UserDelete struct {
	TargetId int `json:"target_id"`
}

type FriendshipRequest struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

type UserData struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}
