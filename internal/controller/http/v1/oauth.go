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
	"strings"
)

var oauthState = "pseudo-random"

type OauthContent struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type OauthParams struct {
	AccessToken  string `json:"access_token"`
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

	oauthParams, trouble := h.setParams(apiName)
	if trouble == ErrInternalServer {
		h.Errors(w, http.StatusInternalServerError)
		return
	} else if trouble == ErrPageNotFound {
		h.Errors(w, http.StatusNotFound)
		return
	}
	fmt.Println(oauthParams)
	var buf bytes.Buffer
	buf.WriteString(oauthParams.AuthURL)
	v := url.Values{"response_type": {"code"}, "client_id": {oauthParams.ClientID}}
	v.Set("redirect_uri", oauthParams.CallbackURL)
	v.Set("scope", oauthParams.Scope)
	v.Set("state", oauthState)
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	url := buf.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) setParams(apiName string) (*OauthParams, string) {
	oauthParams := OauthParams{}

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

func (h *Handler) OauthCallbackGoogleHandler(w http.ResponseWriter, r *http.Request) {
	oauthParams, trouble := h.setParams("github")
	if trouble == ErrInternalServer {
		h.Errors(w, http.StatusInternalServerError)
		return
	} else if trouble == ErrPageNotFound {
		h.Errors(w, http.StatusNotFound)
		return
	}

	err := h.exchageCode(r, oauthParams)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackHandler - exchageCode: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	content, err := h.tokenToCall(oauthParams)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	user := entity.User{}

	oauthContent := OauthContent{}
	oauthContents := []OauthContent{}

	err = json.Unmarshal(content, &oauthContent)
	if err != nil {
		json.Unmarshal(content, &oauthContents)
	}

	if oauthContent.Email != "" {
		user.Email = oauthContent.Email
	} else if oauthContents[0].Email != "" {
		user.Email = oauthContents[0].Email
	}

	user.Email = strings.ToLower(user.Email)

	id, err := h.Usecases.Users.GetIdBy(user)
	if id == 0 {
		h.Usecases.Users.SignUp(user)
		id, err = h.Usecases.Users.GetIdBy(user)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackHandler - GetIdBy#1: %w", err))
			fmt.Println("tuta1")
			h.Errors(w, http.StatusInternalServerError)
			return
		}
	}
	if err != nil {
		if !errors.Is(err, entity.ErrUserNotFound) {
			h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackHandler - GetIdBy#2: %w", err))
			fmt.Println("tuta2")
			h.Errors(w, http.StatusInternalServerError)
			return
		}
	}

	err = h.Usecases.Users.SignIn(user)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackHandler - SignIn: %w", err))
		fmt.Println("tuta3")
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	userWithSession, err := h.Usecases.Users.GetSession(id)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackHandler - GetSession: %w", err))
		fmt.Println("tuta4")
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   userWithSession.SessionToken,
		Expires: userWithSession.SessionTTL,
		Path:    "/",
		Domain:  h.Cfg.Server.Host,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) OauthCallbackGithubHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("zashel")
	oauthParams, trouble := h.setParams("github")
	if trouble == ErrInternalServer {
		h.Errors(w, http.StatusInternalServerError)
		return
	} else if trouble == ErrPageNotFound {
		h.Errors(w, http.StatusNotFound)
		return
	}

	err := h.exchageCode(r, oauthParams)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackGithubHandler - exchageCode: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	content, err := h.tokenToCall(oauthParams)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackGithubHandler - tokenToCall: %w", err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	user := entity.User{}

	oauthContent := OauthContent{}
	oauthContents := []OauthContent{}

	err = json.Unmarshal(content, &oauthContent)
	if err != nil {
		json.Unmarshal(content, &oauthContents)
	}

	if oauthContent.Email != "" {
		user.Email = oauthContent.Email
	} else if oauthContents[0].Email != "" {
		user.Email = oauthContents[0].Email
	}

	user.Email = strings.ToLower(user.Email)

	id, err := h.Usecases.Users.GetIdBy(user)
	if id == 0 {
		h.Usecases.Users.SignUp(user)
		id, err = h.Usecases.Users.GetIdBy(user)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackGithubHandler - GetIdBy#1: %w", err))
			fmt.Println("tuta1")
			h.Errors(w, http.StatusInternalServerError)
			return
		}
	}
	if err != nil {
		if !errors.Is(err, entity.ErrUserNotFound) {
			h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackHandler - GetIdBy#2: %w", err))
			fmt.Println("tuta2")
			h.Errors(w, http.StatusInternalServerError)
			return
		}
	}

	err = h.Usecases.Users.SignIn(user)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackHandler - SignIn: %w", err))
		fmt.Println("tuta3")
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	userWithSession, err := h.Usecases.Users.GetSession(id)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - OauthCallbackHandler - GetSession: %w", err))
		fmt.Println("tuta4")
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   userWithSession.SessionToken,
		Expires: userWithSession.SessionTTL,
		Path:    "/",
		Domain:  h.Cfg.Server.Host,
	})

	http.Redirect(w, r, "/", http.StatusFound)
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
	// response, err := http.Get(h.Cfg.Oauth.ApiURL + "=" + token.AccessToken)
	// var buf bytes.Buffer
	// buf.WriteString(oauthParams.AccessURL)
	// url := buf.String()

	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("NewRequest: %w", err)
	// }
	// req.Header.Set("Authorization", "Bearer "+oauthParams.AccessToken)

	// response, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("do: %w", err)
	// }
	// defer response.Body.Close()
	// contents, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	// }
	// return contents, nil

	response, err := http.Get(oauthParams.AccessToken + "=" + oauthParams.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}
