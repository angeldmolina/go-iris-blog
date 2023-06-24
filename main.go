package main

import (
	database "myapp/pkg"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

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

func main() {
	app := iris.New()
	app.Use(logger.New())

	database.ConnectDb()

	// Define the routes for the API.
	app.Get("/posts", getPosts)
	app.Get("/posts/{id:uint}", getPostByID)
	app.Post("/posts", createPost)
	app.Delete("/posts/{id:uint}", deletePost)

	// Start the server.
	app.Run(iris.Addr(":8080"))
}
