package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adrock-miles/about-me/internal/application/contact"
	"github.com/adrock-miles/about-me/internal/application/project"
	"github.com/adrock-miles/about-me/internal/infrastructure/email"
	"github.com/adrock-miles/about-me/internal/infrastructure/github"
	"github.com/adrock-miles/about-me/internal/interfaces/http/handler"
	"github.com/adrock-miles/about-me/internal/interfaces/http/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "about-me",
	Short: "Personal portfolio site server",
	Long:  "A personal portfolio website backend built with Go, serving a React frontend.",
	RunE:  runServer,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().IntP("port", "p", 8080, "Port to run the server on")
	rootCmd.Flags().StringP("host", "H", "0.0.0.0", "Host to bind the server to")

	viper.BindPFlag("server.port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.host", rootCmd.Flags().Lookup("host"))
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Bind PORT env var (used by Railway and other PaaS providers)
	viper.BindEnv("server.port", "PORT")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: no config file found: %v", err)
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	// Infrastructure layer
	githubClient := github.NewClient(
		viper.GetString("github.username"),
		viper.GetString("github.api_url"),
	)
	emailSender := email.NewSender(
		viper.GetString("smtp.host"),
		viper.GetInt("smtp.port"),
		viper.GetString("smtp.from"),
		viper.GetString("contact.email"),
	)

	// Application layer
	projectService := project.NewService(githubClient)
	contactService := contact.NewService(emailSender, viper.GetString("contact.email"))

	// Interface layer
	projectHandler := handler.NewProjectHandler(projectService)
	contactHandler := handler.NewContactHandler(contactService)

	r := router.New(projectHandler, contactHandler)

	addr := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))
	log.Printf("Starting server on %s", addr)

	return http.ListenAndServe(addr, r)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
