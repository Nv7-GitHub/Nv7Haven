package eod

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

type Data struct {
	Channel string
	User    string
}

func (b *EoD) restart(m types.Msg, rsp types.Rsp) {
	f, err := os.Create("restartinfo.gob")
	if rsp.Error(err) {
		return
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(Data{m.ChannelID, m.Author.ID})
	if rsp.Error(err) {
		return
	}

	rsp.Message("Restarting...")

	os.Exit(2)
}

func (b *EoD) update(m types.Msg, rsp types.Rsp) {
	ping := fmt.Sprintf("<@%s> ", m.Author.ID)
	b.dg.ChannelMessageSend(m.ChannelID, ping+"Downloading updates...")

	cmd := exec.Command("git", "pull")
	err := cmd.Run()
	if rsp.Error(err) {
		return
	}

	cmdStr := `go build -o main -ldflags="-s -w"`
	if strings.HasPrefix(runtime.GOARCH, "arm") {
		cmdStr += ` -tags="arm_logs"`
	}

	b.dg.ChannelMessageSend(m.ChannelID, ping+"Installing updates...")
	cmd = exec.Command("sh", "-c", cmdStr)
	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf
	cmd.Stderr = buf
	err = cmd.Run()
	if err != nil {
		rsp.ErrorMessage(buf.String())
		return
	}

	// Clear logs file
	f, err := os.Create("logs.txt")
	if rsp.Error(err) {
		return
	}
	f.Close()

	b.restart(m, rsp)
}

func (b *EoD) start() {
	_, err := os.Stat("restartinfo.gob")
	if err == nil {
		f, err := os.Open("restartinfo.gob")
		if err != nil {
			return
		}

		dec := gob.NewDecoder(f)
		var data Data
		err = dec.Decode(&data)
		if err != nil {
			return
		}

		b.dg.ChannelMessageSend(data.Channel, fmt.Sprintf("<@%s> Restarted!", data.User))
	}
}
