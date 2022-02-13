package eod

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
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
	if os.IsNotExist(err) {
		// File doesn't exist, send logs
		logs, err := os.ReadFile("logs.txt")
		if err != nil {
			return
		}
		b.dg.ChannelMessageSendEmbed("840344139870371920", &discordgo.MessageEmbed{
			Title:       "Bot Crash!",
			Description: fmt.Sprintf("```\n%s\n```", string(logs)),
		})
		os.Create("logs.txt") // Reset logs
	}

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

		os.Remove("restartinfo.gob")
	}
}

func (b *EoD) guildupdate(m types.Msg, rsp types.Rsp, optimize bool) {
	b.Data.RLock()
	defer b.Data.RUnlock()

	msgContinuous := "Optimizing"
	msgPast := "Optimized"
	if !optimize {
		msgContinuous = "Recalculating"
		msgPast = "Recalculated"
	}
	id := rsp.Message(fmt.Sprintf("%s [0/%d]...", msgContinuous, len(b.DB)))

	taken := time.Duration(0)
	i := 0
	lastUpdated := 0
	for _, db := range b.DB {
		if (len(db.Elements) > 100) || (i-lastUpdated > 10) { // If it has enough elements to take a significant amount of time
			hasEdited := false
			gld, err := b.dg.Guild(db.Guild)
			if err == nil {
				isCommunity := false
				for _, feature := range gld.Features {
					if feature == "COMMUNITY" {
						isCommunity = true
						break
					}
				}

				if isCommunity {
					b.dg.ChannelMessageEdit(m.ChannelID, id, fmt.Sprintf("<@%s> %s **%s**... [%d/%d]", m.Author.ID, msgContinuous, gld.Name, i+1, len(b.DB)))
					hasEdited = true
				}
			}

			if !hasEdited {
				b.dg.ChannelMessageEdit(m.ChannelID, id, fmt.Sprintf("<@%s> %s... [%d/%d]", m.Author.ID, msgContinuous, i+1, len(b.DB)))
			}

			lastUpdated = i
		}

		start := time.Now()
		var err error
		if optimize {
			err = db.Optimize()
		} else {
			err = db.Recalc()
		}
		rsp.Error(err)
		/*if rsp.Error(err) {
			return
		}*/
		taken += time.Since(start)

		i++
	}

	b.dg.ChannelMessageEdit(m.ChannelID, id, fmt.Sprintf("<@%s> %s in **%s**.", m.Author.ID, msgPast, taken.String()))
}
