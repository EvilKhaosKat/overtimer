package main

import (
	"fmt"
	"log"
	"time"
)

const (
	DefaultWorkStartHour = 10
	DefaultWorkEndHour   = 19
)

type App struct {
	Settings        *Settings       `json:"settings"`
	ActiveEntry     *Entry          `json:"active_entry,omitempty"`
	FinishedEntries []FinishedEntry `json:"finished_entries"`
}

func NewApp(currentSettings *Settings, activeEntry *Entry, finishedEntries []FinishedEntry) *App {
	return &App{Settings: currentSettings, ActiveEntry: activeEntry, FinishedEntries: finishedEntries}
}

func NewAppDefault() *App {
	return &App{Settings: &Settings{
		WorkStartHour: DefaultWorkStartHour,
		WorkEndHour:   DefaultWorkEndHour,
	}}
}

func (app *App) String() string {
	total, totalType := app.calculateTotal()
	if totalType == hasOvertime {
		return fmt.Sprintf("overtime: %v", total)
	} else {
		return fmt.Sprintf("debt: %v", total)
	}
}

func main() {
	command, err := readCommand()
	if err != nil {
		log.Fatal(err)
	}

	app, err := getApp()
	if err != nil {
		log.Fatal(err)
	}

	switch command {
	case start:
		started := app.start()
		if started {
			fmt.Println("overtime started at", nowTimeFormatted())
		} else {
			fmt.Println("overtime is already in progress")
		}
	case stop:
		entryType, err := app.stop()
		if err != nil {
			log.Fatal("error occurred during stopping: ", err)
		}
		fmt.Printf("%v stopped at %v", entryType, nowTimeFormatted())
	case spend:
		started, err := app.spend()
		if err != nil {
			log.Fatal("error occurred during spending: ", err)
		}

		if started {
			fmt.Println("time spending started at", nowTimeFormatted())
		} else {
			fmt.Println("spending is already in progress")
		}
	case routine:
		result, err := app.routine()
		if err != nil {
			log.Fatal("error occurred during stopping: ", err)
		}
		fmt.Printf("%s", result)
	case status:
		fmt.Println(app)
	case lunch:
		started := app.lunch()
		if started {
			fmt.Println("lunch started at", nowTimeFormatted())
		} else {
			fmt.Println("lunch is already in progress")
		}
	default:
		log.Fatal("unsupported command $s", command)
	}

	err = app.save()
	if err != nil {
		log.Fatal(err)
	}
}

func nowTimeFormatted() string {
	return time.Now().Format(time.RFC3339)
}
