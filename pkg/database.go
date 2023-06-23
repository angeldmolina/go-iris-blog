package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Post represents a blog post.
type Post struct {
	gorm.Model
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Database connection
var DB *gorm.DB

// ConnectDb initializes database connection.
func ConnectDb() {
	var err error
	DB, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB.AutoMigrate(&Post{})
}

// GetPosts returns all blog posts.
func GetPosts() (posts []Post) {
	DB.Find(&posts)
	return
}

// GetPostByID returns a specific blog post by its ID.
func GetPostByID(id uint) (post Post, err error) {
	result := DB.First(&post, id)
	if result.Error != nil {
		err = result.Error
	}
	return
}

// CreatePost creates a new blog post.
func CreatePost(post *Post) (err error) {
	result := DB.Create(&post)
	if result.Error != nil {
		err = result.Error
	}
	return
}

// DeletePost deletes a blog post by its ID.
func DeletePost(id uint) (err error) {
	var post Post
	result := DB.Delete(&post, id)
	if result.Error != nil || result.RowsAffected == 0 {
		err = result.Error
	}
	return
}
