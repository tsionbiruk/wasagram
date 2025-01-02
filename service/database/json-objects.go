package database

type UserProfileInfo struct {
	Followers []string     `json:"followers"`
	Following []string     `json:"following"`
	Banned    []string     `json:"banned"`
	Photos    []StreamPost `json:"photos"`
}

type CommentData struct {
	Id        int64  `json:"id"`
	User      string `json:"user"`
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content"`
}

type StreamPost struct {
	Id        int64         `json:"id"`
	User      string        `json:"user"`
	Timestamp int64         `json:"timestamp"`
	Likes     []string      `json:"likes"`
	Comments  []CommentData `json:"comments"`
}
