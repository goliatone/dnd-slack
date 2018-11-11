package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

const (
	enableCmd  = "enable"
	disableCmd = "disable"
)

var token *string
var text *string
var emoji *string

func main() {

	flowCommand := flag.NewFlagSet(enableCmd, flag.ExitOnError)
	inactiveCommand := flag.NewFlagSet(disableCmd, flag.ExitOnError)

	var Usage = func() {
		fmt.Println("Utility to toggle Slack Do Not Disturb using notification center")
		fmt.Fprintf(os.Stderr, "Usage: %s [subcommand] {enable|disable}\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.Parse()

	if len(flag.Args()) == 0 {
		Usage()
	}

	switch os.Args[1] {
	case enableCmd:
		token = flowCommand.String("token", "", "Slack API auth <token>.")
		text = flowCommand.String("status-text", "I'm flowing", "Slack status text.")
		emoji = flowCommand.String("status-emoji", ":octagonal_sign:", "Slack status emoji.")
		flowCommand.Parse(os.Args[2:])
	case disableCmd:
		token = inactiveCommand.String("token", "", "Slack API auth <token>.")
		text = inactiveCommand.String("status-text", "", "Slack status text.")
		emoji = inactiveCommand.String("status-emoji", "", "Slack status emoji.")
		inactiveCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("print %s\n", *token)
	api := slack.New(*token)

	if flowCommand.Parsed() {
		enable(api, *text, *emoji)
	}

	if inactiveCommand.Parsed() {
		disable(api, *text, *emoji)
	}
}

func enable(api *slack.Client, statusText string, statusEmoji string) {
	if err := api.SetUserCustomStatus(statusText, statusEmoji); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	const (
		minutes = 20
	)

	snoozeResponse, err := api.SetSnooze(minutes)

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Enabled: %t, StartTime: %b, EndTime: %b\n",
		snoozeResponse.Enabled,
		snoozeResponse.NextStartTimestamp,
		snoozeResponse.NextEndTimestamp)
}

func disable(api *slack.Client, statusText string, statusEmoji string) {
	if err := api.EndDND(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if err := api.SetUserCustomStatus(statusEmoji, statusEmoji); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	fmt.Println("DND macOS status disabled")
}
