package v1

import "forum/internal/entity"

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

type Key string

const (
	SessionDomain          = "localhost"
	UserNotExist           = "Такого пользователя не существует"
	UserPassWrong          = "Неверный пароль, попробуйте ещё раз"
	PasswordsNotSame       = "Пароли не совпадают"
	EmailFormatWrong       = "Неправильный формат почты"
	UserEmailExist         = "Пользователь с такой почтой уже существует"
	UserNameExist          = "Пользователь с таким именем уже существует"
	PostCategoryRequired   = "Выберите хотя бы одну тему"
	ErrPageNotFound        = "Страница не найдена"
	ErrBadRequest          = "Некорректный запрос"
	ErrInternalServer      = "Ошибка сервера"
	ErrMethodNotAllowed    = "Метод не разрешен"
	ErrStatusNotAuthorized = "Вы не авторизованы"
	ErrLowAccessLevel      = "Низкий уровень доступа"
	QueryPost              = "post"
	QueryComment           = "comment"
	QueryLiked             = "liked"
	QueryDisliked          = "disliked"
	CommandPutLike         = "like"
	CommandPutDislike      = "dislike"
	UpdateQueryInfo        = "info"
	ReactionMessageLike    = "\"лайк\""
	ReactionMessageDislike = "\"дизлайк\""
)

var errors ErrMessage
