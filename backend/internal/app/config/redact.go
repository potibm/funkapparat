package config

import "net/url"

const redacted = "***REDACTED***"

func (c Config) RedactConfigForDisplay() Config {
	result := c

	result.Sentry.DSN = redacted
	if result.S3Client != nil {
		result.S3Client = &S3ClientConfig{
			AccessKeyID:     redacted,
			SecretAccessKey: redacted,
			Region:          result.S3Client.Region,
			Endpoint:        result.S3Client.Endpoint,
			UsePathStyle:    result.S3Client.UsePathStyle,
		}
	}

	result.App.RedisURL = result.App.RedisURL.Redacted()

	return result
}

func (r RedisURL) Redacted() RedisURL {
	return RedisURL(redactURLPassword(string(r)))
}

func redactURLPassword(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	if parsedURL.User != nil {
		parsedURL.User = url.UserPassword(parsedURL.User.Username(), redacted)

		return parsedURL.String()
	}

	return rawURL
}
