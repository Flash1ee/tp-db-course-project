package utilits

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	"tp-db-project/configs"
)

type ExpectedConnections struct {
	SqlConnection *pgxpool.Pool
	PathFiles     string
}

func NewLogger(config *configs.Config) (log *logrus.Logger, closeResource func() error) {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	currentTime := time.Now().In(time.UTC)

	formatted := config.LogAddr + fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second()) + ".log"

	f, err := os.OpenFile(formatted, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}

	logger.SetOutput(f)
	logger.Writer()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger, f.Close
}

func NewPostgresConnection(config *configs.RepositoryConnections) (db *pgxpool.Pool, closeResource func()) {
	conn, err := pgxpool.Connect(context.Background(), config.DataBaseUrl)
	if err != nil {
		logrus.Fatal(err)
	}

	return conn, conn.Close
}
