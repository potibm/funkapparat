package initializer

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/potibm/funkapparat/internal/app/config"
	"github.com/potibm/funkapparat/internal/app/exporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

func TestBuildFormatter(t *testing.T) {
	feedConfig := &config.FeedConfig{
		FeedTitle:       "Test Feed",
		FeedDescription: "A test feed",
		FeedLink:        "https://example.com",
		AuthorName:      "Test Author",
		AuthorEmail:     "test@example.com",
	}

	tests := []struct {
		name        string
		cfg         config.ExporterConfig
		feedConfig  *config.FeedConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "rss formatter",
			cfg: config.ExporterConfig{
				Name: "rss-export",
				Type: "rss",
			},
			feedConfig:  feedConfig,
			expectError: false,
		},
		{
			name: "json formatter",
			cfg: config.ExporterConfig{
				Name: "json-export",
				Type: "json",
			},
			feedConfig:  feedConfig,
			expectError: false,
		},
		{
			name: "atom formatter",
			cfg: config.ExporterConfig{
				Name: "atom-export",
				Type: "atom",
			},
			feedConfig:  feedConfig,
			expectError: false,
		},
		{
			name: "unknown type",
			cfg: config.ExporterConfig{
				Name: "unknown-export",
				Type: "csv",
			},
			feedConfig:  feedConfig,
			expectError: true,
			errorMsg:    "unknown exporter type",
		},
		{
			name: "feed formatter without feed config",
			cfg: config.ExporterConfig{
				Name: "rss-export",
				Type: "rss",
			},
			feedConfig:  nil,
			expectError: true,
			errorMsg:    "requires a feed config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := buildFormatter(tt.cfg, tt.feedConfig)
			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, f)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, f)
			}
		})
	}
}

func TestBuildWriter(t *testing.T) {
	s3Client := &s3.Client{}

	tests := []struct {
		name        string
		cfg         config.ExporterConfig
		s3Client    *s3.Client
		expectError bool
		errorMsg    string
	}{
		{
			name: "s3 writer with client and bucket",
			cfg: config.ExporterConfig{
				Name:        "s3-export",
				Destination: "s3",
				Options:     map[string]string{"bucket": "my-bucket"},
			},
			s3Client:    s3Client,
			expectError: false,
		},
		{
			name: "s3 writer without client",
			cfg: config.ExporterConfig{
				Name:        "s3-export",
				Destination: "s3",
				Options:     map[string]string{"bucket": "my-bucket"},
			},
			s3Client:    nil,
			expectError: true,
			errorMsg:    "s3client is not configured",
		},
		{
			name: "s3 writer without bucket option",
			cfg: config.ExporterConfig{
				Name:        "s3-export",
				Destination: "s3",
				Options:     map[string]string{},
			},
			s3Client:    s3Client,
			expectError: true,
			errorMsg:    "requires 'bucket' option",
		},
		{
			name: "file writer with dir",
			cfg: config.ExporterConfig{
				Name:        "file-export",
				Destination: "file",
				Options:     map[string]string{"dir": "/tmp/exports"},
			},
			s3Client:    nil,
			expectError: false,
		},
		{
			name: "file writer without dir option",
			cfg: config.ExporterConfig{
				Name:        "file-export",
				Destination: "file",
				Options:     map[string]string{},
			},
			s3Client:    nil,
			expectError: true,
			errorMsg:    "requires 'dir' option",
		},
		{
			name: "unknown destination",
			cfg: config.ExporterConfig{
				Name:        "ftp-export",
				Destination: "ftp",
			},
			s3Client:    nil,
			expectError: true,
			errorMsg:    "unknown destination",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := buildWriter(tt.cfg, tt.s3Client)
			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, w)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, w)
			}
		})
	}
}

