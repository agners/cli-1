package cmd

import (
	"errors"
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addonsChangelogCmd = &cobra.Command{
	Use:     "changelog [slug]",
	Aliases: []string{"cl", "ch"},
	Short:   "Show changelog of an Hass.io add-on",
	Long: `
This command shows the changelog of an add-on. It gives you what has been
changed in the latest version and tell you about possible breaking changes.`,
	Example: `
hassio addons changelog core_ssh
hassio addons changelog core_mosquitto`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons changelog")

		section := "addons"
		command := "{slug}/changelog"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)

		if err != nil {
			fmt.Println(err)
			return
		}

		request := helper.GetRequest()

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		resp, err := request.SetHeader("Accept", "text/plain").Get(url)

		// returns 200 OK or 400, everything else is wrong
		if err == nil && resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			err = errors.New("Unexpected server response")
			log.Error(err)
		}

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			fmt.Println(string(resp.Body()))
		}

		return
	},
}

func init() {
	addonsCmd.AddCommand(addonsChangelogCmd)
}