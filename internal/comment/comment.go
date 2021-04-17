package comment

import "github.com/jinzhu/gorm"

type Service struct {
	DB *gorm.DB
}

type Comment struct {
	gorm.Model
	Slug   string
	Body   string
	Author string
}

type CommentService interface {
	GetCommnet(ID uint) (Comment, error)
	GetCommentsBySlug(slug string) ([]Comment, error)
	PostComment(comment Comment) (Comment, error)
	UpdateComment(ID uint, newComment Comment) (Comment, error)
	DeleteComment(ID uint) error
	GetAllComment() ([]Comment, error)
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}
