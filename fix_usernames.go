package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"bitbucket.org/d3dev/parse_pikabu/core/logger"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	logging "github.com/op/go-logging"

	"bitbucket.org/d3dev/parse_pikabu/models"
	"github.com/go-pg/pg"
	"golang.org/x/sync/semaphore"
)

var startTimestamp = time.Now().Unix()

func printTimeSinceStart() {
	fmt.Printf("time since start: %v\n", (time.Now().Unix() - startTimestamp))
}

func fixUsernames() {
	file, err := os.OpenFile("logs/fix_usernames.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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
	limit := 1024 * 4
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

	wg.Wait()
}

func fixUsernameProcessUser(user *models.PikabuUser) {
	var usernameVersions []models.PikabuUserUsernameVersion
	err := models.Db.Model(&usernameVersions).
		Where("item_id = ?", user.PikabuId).
		Order("timestamp").
		Select()
	helpers.PanicOnError(err)

	for i, usernameVersion := range usernameVersions {
		// delete versions
		deleteVersions(usernameVersion.ItemId, usernameVersion.Timestamp)
		if i == 0 {
			continue
		}

		previousVersion := usernameVersions[i-1]

		deleteVersions(previousVersion.ItemId, previousVersion.Timestamp)

		if strings.ToLower(previousVersion.Value) == strings.ToLower(usernameVersion.Value) { // &&
			// previousVersion.Value != usernameVersion.Value {

			fmt.Printf("prev version %v curr version %v\n", previousVersion, usernameVersion)
			previousVersion.Value = usernameVersion.Value
			err := models.Db.Update(&previousVersion)
			helpers.PanicOnError(err)
			err = models.Db.Delete(&usernameVersion)
			usernameVersions[i] = previousVersion
			helpers.PanicOnError(err)
		}
	}
}

func deleteVersions(itemId uint64, timestamp models.TimestampType) {
	deleteVersionsTable(&models.PikabuUserRatingVersion{}, itemId, timestamp)
	deleteVersionsTable(&models.PikabuUserNumberOfSubscribersVersion{}, itemId, timestamp)
	deleteVersionsTable(&models.PikabuUserNumberOfCommentsVersion{}, itemId, timestamp)
	deleteVersionsTable(&models.PikabuUserNumberOfStoriesVersion{}, itemId, timestamp)
	deleteVersionsTable(&models.PikabuUserNumberOfHotStoriesVersion{}, itemId, timestamp)
	deleteVersionsTable(&models.PikabuUserNumberOfPlusesVersion{}, itemId, timestamp)
	deleteVersionsTable(&models.PikabuUserNumberOfMinusesVersion{}, itemId, timestamp)

}

func deleteVersionsTable(table interface{}, itemId uint64, timestamp models.TimestampType) {
	_, err := models.Db.Model(table).
		Where("item_id = ? AND timestamp = ?", itemId, timestamp).
		Delete()
	helpers.PanicOnError(err)
}
