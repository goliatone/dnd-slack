package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	// Subcommands
	flowCommand := flag.NewFlagSet("flow", flag.ExitOnError)
	inactiveCommand := flag.NewFlagSet("list", flag.ExitOnError)

	token := flowCommand.String("token", "", "Slack API auth <token>.")
	statusText := flowCommand.String("status-text", "I'm flowing", "Slack status text.")
	emoji := flowCommand.String("status-emoji", ":octagonal_sign:", "Slack status emoji.")

	if len(os.Args) == 1 {
		flag.PrintDefaults()
		flowCommand.PrintDefaults()
		inactiveCommand.PrintDefaults()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "flow":
		flowCommand.Parse(os.Args[2:])
	case "inactive":
		inactiveCommand.Parse(os.Args[2:])
		fmt.Print("inactive execute")
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("print %s\n", *token)
	api := slack.New(*token)

	if flowCommand.Parsed() {
		flow(api, *statusText, *emoji)
	}

	if inactiveCommand.Parsed() {
		inactive(api, "", "")
	}
}

func flow(api *slack.Client, statusText string, statusEmoji string) {
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

func inactive(api *slack.Client, statusText string, statusEmoji string) {
	if err := api.EndDND(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if err := api.SetUserCustomStatus(statusEmoji, statusEmoji); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	fmt.Print("DND status disabled!")
}
