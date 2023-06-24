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

// Comment represents a comment on a blog post
type Comment struct {
	gorm.Model
	PostID uint
	Body   string `json:"body"`
}

// GetCommentsForPost returns all comments for a post
func GetCommentsForPost(postID uint) ([]Comment, error) {
	var comments []Comment
	result := DB.Where("post_id = ?", postID).Find(&comments)
	return comments, result.Error
}

// CreateComment creates a new comment for a post
func CreateComment(comment *Comment) (err error) {
	result := DB.Create(comment)
	return result.Error
}

// Like represents a like on a blog post
type Like struct {
	gorm.Model
	PostID uint
	UserID uint
}

// GetLikesForPost returns the number of likes for a post
func GetLikesForPost(postID uint) (int, error) {
	var count int64
	result := DB.Model(&Like{}).Where("post_id = ?", postID).Count(&count)
	return int(count), result.Error
}

// AddLike adds a like for a post from a user
func AddLike(postID, userID uint) (err error) {
	like := Like{PostID: postID, UserID: userID}
	result := DB.Create(&like)
	return result.Error
}
