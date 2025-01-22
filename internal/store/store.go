// Package store provides utilities for managing game scores persistently
// across sessions using a JSON file.
package store

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

// NoScores is the message displayed when no scores are available. It is
// public for testing purposes.
const NoScores = "No scores yet. Please guess the number to start scoring."

// Score represents a player's game performance, including their name,
// difficulty level, number of attempts, and time taken for the session.
type Score struct {
	Player   string        `json:"player"`
	Level    string        `json:"level"`
	Attempts int           `json:"attempts"`
	Time     time.Duration `json:"time"`
}

// Scores is a collection of Score entries, providing a method to format
// them as an ASCII table.
type Scores []Score

// String formats the Scores collection into an ASCII table. If no scores
// are present, it returns a message indicating that no scores are available.
func (s Scores) String() string {
	if len(s) == 0 {
		return NoScores
	}

	var buffer bytes.Buffer
	table := tablewriter.NewWriter(&buffer)
	table.SetHeader([]string{"Player", "Level", "Attempts", "Time"})

	for _, score := range s {
		table.Append([]string{
			score.Player,
			score.Level,
			strconv.Itoa(score.Attempts),
			score.Time.String(),
		})
	}

	table.Render()
	return buffer.String()
}

// Store defines methods for loading and adding scores, facilitating testing.
type Store interface {
	Load() Scores
	Add(score Score) (Scores, error)
}

// ScoresStore manages the file path for storing scores, which must be a
// JSON file.
type ScoresStore struct {
	FilePath string
}

// Load retrieves previously saved scores. If no scores exist, it returns
// an empty Scores collection.
func (s *ScoresStore) Load() Scores {
	byt, err := os.ReadFile(s.FilePath)
	if err != nil {
		return Scores{}
	}

	var scores Scores
	if err = json.Unmarshal(byt, &scores); err != nil {
		return Scores{}
	}

	return scores
}

// Add inserts a new score into the collection, returning the top 10 scores
// sorted by difficulty, attempts, and time. It handles errors from file
// operations and JSON marshaling.
func (s *ScoresStore) Add(score Score) (Scores, error) {
	file, err := os.OpenFile(s.FilePath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return Scores{}, nil
	}
	defer file.Close()

	scores := s.Load()
	for _, s := range scores {
		if s == score {
			return scores, nil
		}
	}

	scores = append(scores, score)

	byt, err := json.Marshal(scores)
	if err != nil {
		return Scores{}, err
	}

	_, err = file.Write(byt)
	if err != nil {
		return Scores{}, err
	}

	scores.sort()

	if len(scores) > 10 {
		return scores[0:10], nil
	}

	return scores, nil
}

func (s *Scores) sort() {
	levelOrder := map[string]int{
		"Hard":   1,
		"Medium": 2,
		"Easy":   3,
	}

	sort.Slice(*s, func(i, j int) bool {
		if levelOrder[(*s)[i].Level] != levelOrder[(*s)[j].Level] {
			return levelOrder[(*s)[i].Level] < levelOrder[(*s)[j].Level]
		}
		if (*s)[i].Attempts != (*s)[j].Attempts {
			return (*s)[i].Attempts < (*s)[j].Attempts
		}
		return (*s)[i].Time < (*s)[j].Time
	})
}
