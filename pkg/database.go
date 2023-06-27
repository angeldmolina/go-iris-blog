package database

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Post represents a blog post.
type Post struct {
	gorm.Model
	Title string `json:"title"`
	Body  string `json:"body"`
}

// User represents a user in the system.
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
}

// Database connection
var DB *gorm.DB

var jwtSecret = []byte("your-secret-key")

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

// GenerateToken generates a new JWT token for a user.
func GenerateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// VerifyToken verifies the authenticity of a JWT token and returns the associated user ID.
func VerifyToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, jwt.ErrInvalidKey
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return 0, jwt.ErrInvalidKey
	}

	return uint(userID), nil
}
