package usecase

import (
	"encoding/json"
	"forum/internal/entity"
)

type UserInteractor struct {
	repo Repository
}

func NewInteractor(repo Repository) *UserInteractor {
	return &UserInteractor{repo}
}

type Request struct {
	Header `json:"request_header"`
	Body   `json:"request_body"`
}

type Header struct {
	Entity string `json:"entity_type"`
}

type Body struct {
	User            entity.User            `json:"user_attributes,omitempty"`
	UserFilter      entity.UserFilter      `json:"user_filter,omitempty"`
	Post            entity.Post            `json:"post_attributes,omitempty"`
	PostFilter      entity.PostFilter      `json:"post_filter,omitempty"`
	Category        entity.Category        `json:"category_attributes,omitempty"`
	Comment         entity.Comment         `json:"comment_attributes,omitempty"`
	PostReaction    entity.PostReaction    `json:"post_reaction_attributes,omitempty"`
	CommentReaction entity.CommentReaction `json:"comment_reaction_attributes,omitempty"`
}

func (i *UserInteractor) CreateUser(jsonBlob []byte) ([]byte, error) {
	user_attr := entity.User{}
	json.Unmarshal(jsonBlob, &user_attr)
	header := Header{Entity: "User"}
	body := Body{User: user_attr}
	return i.sendCreateRequest(Request{header, body})
}

func (i *UserInteractor) CreatePost(jsonBlob []byte) ([]byte, error) {
	post_attr := entity.Post{}
	json.Unmarshal(jsonBlob, &post_attr)
	header := Header{Entity: "Post"}
	body := Body{Post: post_attr}
	return i.sendCreateRequest(Request{header, body})
}

func (i *UserInteractor) CreateComment(jsonBlob []byte) ([]byte, error) {
	comment_attrs := entity.Comment{}
	json.Unmarshal(jsonBlob, &comment_attrs)
	header := Header{Entity: "Comment"}
	body := Body{Comment: comment_attrs}
	return i.sendCreateRequest(Request{header, body})
}

func (i *UserInteractor) CreatePostReaction(jsonBlob []byte) ([]byte, error) {
	post_reaction_attrs := entity.PostReaction{}
	json.Unmarshal(jsonBlob, &post_reaction_attrs)
	header := Header{Entity: "PostReaction"}
	body := Body{PostReaction: post_reaction_attrs}
	return i.sendCreateRequest(Request{header, body})
}

func (i *UserInteractor) CreateCommentReaction(jsonBlob []byte) ([]byte, error) {
	comment_reaction_attrs := entity.CommentReaction{}
	json.Unmarshal(jsonBlob, &comment_reaction_attrs)
	header := Header{Entity: "CommentReaction"}
	body := Body{CommentReaction: comment_reaction_attrs}
	return i.sendCreateRequest(Request{header, body})
}

func (i *UserInteractor) sendCreateRequest(r Request) ([]byte, error) {
	jsonBlob, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return i.repo.Create(jsonBlob)
}

func (i *UserInteractor) GetUser(jsonBlob []byte) ([]byte, error) {
	user_attr := entity.UserFilter{}
	json.Unmarshal(jsonBlob, &user_attr)
	header := Header{Entity: "UserFilter"}
	body := Body{UserFilter: user_attr}
	return i.sendGetRequest(Request{header, body})
}

func (i *UserInteractor) GetPost(jsonBlob []byte) ([]byte, error) {
	post_attr := entity.PostFilter{}
	json.Unmarshal(jsonBlob, &post_attr)
	header := Header{Entity: "PostFilter"}
	body := Body{PostFilter: post_attr}
	return i.sendGetRequest(Request{header, body})
}

func (i *UserInteractor) GetCategory(jsonBlob []byte) ([]byte, error) {
	category_attr := entity.Category{}
	json.Unmarshal(jsonBlob, &category_attr)
	header := Header{Entity: "Category"}
	body := Body{Category: category_attr}
	return i.sendGetRequest(Request{header, body})
}

func (i *UserInteractor) GetPostReactions(jsonBlob []byte) ([]byte, error) {
	postReactions_attr := entity.PostReaction{}
	json.Unmarshal(jsonBlob, &postReactions_attr)
	header := Header{Entity: "PostReaction"}
	body := Body{PostReaction: postReactions_attr}
	return i.sendGetRequest(Request{header, body})
}

