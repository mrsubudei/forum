package v1

import (
	"forum/internal/entity"
	"forum/internal/usecase"
	"net/http"
	"text/template"
)

type Handler struct {
	usecases *usecase.UseCases
}

type Content struct {
	Authorized   bool
	Unauthorized bool
	Admin        bool
	User         entity.User
	Post         entity.Post
	Posts        []entity.Post
	Users        []entity.User
	Message      string
	ErrorMsg     ErrMessage
}

type ErrMessage struct {
	Code    int
	Message string
}

const (
	SessionDomain          = "localhost"
	userNotExist           = "Такого пользователя не существует"
	userPassWrong          = "Неверный пароль, попробуйте ещё раз"
	passwordsNotSame       = "Пароли не совпадают"
	emailFormatWrong       = "Неправильный формат почты"
	userEmailExist         = "Пользователь с такой почтой уже существует"
	userNameExist          = "Пользователь с таким именем уже существует"
	postCategoryRequired   = "Выберите хотя бы одну тему"
	errPageNotFound        = "Страница не найдена"
	errBadRequest          = "Некорректный запрос"
	ErrInternalServer      = "Ошибка сервера"
	errMethodNotAllowed    = "Метод не разрешен"
	errStatusNotAuthorized = "Вы не авторизованы"
	errLowAccessLevel      = "Низкий уровень доступа"
	queryPost              = "post"
	queryComment           = "comment"
	queryLiked             = "liked"
	queryDisliked          = "disliked"
	commandPutLike         = "like"
	commandPutDislike      = "dislike"
	updateQueryInfo        = "info"
	reactionMessageLike    = "\"лайк\""
	reactionMessageDislike = "\"дизлайк\""
)

var errors ErrMessage

func NewHandler(usecases *usecase.UseCases) *Handler {
	return &Handler{
		usecases: usecases,
	}
}

func (h *Handler) Errors(w http.ResponseWriter, errors ErrMessage) {
	html, err := template.ParseFiles("templates/errors.html")
	if err != nil {
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	html.Execute(w, errors)
}
