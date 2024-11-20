package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/Braendie/RestAPI/internal/user"
	"github.com/Braendie/RestAPI/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// db is a struct, that uses collection from mongodb and logger for logging data.
type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

// NewStorage creates storage for using mongodb by collection.
func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {

	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

// Create adds new user in the database mongodb.
// It returns user id, that need to write on user object.
func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("Create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}

	d.logger.Debug("convert InsertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert InsertedID to ObjectID, probably oid: %s", oid)
}

// FindAll searches all users in the database and returns them.
func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return u, fmt.Errorf("failed to find all users due to error: %v", err)
	}

	if err := cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor. error: %v", err)
	}

	return u, nil
}

// FindOne searches a user by id in the database.
// If it finds a user, it returns a user object with this id.
// Otherwise it returns error 404.
func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to ObjectID, hex: %s", id)
	}

	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			// TODO ErrEntityNotFound
			return u, fmt.Errorf("not found")
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user (id:%s) from DB due to error: %v", id, err)
	}

	return u, nil
}

// Update fuction updates a user, that gives into database.
// If this user is not in the database it returns error 404.
func (d *db) Update(ctx context.Context, user user.User) error {
	d.logger.Debug("Update user")
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", user.ID)
	}

	filter := bson.M{"_id": objectID}
	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user. error: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user bytes. error: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query. error: %v", err)
	}

	if result.MatchedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}

	d.logger.Tracef("Matched %d, documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

// Delete removes a user from database by his id.
// If this user is not in the database it returns error 404.
func (d *db) Delete(ctx context.Context, id string) error {
	d.logger.Debug("Delete user")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", id)
	}

	filter := bson.M{"_id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %v", err)
	}

	if result.DeletedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}

	d.logger.Tracef("Deleted %d, documents", result.DeletedCount)

	return nil
}
