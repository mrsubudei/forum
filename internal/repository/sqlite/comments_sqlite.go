package sqlite

import (
	"forum/internal/entity"
	"forum/pkg/sqlite3"
)

type CommentsRepo struct {
	*sqlite3.Sqlite
}

func NewCommentsRepo(sq *sqlite3.Sqlite) *CommentsRepo {
	return &CommentsRepo{sq}
}

func (cr *CommentsRepo) Store(comment entity.Comment) error {
	return nil
}

func (cr *CommentsRepo) Fetch() ([]entity.Comment, error) {
	var commets []entity.Comment
	return commets, nil
}

func (cr *CommentsRepo) GetById(id int64) (entity.Comment, error) {
	var comment entity.Comment
	return comment, nil
}

func (cr *CommentsRepo) Update(comment entity.Comment) (entity.Comment, error) {
	return comment, nil
}

func (cr *CommentsRepo) Delete(comment entity.Comment) error {
	return nil
}

func (pr *CommentsRepo) ThumbsUp(comment entity.Comment) error {
	return nil
}

func (pr *CommentsRepo) ThumbsDown(comment entity.Comment) error {
	return nil
}
