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
	OwnerId      int64
	ErrorMsg     ErrMessage
}

type ErrMessage struct {
	Code    int
	Message string
}

type Key string

const (
	UserNotExist          = "Такого пользователя не существует"
	UserPassWrong         = "Неверный пароль, попробуйте ещё раз"
	PasswordsNotSame      = "Пароли не совпадают"
	EmailFormatWrong      = "Неправильный формат почты"
	UserEmailAlreadyExist = "Пользователь с такой почтой уже существует"
	UserNameAlreadyExist  = "Пользователь с таким именем уже существует"
	PostCategoryRequired  = "Выберите хотя бы одну тему"
)

const (
	ErrNoRowsInResult      = "no rows in result set"
	ErrPageNotFound        = "Страница не найдена"
	ErrBadRequest          = "Некорректный запрос"
	ErrInternalServer      = "Ошибка сервера"
	ErrMethodNotAllowed    = "Метод не разрешен"
	ErrStatusNotAuthorized = "Вы не авторизованы"
	ErrLowAccessLevel      = "Низкий уровень доступа"
)

const (
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

const (
	GoogleAuthURL   = "https://accounts.google.com/o/oauth2/auth"
	GoogleTokenURL  = "https://oauth2.googleapis.com/token"
	GoogleAccessURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token"
	GoogleScope     = "https://www.googleapis.com/auth/userinfo.email " +
		"https://www.googleapis.com/auth/userinfo.profile"
	GoogleCallbackURL = "http://localhost:8087/oauth2_callback_google"

	GithubAuthURL     = "https://github.com/login/oauth/authorize"
	GithubTokenURL    = "https://github.com/login/oauth/access_token"
	GithubAccessURL   = "https://api.github.com/user/emails"
	GithubScope       = "user"
	GithubCallbackURL = "http://localhost:8087/oauth2_callback_github"
)
