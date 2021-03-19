package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"server-updater/aptclient"
	"server-updater/discordclient"

	"github.com/joho/godotenv"
)

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
	if updateCount == 0 {
		discord.LogInfoIgnoreError("checked for updates but no packages need to be updated. bye :)")
		return
	}
	log.Printf("%d packages need to be updated", updateCount)

	log.Println("start: upgrade")
	if err := apt.Upgrade(); err != nil {
		discord.LogFatalIgnoreError(fmt.Sprintf("could not upgrade. err: %v", err))
	}

	log.Println("start: autoremove")
	if err := apt.Autoremove(); err != nil {
		discord.LogFatalIgnoreError(fmt.Sprintf("could not autoremove. err: %v", err))
	}

	rebootRequired := false
	log.Println("start: reboot check")
	if _, err := os.Stat("/var/run/reboot-required"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("reboot is not required")
		} else {
			discord.LogFatalIgnoreError(fmt.Sprintf("could not read \"/var/run/reboot-required\" to check if restart is required. but update/upgrade/autoremove run without a problem. err: %v", err))
		}
	} else {
		log.Println("reboot is required")
		rebootRequired = true
	}

	discord.LogInfoIgnoreError(fmt.Sprintf("upgraded %d packages. reboot required: %t", updateCount, rebootRequired))
	log.Println("finished node-update")

}
