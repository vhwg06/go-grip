package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App          app
		HTTP         http
		Log          log
		PG           pg
		JWT          jwt
		Metrics      metrics
		Swagger      swagger
		Ecommerce    ecommerce
		Notification notification
		Admin        admin
		Payment      payment
		R2           r2
	}

	// App -.
	app struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
		BaseURL string `env:"APP_BASE_URL" envDefault:"http://localhost:8080"`
	}

	// HTTP -.
	http struct {
		Port                 string `env:"HTTP_PORT,required"`
		UsePreforkMode       bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
		CORSAllowedOrigins   string `env:"HTTP_CORS_ALLOWED_ORIGINS" envDefault:"http://localhost:3000,http://127.0.0.1:3000"`
		CORSAllowCredentials bool   `env:"HTTP_CORS_ALLOW_CREDENTIALS" envDefault:"false"`
	}

	// Log -.
	log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG -.
	pg struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	// JWT -.
	jwt struct {
		Secret      string        `env:"JWT_SECRET,required"`
		TokenExpiry time.Duration `env:"JWT_TOKEN_EXPIRY" envDefault:"24h"`
	}

	// Metrics -.
	metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	// Ecommerce -.
	ecommerce struct {
		ThemeColor        string        `env:"ECOMMERCE_THEME_COLOR" envDefault:"#0f766e"`
		MediaMaxBytes     int64         `env:"ECOMMERCE_MEDIA_MAX_BYTES" envDefault:"5242880"`
		InitialImportMax  int           `env:"ECOMMERCE_INITIAL_IMPORT_MAX" envDefault:"25"`
		SchedulerInterval time.Duration `env:"ECOMMERCE_SCHEDULER_INTERVAL" envDefault:"1m"`
	}

	// Notification -.
	notification struct {
		Enabled bool   `env:"NOTIFICATION_ENABLED" envDefault:"false"`
		From    string `env:"NOTIFICATION_FROM" envDefault:"noreply@example.com"`
	}

	// Admin contains admin access configuration.
	admin struct {
		Users string `env:"ADMIN_USERS" envDefault:""`
	}

	// Payment contains gateway configuration.
	payment struct {
		Provider           string `env:"PAYMENT_PROVIDER" envDefault:"epay"`
		MerchantID         string `env:"PAYMENT_MERCHANT_ID" envDefault:""`
		SecretKey          string `env:"PAYMENT_SECRET_KEY" envDefault:""`
		BaseURL            string `env:"PAYMENT_BASE_URL" envDefault:""`
		NotifyURL          string `env:"PAYMENT_NOTIFY_URL" envDefault:""`
		ReturnURL          string `env:"PAYMENT_RETURN_URL" envDefault:""`
		OrderTimeoutMinute int    `env:"PAYMENT_ORDER_TIMEOUT_MINUTES" envDefault:"5"`
	}

	// r2 -.
	r2 struct {
		AccountID   string `env:"R2_ACCOUNT_ID" envDefault:""`
		AccessKeyID string `env:"R2_ACCESS_KEY_ID" envDefault:""`
		SecretKey   string `env:"R2_SECRET_ACCESS_KEY" envDefault:""`
		BucketName  string `env:"R2_BUCKET_NAME" envDefault:""`
		PublicURL   string `env:"R2_PUBLIC_URL" envDefault:""`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
