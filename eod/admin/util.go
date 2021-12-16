package admin

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

const uidLength = 16

func (a *Admin) CheckUID(uid string) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	_, exists := a.uids[uid]
	return exists
}

func (a *Admin) GetUID() string {
	a.lock.RLock()
	defer a.lock.RUnlock()

	exists := true
	var txt string
	for exists {
		rand.Seed(time.Now().UnixNano())
		b := make([]byte, uidLength)
		rand.Read(b)
		txt = fmt.Sprintf("%x", b)[:uidLength]
		_, exists = a.uids[txt]
	}

	return txt
}

func (a *Admin) Login(c *fiber.Ctx) error {
	password := string(c.Body())
	if password == os.Getenv("PASSWORD") {
		uid := a.GetUID()

		a.lock.Lock()
		a.uids[uid] = empty{}
		a.lock.Unlock()

		c.WriteString(uid)
		return nil
	}
	return errors.New("admin: invalid password")
}