func TestBootstrapExporters(t *testing.T) {
	ctx := context.Background()
	feedConfig := &config.FeedConfig{
		FeedTitle:       "Test Feed",
		FeedDescription: "A test feed",
		FeedLink:        "https://example.com",
		AuthorName:      "Test Author",
		AuthorEmail:     "test@example.com",
	}

	t.Run("empty configs", func(t *testing.T) {
		exporters, err := BootstrapExporters(ctx, "1.0.0", feedConfig, nil, nil, discardLogger())
		require.NoError(t, err)
		assert.Empty(t, exporters)
	})

	t.Run("disabled exporter is skipped", func(t *testing.T) {
		configs := []config.ExporterConfig{
			{
				Name:    "disabled-rss",
				Type:    "rss",
				Enabled: false,
			},
		}
		exporters, err := BootstrapExporters(ctx, "1.0.0", feedConfig, configs, nil, discardLogger())
		require.NoError(t, err)
		assert.Empty(t, exporters)
	})

	t.Run("unknown type is skipped", func(t *testing.T) {
		configs := []config.ExporterConfig{
			{
				Name:        "bad-type",
				Type:        "csv",
				Destination: "file",
				Enabled:     true,
				Options:     map[string]string{"dir": "/tmp"},
			},
		}
		exporters, err := BootstrapExporters(ctx, "1.0.0", feedConfig, configs, nil, discardLogger())
		require.NoError(t, err)
		assert.Empty(t, exporters)
	})

	t.Run("unknown destination is skipped", func(t *testing.T) {
		configs := []config.ExporterConfig{
			{
				Name:        "bad-dest",
				Type:        "rss",
				Destination: "ftp",
				Enabled:     true,
			},
		}
		exporters, err := BootstrapExporters(ctx, "1.0.0", feedConfig, configs, nil, discardLogger())
		require.NoError(t, err)
		assert.Empty(t, exporters)
	})

	t.Run("s3 without client returns error", func(t *testing.T) {
		configs := []config.ExporterConfig{
			{
				Name:        "s3-rss",
				Type:        "rss",
				Destination: "s3",
				Enabled:     true,
				Options:     map[string]string{"bucket": "my-bucket"},
			},
		}
		exporters, err := BootstrapExporters(ctx, "1.0.0", feedConfig, configs, nil, discardLogger())
		require.Error(t, err)
		assert.Nil(t, exporters)
		assert.Contains(t, err.Error(), "s3client is not configured")
	})

	t.Run("s3 without bucket option returns error", func(t *testing.T) {
		configs := []config.ExporterConfig{
			{
				Name:        "s3-rss",
				Type:        "rss",
				Destination: "s3",
				Enabled:     true,
				Options:     map[string]string{},
			},
		}
		s3Client := &s3.Client{}
		exporters, err := BootstrapExporters(ctx, "1.0.0", feedConfig, configs, s3Client, discardLogger())
		require.Error(t, err)
		assert.Nil(t, exporters)
		assert.Contains(t, err.Error(), "requires 'bucket' option")
	})

	t.Run("file without dir option returns error", func(t *testing.T) {
		configs := []config.ExporterConfig{
			{
				Name:        "file-rss",
				Type:        "rss",
				Destination: "file",
				Enabled:     true,
				Options:     map[string]string{},
			},
		}
		exporters, err := BootstrapExporters(ctx, "1.0.0", feedConfig, configs, nil, discardLogger())
		require.Error(t, err)
		assert.Nil(t, exporters)
		assert.Contains(t, err.Error(), "requires 'dir' option")
	})

	t.Run("successful file rss exporter", func(t *testing.T) {
		configs := []config.ExporterConfig{
			{
				Name:        "file-rss",
				Type:        "rss",
				Destination: "file",
				Filename:    "feed.xml",
				Enabled:     true,
				Options:     map[string]string{"dir": "/tmp/exports"},
			},
		}
		exporters, err := BootstrapExporters(ctx, "1.0.0", feedConfig, configs, nil, discardLogger())
		require.NoError(t, err)
		require.Len(t, exporters, 1)
		assert.Equal(t, "file-rss", exporters[0].Name())
		assert.IsType(t, &exporter.UniversalExporter{}, exporters[0])
	})

	t.Run("successful s3 json exporter", func(t *testing.T) {
		configs := []config.ExporterConfig{
			{
				Name:        "s3-json",
				Type:        "json",
				Destination: "s3",
				Filename:    "feed.json",
				Enabled:     true,
				Options:     map[string]string{"bucket": "my-bucket"},
			},
		}
		s3Client := &s3.Client{}
		exporters, err := BootstrapExporters(ctx, "1.0.0", feedConfig, configs, s3Client, discardLogger())
		require.NoError(t, err)
		require.Len(t, exporters, 1)
		assert.Equal(t, "s3-json", exporters[0].Name())
		assert.IsType(t, &exporter.UniversalExporter{}, exporters[0])
	})

	t.Run("feed formatter without feed config returns error", func(t *testing.T) {
		configs := []config.ExporterConfig{
			{
				Name:        "file-rss",
				Type:        "rss",
				Destination: "file",
				Filename:    "feed.xml",
				Enabled:     true,
				Options:     map[string]string{"dir": "/tmp/exports"},
			},
		}
		exporters, err := BootstrapExporters(ctx, "1.0.0", nil, configs, nil, discardLogger())
		require.NoError(t, err)
		assert.Empty(t, exporters)
	})

	t.Run("mixed valid and invalid exporters", func(t *testing.T) {
		configs := []config.ExporterConfig{
			{
				Name:    "disabled",
				Type:    "rss",
				Enabled: false,
			},
			{
				Name:        "bad-type",
				Type:        "csv",
				Destination: "file",
				Enabled:     true,
				Options:     map[string]string{"dir": "/tmp"},
			},
			{
				Name:        "valid-atom",
				Type:        "atom",
				Destination: "file",
				Filename:    "feed.atom",
				Enabled:     true,
				Options:     map[string]string{"dir": "/tmp/exports"},
			},
		}
		exporters, err := BootstrapExporters(ctx, "1.0.0", feedConfig, configs, nil, discardLogger())
		require.NoError(t, err)
		require.Len(t, exporters, 1)
		assert.Equal(t, "valid-atom", exporters[0].Name())
	})
}
