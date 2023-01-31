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

var OauthState = "pseudo-random-fs3f#ds38A@f"

type OauthContent struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type OauthParams struct {
	AccessToken  string `json:"access_token"`
	ApiName      string
	ClientID     string
	ClientSecret string
	OauthURLS    OauthURLs
}

type OauthURLs struct {
	Auth     string
	Token    string
	Access   string
	Callback string
	Scope    string
}

var GoogleOauthURLs = OauthURLs{
	Auth:     "https://accounts.google.com/o/oauth2/auth",
	Token:    "https://oauth2.googleapis.com/token",
	Access:   "https://www.googleapis.com/oauth2/v2/userinfo?access_token",
	Callback: "http://localhost:8087/oauth2_callback_google",
	Scope: "https://www.googleapis.com/auth/userinfo.email " +
		"https://www.googleapis.com/auth/userinfo.profile",
}

var GithubOauthURLs = OauthURLs{
	Auth:     "https://github.com/login/oauth/authorize",
	Token:    "https://github.com/login/oauth/access_token",
	Access:   "https://api.github.com/user/emails",
	Callback: "http://localhost:8087/oauth2_callback_github",
	Scope:    "user",
}

var MailruOauthURLs = OauthURLs{
	Auth:     "https://oauth.mail.ru/login",
	Token:    "https://oauth.mail.ru/token",
	Access:   "https://oauth.mail.ru/userinfo?access_token",
	Callback: "http://localhost:8087/oauth2_callback_mailru",
	Scope:    "userinfo",
}
