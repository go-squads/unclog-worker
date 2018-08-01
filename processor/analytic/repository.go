package analytic

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"

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
	conn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"),
		viper.GetString("DB_HOST"), viper.GetString("DB_PORT"), viper.GetString("DB_SSLMODE"))

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
	insertStatement := "insert into log_metrics(timestamp, log_level, log_level_id, quantity) values($1,$2,$3,$4)"
	resp, err := r.dbClient.Exec(insertStatement, "now()", metric.logLevel, metric.logLevelId, metric.count)
	if err != nil {
		log.Print("\n\n", err.Error(), "\n\n")
		return err
	}
	resp.LastInsertId()
	return nil
}
