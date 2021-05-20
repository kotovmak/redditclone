package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// databaseURL = "host=localhost dbname=reddit_test port=5433 sslmode=disable user=user password=password"
		databaseURL = "postgresql://user:password@localhost:5433/reddit_test?sslmode=disable"
	}

	os.Exit(m.Run())
}
