package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"forum/internal/config"
	"forum/internal/entity"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var oauthState = "pseudo-random"

type OauthContent struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type OauthParams struct {
	AccessToken  string `json:"access_token"`
	ApiName      string
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	AccessURL    string
	CallbackURL  string
	Scope        string
}

func (h *Handler) OauthSignHandler(w http.ResponseWriter, r *http.Request) {
	err := config.ReadEnv(".env")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthSignHandler - ReadEnv: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	path := strings.Split(r.URL.Path, "/")
	apiName := path[len(path)-1]
	if r.URL.Path != "/oauth2_signin/"+apiName {
		h.Errors(w, http.StatusNotFound)
		return
	}

	// setting parametrs depends on api name
	oauthParams, trouble := h.setParams(apiName)
	if trouble == ErrInternalServer {
		h.Errors(w, http.StatusInternalServerError)
		return
	} else if trouble == ErrPageNotFound {
		h.Errors(w, http.StatusNotFound)
		return
	}

	// oauth2 consists of three streps:
	// 1. request to get code
	// 2. exchange code for token
	// 3. use token to call API

	// this is the first step
	var buf bytes.Buffer
	buf.WriteString(oauthParams.AuthURL)
	v := url.Values{"response_type": {"code"}, "client_id": {oauthParams.ClientID}}
	v.Set("redirect_uri", oauthParams.CallbackURL)
	v.Set("scope", oauthParams.Scope)
	v.Set("state", oauthState)
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	url := buf.String()

	// respone of this request should call callback handler func
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) OauthCallbackGoogleHandler(w http.ResponseWriter, r *http.Request) {
	oauthParams, trouble := h.setParams("google")

	if trouble == ErrInternalServer {
		h.Errors(w, http.StatusInternalServerError)
		return
	} else if trouble == ErrPageNotFound {
		h.Errors(w, http.StatusNotFound)
		return
	}

	h.OauthSignIn(w, r, oauthParams)
}

func (h *Handler) OauthCallbackGithubHandler(w http.ResponseWriter, r *http.Request) {
	oauthParams, trouble := h.setParams("github")
	if trouble == ErrInternalServer {
		h.Errors(w, http.StatusInternalServerError)
		return
	} else if trouble == ErrPageNotFound {
		h.Errors(w, http.StatusNotFound)
		return
	}

	h.OauthSignIn(w, r, oauthParams)
}

