package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthcheckNoEnvVarSet(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHealthcheckSendUsingInvalidValue(t *testing.T) {
	// Setup
	t.Setenv("CALENVITE_SVC_SEND_USING", "invalid_value")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHealthcheckSendUsingMailgunOK(t *testing.T) {
	// Setup
	t.Setenv("CALENVITE_SVC_SEND_USING", "MAILGUN")
	t.Setenv("CALENVITE_SVC_MAILGUN_KEY", "key")
	t.Setenv("CALENVITE_SVC_MAILGUN_DOMAIN", "domain")

	t.Setenv("CALENVITE_SVC_EMAIL_SENDER_ADDRESS", "me@mail.com")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHealthcheckSendUsingSMTPOK(t *testing.T) {
	// Setup
	t.Setenv("CALENVITE_SVC_SEND_USING", "SMTP")
	t.Setenv("CALENVITE_SVC_SMTP_HOST", "host")
	t.Setenv("CALENVITE_SVC_SMTP_PORT", "123")
	t.Setenv("CALENVITE_SVC_SMTP_USERNAME", "user")
	t.Setenv("CALENVITE_SVC_SMTP_PASSWORD", "pass")

	t.Setenv("CALENVITE_SVC_EMAIL_SENDER_ADDRESS", "me@mail.com")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHealthcheckSendUsingMailgunMissingKeyError(t *testing.T) {
	// Setup
	t.Setenv("CALENVITE_SVC_SEND_USING", "MAILGUN")
	t.Setenv("CALENVITE_SVC_MAILGUN_DOMAIN", "domain")

	t.Setenv("CALENVITE_SVC_EMAIL_SENDER_ADDRESS", "me@mail.com")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHealthcheckSendUsingMailgunMissingDomainError(t *testing.T) {
	// Setup
	t.Setenv("CALENVITE_SVC_SEND_USING", "MAILGUN")
	t.Setenv("CALENVITE_SVC_MAILGUN_KEY", "key")

	t.Setenv("CALENVITE_SVC_EMAIL_SENDER_ADDRESS", "me@mail.com")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHealthcheckSendUsingSMTPMissingHostError(t *testing.T) {
	// Setup
	t.Setenv("CALENVITE_SVC_SEND_USING", "SMTP")
	t.Setenv("CALENVITE_SVC_SMTP_PORT", "123")
	t.Setenv("CALENVITE_SVC_SMTP_USERNAME", "user")
	t.Setenv("CALENVITE_SVC_SMTP_PASSWORD", "pass")

	t.Setenv("CALENVITE_SVC_EMAIL_SENDER_ADDRESS", "me@mail.com")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHealthcheckSendUsingSMTPMissingUserError(t *testing.T) {
	// Setup
	t.Setenv("CALENVITE_SVC_SEND_USING", "SMTP")
	t.Setenv("CALENVITE_SVC_SMTP_HOST", "host")
	t.Setenv("CALENVITE_SVC_SMTP_PORT", "123")
	t.Setenv("CALENVITE_SVC_SMTP_PASSWORD", "pass")

	t.Setenv("CALENVITE_SVC_EMAIL_SENDER_ADDRESS", "me@mail.com")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestHealthcheckSendUsingSMTPMissingPortError(t *testing.T) {
	// Setup
	t.Setenv("CALENVITE_SVC_SEND_USING", "SMTP")
	t.Setenv("CALENVITE_SVC_SMTP_HOST", "host")
	t.Setenv("CALENVITE_SVC_SMTP_USERNAME", "user")
	t.Setenv("CALENVITE_SVC_SMTP_PASSWORD", "pass")

	t.Setenv("CALENVITE_SVC_EMAIL_SENDER_ADDRESS", "me@mail.com")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHealthcheckSendUsingSMTPMissingPasswordError(t *testing.T) {
	// Setup
	t.Setenv("CALENVITE_SVC_SEND_USING", "SMTP")
	t.Setenv("CALENVITE_SVC_SMTP_HOST", "host")
	t.Setenv("CALENVITE_SVC_SMTP_PORT", "123")
	t.Setenv("CALENVITE_SVC_SMTP_USERNAME", "user")

	t.Setenv("CALENVITE_SVC_EMAIL_SENDER_ADDRESS", "me@mail.com")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHealthcheckMissingEmailSenderAddressError(t *testing.T) {
	// Setup
	t.Setenv("CALENVITE_SVC_SEND_USING", "MAILGUN")
	t.Setenv("CALENVITE_SVC_MAILGUN_KEY", "key")
	t.Setenv("CALENVITE_SVC_MAILGUN_DOMAIN", "domain")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.GET("/healthcheck", HealthcheckHandler)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
