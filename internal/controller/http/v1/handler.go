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
	Posts        []entity.Post
	Users        []entity.User
	ErrorMsg     ErrMessage
}

type ContentSingle struct {
	Authorized   bool
	Unauthorized bool
	Admin        bool
	Post         entity.Post
	User         entity.User
	ErrorMsg     ErrMessage
}

type ErrMessage struct {
	Code    int
	Message string
}

var (
	// DateFormat             = "2006-01-02"
	sessionDomain          = "localhost"
	userNotExist           = "Такого пользователя не существует"
	userPassWrong          = "Неверный пароль, попробуйте ещё раз"
	passwordsNotSame       = "Пароли не совпадают"
	emailFormatWrong       = "Неправильный формат почты"
	userEmailExist         = "Пользователь с такой почтой уже существует"
	userNameExist          = "Пользователь с таким именем уже существует"
	postCategoryRequired   = "Выберите хотя бы одну тему"
	errPageNotFound        = "Страница не найдена"
	errBadRequest          = "Некорректный запрос"
	errInternalServer      = "Ошибка сервера"
	errMethodNotAllowed    = "Метод не разрешен"
	errStatusNotAuthorized = "Вы не авторизованы"
	commandPutLike         = "like"
	commandPutDislike      = "dislike"
	errors                 ErrMessage
)

func NewHandler(services *usecase.UseCases) *Handler {
	return &Handler{
		usecases: services,
	}
}

func (h *Handler) Errors(w http.ResponseWriter, errors ErrMessage) {
	html, err := template.ParseFiles("templates/errors.html")
	if err != nil {
		http.Error(w, errInternalServer, http.StatusInternalServerError)
		return
	}

	html.Execute(w, errors)
}
