package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/r3labs/sse/v2"
)

type Service struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Running  bool   `json:"running"`
	Building bool   `json:"building"`

	Cmd    *exec.Cmd `json:"-"`
	Output *Output   `json:"-"`
}

type Output struct {
	Content *strings.Builder

	Cond   *sync.Cond
	Data   []byte
	Finish chan struct{}
}

func (s *Output) Write(p []byte) (n int, err error) {
	s.Data = p
	s.Cond.L.Lock()
	s.Cond.Broadcast()
	s.Cond.L.Unlock()

	return s.Content.Write(p)
}

var lock = &sync.Mutex{}
var services = map[string]*Service{
	"test": {
		ID:   "test",
		Name: "Test",
		Output: &Output{
			Content: &strings.Builder{},
			Cond:    sync.NewCond(&sync.Mutex{}),
			Finish:  make(chan struct{}),
		},
	},
}

func marshalServices() []byte {
	arr := make([]*Service, 0, len(services))
	lock.Lock()
	for _, s := range services {
		arr = append(arr, s)
	}
	lock.Unlock()
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Name < arr[j].Name
	})

	v, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	return v
}

func PublishServices() {
	events.Publish("services", &sse.Event{
		Data: marshalServices(),
	})
}

func Build(s *Service) error {
	s.Building = true
	PublishServices()
	defer func() {
		s.Building = false
		PublishServices()
	}()

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Pull
	err = exec.Command("git", "pull").Run()
	if err != nil {
		return err
	}

	// Build
	cmd := exec.Command("go", "build", "-o", filepath.Join(wd, "build", s.ID))
	cmd.Dir = filepath.Join(wd, "run", s.ID)
	cmd.Stderr = &strings.Builder{}
	return cmd.Run()
}

func autoRestart(s *Service, wd string) {
	err := s.Cmd.Wait()
	if s.Running {
		s.Output.Write([]byte("Server crashed with error: " + err.Error() + "\n"))
		s.Cmd = exec.Command(filepath.Join(wd, "build", s.ID))
		s.Cmd.Stdout = s.Output
		s.Cmd.Stderr = s.Output
		err = s.Cmd.Start()
		if err != nil {
			s.Output.Write([]byte("Server couldn't start with error: " + err.Error() + "\n"))
			s.Running = false
			PublishServices()
		} else {
			go autoRestart(s, wd)
		}
	} else {
		s.Output.Finish <- struct{}{}
	}
}

func Run(s *Service) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	s.Cmd = exec.Command(filepath.Join(wd, "build", s.ID))
	s.Output.Content = &strings.Builder{}
	s.Cmd.Stdout = s.Output
	s.Cmd.Stderr = s.Output

	err = s.Cmd.Start()
	if err != nil {
		return err
	}
	s.Running = true
	PublishServices()

	go autoRestart(s, wd)
	return nil
}

func Stop(s *Service) error {
	s.Running = false
	err := s.Cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}
	<-s.Output.Finish
	s.Cmd = nil
	PublishServices()
	return nil
}
