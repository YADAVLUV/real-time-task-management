package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string            `json:"username" bson:"username"`
	Password string            `json:"-" bson:"password"`
}

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string            `json:"title" bson:"title"`
	Description string            `json:"description" bson:"description"`
	Priority    int               `json:"priority" bson:"priority"`
	Completed   bool              `json:"completed" bson:"completed"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	CreatedAt   time.Time         `json:"createdAt" bson:"createdAt"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    int    `json:"priority" binding:"min=1,max=5"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *int    `json:"priority" binding:"omitempty,min=1,max=5"`
	Completed   *bool   `json:"completed"`
}
