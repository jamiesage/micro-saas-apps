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
// WARNING: Do not log this struct or its nested fields directly.
// Several fields contain secrets (API keys, tokens, passwords) loaded from
// environment variables that must not appear in logs or error messages.
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
	Recipients  RecipientsConfig // Loaded from env only
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
	DailyReport string   // Loaded from env only
	ErrorAlerts []string // Loaded from env only
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
// path traversal risk: this function expects a trusted configPath input
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

	// Load and validate secrets from environment
	var envErrs []string

	// Required API keys
	cfg.Anthropic.APIKey = requireEnv("ANTHROPIC_API_KEY", &envErrs)
	cfg.Instagram.AccountID = requireEnv("INSTAGRAM_ACCOUNT_ID", &envErrs)
	cfg.Instagram.AccessToken = requireEnv("INSTAGRAM_ACCESS_TOKEN", &envErrs)

	// Image generation (at least one required)
	cfg.ImageGeneration.BFLAPIKey = optionalEnv("BFL_API_KEY")
	cfg.ImageGeneration.ReplicateToken = optionalEnv("REPLICATE_API_TOKEN")
	if cfg.ImageGeneration.BFLAPIKey == "" && cfg.ImageGeneration.ReplicateToken == "" {
		envErrs = append(envErrs, "BFL_API_KEY or REPLICATE_API_TOKEN environment variable is required")
	}

	// Email configuration (conditionally required)
	if cfg.Email.Enabled {
		switch cfg.Email.Provider {
		case "smtp":
			cfg.Email.SMTP.User = requireEnv("SMTP_USER", &envErrs)
			cfg.Email.SMTP.Password = requireEnv("SMTP_PASSWORD", &envErrs)
		case "sendgrid":
			cfg.Email.SendGrid.APIKey = requireEnv("SENDGRID_API_KEY", &envErrs)
		}
		// Load email recipients
		cfg.Email.Recipients.DailyReport = requireEnv("EMAIL_RECIPIENT_DAILY_REPORT", &envErrs)
		errorAlertsStr := requireEnv("EMAIL_RECIPIENT_ERROR_ALERTS", &envErrs)
		if errorAlertsStr != "" {
			// Split comma-separated list and trim whitespace
			cfg.Email.Recipients.ErrorAlerts = strings.Split(errorAlertsStr, ",")
			for i := range cfg.Email.Recipients.ErrorAlerts {
				cfg.Email.Recipients.ErrorAlerts[i] = strings.TrimSpace(cfg.Email.Recipients.ErrorAlerts[i])
			}
		}
	} else {
		// Load optionally if email is disabled
		cfg.Email.SMTP.User = optionalEnv("SMTP_USER")
		cfg.Email.SMTP.Password = optionalEnv("SMTP_PASSWORD")
		cfg.Email.SendGrid.APIKey = optionalEnv("SENDGRID_API_KEY")
		cfg.Email.Recipients.DailyReport = optionalEnv("EMAIL_RECIPIENT_DAILY_REPORT")
		errorAlertsStr := optionalEnv("EMAIL_RECIPIENT_ERROR_ALERTS")
		if errorAlertsStr != "" {
			cfg.Email.Recipients.ErrorAlerts = strings.Split(errorAlertsStr, ",")
			for i := range cfg.Email.Recipients.ErrorAlerts {
				cfg.Email.Recipients.ErrorAlerts[i] = strings.TrimSpace(cfg.Email.Recipients.ErrorAlerts[i])
			}
		}
	}

	// Required monitoring
	cfg.Monitoring.Healthchecks.PingURL = requireEnv("HEALTHCHECKS_PING_URL", &envErrs)

	// Check for environment variable errors
	if len(envErrs) > 0 {
		return nil, fmt.Errorf("environment variable errors:\n  - %s", strings.Join(envErrs, "\n  - "))
	}

	// Validate logical value constraints
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// checks logical constraints on configuration values.
// Environment variable presence is validated during Load().
func (c *Config) Validate() error {
	var errs []string

	// Validate each section
	errs = append(errs, c.validateAnthropic()...)
	errs = append(errs, c.validateInstagram()...)
	errs = append(errs, c.validateImageGeneration()...)
	errs = append(errs, c.validateEmail()...)
	errs = append(errs, c.validatePipeline()...)
	errs = append(errs, c.validatePaths()...)
	errs = append(errs, c.validateMonitoring()...)
	errs = append(errs, c.validateLogging()...)

	if len(errs) > 0 {
		return fmt.Errorf("configuration errors:\n  - %s", strings.Join(errs, "\n  - "))
	}

	return nil
}