func (h *Handler) OauthSignIn(w http.ResponseWriter, r *http.Request, oauthParams *OauthParams) {
	// this is the second step of oauth2
	err := h.exchageCode(r, oauthParams)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthSignIn %v - exchageCode: %w", oauthParams.ApiName, err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// this is the third step of oauth2
	content, err := h.tokenToCall(oauthParams)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthSignIn %v - exchageCode: %w", oauthParams.ApiName, err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	user := entity.User{}
	oauthContent := OauthContent{}
	oauthContents := []OauthContent{}

	// parsing data
	err = json.Unmarshal(content, &oauthContent)
	if err != nil {
		if err = json.Unmarshal(content, &oauthContents); err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - OauthSignIn %v - json.Unmarshal: %w", oauthParams.ApiName, err))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	}

	// different apis give response in a different way
	// and parsing information gets a bit messy
	if oauthContent.Email != "" {
		user.Email = oauthContent.Email
	} else if oauthContents[0].Email != "" {
		user.Email = strings.ToLower(oauthContents[0].Email)
	}

	// if there is no user with such email in db, register it
	id, err := h.Usecases.Users.GetIdBy(user)
	if err != nil && !errors.Is(err, entity.ErrUserNotFound) {
		h.l.WriteLog(fmt.Errorf("v1 - OauthSignIn GetIdBy: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	if errors.Is(err, entity.ErrUserNotFound) {
		if oauthContent.Name != "" {
			user.Name = oauthContent.Name
		} else {
			// if in data recieved from api user's name is empty,
			// name will be chars before '@' from email
			user.Name = getNameFromEmail(user.Email)
		}
		if err := h.Usecases.Users.SignUp(user); err == entity.ErrUserNameAlreadyExists {
			suffix := 0
			// if there is already user with that name
			// add incrementing integer suffix to it, until register is ok
			for err == entity.ErrUserNameAlreadyExists {
				suffix++
				user.Name = user.Name + strconv.Itoa(suffix)
				err = h.Usecases.Users.SignUp(user)
			}
		}

		// getting new registered user's id
		id, err = h.Usecases.Users.GetIdBy(user)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - OauthSignIn %v - exchageCode: %w", oauthParams.ApiName, err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}
	}

	// generating session token
	err = h.Usecases.Users.SignIn(user)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthSignIn %v - exchageCode: %w", oauthParams.ApiName, err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	// getting generated token from db for saving in cookie
	userWithSession, err := h.Usecases.Users.GetSession(id)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthSignIn %v - exchageCode: %w", oauthParams.ApiName, err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	// saving session token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   userWithSession.SessionToken,
		Expires: userWithSession.SessionTTL,
		Path:    "/",
		Domain:  h.Cfg.Server.Host,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) setParams(apiName string) (*OauthParams, string) {
	oauthParams := OauthParams{}

	// getting token from environment and const variables
	oauthParams.ApiName = apiName
	switch apiName {
	case "google":
		if clientId, ok := os.LookupEnv("GOOGLE_CLIENT_ID"); ok {
			oauthParams.ClientID = clientId
		} else {
			h.l.WriteLog(fmt.Errorf("v1 - OauthSignHandler - setParams - LookupEnv CLIENT_ID"))
			return nil, ErrInternalServer
		}
		if clientSecret, ok := os.LookupEnv("GOOGLE_CLIENT_SECRET"); ok {
			oauthParams.ClientSecret = clientSecret
		} else {
			h.l.WriteLog(fmt.Errorf("v1 - OauthSignHandler - setParams - LookupEnv CLIENT_SECRET"))
			return nil, ErrInternalServer
		}
		oauthParams.AuthURL = GoogleAuthURL
		oauthParams.TokenURL = GoogleTokenURL
		oauthParams.AccessURL = GoogleAccessURL
		oauthParams.Scope = GoogleScope
		oauthParams.CallbackURL = GoogleCallbackURL
	case "github":
		if clientId, ok := os.LookupEnv("GITHUB_CLIENT_ID"); ok {
			oauthParams.ClientID = clientId
		} else {
			h.l.WriteLog(fmt.Errorf("v1 - OauthSignHandler - setParams - LookupEnv CLIENT_ID"))
			return nil, ErrInternalServer
		}
		if clientSecret, ok := os.LookupEnv("GITHUB_CLIENT_SECRET"); ok {
			oauthParams.ClientSecret = clientSecret
		} else {
			h.l.WriteLog(fmt.Errorf("v1 - OauthSignHandler - setParams - LookupEnv CLIENT_SECRET"))
			return nil, ErrInternalServer
		}
		oauthParams.AuthURL = GithubAuthURL
		oauthParams.TokenURL = GithubTokenURL
		oauthParams.AccessURL = GithubAccessURL
		oauthParams.Scope = GithubScope
		oauthParams.CallbackURL = GithubCallbackURL
	default:
		return nil, ErrPageNotFound
	}
	return &oauthParams, ""
}

func (h *Handler) exchageCode(r *http.Request, oauthParams *OauthParams) error {
	state := r.FormValue("state")
	code := r.FormValue("code")
	if state != oauthState {
		return fmt.Errorf("invalid oauth state")
	}

	var buf bytes.Buffer
	buf.WriteString(oauthParams.TokenURL + "?")
	v := url.Values{"grant_type": {"authorization_code"}, "code": {code}}
	v.Set("redirect_uri", oauthParams.CallbackURL)
	v.Set("client_id", oauthParams.ClientID)
	v.Set("client_secret", oauthParams.ClientSecret)
	buf.WriteString(v.Encode())
	url := buf.String()
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bytes, &oauthParams)
	return nil
}

func (h *Handler) tokenToCall(oauthParams *OauthParams) ([]byte, error) {
	var contents []byte
	var err error
	var response *http.Response

	switch oauthParams.ApiName {
	case "google":
		response, err = http.Get(oauthParams.AccessURL + "=" + oauthParams.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("failed getting user info: %s", err.Error())
		}
	case "github":
		req, err := http.NewRequest("GET", oauthParams.AccessURL, nil)
		if err != nil {
			return nil, fmt.Errorf("newRequest: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+oauthParams.AccessToken)
		response, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("do: %w", err)
		}
	default:
		return nil, fmt.Errorf("wrong api name")
	}

	defer response.Body.Close()
	contents, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

func getNameFromEmail(str string) string {
	sl := strings.Split(str, "@")
	if sl[0] != "" {
		return sl[0]
	}
	return "xioa"
}
