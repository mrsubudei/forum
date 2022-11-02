package usecase

import (
	"encoding/json"
	"forum/internal/entity"
	"testing"
)

var mockRepo = &DeleteMockRepository{
	Users: map[int]entity.User{
		1: user1,
		2: user2,
	},
	Posts: map[int]entity.Post{
		1: post1,
		2: post2,
	},
	Comments: map[int]entity.Comment{
		1: comment1,
		2: comment2,
	},
}

func TestDeleteEntity(t *testing.T) {
	userItr := NewInteractor(mockRepo)
	//delete user
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(user1)
		err := userItr.DeleteUser(jsonBlob)
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("user doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(user3)
		err := userItr.DeleteUser(jsonBlob)
		if err != USER_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", USER_NOT_EXISTS, err)
		}
	})
	//delete post
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(post1)
		err := userItr.DeletePost(jsonBlob)
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("post doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(post3)
		err := userItr.DeletePost(jsonBlob)
		if err != POST_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", POST_NOT_EXISTS, err)
		}
	})
	//delete comment
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(comment1)
		err := userItr.DeleteComment(jsonBlob)
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("comment doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(comment3)
		err := userItr.DeleteComment(jsonBlob)
		if err != COMMENT_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", COMMENT_NOT_EXISTS, err)
		}
	})
	//delete post reaction
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(postReaction1)
		err := userItr.DeletePostReaction(jsonBlob)
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("post doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(postReaction3)
		err := userItr.DeletePostReaction(jsonBlob)
		if err != POST_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", POST_NOT_EXISTS, err)
		}
	})
	t.Run("user doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(postReaction4)
		err := userItr.DeletePostReaction(jsonBlob)
		if err != USER_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", USER_NOT_EXISTS, err)
		}
	})
	//delete comment reaction
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(postReaction1)
		err := userItr.DeletePostReaction(jsonBlob)
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("comment doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(commentReaction3)
		err := userItr.DeleteCommentReaction(jsonBlob)
		if err != COMMENT_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", COMMENT_NOT_EXISTS, err)
		}
	})
	t.Run("user doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(commentReaction4)
		err := userItr.DeleteCommentReaction(jsonBlob)
		if err != USER_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", USER_NOT_EXISTS, err)
		}
	})

}

type DeleteMockRepository struct {
	Users    map[int]entity.User
	Posts    map[int]entity.Post
	Comments map[int]entity.Comment
}

func (m *DeleteMockRepository) Delete(jsonBlob []byte) error {
	r := Request{}
	json.Unmarshal(jsonBlob, &r)
	switch r.Header.Entity {
	case "User":
		id := r.Body.User.Id
		if _, ok := m.Users[id]; ok {
			return nil
		}
		return USER_NOT_EXISTS
	case "Post":
		id := r.Body.Post.Id
		if _, ok := m.Posts[id]; ok {
			return nil
		}
		return POST_NOT_EXISTS
	case "Comment":
		id := r.Body.Comment.Id
		if _, ok := m.Comments[id]; ok {
			return nil
		}
		return COMMENT_NOT_EXISTS

	case "PostReaction":
		postId := r.Body.PostReaction.Post.Id
		userId := r.Body.PostReaction.Reaction.User.Id
		if _, ok := m.Posts[postId]; ok {
			if _, ok := m.Users[userId]; ok {
				return nil
			}
			return USER_NOT_EXISTS
		}
		return POST_NOT_EXISTS
	case "CommentReaction":
		commentId := r.Body.CommentReaction.Comment.Id
		userId := r.Body.CommentReaction.Reaction.User.Id
		if _, ok := m.Comments[commentId]; ok {
			if _, ok := m.Users[userId]; ok {
				return nil
			}
			return USER_NOT_EXISTS
		}
		return COMMENT_NOT_EXISTS
	}
	return nil
}

func (m *DeleteMockRepository) Get([]byte) ([]byte, error)    { return nil, nil }
func (m *DeleteMockRepository) Change([]byte) ([]byte, error) { return nil, nil }
func (m *DeleteMockRepository) Create([]byte) ([]byte, error) { return nil, nil }
