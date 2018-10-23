package main

import (
	botcli "cli"
	"discord/anime/mal"
	messageMal "discord/message/anime/mal"
	messageMiddleware "discord/message/middleware"
	"discord/writer"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

func init() {
	cli.AppHelpTemplate = `NAME:
	{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

 USAGE:
	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

 VERSION:
	{{.Version}}{{end}}{{end}}{{if .Description}}

 DESCRIPTION:
	{{.Description}}{{end}}{{if len .Authors}}

 AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
	{{range $index, $author := .Authors}}{{if $index}}
	{{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

 COMMANDS:{{range .VisibleCategories}}{{if .Name}}

	{{.Name}}:{{end}}{{range .VisibleCommands}}
	  {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

 GLOBAL OPTIONS:
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

 COPYRIGHT:
	{{.Copyright}}{{end}}
 ` + fmt.Sprintf("\n%s", os.Getenv("EOF_DELIM"))

	cli.CommandHelpTemplate = `NAME:
 {{.HelpName}} - {{.Usage}}

USAGE:
 {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

CATEGORY:
 {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
 {{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
 {{range .VisibleFlags}}{{.}}
 {{end}}{{end}}
` + fmt.Sprintf("\n%s", os.Getenv("EOF_DELIM"))
	cli.SubcommandHelpTemplate = `NAME:
 {{.HelpName}} - {{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}

USAGE:
 {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
 {{.Name}}:{{end}}{{range .VisibleCommands}}
   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
OPTIONS:
 {{range .VisibleFlags}}{{.}}
 {{end}}{{end}}
` + fmt.Sprintf("\n%s", os.Getenv("EOF_DELIM"))
}

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if (!messageMiddleware.IgnoreOwnMessagesMiddleware{}.ProcessMessage(session, message)) {
		return
	}

	cmdPrefix := os.Getenv("CMD_PREFIX")
	if strings.HasPrefix(message.Content, cmdPrefix) {
		cmdStr := strings.TrimPrefix(message.Content, cmdPrefix)
		args := strings.Split(cmdStr, " ")
		cliApp := botcli.CreateCLI(os.Getenv("APP_NAME"), cmdPrefix)

		writer := writer.CreateDiscordWriter(
			session, message.ChannelID, os.Getenv("EOF_DELIM"))
		cliApp.Writer = writer

		cliApp.Action = func(c *cli.Context) error {
			return nil
		}
		cliApp.Commands = []cli.Command{
			{
				Name:  "anime",
				Usage: "List anime commands",
				Action: func(c *cli.Context) error {
					cli.ShowCommandHelp(c, "anime")

					return nil
				},
				Subcommands: cli.Commands{
					cli.Command{
						Name:  "profile",
						Usage: "Displays a user's anime profile",
						Flags: []cli.Flag{
							cli.StringFlag{
								Name:  "name, n",
								Usage: "My Anime List Username",
							},
						},
						Action: func(c *cli.Context) error {
							userProfile, err := mal.GetProfile(c.String("name"))
							if err != nil {
								log.Panicln(err)
								return nil
							}
							options := messageMal.CreateAnimeProfileEmbeddedOptions{
								AnimeProfile: userProfile,
							}
							embeddedMessage := messageMal.CreateAnimeProfileEmbedded(
								options)

							_, error := session.ChannelMessageSendEmbed(message.ChannelID, embeddedMessage)
							if error != nil {
								log.Panic(error)
							}
							return nil
						},
					},
				},
			},
		}

		err := cliApp.Run(args)
		if err != nil {
			log.Println(err)
		}

	}
}