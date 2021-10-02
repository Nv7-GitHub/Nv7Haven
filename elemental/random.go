package elemental

import (
	"context"
	"fmt"
	"time"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const randomQuery = `SELECT name FROM suggestions a, (SELECT found FROM users WHERE uid=? LIMIT 1) b WHERE %s AND JSON_CONTAINS(b.found, CONCAT('"', (SELECT elem1 FROM sugg_combos WHERE elem3=a.name LIMIT 1) ,'"'), "$") AND JSON_CONTAINS(b.found, CONCAT('"', (SELECT elem2 FROM sugg_combos WHERE elem3=a.name LIMIT 1) ,'"'), "$") ORDER BY RAND() LIMIT 1`

func (e *Elemental) randomSuggestion(where string, uid string) ([]string, error) {
	isAnarchy := time.Now().Weekday() == anarchyDay
	params := []interface{}{maxVotes, "%\"" + uid + "\"%"}
	if isAnarchy {
		where = "1"
		params = []interface{}{}
	}
	params = append([]interface{}{uid}, params...)

	row := e.db.QueryRow(fmt.Sprintf(randomQuery, where), params...)
	var elem3 string
	err := row.Scan(&elem3)
	if err != nil {
		return []string{}, err
	}

	var elem1, elem2 string
	row = e.db.QueryRow("SELECT elem1, elem2 FROM sugg_combos WHERE elem3 LIKE ?", elem3)
	err = row.Scan(&elem1, &elem2)
	if err != nil {
		return []string{}, err
	}

	return []string{elem1, elem2}, nil
}

func (e *Elemental) RandomLonely(_ context.Context, req *wrapperspb.StringValue) (*pb.RandomCombinationResponse, error) {
	ans, err := e.RandomLonelySuggestion(req.Value)
	return &pb.RandomCombinationResponse{
		Elements: ans,
	}, err
}

// Pretty much the same, just different first line
func (e *Elemental) UpAndComing(_ context.Context, req *wrapperspb.StringValue) (*pb.RandomCombinationResponse, error) {
	ans, err := e.UpAndComingSuggestion(req.Value)
	return &pb.RandomCombinationResponse{
		Elements: ans,
	}, err
}

// RandomLonelySuggestion gets a random lonely suggestion
func (e *Elemental) RandomLonelySuggestion(uid string) ([]string, error) {
	return e.randomSuggestion("votes<? AND voted NOT LIKE ?", uid)
}

// UpAndComingSuggestion suggestion gets a suggestion that needs one vote
func (e *Elemental) UpAndComingSuggestion(uid string) ([]string, error) {
	return e.randomSuggestion("votes=(?-1) AND voted NOT LIKE ?", uid)
}
