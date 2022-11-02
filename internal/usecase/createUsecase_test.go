package usecase

import (
	"encoding/json"
	"errors"
	"forum/internal/entity"
	"testing"
)

var (
	USER_EXISTS        = errors.New("user with a given email exists")
	USER_NOT_EXISTS    = errors.New("user with a given id doesn't exists")
	POST_NOT_EXISTS    = errors.New("post with a given id doesn't exists")
	COMMENT_NOT_EXISTS = errors.New("comment with a given id doesn't exists")
)

var (
	user1 = entity.User{
		Id:          1,
		Name:        "User1",
		Email:       "user1@mail.com",
		Password:    "password",
		RegDate:     "01-11-2022",
		DateOfBirth: "01-11-1990",
		City:        "Astana",
		Sex:         "male",
	}
	user2 = entity.User{
		Id:          2,
		Name:        "User2",
		Email:       "user2@mail.com",
		Password:    "password",
		RegDate:     "01-11-2022",
		DateOfBirth: "01-11-1990",
		City:        "Almaty",
		Sex:         "female",
	}
	user3 = entity.User{
		Id:          3,
		Name:        "User3",
		Email:       "user3@mail.com",
		Password:    "password",
		RegDate:     "01-11-2022",
		DateOfBirth: "01-11-1990",
		City:        "Taiynsha",
		Sex:         "male",
	}
	post1 = entity.Post{
		Id:       1,
		User:     user1,
		Date:     "01-11-2022",
		Title:    "Post1",
		Content:  "Post1 Content",
		Category: []entity.Category{category1, category2},
	}
	post2 = entity.Post{
		Id:       2,
		User:     user2,
		Date:     "01-11-2022",
		Title:    "Post2",
		Content:  "Post2 Content",
		Category: []entity.Category{category1, category3},
	}
	post3 = entity.Post{
		Id:       3,
		User:     user3,
		Date:     "01-11-2022",
		Title:    "Post2",
		Content:  "Post2 Content",
		Category: []entity.Category{category2, category3},
	}
	postReaction1 = entity.PostReaction{
		Post:     post1,
		Reaction: entity.Reaction{Like: true, Date: "01-11-2022", User: user2},
	}
	postReaction2 = entity.PostReaction{
		Post:     post2,
		Reaction: entity.Reaction{Like: false, Date: "01-11-2022", User: user1},
	}
	postReaction3 = entity.PostReaction{
		Post:     post3,
		Reaction: entity.Reaction{Like: true, Date: "01-11-2022", User: user1},
	}
	postReaction4 = entity.PostReaction{
		Post:     post2,
		Reaction: entity.Reaction{Like: false, Date: "01-11-2022", User: user3},
	}
	comment1 = entity.Comment{
		Id:      1,
		Post:    post1,
		User:    user2,
		Date:    "01-11-22",
		Content: "comment1",
	}
	comment2 = entity.Comment{
		Id:      2,
		Post:    post2,
		User:    user1,
		Date:    "01-11-22",
		Content: "comment2",
	}
	comment3 = entity.Comment{
		Id:      3,
		Post:    post3,
		User:    user2,
		Date:    "01-11-22",
		Content: "comment3",
	}
	comment4 = entity.Comment{
		Id:      4,
		Post:    post1,
		User:    user3,
		Date:    "01-11-22",
		Content: "comment4",
	}
	commentReaction1 = entity.CommentReaction{
		Comment:  comment1,
		Reaction: entity.Reaction{Like: true, Date: "01-11-2022", User: user2},
	}
	commentReaction2 = entity.CommentReaction{
		Comment:  comment2,
		Reaction: entity.Reaction{Like: false, Date: "01-11-2022", User: user1},
	}
	commentReaction3 = entity.CommentReaction{
		Comment:  comment3,
		Reaction: entity.Reaction{Like: true, Date: "01-11-2022", User: user1},
	}
	commentReaction4 = entity.CommentReaction{
		Comment:  comment2,
		Reaction: entity.Reaction{Like: false, Date: "01-11-2022", User: user3},
	}
	category1 = entity.Category{
		Title: "Golang",
	}
	category2 = entity.Category{
		Title: "Python",
	}
	category3 = entity.Category{
		Title: "Java",
	}
)

