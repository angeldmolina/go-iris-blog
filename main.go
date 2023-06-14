package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

// Post represents a blog post.
type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Posts is an in-memory store for blog posts.
var Posts = []Post{
	{ID: 1, Title: "First Post", Body: "Hello, world!"},
	{ID: 2, Title: "Second Post", Body: "Iris is awesome!"},
}

// getPosts returns all blog posts.
func getPosts(ctx iris.Context) {
	ctx.JSON(Posts)
}

// getPostByID returns a specific blog post by its ID.
func getPostByID(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	for _, post := range Posts {
		if post.ID == id {
			ctx.JSON(post)
			return
		}
	}

	ctx.StatusCode(iris.StatusNotFound)
}

// createPost creates a new blog post.
func createPost(ctx iris.Context) {
	var post Post
	err := ctx.ReadJSON(&post)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	// Generate a new unique ID for the post.
	post.ID = len(Posts) + 1

	Posts = append(Posts, post)

	ctx.StatusCode(iris.StatusCreated)
}

// deletePost deletes a blog post by its ID.
func deletePost(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	for i, post := range Posts {
		if post.ID == id {
			Posts = append(Posts[:i], Posts[i+1:]...)
			ctx.StatusCode(iris.StatusNoContent)
			return
		}
	}

	ctx.StatusCode(iris.StatusNotFound)
}

func main() {
	app := iris.New()
	app.Use(logger.New())

	// Define the routes for the API.
	app.Get("/posts", getPosts)
	app.Get("/posts/{id:int}", getPostByID)
	app.Post("/posts", createPost)
	app.Delete("/posts/{id:int}", deletePost)

	// Start the server.
	app.Run(iris.Addr(":8080"))
}
