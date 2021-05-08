package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

func traceHandler(stack interface{}) {
	body := map[string]string{
		"avatar_url": "",
		"username":   "Nv7 Server Status",
		"content":    fmt.Sprintf("**Bot Error!**\nStack Trace:\n```\n%s\n```", stack),
	}
	dat, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post("https://discord.com/api/webhooks/840345330352652309/Y-oZf-Riw344TkR_gMNELytJgN3nL2P2teIR9__iQ1zRqcQwDXdmDHfZUoobftsOy3th", "application/json", bytes.NewBuffer(dat))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func recoverer() {
	r := recover()
	if r != nil {
		stack := string(debug.Stack())
		traceHandler(stack)
	}
}
