package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"forum/internal/entity"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var oauthStateString = "pseudo-random"

type GoogleContent struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GoogleTokens struct {
	AccessToken string `json:"access_token,omitempty"`
	Tokens      struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"google_tokens,omitempty"`
}

func (h *Handler) SignInGoogleHandler(w http.ResponseWriter, r *http.Request) {
	configFile, err := os.Open("tokens.json")
	if err != nil {
		log.Fatal(err)
	}
	data, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()
	tokens := GoogleTokens{}
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		log.Fatal(err)
	}
	h.Cfg.GoogleTokens.ClientID = tokens.Tokens.ClientID
	h.Cfg.GoogleTokens.ClientSecret = tokens.Tokens.ClientSecret

	var buf bytes.Buffer
	buf.WriteString("https://accounts.google.com/o/oauth2/auth")
	v := url.Values{"response_type": {"code"}, "client_id": {h.Cfg.GoogleTokens.ClientID}}
	v.Set("redirect_uri", h.Cfg.Server.Host+":8087/google_callback")
	v.Set("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
	v.Set("state", oauthStateString)
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	url := buf.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	content, err := h.getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	userData := GoogleContent{}
	json.Unmarshal(content, &userData)
	user := entity.User{
		Email: userData.Email,
	}

	id, err := h.Usecases.Users.GetIdBy(user)
	if id == 0 {
		user.Name = userData.Name
		h.Usecases.Users.SignUp(user)
		id, err = h.Usecases.Users.GetIdBy(user)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - GoogleCallbackHandler - GetIdBy#1: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}
	}
	if err != nil {
		if !errors.Is(err, entity.ErrUserNotFound) {
			h.l.WriteLog(fmt.Errorf("v1 - GoogleCallbackHandler - GetIdBy#2: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}
	}

	err = h.Usecases.Users.SignIn(user)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - GoogleCallbackHandler - SignIn: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return
	}

	userWithSession, err := h.Usecases.Users.GetSession(id)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - GoogleCallbackHandler - GetSession: %w", err))
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

func (h *Handler) getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := h.exchange(code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}

func (h *Handler) exchange(code string) (GoogleTokens, error) {
	var buf bytes.Buffer
	buf.WriteString("https://oauth2.googleapis.com/token?")
	v := url.Values{"grant_type": {"authorization_code"}, "code": {code}}
	v.Set("redirect_uri", h.Cfg.Server.Host+":8087/google_callback")
	v.Set("client_id", h.Cfg.GoogleTokens.ClientID)
	v.Set("client_secret", h.Cfg.GoogleTokens.ClientSecret)
	buf.WriteString(v.Encode())
	url := buf.String()
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return GoogleTokens{}, err
	}
	defer resp.Body.Close()
	var token GoogleTokens
	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &token)
	return token, nil
}