// validateAnthropic validates Anthropic configuration.
func (c *Config) validateAnthropic() []string {
	var errs []string

	if c.Anthropic.PrimaryModel == "" {
		errs = append(errs, "anthropic.primary_model is required")
	}

	if c.Anthropic.FastModel == "" {
		errs = append(errs, "anthropic.fast_model is required")
	}

	if c.Anthropic.MaxTokens <= 0 {
		errs = append(errs, "anthropic.max_tokens must be greater than 0")
	}

	// Validate thinking budget bounds
	if c.Anthropic.Thinking.Enabled {
		if c.Anthropic.Thinking.BudgetTokens < 1024 {
			errs = append(errs, "anthropic.thinking.budget_tokens must be at least 1024")
		}
		if c.Anthropic.Thinking.BudgetTokens > 64000 {
			errs = append(errs, "anthropic.thinking.budget_tokens must not exceed 64000")
		}
	}

	// Validate retry configuration
	if c.Anthropic.Retry.MaxAttempts <= 0 {
		errs = append(errs, "anthropic.retry.max_attempts must be greater than 0")
	}
	if c.Anthropic.Retry.InitialDelayMs <= 0 {
		errs = append(errs, "anthropic.retry.initial_delay_ms must be greater than 0")
	}
	if c.Anthropic.Retry.MaxDelayMs <= 0 {
		errs = append(errs, "anthropic.retry.max_delay_ms must be greater than 0")
	}
	if c.Anthropic.Retry.MaxDelayMs < c.Anthropic.Retry.InitialDelayMs {
		errs = append(errs, "anthropic.retry.max_delay_ms must be greater than or equal to initial_delay_ms")
	}
	if c.Anthropic.Retry.Multiplier <= 0 {
		errs = append(errs, "anthropic.retry.multiplier must be greater than 0")
	}

	return errs
}

// validateInstagram validates Instagram configuration.
func (c *Config) validateInstagram() []string {
	var errs []string

	if c.Instagram.CommentsToFetch < 1 || c.Instagram.CommentsToFetch > 100 {
		errs = append(errs, "instagram.comments_to_fetch must be between 1 and 100")
	}

	if c.Instagram.MinLikesThreshold < 0 {
		errs = append(errs, "instagram.min_likes_threshold must be greater than or equal to 0")
	}

	if c.Instagram.TopCommentsForFiltering <= 0 {
		errs = append(errs, "instagram.top_comments_for_filtering must be greater than 0")
	}

	if c.Instagram.TopCommentsForFiltering > c.Instagram.CommentsToFetch {
		errs = append(errs, "instagram.top_comments_for_filtering must not exceed comments_to_fetch")
	}

	// Validate rate limits
	if c.Instagram.RateLimit.RequestsPerHour <= 0 {
		errs = append(errs, "instagram.rate_limit.requests_per_hour must be greater than 0")
	}
	if c.Instagram.RateLimit.RetryAfterSeconds <= 0 {
		errs = append(errs, "instagram.rate_limit.retry_after_seconds must be greater than 0")
	}

	return errs
}

// validateImageGeneration validates image generation configuration.
func (c *Config) validateImageGeneration() []string {
	var errs []string

	if c.ImageGeneration.Provider != "bfl" && c.ImageGeneration.Provider != "replicate" {
		errs = append(errs, "image_generation.provider must be 'bfl' or 'replicate'")
	}

	if c.ImageGeneration.Model == "" {
		errs = append(errs, "image_generation.model is required")
	}

	if c.ImageGeneration.ImagesPerChapter < 1 || c.ImageGeneration.ImagesPerChapter > 4 {
		errs = append(errs, "image_generation.images_per_chapter must be between 1 and 4")
	}

	if c.ImageGeneration.Width <= 0 {
		errs = append(errs, "image_generation.width must be greater than 0")
	}

	if c.ImageGeneration.Height <= 0 {
		errs = append(errs, "image_generation.height must be greater than 0")
	}

	if c.ImageGeneration.GuidanceScale <= 0 {
		errs = append(errs, "image_generation.guidance_scale must be greater than 0")
	}

	if c.ImageGeneration.NumInferenceSteps <= 0 {
		errs = append(errs, "image_generation.num_inference_steps must be greater than 0")
	}

	if c.ImageGeneration.TimeoutSeconds <= 0 {
		errs = append(errs, "image_generation.timeout_seconds must be greater than 0")
	}

	// Validate retry configuration
	if c.ImageGeneration.Retry.MaxAttempts <= 0 {
		errs = append(errs, "image_generation.retry.max_attempts must be greater than 0")
	}
	if c.ImageGeneration.Retry.DelaySeconds <= 0 {
		errs = append(errs, "image_generation.retry.delay_seconds must be greater than 0")
	}

	return errs
}

