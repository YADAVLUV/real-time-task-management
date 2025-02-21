// main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "golang.org/x/crypto/bcrypt"
    "github.com/gin-contrib/cors"
)

// Models
type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Email    string            `bson:"email" json:"email" binding:"required"`
    Password string            `bson:"password" json:"password" binding:"required"`
}

type Task struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Title       string            `bson:"title" json:"title" binding:"required"`
    Description string            `bson:"description" json:"description"`
    DueDate     time.Time         `bson:"due_date" json:"due_date"`
    Status      string            `bson:"status" json:"status"`
    UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
    AssigneeID  primitive.ObjectID `bson:"assignee_id,omitempty" json:"assignee_id,omitempty"`
    CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

// JWT claims struct
type Claims struct {
    UserID primitive.ObjectID `json:"user_id"`
    jwt.RegisteredClaims
}

var client *mongo.Client
var userCollection *mongo.Collection
var taskCollection *mongo.Collection
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func CORSMiddleware() gin.HandlerFunc {
    allowedOrigin := os.Getenv("FRONTEND_URL")
    if allowedOrigin == "" {
        allowedOrigin = "http://localhost:5173" // Default fallback
    }

    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
        c.Header("Access-Control-Max-Age", "86400")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString, err := c.Cookie("token") // ✅ Read token from cookie
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("userID", claims.UserID)
        c.Next()
    }
}


func register(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    user.Password = string(hashedPassword)
    
    // Create user
    result, err := userCollection.InsertOne(context.Background(), user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }
    
    user.ID = result.InsertedID.(primitive.ObjectID)
    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "id": user.ID})
}

func login(c *gin.Context) {
    var credentials struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
        return
    }

    var user User
    err := userCollection.FindOne(context.Background(), bson.M{"email": credentials.Email}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // ✅ Generate JWT token
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // ✅ Store JWT in HttpOnly Cookie
    c.SetCookie("token", tokenString, 86400, "/", "localhost", false, true)

    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
func logout(c *gin.Context) {
    c.SetCookie("token", "", -1, "/", "", false, true) // ✅ Expire cookie immediately
    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}


func getTasks(c *gin.Context) {
    userID := c.MustGet("userID").(primitive.ObjectID)
    
    cursor, err := taskCollection.Find(context.Background(), bson.M{"user_id": userID})
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
    var task Task
    if err := c.ShouldBindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    task.UserID = c.MustGet("userID").(primitive.ObjectID)
    task.Status = "pending"
    task.CreatedAt = time.Now()
    task.UpdatedAt = time.Now()
    
    result, err := taskCollection.InsertOne(context.Background(), task)
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
    
    userID := c.MustGet("userID").(primitive.ObjectID)
    
    var task Task
    err = taskCollection.FindOne(context.Background(), bson.M{
        "_id": taskID,
        "user_id": userID,
    }).Decode(&task)
    
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    
    var updateData Task
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    updateData.UpdatedAt = time.Now()
    updateData.UserID = userID
    
    _, err = taskCollection.UpdateOne(
        context.Background(),
        bson.M{"_id": taskID},
        bson.M{"$set": updateData},
    )
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
        return
    }
    
    c.JSON(http.StatusOK, updateData)
}

func deleteTask(c *gin.Context) {
    taskID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }
    
    userID := c.MustGet("userID").(primitive.ObjectID)
    
    result, err := taskCollection.DeleteOne(context.Background(), bson.M{
        "_id": taskID,
        "user_id": userID,
    })
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
        return
    }
    
    if result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func protected(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "You are authenticated!", "userID": userID})
}


func main() {
    // Load .env file
    godotenv.Load()
    
    // MongoDB connection
    ctx := context.Background()
    clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
    var err error
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    
    // Ping the database
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }
    
    db := client.Database("taskmanager")
    userCollection = db.Collection("users")
    taskCollection = db.Collection("tasks")
    
    // Create indexes
    _, err = userCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys:    bson.D{{Key: "email", Value: 1}},
        Options: options.Index().SetUnique(true),
    })
    if err != nil {
        log.Fatal("Failed to create index:", err)
    }
    
    // Initialize router
    r := gin.Default()
    
    // Middleware
    // r.Use(CORSMiddleware())
    r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Change to your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Allow cookies and authentication headers
		MaxAge:           12 * time.Hour,
	}))
    
    // Routes
    auth := r.Group("/auth")
    {
        auth.POST("/register", register)
        auth.POST("/login", login)
		auth.GET("/ping",ping)
        auth.POST("/logout", logout)
        auth.GET("/protected",AuthMiddleware(),protected)
    }
    
    api := r.Group("/api")
    api.Use(AuthMiddleware())
    {
        api.GET("/gettasks", getTasks)
        api.POST("/tasks", createTask)
        api.PUT("/tasks/:id", updateTask)
        api.DELETE("/tasks/:id", deleteTask)
    }

    
    r.Run(":8080")
}
