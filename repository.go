package x_clone_post_svc

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository follows gORM convention for the method namings
type Repository interface {
	Create(ctx context.Context, post Post) (Post, error)
	Find(ctx context.Context) ([]Post, error)
	FirstByID(ctx context.Context, id string) (Post, error)
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepository{
		coll: db.Collection("posts"),
	}
}

type mongoRepository struct {
	coll *mongo.Collection
}

type repoPost struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Content   string              `bson:"content"`
	CreatedAt primitive.Timestamp `bson:"created_at"`
	UserID    primitive.ObjectID  `bson:"user_id"`
}

func (r *mongoRepository) Create(ctx context.Context, post Post) (Post, error) {
	post.CreatedAt = uint32(time.Now().Unix())
	userObjectID, _ := primitive.ObjectIDFromHex(post.User.ID)
	repoPost := repoPost{
		Content:   post.Content,
		CreatedAt: primitive.Timestamp{T: post.CreatedAt},
		UserID:    userObjectID,
	}
	result, err := r.coll.InsertOne(ctx, repoPost)
	if err != nil {
		return Post{}, err
	}

	post.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return post, nil
}

func (r *mongoRepository) Find(ctx context.Context) (posts []Post, err error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repoPost repoPost
	for cursor.Next(ctx) {
		if err := cursor.Decode(&repoPost); err != nil {
			return nil, err
		}
		posts = append(posts, Post{
			ID:        repoPost.ID.Hex(),
			Content:   repoPost.Content,
			CreatedAt: repoPost.CreatedAt.T,
			User: User{
				ID: repoPost.UserID.Hex(),

				// TODO Change with the real one
				Username: "dummyusername_" + repoPost.UserID.Hex(),
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *mongoRepository) FirstByID(ctx context.Context, id string) (post Post, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return post, err
	}

	var repoPost repoPost
	err = r.coll.FindOne(ctx, bson.M{"_id": objectID}).Decode(&repoPost)
	if err == mongo.ErrNoDocuments {
		return post, errors.New("post not found")
	}
	if err != nil {
		return post, err
	}
	return Post{
		ID:        repoPost.ID.Hex(),
		Content:   repoPost.Content,
		CreatedAt: repoPost.CreatedAt.T,
		User: User{
			ID:       repoPost.UserID.Hex(),
			Username: "dummyusername_" + repoPost.UserID.Hex(),
		},
	}, nil
}
