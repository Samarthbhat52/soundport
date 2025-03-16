package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	accent = lipgloss.NewStyle().Foreground(lipgloss.Color("163"))
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("161"))
)

func init() {
	rootCmd.AddCommand(spotifyCmd)
	spotifyCmd.AddCommand(spotifyLoginCmd)
	spotifyCmd.AddCommand(spotifyPlaylistsCmd)
}

type listOptions struct {
	options []string
}

var spotifyCmd = &cobra.Command{
	Use:   "spotify",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
}

var spotifyLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var message strings.Builder
		var status strings.Builder

		creds := spotify.NewCredentials()

		message.WriteString("Click on " + accent.Render("Accept") + " in the browser popup\n")
		fmt.Println(message.String())

		ch := make(chan int)
		state := spotify.RandStringBytes(16)

		url := creds.GetAuthURL(state)
		go creds.StartHttpServer(ch, state)
		go spotify.OpenBrowser(url)

		val := <-ch
		if val == 0 {
			status.WriteString(green.Render("Login successful\n"))
			fmt.Println(status.String())
		} else {
			status.WriteString(red.Render("Login failed\n"))
			fmt.Println(status.String())
		}
		fmt.Println("Browser window/tab can be closed.")
	},
}

var spotifyPlaylistsCmd = &cobra.Command{
	Use: "get",
	PreRun: func(cmd *cobra.Command, args []string) {
		a, err := spotify.NewAuth()
		if err != nil {
			log.Fatal("Please login to continue")
		}
		err = a.RefreshSession()
		if err != nil {
			log.Fatal("error processing auth request")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		a, _ := spotify.NewAuth()
		resp, err := a.GetPlaylists()
		if err != nil {
			log.Fatal(err)
		}

		l := list.New(resp.ItemPlaylists, list.NewDefaultDelegate(), 0, 0)
	},
}
