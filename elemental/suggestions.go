package elemental

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const minVotes = -1
const maxVotes = 2
const anarchyDay = time.Friday

func (e *Elemental) getSugg(id string) (*pb.Suggestion, error) {
	row := e.db.QueryRow("SELECT * FROM suggestions WHERE name=?", id)
	suggestion := &pb.Suggestion{}
	var color string
	var voted string
	err := row.Scan(&suggestion.Name, &color, &suggestion.Creator, &voted, &suggestion.Votes)
	if err != nil {
		return &pb.Suggestion{}, err
	}

	colors := strings.Split(color, "_")
	sat, err := strconv.ParseFloat(colors[1], 32)
	if err != nil {
		return &pb.Suggestion{}, err
	}
	light, err := strconv.ParseFloat(colors[2], 32)
	if err != nil {
		return &pb.Suggestion{}, err
	}
	suggestion.Color = &pb.Color{
		Base:       colors[0],
		Saturation: float32(sat),
		Lightness:  float32(light),
	}

	var votedData []string
	err = json.Unmarshal([]byte(voted), &votedData)
	if err != nil {
		return &pb.Suggestion{}, err
	}
	suggestion.Voted = votedData

	return suggestion, nil
}

func (e *Elemental) GetSuggestion(ctx context.Context, id *wrapperspb.StringValue) (*pb.Suggestion, error) {
	suggestion, err := e.getSugg(id.Value)
	if err != nil {
		if err.Error() == "null" {
			return &pb.Suggestion{}, errors.New("null")
		}
	}
	return suggestion, err
}

func (e *Elemental) GetSuggestionCombos(ctx context.Context, req *pb.Combination) (*pb.SuggestionCombinationResponse, error) {
	data, err := e.GetSuggestions(req.Elem1, req.Elem2)
	return &pb.SuggestionCombinationResponse{
		Suggestions: data,
	}, err
}

func (e *Elemental) DownSuggestion(ctx context.Context, req *pb.SuggestionRequest) (*pb.VoteResponse, error) {
	suc, msg := e.DownvoteSuggestion(req.Element, req.Uid)
	if !suc {
		return &pb.VoteResponse{}, errors.New(msg)
	}

	return &pb.VoteResponse{
		Create: false,
	}, nil
}

// DownvoteSuggestion downvotes a suggestion
func (e *Elemental) DownvoteSuggestion(id, uid string) (bool, string) {
	existing, err := e.getSugg(id)
	if err != nil {
		return false, err.Error()
	}
	for _, voted := range existing.Voted {
		if voted == uid {
			return false, "You already voted!"
		}
	}
	existing.Votes--
	if existing.Votes < minVotes {
		e.db.Exec("DELETE FROM suggestions WHERE name=?", id)
		return true, ""
	}
	existing.Voted = append(existing.Voted, uid)
	data, err := json.Marshal(existing.Voted)
	if err != nil {
		return false, err.Error()
	}
	_, err = e.db.Exec("UPDATE suggestions SET voted=?, votes=? WHERE name=?", data, existing.Votes, existing.Name)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func (e *Elemental) UpSuggestion(ctx context.Context, req *pb.SuggestionRequest) (*pb.VoteResponse, error) {
	create, suc, msg := e.UpvoteSuggestion(req.Element, req.Uid)
	if !suc {
		return &pb.VoteResponse{}, errors.New(msg)
	}

	return &pb.VoteResponse{
		Create: create,
	}, nil
}

// UpvoteSuggestion upvotes a suggestion
func (e *Elemental) UpvoteSuggestion(id, uid string) (bool, bool, string) {
	existing, err := e.getSugg(id)
	if err != nil {
		return false, false, err.Error()
	}

	isAnarchy := time.Now().Weekday() == anarchyDay
	if !(isAnarchy) {
		for _, voted := range existing.Voted {
			if voted == uid {
				return false, false, "You already voted!"
			}
		}
	}

	existing.Votes++
	existing.Voted = append(existing.Voted, uid)
	data, err := json.Marshal(existing.Voted)
	if err != nil {
		return false, false, err.Error()
	}
	_, err = e.db.Exec("UPDATE suggestions SET votes=?, voted=? WHERE name=?", existing.Votes, data, existing.Name)
	if err != nil {
		return false, false, err.Error()
	}
	if (existing.Votes >= maxVotes) || isAnarchy {
		return true, true, ""
	}
	return false, true, ""
}

func (e *Elemental) NewSugg(ctx context.Context, req *pb.NewSuggestionRequest) (*pb.VoteResponse, error) {
	create, err := e.NewSuggestion(req.Elem1, req.Elem2, req.Suggestion)
	return &pb.VoteResponse{
		Create: create,
	}, err
}

// NewSuggestion makes a new suggestion
func (e *Elemental) NewSuggestion(elem1, elem2 string, suggestion *pb.Suggestion) (bool, error) {
	voted, _ := json.Marshal(suggestion.Voted)
	color := fmt.Sprintf("%s_%f_%f", suggestion.Color.Base, suggestion.Color.Saturation, suggestion.Color.Lightness)
	_, err := e.db.Exec("INSERT INTO suggestions VALUES( ?, ?, ?, ?, ? )", suggestion.Name, color, suggestion.Creator, voted, suggestion.Votes)
	if err != nil {
		return false, err
	}

	_, err = e.db.Exec("INSERT INTO sugg_combos VALUES ( ?, ?, ? )", elem1, elem2, suggestion.Name)
	if err != nil {
		return false, err
	}
	if time.Now().Weekday() == anarchyDay {
		return true, nil
	}
	return false, nil
}
