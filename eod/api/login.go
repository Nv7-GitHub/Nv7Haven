package api

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Nv7-Github/Nv7Haven/eod/api/data"
	"golang.org/x/oauth2"
)

func (a *API) MethodGuild(vals map[string]any, id string) data.Response {
	// Process params
	gldV, ok := vals["gld"]
	if !ok {
		return data.RSPBadRequest
	}
	gld, ok := gldV.(string)
	if !ok {
		return data.RSPBadRequest
	}

	// Check
	db, res := a.GetDB(gld)
	if !res.Exists {
		return data.RSPError(res.Message)
	}
	db.RLock()
	_, exists := db.Invs()[id]
	db.RUnlock()
	if !exists {
		return data.RSPError(fmt.Sprintf("User %s doesn't have an inventory in guild %s!", id, gld))
	}
	return data.RSPSuccess(map[string]any{})
}

var conf = &oauth2.Config{
	RedirectURL:  "https://http.nv7haven.com/eode/oauth",
	ClientID:     "964274065508556800",
	ClientSecret: "dRDvGpuHZAgH6u-F5_UHakTxZgLewhe4",
	Scopes:       []string{"identify"},
	Endpoint: oauth2.Endpoint{
		AuthURL:   "https://discord.com/api/oauth2/authorize",
		TokenURL:  "https://discord.com/api/oauth2/token",
		AuthStyle: oauth2.AuthStyleInParams,
	},
}

func (a *API) GenURL() (data.Response, string) {
	// Gen state
	var state string
	exists := true
	a.loginLock.RLock()
	var err error
	for exists {
		state, err = genState(stateLength)
		if err != nil {
			return data.RSPError(err.Error()), ""
		}
		_, exists = a.loginLinks[state]
	}
	a.loginLock.RUnlock()

	// Gen chan
	resp := make(chan string)
	a.loginLock.Lock()
	a.loginLinks[state] = resp
	a.loginLock.Unlock()

	// Gen URL
	url := conf.AuthCodeURL(state)
	return data.RSPSuccess(map[string]any{"url": url}), state
}

const stateLength = 10

// https://stackoverflow.com/questions/35558166/when-to-randomize-auth-code-state-in-oauth2
func genState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(data), nil
}

func http_err(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}

func (a *API) HandleOAuth(w http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")

	a.loginLock.RLock()
	ch, exists := a.loginLinks[state]
	a.loginLock.RUnlock()
	if !exists {
		http_err(w, "Invalid state")
		return
	}

	// Exchange code for access token
	token, err := conf.Exchange(context.Background(), req.FormValue("code"))
	if err != nil {
		http_err(w, err.Error())
		return
	}

	// Get info
	res, err := conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
	if err != nil || res.StatusCode != 200 {
		http_err(w, "Failed to get user info")
		return
	}
	defer res.Body.Close()

	// Parse user info
	var data map[string]interface{}
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&data); err != nil {
		http_err(w, err.Error())
		return
	}

	// Login
	ch <- data["id"].(string)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged in successfully."))
}
