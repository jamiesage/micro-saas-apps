// USAGE EXAMPLES
//
// Load the config:
// cfg, err := config.Load("./config/pipeline.yaml")
// if err != nil {
//     log.Fatalf("Failed to load config: %v", err)
// }
//
// Use the config
// Get a config value:
// fmt.Println(cfg.Anthropic.PrimaryModel)
//
// Config functions:
// fmt.Println(cfg.GetTimezone())
// fmt.Println(cfg.GetModelForAgent("story_planner"))

package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config is the root configuration struct for the story engine.
type Config struct {
	Anthropic       AnthropicConfig       `mapstructure:"anthropic"`
	Instagram       InstagramConfig       `mapstructure:"instagram"`
	ImageGeneration ImageGenerationConfig `mapstructure:"image_generation"`
	Email           EmailConfig           `mapstructure:"email"`
	Pipeline        PipelineConfig        `mapstructure:"pipeline"`
	Agents          AgentsConfig          `mapstructure:"agents"`
	Paths           PathsConfig           `mapstructure:"paths"`
	Monitoring      MonitoringConfig      `mapstructure:"monitoring"`
	Logging         LoggingConfig         `mapstructure:"logging"`
}

type AnthropicConfig struct {
	PrimaryModel string         `mapstructure:"primary_model"`
	FastModel    string         `mapstructure:"fast_model"`
	MaxTokens    int            `mapstructure:"max_tokens"`
	Thinking     ThinkingConfig `mapstructure:"thinking"`
	Retry        RetryConfig    `mapstructure:"retry"`
	APIKey       string         // Loaded from env only
}

type ThinkingConfig struct {
	Enabled      bool `mapstructure:"enabled"`
	BudgetTokens int  `mapstructure:"budget_tokens"`
}

type RetryConfig struct {
	MaxAttempts    int     `mapstructure:"max_attempts"`
	InitialDelayMs int     `mapstructure:"initial_delay_ms"`
	MaxDelayMs     int     `mapstructure:"max_delay_ms"`
	Multiplier     float64 `mapstructure:"multiplier"`
}

type InstagramConfig struct {
	CommentsToFetch         int             `mapstructure:"comments_to_fetch"`
	MinLikesThreshold       int             `mapstructure:"min_likes_threshold"`
	TopCommentsForFiltering int             `mapstructure:"top_comments_for_filtering"`
	RateLimit               RateLimitConfig `mapstructure:"rate_limit"`
	AccountID               string          // Loaded from env only
	AccessToken             string          // Loaded from env only
}

type RateLimitConfig struct {
	RequestsPerHour   int `mapstructure:"requests_per_hour"`
	RetryAfterSeconds int `mapstructure:"retry_after_seconds"`
}

type ImageGenerationConfig struct {
	Provider          string           `mapstructure:"provider"`
	Model             string           `mapstructure:"model"`
	ImagesPerChapter  int              `mapstructure:"images_per_chapter"`
	Width             int              `mapstructure:"width"`
	Height            int              `mapstructure:"height"`
	GuidanceScale     float64          `mapstructure:"guidance_scale"`
	NumInferenceSteps int              `mapstructure:"num_inference_steps"`
	StyleAnchors      []string         `mapstructure:"style_anchors"`
	TimeoutSeconds    int              `mapstructure:"timeout_seconds"`
	Retry             ImageRetryConfig `mapstructure:"retry"`
	BFLAPIKey         string           // Loaded from env only
	ReplicateToken    string           // Loaded from env only
}

type ImageRetryConfig struct {
	MaxAttempts  int `mapstructure:"max_attempts"`
	DelaySeconds int `mapstructure:"delay_seconds"`
}

