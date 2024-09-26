package x_clone_post_srv

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Post      string             `bson:"post" json:"post"`
	CreatedAt int64              `json:"created_at"`
	User      User               `json:"user"`
}
