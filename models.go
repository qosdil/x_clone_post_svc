package x_clone_post_srv

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Content   string              `bson:"content" json:"content"`
	CreatedAt primitive.Timestamp `bson:"created_at" json:"created_at"`
	UserID    primitive.ObjectID  `bson:"user_id" json:"user_id"`
}

type PostResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt uint32 `json:"created_at"`
	User      User   `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