func (i *UserInteractor) GetCommentReactions(jsonBlob []byte) ([]byte, error) {
	commentReactions_attr := entity.CommentReaction{}
	json.Unmarshal(jsonBlob, &commentReactions_attr)
	header := Header{Entity: "CommentReaction"}
	body := Body{CommentReaction: commentReactions_attr}
	return i.sendGetRequest(Request{header, body})
}

func (i *UserInteractor) GetComments(jsonBlob []byte) ([]byte, error) {
	comments_attr := entity.Comment{}
	json.Unmarshal(jsonBlob, &comments_attr)
	header := Header{Entity: "Comment"}
	body := Body{Comment: comments_attr}
	return i.sendGetRequest(Request{header, body})
}

func (i *UserInteractor) sendGetRequest(r Request) ([]byte, error) {
	jsonBlob, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return i.repo.Get(jsonBlob)
}

func (i *UserInteractor) ChangePost(jsonBlob []byte) ([]byte, error) {
	post_attr := entity.Post{}
	json.Unmarshal(jsonBlob, &post_attr)
	header := Header{Entity: "Post"}
	body := Body{Post: post_attr}
	return i.sendChangeRequest(Request{header, body})
}

func (i *UserInteractor) ChangeComment(jsonBlob []byte) ([]byte, error) {
	comment_attr := entity.Comment{}
	json.Unmarshal(jsonBlob, &comment_attr)
	header := Header{Entity: "Comment"}
	body := Body{Comment: comment_attr}
	return i.sendChangeRequest(Request{header, body})
}

func (i *UserInteractor) ChangePostReaction(jsonBlob []byte) ([]byte, error) {
	post_reaction_attr := entity.PostReaction{}
	json.Unmarshal(jsonBlob, &post_reaction_attr)
	header := Header{Entity: "PostReaction"}
	body := Body{PostReaction: post_reaction_attr}
	return i.sendChangeRequest(Request{header, body})
}

func (i *UserInteractor) ChangeCommentReaction(jsonBlob []byte) ([]byte, error) {
	comment_reaction_attr := entity.CommentReaction{}
	json.Unmarshal(jsonBlob, &comment_reaction_attr)
	header := Header{Entity: "CommentReaction"}
	body := Body{CommentReaction: comment_reaction_attr}
	return i.sendChangeRequest(Request{header, body})
}
func (i *UserInteractor) sendChangeRequest(r Request) ([]byte, error) {
	jsonBlob, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return i.repo.Change(jsonBlob)
}
func (i *UserInteractor) DeleteUser(jsonBlob []byte) error {
	user_attr := entity.User{}
	json.Unmarshal(jsonBlob, &user_attr)
	header := Header{Entity: "User"}
	body := Body{User: user_attr}
	return i.sendDeleteRequest(Request{header, body})
}

func (i *UserInteractor) DeletePost(jsonBlob []byte) error {
	post_attr := entity.Post{}
	json.Unmarshal(jsonBlob, &post_attr)
	header := Header{Entity: "Post"}
	body := Body{Post: post_attr}
	return i.sendDeleteRequest(Request{header, body})
}

func (i *UserInteractor) DeleteComment(jsonBlob []byte) error {
	comment_attr := entity.Comment{}
	json.Unmarshal(jsonBlob, &comment_attr)
	header := Header{Entity: "Comment"}
	body := Body{Comment: comment_attr}
	return i.sendDeleteRequest(Request{header, body})
}

func (i *UserInteractor) DeletePostReaction(jsonBlob []byte) error {
	post_reaction_attr := entity.PostReaction{}
	json.Unmarshal(jsonBlob, &post_reaction_attr)
	header := Header{Entity: "PostReaction"}
	body := Body{PostReaction: post_reaction_attr}
	return i.sendDeleteRequest(Request{header, body})
}

func (i *UserInteractor) DeleteCommentReaction(jsonBlob []byte) error {
	comment_reaction_attr := entity.CommentReaction{}
	json.Unmarshal(jsonBlob, &comment_reaction_attr)
	header := Header{Entity: "CommentReaction"}
	body := Body{CommentReaction: comment_reaction_attr}
	return i.sendDeleteRequest(Request{header, body})
}

func (i *UserInteractor) sendDeleteRequest(r Request) error {
	jsonBlob, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return i.repo.Delete(jsonBlob)
}

// TODO:
// [ ] Get methods
//    [ ] Convert errors to show or not
//    [ ] Convert data from repo (e.g. count posts for users) i.e. from json to struct > do some operation > convert back > send further
// [ ] Unmarshal interface for empty fields (e.g. `User` field in `Post` )
