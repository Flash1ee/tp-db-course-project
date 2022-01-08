package utilits

import (
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"io"
	"tp-db-project/configs"
)
const (
	MaxCountConnections = 100
)
type ExpectedConnections struct {
	SqlConnection *pgx.ConnPool
	PathFiles     string
}

func NewLogger(config *configs.Config) (log *logrus.Logger) {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	//currentTime := time.Now().In(time.UTC)

	//formatted := config.LogAddr + fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
	//	currentTime.Year(), currentTime.Month(), currentTime.Day(),
	//	currentTime.Hour(), currentTime.Minute(), currentTime.Second()) + ".log"

	//f, err := os.OpenFile(formatted, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	//if err != nil {
	//	logrus.Fatalf("error opening file: %v", err)
	//}

	logger.SetOutput(io.Discard)
	logger.Writer()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}

func NewPostgresConnection(config *configs.RepositoryConnections) (db *pgx.ConnPool, closeResource func()) {
	conn, err := pgx.ParseConnectionString(config.DataBaseUrl)
	if err != nil {
		logrus.Fatal(err)
	}
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: MaxCountConnections,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	return pool, pool.Close
}
