package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"server-updater/aptclient"
	"server-updater/discordclient"

	"github.com/joho/godotenv"
)

//go:embed start-screen.txt
var startScreen string

var version = "unset"
var commit = "unset"
var buildDate = "unset"

func getEnvOrFail(envName string) string {
	value := os.Getenv(envName)
	if len(value) == 0 {
		log.Fatalf("environment variable could not be loaded. variable: %s", envName)
	}
	return value
}

func loadDotenvOrFail() {
	envFileLocation := flag.String("env-file", "", "env file location")
	flag.Parse()

	if *envFileLocation == "" {
		log.Println("not loading .env file. using process env")
		return
	}

	log.Printf("loading env file from %s", *envFileLocation)

	err := godotenv.Load(*envFileLocation)
	if err != nil {
		log.Fatalf("err could not load .env file: %s with godotenv. err: %v", *envFileLocation, err)
	}
}

func main() {
	start := time.Now()
	log.Printf("\n%s\nversion: %s(%s) built: %s\n", startScreen, version, commit, buildDate)

	loadDotenvOrFail()

	webhookURLInfo := getEnvOrFail("WEBHOOK_URL_INFO")
	webhookURLError := getEnvOrFail("WEBHOOK_URL_ERROR")

	apt := aptclient.New()
	discord := discordclient.New("Update Bot", webhookURLError, webhookURLInfo)

	log.Println("start: update")
	updateCount, err := apt.Update()
	if err != nil {
		discord.LogFatalIgnoreError(fmt.Sprintf("could not update. err: %v", err))
	}
	log.Printf("%d packages need to be updated", updateCount)

	log.Println("start: upgrade")
	upgradeInfo, err := apt.Upgrade()
	if err != nil {
		discord.LogFatalIgnoreError(fmt.Sprintf("could not upgrade. err: %v", err))
	}

	log.Println("start: autoremove")
	if err := apt.Autoremove(); err != nil {
		discord.LogFatalIgnoreError(fmt.Sprintf("could not autoremove. err: %v", err))
	}

	rebootRequired := isRebootRequired(discord)

	discord.LogInfoIgnoreError(
		fmt.Sprintf("reboot required: %t, execution time: %s. %s\n",
			rebootRequired, time.Since(start), upgradeInfo))
	log.Println("finished server-update")
}

func isRebootRequired(discord *discordclient.DiscordClient) bool {
	log.Println("start: reboot check")
	if _, err := os.Stat("/var/run/reboot-required"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		discord.LogFatalIgnoreError(fmt.Sprintf("could not read \"/var/run/reboot-required\" to check if restart is required. but update/upgrade/autoremove run without a problem. err: %v", err))
	}
	return true
}
