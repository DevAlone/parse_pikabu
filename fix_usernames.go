package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	logging "github.com/op/go-logging"

	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
	"golang.org/x/sync/semaphore"
)

func fixUsernames() {
	file, err := os.OpenFile("logs/load_from_old_db.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		helpers.PanicOnError(err)
	}
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, logger.LogFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)

	if config.Settings.Debug {
		logging.SetLevel(logging.DEBUG, "parse_pikabu")
	} else {
		logging.SetLevel(logging.WARNING, "parse_pikabu")
	}

	logger.Log.Debug("load_from_old_db started")

	fmt.Println("fixing usernames...")
	printTimeSinceStart()

	err = models.InitDb()
	if err != nil {
		helpers.PanicOnError(err)
	}

	fmt.Println("processing users...")
	printTimeSinceStart()

	var (
		maxWorkers = 64 // 128
		sem        = semaphore.NewWeighted(int64(maxWorkers))
	)
	ctx := context.TODO()
	var wg sync.WaitGroup

	offset := 0
	limit := 32
	for true {
		fmt.Printf("processing users, offset=%v, limit=%v\n", offset, limit)
		printTimeSinceStart()

		var users []models.PikabuUser
		err := models.Db.Model(&users).
			Where("pikabu_id > ?", offset).
			Order("pikabu_id").
			Limit(limit).
			Select()

		if err == pg.ErrNoRows || len(users) == 0 {
			break
		}
		helpers.PanicOnError(err)

		for _, user := range users {
			offset = int(math.Max(float64(offset), float64(user.PikabuId)))
			helpers.PanicOnError(sem.Acquire(ctx, 1))
			wg.Add(1)
			go func(u models.PikabuUser) {
				defer sem.Release(1)
				fixUsernameProcessUser(&u)
				wg.Done()
			}(user)
		}
	}

	// TODO: process skipped
	wg.Wait()
}

func fixUsernameProcessUser(user *models.PikabuUser) {
	var usernameVersions []models.PikabuUserUsernameVersion
	err := models.Db.Model(&usernameVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	usernames := map[string]bool{}

	for _, usernameVersion := range usernameVersions {
		usernames[strings.ToLower(usernameVersion.Value)] = true
	}

	if len(usernames) == 1 {
		for _, usernameVersion := range usernameVersions {
			helpers.PanicOnError(models.Db.Delete(&usernameVersion))
		}
	} else {
		fmt.Printf("NOT fixing user %v with id %v\n", user.Username, user.PikabuId)
	}
}
