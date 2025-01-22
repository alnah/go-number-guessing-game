package store_test

import (
	"bytes"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	s "github.com/go-number-guessing-game/internal/store"
	"github.com/olekukonko/tablewriter"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationScoresString(t *testing.T) {
	t.Run("return scores table", func(t *testing.T) {
		scores := s.Scores{}
		buffer := bytes.Buffer{}

		table := tablewriter.NewWriter(&buffer)
		table.SetHeader([]string{"Player", "Level", "Attempts", "Time"})
		for i := 1; i <= 10; i++ {
			score := createRandomScore(t)
			table.Append([]string{
				score.Player,
				score.Level,
				strconv.Itoa(score.Attempts),
				score.Time.String(),
			})
			scores = append(scores, score)
		}
		table.Render()

		want := buffer.String()
		got := scores.String()
		assert.Equal(t, want, got)
	})

	t.Run("return message to user when no scores", func(t *testing.T) {
		scores := s.Scores{}

		want := s.NoScores
		got := scores.String()
		assert.Equal(t, want, got)
	})
}

func TestIntegrationScoreStoreLoad(t *testing.T) {
	t.Run("return loaded scores", func(t *testing.T) {
		file := createTempFile(t)
		store := s.ScoreStore{FilePath: file.Name()}
		score := createRandomScore(t)

		want, err := store.Add(score)
		assert.NoError(t, err)

		got := store.Load()
		assert.Equal(t, want, got)
	})

	t.Run("return empty scores when not existing file", func(t *testing.T) {
		store := s.ScoreStore{FilePath: "not_exist.json"}

		want := s.Scores{}
		got := store.Load()
		assert.Equal(t, want, got)
	})
}

func TestIntegrationScoreStoreAdd(t *testing.T) {
	t.Run("add score and return sorted scores list", func(t *testing.T) {
		file := createTempFile(t)
		store := s.ScoreStore{FilePath: file.Name()}
		score1 := s.Score{
			Player:   "Test1",
			Level:    "Hard",
			Attempts: 3,
			Time:     30 * time.Second,
		}
		score2 := s.Score{
			Player:   "Test2",
			Level:    "Hard",
			Attempts: 3,
			Time:     20 * time.Second,
		}
		score3 := s.Score{
			Player:   "Test3",
			Level:    "Hard",
			Attempts: 2,
			Time:     30 * time.Second,
		}
		score4 := s.Score{
			Player:   "Test4",
			Level:    "Medium",
			Attempts: 3,
			Time:     30 * time.Second,
		}
		score5 := s.Score{
			Player:   "Test5",
			Level:    "Easy",
			Attempts: 3,
			Time:     30 * time.Second,
		}

		// Want sorted scores by level, attempts, and time
		want := s.Scores{score3, score2, score1, score4, score5}
		var got s.Scores
		var err error
		for _, s := range want {
			got, err = store.Add(s)
			assert.NoError(t, err)
		}
		assert.Equal(t, want, got)
	})

	t.Run("can't add twice the same score", func(t *testing.T) {
		file := createTempFile(t)
		store := s.ScoreStore{FilePath: file.Name()}
		score := createRandomScore(t)

		want := s.Scores{score}
		var got s.Scores
		var err error
		for i := 1; i <= 2; i++ {
			got, err = store.Add(score)
			assert.NoError(t, err)
		}
		assert.Equal(t, want, got)
	})

	t.Run("add twice same score for different players", func(t *testing.T) {
		file := createTempFile(t)
		store := s.ScoreStore{FilePath: file.Name()}
		score1 := s.Score{
			Player:   "Test1",
			Level:    "Hard",
			Attempts: 3,
			Time:     30 * time.Second,
		}
		score2 := s.Score{
			Player:   "Test2",
			Level:    "Hard",
			Attempts: 3,
			Time:     30 * time.Second,
		}

		want := s.Scores{score1, score2}
		var got s.Scores
		var err error
		for _, s := range want {
			got, err = store.Add(s)
			assert.NoError(t, err)
		}
		assert.Equal(t, want, got)
	})

	t.Run("keep only the 10 best scores", func(t *testing.T) {
		file := createTempFile(t)
		store := s.ScoreStore{FilePath: file.Name()}

		var got s.Scores
		var err error
		for i := 1; i <= 11; i++ {
			score := createRandomScore(t)
			got, err = store.Add(score)
			assert.NoError(t, err)
		}
		assert.Len(t, got, 10)
	})
}

func createTempFile(t *testing.T) *os.File {
	t.Helper()

	file, err := os.CreateTemp("", "data_test.json")
	assert.NoError(t, err)

	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})

	return file
}

func createRandomScore(t *testing.T) s.Score {
	t.Helper()

	player := getRandomPlayer(4)
	level := getRandomLevel()
	attempts := getRandomAttempts(level)
	time := getRandomTime(attempts)

	return s.Score{Player: player, Level: level, Attempts: attempts, Time: time}
}

func getRandomPlayer(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getRandomLevel() string {
	levels := []string{"Hard", "Medium", "Easy"}
	index := rand.Intn(3)
	return levels[index]
}

func getRandomAttempts(level string) int {
	levelToMaxAttempts := map[string]int{
		"Hard":   3,
		"Medium": 5,
		"Easy":   10,
	}

	maxAttempts := levelToMaxAttempts[level]
	return rand.Intn(maxAttempts) + 1
}

func getRandomTime(attempts int) time.Duration {
	return time.Duration(attempts) * 10 * time.Second
}
