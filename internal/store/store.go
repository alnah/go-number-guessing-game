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

const NoScores = "No scores yet. Please guess the number to start scoring."

type Score struct {
	Player   string        `json:"player"`
	Level    string        `json:"level"`
	Attempts int           `json:"attempts"`
	Time     time.Duration `json:"time"`
}

type Scores []Score

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

type Store interface {
	Load() Scores
	Add(score Score) (Scores, error)
}

type ScoresStore struct {
	FilePath string
}

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