type EmailConfig struct {
	Enabled     bool             `mapstructure:"enabled"`
	Provider    string           `mapstructure:"provider"`
	SMTP        SMTPConfig       `mapstructure:"smtp"`
	SendGrid    SendGridConfig   `mapstructure:"sendgrid"`
	Recipients  RecipientsConfig `mapstructure:"recipients"`
	FromAddress string           `mapstructure:"from_address"`
	FromName    string           `mapstructure:"from_name"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string // Loaded from env only
	Password string // Loaded from env only
}

type SendGridConfig struct {
	FromName string `mapstructure:"from_name"`
	APIKey   string // Loaded from env only
}

type RecipientsConfig struct {
	DailyReport string   `mapstructure:"daily_report"`
	ErrorAlerts []string `mapstructure:"error_alerts"`
}

type PipelineConfig struct {
	Schedule    string            `mapstructure:"schedule"`
	Timezone    string            `mapstructure:"timezone"`
	DryRun      bool              `mapstructure:"dry_run"`
	Story       StoryConfig       `mapstructure:"story"`
	Hashtags    HashtagsConfig    `mapstructure:"hashtags"`
	Context     ContextConfig     `mapstructure:"context"`
	Validation  ValidationConfig  `mapstructure:"validation"`
	Checkpoints CheckpointsConfig `mapstructure:"checkpoints"`
}

type StoryConfig struct {
	MaxChapterLength    int `mapstructure:"max_chapter_length"`
	MinChapterLength    int `mapstructure:"min_chapter_length"`
	TargetChapterLength int `mapstructure:"target_chapter_length"`
}

type HashtagsConfig struct {
	CountMin           int `mapstructure:"count_min"`
	CountMax           int `mapstructure:"count_max"`
	MaxTotalCharacters int `mapstructure:"max_total_characters"`
}

type ContextConfig struct {
	RecentChaptersCount int `mapstructure:"recent_chapters_count"`
	FullTextChapters    int `mapstructure:"full_text_chapters"`
	MaxEntities         int `mapstructure:"max_entities"`
}

type ValidationConfig struct {
	StrictLengthCheck   bool `mapstructure:"strict_length_check"`
	LengthRetryAttempts int  `mapstructure:"length_retry_attempts"`
	CanonCheckEnabled   bool `mapstructure:"canon_check_enabled"`
}

type CheckpointsConfig struct {
	Enabled       bool `mapstructure:"enabled"`
	RetentionDays int  `mapstructure:"retention_days"`
}

type AgentsConfig struct {
	CommentFilter        AgentConfig `mapstructure:"comment_filter"`
	StoryPlanner         AgentConfig `mapstructure:"story_planner"`
	StoryWriter          AgentConfig `mapstructure:"story_writer"`
	HashtagGenerator     AgentConfig `mapstructure:"hashtag_generator"`
	ImagePromptGenerator AgentConfig `mapstructure:"image_prompt_generator"`
}

type AgentConfig struct {
	Model                string  `mapstructure:"model"`
	UseThinking          bool    `mapstructure:"use_thinking"`
	ThinkingBudget       int     `mapstructure:"thinking_budget"`
	Temperature          float64 `mapstructure:"temperature"`
	MaxSuggestions       int     `mapstructure:"max_suggestions"`
	MaxBankedIdeas       int     `mapstructure:"max_banked_ideas"`
	IncludeStoryTags     bool    `mapstructure:"include_story_tags"`
	IncludeGenreTags     bool    `mapstructure:"include_genre_tags"`
	IncludeGeneralTags   bool    `mapstructure:"include_general_tags"`
	IncludeCharacterRefs bool    `mapstructure:"include_character_refs"`
	MaxPromptLength      int     `mapstructure:"max_prompt_length"`
}

type PathsConfig struct {
	DataDir      string `mapstructure:"data_dir"`
	StoryBible   string `mapstructure:"story_bible"`
	EntitiesDir  string `mapstructure:"entities_dir"`
	ChaptersDir  string `mapstructure:"chapters_dir"`
	RunsDir      string `mapstructure:"runs_dir"`
	PromptsDir   string `mapstructure:"prompts_dir"`
	TemplatesDir string `mapstructure:"templates_dir"`
}

type MonitoringConfig struct {
	Healthchecks HealthchecksConfig `mapstructure:"healthchecks"`
	HTTP         HTTPConfig         `mapstructure:"http"`
	CostTracking CostTrackingConfig `mapstructure:"cost_tracking"`
}

type HealthchecksConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
	PingURL        string // Loaded from env only
}

type HTTPConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Port    int  `mapstructure:"port"`
}

type CostTrackingConfig struct {
	Enabled             bool    `mapstructure:"enabled"`
	DailyAlertThreshold float64 `mapstructure:"daily_alert_threshold"`
	LogFile             string  `mapstructure:"log_file"`
}

type LoggingConfig struct {
	Level           string     `mapstructure:"level"`
	Format          string     `mapstructure:"format"`
	File            FileConfig `mapstructure:"file"`
	IncludeThinking bool       `mapstructure:"include_thinking"`
}

type FileConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	FilenamePattern string `mapstructure:"filename_pattern"`
}

// Load reads configuration from the pipeline.yaml and .env files
// Environment variables take precedence over pipeline.yaml values with the same name
// see https://github.com/spf13/viper
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Enable environment variable overrides
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Load sensitive values from environment
	cfg.Anthropic.APIKey = os.Getenv("ANTHROPIC_API_KEY")
	cfg.Instagram.AccountID = os.Getenv("INSTAGRAM_ACCOUNT_ID")
	cfg.Instagram.AccessToken = os.Getenv("INSTAGRAM_ACCESS_TOKEN")
	cfg.ImageGeneration.BFLAPIKey = os.Getenv("BFL_API_KEY")
	cfg.ImageGeneration.ReplicateToken = os.Getenv("REPLICATE_API_TOKEN")
	cfg.Email.SMTP.User = os.Getenv("SMTP_USER")
	cfg.Email.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	cfg.Email.SendGrid.APIKey = os.Getenv("SENDGRID_API_KEY")
	cfg.Monitoring.Healthchecks.PingURL = os.Getenv("HEALTHCHECKS_PING_URL")

	// Environment variable overrides for YAML values - i.e. for diff configs in dev / prod
	// currently not used, example:
	// if accountID := os.Getenv("INSTAGRAM_ACCOUNT_ID"); accountID != "" {
	// 	cfg.Instagram.AccountID = accountID
	// }

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// Validate checks that all required configuration values are present.
func (c *Config) Validate() error {
	var errs []string

	// Required API keys
	if c.Anthropic.APIKey == "" {
		errs = append(errs, "ANTHROPIC_API_KEY environment variable is required")
	}
	if c.Instagram.AccessToken == "" {
		errs = append(errs, "INSTAGRAM_ACCESS_TOKEN environment variable is required")
	}
	if c.Instagram.AccountID == "" {
		errs = append(errs, "INSTAGRAM_ACCOUNT_ID environment variable is required")
	}

	// Image generation API key (one of them required)
	if c.ImageGeneration.BFLAPIKey == "" && c.ImageGeneration.ReplicateToken == "" {
		errs = append(errs, "BFL_API_KEY or REPLICATE_API_TOKEN environment variable is required")
	}

	// Email configuration (if enabled)
	if c.Email.Enabled {
		switch c.Email.Provider {
		case "smtp":
			if c.Email.SMTP.User == "" || c.Email.SMTP.Password == "" {
				errs = append(errs, "SMTP_USER and SMTP_PASSWORD are required when using SMTP")
			}
		case "sendgrid":
			if c.Email.SendGrid.APIKey == "" {
				errs = append(errs, "SENDGRID_API_KEY is required when using SendGrid")
			}
		}
	}

	// Validate thinking budget bounds
	if c.Anthropic.Thinking.Enabled {
		if c.Anthropic.Thinking.BudgetTokens < 1024 {
			errs = append(errs, "thinking.budget_tokens must be at least 1024")
		}
		if c.Anthropic.Thinking.BudgetTokens > 64000 {
			errs = append(errs, "thinking.budget_tokens must not exceed 64000")
		}
	}

	// Validate story constraints
	if c.Pipeline.Story.MaxChapterLength <= c.Pipeline.Story.MinChapterLength {
		errs = append(errs, "story.max_chapter_length must be greater than min_chapter_length")
	}

	if len(errs) > 0 {
		return fmt.Errorf("configuration errors:\n  - %s", strings.Join(errs, "\n  - "))
	}

	return nil
}

// GetTimezone returns the configured timezone as a *time.Location.
func (c *Config) GetTimezone() (*time.Location, error) {
	return time.LoadLocation(c.Pipeline.Timezone)
}

// GetModelForAgent returns the appropriate model for a given agent,
// falling back to defaults if not specified.
func (c *Config) GetModelForAgent(agentName string) string {
	switch agentName {
	case "comment_filter":
		if c.Agents.CommentFilter.Model != "" {
			return c.Agents.CommentFilter.Model
		}
		return c.Anthropic.FastModel
	case "story_planner":
		if c.Agents.StoryPlanner.Model != "" {
			return c.Agents.StoryPlanner.Model
		}
		return c.Anthropic.PrimaryModel
	case "story_writer":
		if c.Agents.StoryWriter.Model != "" {
			return c.Agents.StoryWriter.Model
		}
		return c.Anthropic.PrimaryModel
	case "hashtag_generator":
		if c.Agents.HashtagGenerator.Model != "" {
			return c.Agents.HashtagGenerator.Model
		}
		return c.Anthropic.FastModel
	case "image_prompt_generator":
		if c.Agents.ImagePromptGenerator.Model != "" {
			return c.Agents.ImagePromptGenerator.Model
		}
		return c.Anthropic.PrimaryModel
	default:
		return c.Anthropic.PrimaryModel
	}
}
