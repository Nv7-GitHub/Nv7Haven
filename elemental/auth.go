package elemental

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// AuthResponse is the response to CreateUser, LoginUser, and NewAnonymous
type AuthResponse struct {
	Success bool
	Data    string
}

// CreateUser creates a user
func (e *Elemental) CreateUser(name string, password string) AuthResponse {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return AuthResponse{
			Success: false,
			Data:    err.Error(),
		}
	}

	var uid string
	count := 1
	for count != 0 {
		uid, err = GenerateRandomStringURLSafe(16)
		if err != nil {
			return AuthResponse{
				Success: false,
				Data:    err.Error(),
			}
		}

		// Check if name taken
		res, err := e.db.Query("SELECT COUNT(1) FROM users WHERE uid=? LIMIT 1", name)
		if err != nil {
			return AuthResponse{
				Success: false,
				Data:    err.Error(),
			}
		}
		defer res.Close()
		res.Next()
		err = res.Scan(&count)
		if err != nil {
			return AuthResponse{
				Success: false,
				Data:    err.Error(),
			}
		}
	}

	// Check if name taken
	res, err := e.db.Query("SELECT COUNT(1) FROM users WHERE name=? LIMIT 1", name)
	if err != nil {
		return AuthResponse{
			Success: false,
			Data:    err.Error(),
		}
	}
	defer res.Close()
	res.Next()
	err = res.Scan(&count)
	if err != nil {
		return AuthResponse{
			Success: false,
			Data:    err.Error(),
		}
	}
	if count == 1 {
		return AuthResponse{
			Success: false,
			Data:    "Account already exists!",
		}
	}

	_, err = e.db.Exec("INSERT INTO users VALUES( ?, ?, ?, ? )", name, uid, string(hashedPassword), `["Air", "Earth", "Fire", "Water"]`)
	if err != nil {
		return AuthResponse{
			Success: false,
			Data:    err.Error(),
		}
	}

	return AuthResponse{
		Success: true,
		Data:    uid,
	}
}

func (e *Elemental) LoginUser(name string, password string) AuthResponse {
	// Check if user exists
	res, err := e.db.Query("SELECT COUNT(1) FROM users WHERE name=?", name)
	if err != nil {
		return AuthResponse{
			Success: false,
			Data:    err.Error(),
		}
	}
	defer res.Close()
	var count int
	res.Next()
	err = res.Scan(&count)
	if err != nil {
		return AuthResponse{
			Success: false,
			Data:    err.Error(),
		}
	}
	if count == 0 {
		return AuthResponse{
			Success: false,
			Data:    "Invalid username",
		}
	}

	res, err = e.db.Query("SELECT uid, password FROM users WHERE name=? LIMIT 1", name)
	if err != nil {
		return AuthResponse{
			Success: false,
			Data:    err.Error(),
		}
	}
	defer res.Close()
	var uid string
	var pwd string
	res.Next()
	err = res.Scan(&uid, &pwd)
	if err != nil {
		return AuthResponse{
			Success: false,
			Data:    err.Error(),
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(password))
	if err != nil {
		return AuthResponse{
			Success: false,
			Data:    "Invalid password",
		}
	}

	return AuthResponse{
		Success: true,
		Data:    uid,
	}
}

func (e *Elemental) NewAnonymousUser() AuthResponse {
	count := 1
	var name string
	var err error
	for count != 0 {
		name, err = GenerateRandomStringURLSafe(8)
		if err != nil {
			return AuthResponse{
				Success: false,
				Data:    err.Error(),
			}
		}

		// Check if name taken
		res, err := e.db.Query("SELECT COUNT(1) FROM users WHERE name=? LIMIT 1", name)
		if err != nil {
			return AuthResponse{
				Success: false,
				Data:    err.Error(),
			}
		}
		defer res.Close()
		res.Next()
		err = res.Scan(&count)
		if err != nil {
			return AuthResponse{
				Success: false,
				Data:    err.Error(),
			}
		}
	}
	return AuthResponse{
		Success: true,
		Data:    name,
	}
}

// HTTP Handlers

func (e *Elemental) createUser(c *fiber.Ctx) error {
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return c.JSON(map[string]any{
			"success": false,
			"data":    err.Error(),
		})
	}
	password := string(c.Body())

	resp := e.CreateUser(name, password)
	return c.JSON(map[string]any{
		"success": resp.Success,
		"data":    resp.Data,
	})
}

func (e *Elemental) loginUser(c *fiber.Ctx) error {
	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return c.JSON(map[string]any{
			"success": false,
			"data":    err.Error(),
		})
	}
	password := string(c.Body())
	resp := e.LoginUser(name, password)
	return c.JSON(map[string]any{
		"success": resp.Success,
		"data":    resp.Data,
	})
}

func (e *Elemental) newAnonymousUser(c *fiber.Ctx) error {
	resp := e.NewAnonymousUser()
	return c.JSON(map[string]any{
		"success": resp.Success,
		"data":    resp.Data,
	})
}
