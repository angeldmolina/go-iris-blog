package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

// Post represents a blog post.
type Post struct {
	gorm.Model
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Database connection
var db *gorm.DB

// Initialize database connection
func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Post{})
}

// getPosts returns all blog posts.
func getPosts(ctx iris.Context) {
	var posts []Post
	db.Find(&posts)
	ctx.JSON(posts)
}

// getPostByID returns a specific blog post by its ID.
func getPostByID(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	var post Post
	result := db.First(&post, id)
	if result.Error != nil {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	ctx.JSON(post)
}

// createPost creates a new blog post.
func createPost(ctx iris.Context) {
	var post Post
	err := ctx.ReadJSON(&post)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	result := db.Create(&post)
	if result.Error != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	ctx.StatusCode(iris.StatusCreated)
}

// deletePost deletes a blog post by its ID.
func deletePost(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	var post Post
	result := db.Delete(&post, id)
	if result.Error != nil || result.RowsAffected == 0 {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}

func main() {
	app := iris.New()
	app.Use(logger.New())

	initDatabase()

	// Define the routes for the API.
	app.Get("/posts", getPosts)
	app.Get("/posts/{id:uint}", getPostByID)
	app.Post("/posts", createPost)
	app.Delete("/posts/{id:uint}", deletePost)

	// Start the server.
	app.Run(iris.Addr(":8080"))
}