func TestCreateEntity(t *testing.T) {
	userItr := NewInteractor(&CreateMockRepository{})
	//`create user` tests
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(user1)
		result, err := userItr.CreateUser(jsonBlob)
		if string(result) != string(jsonBlob) {
			t.Errorf("result of function fails. wants %v but gets %v", string(jsonBlob), string(result))
		}
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(user2)
		result, err := userItr.CreateUser(jsonBlob)
		if string(result) != string(jsonBlob) {
			t.Errorf("result of function fails. wants %v but gets %v", string(jsonBlob), string(result))
		}
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("user exists", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(user1)
		_, err := userItr.CreateUser(jsonBlob)
		if !errors.Is(err, USER_EXISTS) {
			t.Errorf("error of function fails. wants %v but get %v", USER_EXISTS, err)
		}
	})
	//`create post` tests
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(post1)
		result, err := userItr.CreatePost(jsonBlob)
		if string(result) != string(jsonBlob) {
			t.Errorf("result of function fails. wants %v but gets %v", string(jsonBlob), string(result))
		}
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(post2)
		result, err := userItr.CreatePost(jsonBlob)
		if string(result) != string(jsonBlob) {
			t.Errorf("result of function fails. wants %v but gets %v", string(jsonBlob), string(result))
		}
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("user doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(post3)
		_, err := userItr.CreatePost(jsonBlob)
		if !errors.Is(err, USER_NOT_EXISTS) {
			t.Errorf("error of function fails. wants %v but get %v", USER_NOT_EXISTS, err)
		}
	})
	//`create post reaction` tests
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(postReaction1)
		result, err := userItr.CreatePostReaction(jsonBlob)
		if string(result) != string(jsonBlob) {
			t.Errorf("result of function fails. wants %v but gets %v", string(jsonBlob), string(result))
		}
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(postReaction2)
		result, err := userItr.CreatePostReaction(jsonBlob)
		if string(result) != string(jsonBlob) {
			t.Errorf("result of function fails. wants %v but gets %v", string(jsonBlob), string(result))
		}
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("post doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(postReaction3)
		_, err := userItr.CreatePostReaction(jsonBlob)
		if err != POST_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("user doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(postReaction4)
		_, err := userItr.CreatePostReaction(jsonBlob)
		if err != USER_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	//create comments
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(comment1)
		result, err := userItr.CreateComment(jsonBlob)
		if string(result) != string(jsonBlob) {
			t.Errorf("result of function fails. wants %v but gets %v", string(jsonBlob), string(result))
		}
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(comment2)
		result, err := userItr.CreateComment(jsonBlob)
		if string(result) != string(jsonBlob) {
			t.Errorf("result of function fails. wants %v but gets %v", string(jsonBlob), string(result))
		}
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("post doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(comment3)
		_, err := userItr.CreateComment(jsonBlob)
		if err != POST_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("user doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(comment4)
		_, err := userItr.CreateComment(jsonBlob)
		if err != USER_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	//create comment reactions
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(commentReaction1)
		result, err := userItr.CreateCommentReaction(jsonBlob)
		if string(result) != string(jsonBlob) {
			t.Errorf("result of function fails. wants %v but gets %v", string(jsonBlob), string(result))
		}
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("ok", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(commentReaction2)
		result, err := userItr.CreateCommentReaction(jsonBlob)
		if string(result) != string(jsonBlob) {
			t.Errorf("result of function fails. wants %v but gets %v", string(jsonBlob), string(result))
		}
		if err != nil {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("comment doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(commentReaction3)
		_, err := userItr.CreateCommentReaction(jsonBlob)
		if err != COMMENT_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
	t.Run("user doesn't exist", func(t *testing.T) {
		jsonBlob, _ := json.Marshal(commentReaction4)
		_, err := userItr.CreateCommentReaction(jsonBlob)
		if err != USER_NOT_EXISTS {
			t.Errorf("error of function fails. wants %v but get %v", nil, err)
		}
	})
}

type CreateMockRepository struct {
	Users    []entity.User
	Posts    []entity.Post
	Comments []entity.Comment
}

func (m *CreateMockRepository) Create(jsonBlob []byte) ([]byte, error) {
	r := Request{}
	json.Unmarshal(jsonBlob, &r)
	switch r.Header.Entity {
	case "User":
		if m.userExists(r.Body.User) {
			return nil, USER_EXISTS
		}
		m.Users = append(m.Users, r.Body.User)
		return json.Marshal(m.Users[r.Body.User.Id-1])
	case "Post":
		if !m.userExists(r.Body.Post.User) {
			return nil, USER_NOT_EXISTS
		}
		m.Posts = append(m.Posts, r.Body.Post)
		return json.Marshal(m.Posts[r.Body.Post.Id-1])
	case "PostReaction":
		postReaction := r.Body.PostReaction
		if !m.userExists(postReaction.Reaction.User) {
			return nil, USER_NOT_EXISTS
		}
		if !m.postExists(postReaction.Post) {
			return nil, POST_NOT_EXISTS
		}
		return json.Marshal(postReaction)
	case "Comment":
		comment := r.Body.Comment
		if !m.userExists(comment.User) {
			return nil, USER_NOT_EXISTS
		}
		if !m.postExists(comment.Post) {
			return nil, POST_NOT_EXISTS
		}
		m.Comments = append(m.Comments, comment)
		return json.Marshal(comment)
	case "CommentReaction":
		commentReaction := r.Body.CommentReaction
		if !m.userExists(commentReaction.Reaction.User) {
			return nil, USER_NOT_EXISTS
		}
		if !m.commentExists(commentReaction.Comment) {
			return nil, COMMENT_NOT_EXISTS
		}
		return json.Marshal(commentReaction)
	default:
		return nil, nil
	}
}
func (m *CreateMockRepository) userExists(user entity.User) bool {
	for _, otherUser := range m.Users {
		if otherUser.Email == user.Email {
			return true
		}
	}
	return false
}

func (m *CreateMockRepository) postExists(post entity.Post) bool {
	for _, otherPost := range m.Posts {
		if otherPost.Id == post.Id {
			return true
		}
	}
	return false
}

func (m *CreateMockRepository) commentExists(comment entity.Comment) bool {
	for _, otherComment := range m.Comments {
		if otherComment.Id == comment.Id {
			return true
		}
	}
	return false
}

func (m *CreateMockRepository) Get([]byte) ([]byte, error)    { return nil, nil }
func (m *CreateMockRepository) Change([]byte) ([]byte, error) { return nil, nil }
func (m *CreateMockRepository) Delete([]byte) error           { return nil }
