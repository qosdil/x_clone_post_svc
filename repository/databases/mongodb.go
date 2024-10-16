package databases

import (
	"context"
	"time"
	model "x_clone_post_svc/model"
	repository "x_clone_post_svc/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewMongoRepository(db *mongo.Database) repository.Repository {
	return &mongoRepository{
		coll: db.Collection("posts"),
	}
}

type mongoRepository struct {
	coll *mongo.Collection
}

type mongoRepoPost struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Content   string              `bson:"content"`
	CreatedAt primitive.Timestamp `bson:"created_at"`
	UserID    primitive.ObjectID  `bson:"user_id"`
}

func (r *mongoRepository) Create(ctx context.Context, post model.Post) (model.Post, error) {
	post.CreatedAt = uint32(time.Now().Unix())
	userObjectID, _ := primitive.ObjectIDFromHex(post.User.ID)
	repoPost := mongoRepoPost{
		Content:   post.Content,
		CreatedAt: primitive.Timestamp{T: post.CreatedAt},
		UserID:    userObjectID,
	}
	result, err := r.coll.InsertOne(ctx, repoPost)
	if err != nil {
		return model.Post{}, err
	}

	post.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return post, nil
}

func (r *mongoRepository) Find(ctx context.Context) (posts []model.Post, err error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoRepoPost mongoRepoPost
	for cursor.Next(ctx) {
		if err := cursor.Decode(&mongoRepoPost); err != nil {
			return nil, err
		}
		posts = append(posts, model.Post{
			ID:        mongoRepoPost.ID.Hex(),
			Content:   mongoRepoPost.Content,
			CreatedAt: mongoRepoPost.CreatedAt.T,
			User: model.User{
				ID: mongoRepoPost.UserID.Hex(),

				// TODO Change with the real one
				Username: "username_" + mongoRepoPost.UserID.Hex(),
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *mongoRepository) FirstByID(ctx context.Context, id string) (post model.Post, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return post, model.ErrBadIDFormat
	}

	var mongoRepoPost mongoRepoPost
	err = r.coll.FindOne(ctx, bson.M{"_id": objectID}).Decode(&mongoRepoPost)
	if err == mongo.ErrNoDocuments {
		return post, model.ErrNotFound
	}
	if err != nil {
		return post, err
	}
	return model.Post{
		ID:        mongoRepoPost.ID.Hex(),
		Content:   mongoRepoPost.Content,
		CreatedAt: mongoRepoPost.CreatedAt.T,
		User: model.User{
			ID:       mongoRepoPost.UserID.Hex(),
			Username: "username_" + mongoRepoPost.UserID.Hex(),
		},
	}, nil
}
