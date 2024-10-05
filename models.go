package x_clone_post_svc

type Post struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt uint32 `json:"created_at"`
	User      User   `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