// validateEmail validates email configuration with conditional logic.
func (c *Config) validateEmail() []string {
	var errs []string

	// Only validate if email is enabled
	if !c.Email.Enabled {
		return errs
	}

	// Validate provider choice
	if c.Email.Provider != "smtp" && c.Email.Provider != "sendgrid" {
		errs = append(errs, "email.provider must be 'smtp' or 'sendgrid' when email is enabled")
	}

	// Validate provider-specific configuration
	if c.Email.Provider == "smtp" {
		if c.Email.SMTP.Host == "" {
			errs = append(errs, "email.smtp.host is required when using smtp provider")
		}
		if c.Email.SMTP.Port <= 0 {
			errs = append(errs, "email.smtp.port must be greater than 0 when using smtp provider")
		}
	}

	// Validate common email fields
	if c.Email.FromAddress == "" {
		errs = append(errs, "email.from_address is required when email is enabled")
	}

	if c.Email.Recipients.DailyReport == "" {
		errs = append(errs, "email.recipients.daily_report is required when email is enabled")
	}

	if len(c.Email.Recipients.ErrorAlerts) == 0 {
		errs = append(errs, "email.recipients.error_alerts must have at least one recipient when email is enabled")
	}

	return errs
}

// validatePipeline validates pipeline configuration.
func (c *Config) validatePipeline() []string {
	var errs []string

	if c.Pipeline.Schedule == "" {
		errs = append(errs, "pipeline.schedule is required")
	}

	if c.Pipeline.Timezone == "" {
		errs = append(errs, "pipeline.timezone is required")
	} else {
		// Validate timezone is valid
		if _, err := time.LoadLocation(c.Pipeline.Timezone); err != nil {
			errs = append(errs, fmt.Sprintf("pipeline.timezone '%s' is not a valid timezone", c.Pipeline.Timezone))
		}
	}

	// Validate story constraints
	if c.Pipeline.Story.MaxChapterLength <= 0 {
		errs = append(errs, "pipeline.story.max_chapter_length must be greater than 0")
	}
	if c.Pipeline.Story.MinChapterLength <= 0 {
		errs = append(errs, "pipeline.story.min_chapter_length must be greater than 0")
	}
	if c.Pipeline.Story.TargetChapterLength <= 0 {
		errs = append(errs, "pipeline.story.target_chapter_length must be greater than 0")
	}
	if c.Pipeline.Story.MaxChapterLength <= c.Pipeline.Story.MinChapterLength {
		errs = append(errs, "pipeline.story.max_chapter_length must be greater than min_chapter_length")
	}
	if c.Pipeline.Story.TargetChapterLength < c.Pipeline.Story.MinChapterLength ||
		c.Pipeline.Story.TargetChapterLength > c.Pipeline.Story.MaxChapterLength {
		errs = append(errs, "pipeline.story.target_chapter_length must be between min and max chapter length")
	}

	// Validate hashtag constraints
	if c.Pipeline.Hashtags.CountMin <= 0 {
		errs = append(errs, "pipeline.hashtags.count_min must be greater than 0")
	}
	if c.Pipeline.Hashtags.CountMax <= 0 {
		errs = append(errs, "pipeline.hashtags.count_max must be greater than 0")
	}
	if c.Pipeline.Hashtags.CountMax < c.Pipeline.Hashtags.CountMin {
		errs = append(errs, "pipeline.hashtags.count_max must be greater than or equal to count_min")
	}
	if c.Pipeline.Hashtags.MaxTotalCharacters <= 0 {
		errs = append(errs, "pipeline.hashtags.max_total_characters must be greater than 0")
	}

	// Validate context constraints
	if c.Pipeline.Context.RecentChaptersCount <= 0 {
		errs = append(errs, "pipeline.context.recent_chapters_count must be greater than 0")
	}
	if c.Pipeline.Context.FullTextChapters <= 0 {
		errs = append(errs, "pipeline.context.full_text_chapters must be greater than 0")
	}
	if c.Pipeline.Context.FullTextChapters > c.Pipeline.Context.RecentChaptersCount {
		errs = append(errs, "pipeline.context.full_text_chapters must not exceed recent_chapters_count")
	}
	if c.Pipeline.Context.MaxEntities <= 0 {
		errs = append(errs, "pipeline.context.max_entities must be greater than 0")
	}

	// Validate validation constraints
	if c.Pipeline.Validation.LengthRetryAttempts < 0 {
		errs = append(errs, "pipeline.validation.length_retry_attempts must be greater than or equal to 0")
	}

	// Validate checkpoint constraints
	if c.Pipeline.Checkpoints.Enabled && c.Pipeline.Checkpoints.RetentionDays <= 0 {
		errs = append(errs, "pipeline.checkpoints.retention_days must be greater than 0 when checkpoints are enabled")
	}

	return errs
}

