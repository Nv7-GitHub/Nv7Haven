package elemental

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
)

// https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

// GenerateRandomStringURLSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

// MY FUNCTIONS BELOW

func (e *Elemental) getSuggestions(elem1 string) (map[string][]string, error) {
	res, err := e.db.Query("SELECT combos FROM suggestion_combos WHERE name=?", elem1)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var data string
	res.Next()
	err = res.Scan(&data)
	if err != nil {
		return nil, err
	}
	var out map[string][]string
	err = json.Unmarshal([]byte(data), &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (e *Elemental) addCombo(elem1 string, elem2 string, out string) error {
	res, err := e.db.Query("SELECT combos FROM element_combos WHERE name=?", elem1)
	if err != nil {
		return err
	}
	defer res.Close()
	var data string
	res.Next()
	err = res.Scan(&data)
	if err != nil {
		return err
	}
	var combos map[string]string
	if data == "" {
		combos = make(map[string]string, 0)
	} else {
		err = json.Unmarshal([]byte(data), &out)
		if err != nil {
			return err
		}
	}
	combos[elem2] = out
	dat, err := json.Marshal(combos)
	if err != nil {
		return err
	}
	_, err = e.db.Exec("UPDATE selement_combos SET combos=? WHERE name=?", string(dat), elem1)
	if err != nil {
		return err
	}
	return nil
}
