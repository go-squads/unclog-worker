package analytic

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/go-squads/unclog-worker/models"
	_ "github.com/lib/pq"
)

type (
	LogLevelMetricRepository interface {
		SaveV1(t models.TimberWolf) error
		SaveV2(t models.TimberWolf) error
		GetAlertConfig() ([]AlertConfig, error)
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

func (r *LogLevelRepositoryImpl) SaveV1(t models.TimberWolf) (err error) {
	insertStatement := "insert into log_metrics_v1(timestamp, log_level, quantity)" +
		"values($1, $2, $3)"
	resp, err := r.dbClient.Exec(insertStatement, "now", t.LogLevel, t.Counter)
	if err != nil {
		log.Println(err.Error())
		return
	}
	resp.LastInsertId()
	return
}

func (r *LogLevelRepositoryImpl) SaveV2(t models.TimberWolf) (err error) {
	insertStatement := "insert into log_metrics_v2(timestamp, app_name, node_id, log_level, quantity)" +
		"values($1, $2, $3, $4, $5)"
	resp, err := r.dbClient.Exec(insertStatement, "now", t.ApplicationName, t.NodeId, t.LogLevel, t.Counter)
	if err != nil {
		log.Println(err.Error())
		return
	}
	resp.LastInsertId()
	return
}

func (r *LogLevelRepositoryImpl) GetAlertConfig() (alertsConfigs []AlertConfig, err error) {
	rows, err := r.dbClient.Query("select * from alerts ")

	if err != nil {
		log.Println(err.Error())
		return
	}

	for rows.Next() {
		a := AlertConfig{}

		err := rows.Scan(&a.Id, &a.AppName, &a.LogLevel, &a.Duration, &a.Limit, &a.Callback)

		if err != nil {
			return nil, err
		}

		alertsConfigs = append(alertsConfigs, a)
	}

	return
}
