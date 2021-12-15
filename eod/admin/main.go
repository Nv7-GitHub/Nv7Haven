package admin

import (
	"log"
	"net/http"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/vcomm"
)

type Admin struct {
	*eodb.Data
}

func InitAdmin(data *eodb.Data) {
	a := &Admin{
		Data: data,
	}
	serv := vcomm.NewVComm(a)

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
