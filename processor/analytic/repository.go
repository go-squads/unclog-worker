package analytic

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type (
	LogLevelMetricRepository interface {
		Save(metric LogLevelMetric) error
	}

	LogLevelMetric struct {
		logLevel   string
		logLevelId int
		count      uint32
	}

	LogLevelRepositoryImpl struct {
		dbClient *sql.DB
	}
)

func NewLogLevelRepositoryImpl() *LogLevelRepositoryImpl {
	return &LogLevelRepositoryImpl{
		dbClient: NewDbClient(),
	}
}

func NewDbClient() *sql.DB {

	conn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		"pt.go-jekindonesia", "1", "postgres", "localhost", "5432")

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("Failed to connect db: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping db: %v\n", err)
	}

	return db
}

func (r *LogLevelRepositoryImpl) Save(metric LogLevelMetric) error {
	insertStatement := "insert into logs_schema.log_overview(timestamp,loglevel, loglevelId, count) values($1,$2,$3,$4)"
	resp, err := r.dbClient.Exec(insertStatement, "now()", metric.logLevel, metric.logLevelId, metric.count)
	if err != nil {
		log.Print("\n\n", err.Error(), "\n\n")
		return err
	}
	resp.LastInsertId()
	return nil
}
