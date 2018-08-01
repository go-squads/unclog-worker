package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-squads/unclog-worker/util"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var basepath = util.GetRootFolderPath()

func connectDatabase() (*sql.DB, error) {
	connection := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"),
		viper.GetString("DB_HOST"), viper.GetString("DB_PORT"), viper.GetString("DB_SSLMODE"))

	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigration(c *cli.Context) {
	db, err := connectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	migrationQueryString, err := GetStringFromFile(basepath + "migration/schema.sql")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(*migrationQueryString)

	if err != nil {
		log.Fatal(err)
	}
}

func GetStringFromFile(filename string) (*string, error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, errors.New("File not found")
	}

	queryString := fmt.Sprintf("%s", content)

	return &queryString, nil
}
