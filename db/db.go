package db

import (
	"context"
	"errors"
	"filestorageapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DBManager struct {
	client *mongo.Client
}

func NewDBManager(uri string) (*DBManager, error) {
	dbManager := &DBManager{}
	err := dbManager.Connect(uri)
	if err != nil {
		return nil, err
	}

	return dbManager, nil
}

func (d *DBManager) Connect(uri string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	d.client = client
	return nil
}

func (d *DBManager) Disconnect() {
	if d.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		d.client.Disconnect(ctx)
	}
}

func (d *DBManager) Users() *mongo.Collection {
	return d.client.Database("filestorage").Collection("users")
}

func (d *DBManager) Files() *mongo.Collection {
	return d.client.Database("filestorage").Collection("files")
}

func (d *DBManager) SaveFileMetadata(fileMetadata *models.FileMetadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := d.Files().InsertOne(ctx, fileMetadata)
	return err
}

func (d *DBManager) GetFileMetadata(fileID string) (*models.FileMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": fileID}
	var fileMetadata models.FileMetadata

	err := d.Files().FindOne(ctx, filter).Decode(&fileMetadata)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("File metadata not found")
		}
		return nil, err
	}

	return &fileMetadata, nil
}

func (d *DBManager) DeleteFileMetadata(fileID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": fileID}
	res, err := d.Files().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("File metadata not found")
	}

	return nil
}

func (d *DBManager) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	var user models.User

	err := d.Users().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return &user, nil
}

func (d *DBManager) CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := d.Users().InsertOne(ctx, user)
	return err
}
