package admin

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/vcomm"
	"github.com/Nv7-Github/vcomm/definitions"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Admin struct {
	data *eodb.Data
}

func (a *Admin) Config(guild string) (string, error) {
	v, res := a.data.GetDB(guild)
	if !res.Exists {
		return "", errors.New(res.Message)
	}
	d, err := json.Marshal(v.Config)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func InitAdmin(data *eodb.Data) {
	a := &Admin{
		data: data,
	}
	serv := vcomm.NewVComm(a)
	// if my machine
	if runtime.GOOS == "darwin" {
		home, _ := os.UserHomeDir()
		f, err := os.Create(filepath.Join(home, "Documents", "Coding", "eodmin", "src", "lib", "bindings.ts"))
		if err != nil {
			panic(err)
		}
		def, err := serv.CreateDefinitions()
		if err != nil {
			panic(err)
		}
		ts := definitions.GenerateTypescript(def)
		_, err = f.WriteString(ts)
		if err != nil {
			panic(err)
		}
		f.Close()
	}

	// Websocket handler
	http.HandleFunc("/eodconsole", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			res := serv.Message(string(message))

			err = c.WriteMessage(mt, []byte(res))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	})

}
