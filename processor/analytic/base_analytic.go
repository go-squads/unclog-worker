package analytic

import (
	"github.com/go-squads/unclog-worker/models"
	"github.com/robfig/cron"
)

var c cron.Cron

func resetAll(timberWolves []models.TimberWolf) {
	timberWolves = []models.TimberWolf{}
}
