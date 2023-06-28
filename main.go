package main

import (
	"context"
	"encoding/json"
	"fmt"
	database "myapp/pkg"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

var redisClient *redis.Client

// Post represents a blog post.
type Post database.Post

func getPosts(ctx iris.Context) {
	posts := database.GetPosts()
	ctx.JSON(posts)
}

func getPostByID(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	post, err := database.GetPostByID(id)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.JSON(post)
}

func createPost(ctx iris.Context) {
	var post database.Post
	ctx.ReadJSON(&post)
	err := database.CreatePost(&post)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	ctx.StatusCode(iris.StatusCreated)
}

func deletePost(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	err := database.DeletePost(id)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}
	ctx.StatusCode(iris.StatusNoContent)
}

// authenticate authenticates a user and returns a JWT token upon successful authentication.
func authenticate(ctx iris.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ReadJSON(&credentials); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	var user database.User
	if err := database.DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	token, err := database.GenerateToken(&user)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(map[string]string{
		"token": token,
	})
}

// authenticateMiddleware is a middleware to verify the JWT token in the request headers.
func authenticateMiddleware(ctx iris.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	userID, err := database.VerifyToken(tokenString)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	ctx.Values().Set("userID", userID)
	ctx.Next()
}

// Initialize Redis client
func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",               // Replace with your Redis server password
		DB:       0,                // Replace with your Redis database index
	})
}

// getPostsFromCache retrieves blog posts from the Redis cache if available,
// otherwise retrieves them from the database and stores them in the cache.
func getPostsFromCache() ([]Post, error) {
	ctx := context.Background()
	cacheKey := "posts"

	// Check if posts exist in the cache
	if val, err := redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var posts []Post
		if err := json.Unmarshal([]byte(val), &posts); err != nil {
			return nil, err
		}
		return posts, nil
	}

	// Retrieve posts from the database
	var posts []Post
	if result := database.DB.Find(&posts); result.Error != nil {
		return nil, result.Error
	}

	// Store posts in the cache
	if postsJSON, err := json.Marshal(posts); err == nil {
		if err := redisClient.Set(ctx, cacheKey, postsJSON, 10*time.Minute).Err(); err != nil {
			fmt.Println("Failed to cache posts:", err)
		}
	}

	return posts, nil
}

func main() {
	app := iris.New()
	app.Use(logger.New())

	database.ConnectDb()

	// Authenticate user
	app.Post("/authenticate", authenticate)

	// Group routes that require authentication using the middleware.
	authGroup := app.Party("/api")
	authGroup.Use(authenticateMiddleware)

	// Define the routes for the API.
	app.Get("/posts", getPosts)
	app.Get("/posts/{id:uint}", getPostByID)
	app.Post("/posts", createPost)
	app.Delete("/posts/{id:uint}", deletePost)

	// Start the server.
	app.Run(iris.Addr(":8080"))
}
