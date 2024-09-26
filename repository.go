package x_clone_post_srv

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, post Post) (Post, error)
	Find(ctx context.Context) ([]Post, error)
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
	result, err := r.coll.InsertOne(ctx, post)
	if err != nil {
		return post, err
	}
	insertedID, _ := result.InsertedID.(primitive.ObjectID)
	post.ID = insertedID
	return post, nil
}

func (r *mongoRepository) Find(ctx context.Context) (posts []Post, err error) {
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
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
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
