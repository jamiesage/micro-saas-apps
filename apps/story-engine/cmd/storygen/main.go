// SECRET MANAGEMENT
// Development: loaded from .env
// Production: loaded from google secret manager
//
// DEPLOYMENT TYPE
// cloud run job

package main

import (
	"fmt"
	"log"

	"github.com/jamiesage/micro-saas-apps/apps/story-engine/config"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Job initialising...")

	// load .env file incase in local development, otherwise ignore error
	if err := godotenv.Load("./config/dev.env"); err != nil {
		log.Println("No dev.env file found, defaulting to environment variables")
	}

	// Load config
	cfg, err := config.Load("./config/pipeline.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// test using config
	log.Println("Testing config")
	fmt.Println(cfg.Anthropic.PrimaryModel)
	fmt.Println(cfg.GetTimezone())
	fmt.Println(cfg.GetModelForAgent("story_planner"))

	// Start app
	fmt.Println("Starting story pipeline...")
	// TODO: Implement story pipeline

	fmt.Println("Job completed successfully.")
}
