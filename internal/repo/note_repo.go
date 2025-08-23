package repo

import (
	"context"
	"time"

	"github.com/saurabhraut1212/notes_sharing_api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NoteRepo struct {
	col *mongo.Collection
}

func NewNoteRepo(db *mongo.Database) *NoteRepo {
	return &NoteRepo{
		col: db.Collection("notes"),
	}
}

// Add methods for NoteRepo as needed, e.g., Create, FindByID, Update, Delete, etc.

func (r *NoteRepo) Create(ctx context.Context, n *models.Note) error {
	now := time.Now().UTC()
	n.ID = primitive.NewObjectID()
	n.CreatedAt = now
	n.UpdatedAt = now
	_, err := r.col.InsertOne(ctx, n)
	return err

}

func (r *NoteRepo) FindById(ctx context.Context, id primitive.ObjectID) (*models.Note, error) {
	var n models.Note
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&n)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &n, err
}

func (r *NoteRepo) ListByUser(ctx context.Context, userId primitive.ObjectID, page, limit int) ([]models.Note, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	filter := bson.M{"user_id": userId}
	cur, err := r.col.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit64,
		Sort:  bson.M{"created_at": -1}},
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var notes []models.Note
	for cur.Next(ctx) {
		var note models.Note
		if err := cur.Decode(&note); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, cur.Err()
}

func (r *NoteRepo) ListPublic(ctx context.Context, page, limit int) ([]models.Note, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	filter := bson.M{"is_public": true}
	cur, err := r.col.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit64,
		Sort:  bson.M{"created_at": -1}},
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var notes []models.Note
	for cur.Next(ctx) {
		var note models.Note
		if err := cur.Decode(&note); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, cur.Err()
}

func (r *NoteRepo) Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.Note, error) {
	update["updated_at"] = time.Now().UTC()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var n models.Note
	err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": update}, opts).Decode(&n)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &n, err
}

func (r *NoteRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
