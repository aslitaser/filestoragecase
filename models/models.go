package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

type File struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	Filename   string             `bson:"filename"`
	UploadDate time.Time          `bson:"upload_date"`
	FileSize   int64              `bson:"file_size"`
}
type FileMetadata struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Filename    string `bson:"file_name" json:"file_name"`
	ContentType string `bson:"content_type" json:"content_type"`
	Filesize    int64  `bson:"size" json:"size"`
	UploadDate  int64  `bson:"upload_date" json:"upload_date"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