// validatePaths validates paths configuration.
func (c *Config) validatePaths() []string {
	var errs []string

	if c.Paths.DataDir == "" {
		errs = append(errs, "paths.data_dir is required")
	}

	if c.Paths.StoryBible == "" {
		errs = append(errs, "paths.story_bible is required")
	}

	if c.Paths.EntitiesDir == "" {
		errs = append(errs, "paths.entities_dir is required")
	}

	if c.Paths.ChaptersDir == "" {
		errs = append(errs, "paths.chapters_dir is required")
	}

	if c.Paths.RunsDir == "" {
		errs = append(errs, "paths.runs_dir is required")
	}

	if c.Paths.PromptsDir == "" {
		errs = append(errs, "paths.prompts_dir is required")
	}

	if c.Paths.TemplatesDir == "" {
		errs = append(errs, "paths.templates_dir is required")
	}

	return errs
}

// validateMonitoring validates monitoring configuration.
func (c *Config) validateMonitoring() []string {
	var errs []string

	// Validate healthchecks configuration
	if c.Monitoring.Healthchecks.Enabled {
		if c.Monitoring.Healthchecks.TimeoutSeconds <= 0 {
			errs = append(errs, "monitoring.healthchecks.timeout_seconds must be greater than 0 when healthchecks are enabled")
		}
	}

	// Validate HTTP configuration
	if c.Monitoring.HTTP.Enabled {
		if c.Monitoring.HTTP.Port <= 0 {
			errs = append(errs, "monitoring.http.port must be greater than 0 when HTTP monitoring is enabled")
		}
	}

	// Validate cost tracking configuration
	if c.Monitoring.CostTracking.Enabled {
		if c.Monitoring.CostTracking.DailyAlertThreshold < 0 {
			errs = append(errs, "monitoring.cost_tracking.daily_alert_threshold must be greater than or equal to 0")
		}
		if c.Monitoring.CostTracking.LogFile == "" {
			errs = append(errs, "monitoring.cost_tracking.log_file is required when cost tracking is enabled")
		}
	}

	return errs
}

// validateLogging validates logging configuration.
func (c *Config) validateLogging() []string {
	var errs []string

	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.Logging.Level] {
		errs = append(errs, "logging.level must be one of: debug, info, warn, error")
	}

	validFormats := map[string]bool{"json": true, "text": true}
	if !validFormats[c.Logging.Format] {
		errs = append(errs, "logging.format must be one of: json, text")
	}

	if c.Logging.File.Enabled && c.Logging.File.FilenamePattern == "" {
		errs = append(errs, "logging.file.filename_pattern is required when file logging is enabled")
	}

	return errs
}

// for cron job scheduling
// GetTimezone returns the configured timezone as a *time.Location.
func (c *Config) GetTimezone() (*time.Location, error) {
	loc, err := time.LoadLocation(c.Pipeline.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone %q: %w", c.Pipeline.Timezone, err)
	}
	return loc, nil
}

// requireEnv loads a required environment variable and adds an error if missing.
func requireEnv(key string, errs *[]string) string {
	val := os.Getenv(key)
	if val == "" {
		*errs = append(*errs, fmt.Sprintf("%s environment variable is required", key))
	}
	return val
}

// optionalEnv loads an optional environment variable.
func optionalEnv(key string) string {
	return os.Getenv(key)
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
