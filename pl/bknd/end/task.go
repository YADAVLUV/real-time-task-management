package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupTaskRoutes(r *gin.RouterGroup) {
	tasks := r.Group("/tasks")
	tasks.Use(authMiddleware())
	{
		tasks.GET("", getTasks)
		tasks.POST("", createTask)
		tasks.PATCH("/:id", updateTask)
		tasks.DELETE("/:id", deleteTask)
	}
}

func getTasks(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	tasksCollection := mongoClient.Database("taskmanager").Collection("tasks")
	cursor, err := tasksCollection.Find(context.Background(), bson.M{"userId": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	defer cursor.Close(context.Background())

	var tasks []Task
	if err := cursor.All(context.Background(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := primitive.ObjectIDFromHex(c.GetString("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	task := Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Completed:   false,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	tasksCollection := mongoClient.Database("taskmanager").Collection("tasks")
	result, err := tasksCollection.InsertOne(context.Background(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	task.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, task)
}

func updateTask(c *gin.Context) {
	taskID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(c.GetString("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasksCollection := mongoClient.Database("taskmanager").Collection("tasks")

	// Verify task ownership
	var task Task
	err = tasksCollection.FindOne(context.Background(), bson.M{
		"_id":    taskID,
		"userId": userID,
	}).Decode(&task)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Build update document
	update := bson.M{}
	if req.Title != nil {
		update["title"] = *req.Title
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if req.Priority != nil {
		update["priority"] = *req.Priority
	}
	if req.Completed != nil {
		update["completed"] = *req.Completed
	}

	result := tasksCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": taskID, "userId": userID},
		bson.M{"$set": update},
		&options.FindOneAndUpdateOptions{ReturnDocument: options.After},
	)

	if err := result.Decode(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	taskID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(c.GetString("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	tasksCollection := mongoClient.Database("taskmanager").Collection("tasks")
	result, err := tasksCollection.DeleteOne(context.Background(), bson.M{
		"_id":    taskID,
		"userId": userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
