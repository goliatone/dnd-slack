package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

const (
	enableFlag  = "enable"
	disableFlag = "disable"
)

var defaults struct {
	token string
	text  string
	emoji string
}

var token *string
var text *string
var emoji *string

var usage func()

var command *flag.FlagSet
var enableCmd *flag.FlagSet
var disableCmd *flag.FlagSet

var action func(api *slack.Client, statusText string, statusEmoji string)

func init() {
	defaults.token = os.Getenv("SLACK_TOKEN")

	enableCmd = flag.NewFlagSet(enableFlag, flag.ExitOnError)
	disableCmd = flag.NewFlagSet(disableFlag, flag.ExitOnError)

	usage = func() {
		fmt.Println("Utility to toggle Slack Do Not Disturb using notification center")
		fmt.Fprintf(os.Stderr, "Usage: %s [subcommand] {enable|disable}\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {

	flag.Parse()

	if len(flag.Args()) == 0 {
		usage()
	}

	switch os.Args[1] {
	case enableFlag:
		defaults.emoji = ":octagonal_sign:"
		defaults.text = "I'm flowing"
		command = enableCmd
		action = enable
	case disableFlag:
		defaults.emoji = ""
		defaults.text = ""
		command = disableCmd
		action = disable
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	token = command.String("token", defaults.token, "Slack API auth <token>.")
	text = command.String("status-text", defaults.text, "Slack status text.")
	emoji = command.String("status-emoji", defaults.emoji, "Slack status emoji.")

	command.Parse(os.Args[2:])

	api := slack.New(*token)

	action(api, *text, *emoji)
}

func enable(api *slack.Client, statusText string, statusEmoji string) {
	if err := api.SetUserCustomStatus(statusText, statusEmoji); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	const (
		minutes = 20
	)

	res, err := api.SetSnooze(minutes)

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Enabled: %t, StartTime: %b, EndTime: %b\n",
		res.Enabled,
		res.NextStartTimestamp,
		res.NextEndTimestamp)
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
