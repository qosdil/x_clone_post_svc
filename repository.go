package x_clone_post_svc

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, post Post) (Post, error)
	Find(ctx context.Context) ([]PostResponse, error)
	FindByID(ctx context.Context, id string) (Post, error)
}

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepository{
		coll: db.Collection("posts"),
	}
}

func (r *mongoRepository) Create(ctx context.Context, post Post) (Post, error) {
	post.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
	result, err := r.coll.InsertOne(ctx, post)
	if err != nil {
		return post, err
	}
	insertedID, _ := result.InsertedID.(primitive.ObjectID)
	post.ID = insertedID
	return post, nil
}

func (r *mongoRepository) Find(ctx context.Context) (postResponses []PostResponse, err error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var post Post
	for cursor.Next(ctx) {
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		postResponses = append(postResponses, PostResponse{
			ID:        post.ID.Hex(),
			Content:   post.Content,
			CreatedAt: post.CreatedAt.T,
			User: User{
				ID: post.UserID.Hex(),
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return postResponses, nil
}

func (r *mongoRepository) FindByID(ctx context.Context, id string) (post Post, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return post, err
	}

	err = r.coll.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return post, errors.New("post not found")
		}
		return post, err
	}
	return post, nil
}
