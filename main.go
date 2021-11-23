package main

import (
	"calenvite/models"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mailgun/mailgun-go/v4"
	"gopkg.in/go-playground/validator.v9"
	gomail "gopkg.in/gomail.v2"
)

var settings = models.Settings{
	Mailgun: models.Mailgun{
		Domain:    os.Getenv("CALENVITE_SVC_MAILGUN_DOMAIN"),
		SecretKey: os.Getenv("CALENVITE_SVC_MAILGUN_KEY"),
	},
	SMTP: models.SMTP{
		Host:     os.Getenv("CALENVITE_SVC_SMTP_HOST"),
		Port:     os.Getenv("CALENVITE_SVC_SMTP_PORT"),
		User:     os.Getenv("CALENVITE_SVC_SMTP_USER"),
		Password: os.Getenv("CALENVITE_SVC_SMTP_PASSWORD"),
	},
	SenderAddress: os.Getenv("CALENVITE_SVC_EMAIL_SENDER_ADDRESS"),
	SendUsing:     os.Getenv("CALENVITE_SVC_SEND_USING"),
	ServicePort:   os.Getenv("CALENVITE_SVC_PORT"),
}

var API500Error = models.APIResponse{
	Message:     "ERROR",
	StatusCode:  http.StatusInternalServerError,
	ErrorFields: nil,
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func HealthcheckHandler(c echo.Context) error {

	if value, ok := os.LookupEnv("CALENVITE_SVC_SEND_USING"); ok {
		if value != "MAILGUN" && value != "SMTP" {
			log.Printf("Env var CALENVITE_SVC_SEND_USING value invalid: %s. Valid Values: MAILGUN or SMTP. Check documentation\n", value)
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else {
		log.Println("Env var CALENVITE_SVC_SEND_USING not set. Check documentation")
		return c.JSON(http.StatusInternalServerError, nil)

	}

	if value, ok := os.LookupEnv("CALENVITE_SVC_SEND_USING"); ok {
		if value == "MAILGUN" {
			if _, ok := os.LookupEnv("CALENVITE_SVC_MAILGUN_DOMAIN"); !ok {
				log.Println("Env var CALENVITE_SVC_MAILGUN_DOMAIN not set. Check documentation")
				return c.JSON(http.StatusInternalServerError, nil)
			}

			if _, ok := os.LookupEnv("CALENVITE_SVC_MAILGUN_KEY"); !ok {
				log.Println("Env var CALENVITE_SVC_MAILGUN_KEY not set. Check documentation")
				return c.JSON(http.StatusInternalServerError, nil)

			}
		} else if value == "SMTP" {
			if _, ok := os.LookupEnv("CALENVITE_SVC_SMTP_HOST"); !ok {
				log.Println("Env var CALENVITE_SVC_SMTP_HOST not set. Check documentation")
				return c.JSON(http.StatusInternalServerError, nil)
			}

			if _, ok := os.LookupEnv("CALENVITE_SVC_SMTP_USER"); !ok {
				log.Println("Env var CALENVITE_SVC_SMTP_USER not set. Check documentation")
				return c.JSON(http.StatusInternalServerError, nil)
			}

			if _, ok := os.LookupEnv("CALENVITE_SVC_SMTP_PASSWORD"); !ok {
				log.Println("Env var CALENVITE_SVC_SMTP_PASSWORD not set. Check documentation")
				return c.JSON(http.StatusInternalServerError, nil)
			}

			if _, ok := os.LookupEnv("CALENVITE_SVC_SMTP_PORT"); !ok {
				log.Println("Env var CALENVITE_SVC_SMTP_PORT not set. Check documentation")
				return c.JSON(http.StatusInternalServerError, nil)
			}
		}
	}

	if _, ok := os.LookupEnv("CALENVITE_SVC_EMAIL_SENDER_ADDRESS"); !ok {
		log.Println("Env var CALENVITE_SVC_EMAIL_SENDER_ADDRESS is not set. Check documentation")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)

}

func InviteHandler(c echo.Context) error {

	payload := new(models.RequestPayload)

	if err := c.Bind(payload); err != nil {
		return err
	}

	// validate payload
	validate = validator.New()

	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.Struct(payload)

	if err != nil {

		var fieldErr = []models.ErrorFields{}

		for _, err := range err.(validator.ValidationErrors) {

			fieldErr = append(fieldErr, models.ErrorFields{
				Field:   err.Field(),
				Message: "",
				Code:    err.Tag(),
			})
		}

		var res models.APIResponse
		res.Message = "INVALID_PAYLOAD"
		res.ErrorFields = fieldErr
		res.StatusCode = http.StatusBadRequest

		return c.JSON(http.StatusBadRequest, res)
	}

	var emailsAddress []string

	for _, u := range payload.Users {
		emailsAddress = append(emailsAddress, u.Email)
	}

	var calendarICSFileAttendees string
	var calendarICSFileOrganizer string

	// create ics file
	if !reflect.ValueOf(payload.Invitation).IsNil() {

		startAt, err := time.Parse(time.RFC3339, payload.Invitation.StartAt)

		if err != nil {
			log.Printf("Error parsing StartAt value: %s", payload.Invitation.StartAt)
			return c.JSON(http.StatusInternalServerError, API500Error)
		}

		endAt, err := time.Parse(time.RFC3339, payload.Invitation.EndAt)

		if err != nil {
			log.Printf("Error parsing EndAt value: %s", payload.Invitation.EndAt)
			return c.JSON(http.StatusInternalServerError, API500Error)
		}

		var attendees = make(map[string]string)

		for _, u := range payload.Users {
			attendees[u.Email] = u.FullName
		}

		// we need to create two separate ics files
		// otherwise the organizator won't be able
		// to add the event automatically to his calendar
		// this is why we usee ics.MethodRequest for attendees
		// and ics.MethodPublish for the organizer
		// https://datatracker.ietf.org/doc/html/rfc2446#section-3.2
		icsAttendees := createICS(startAt, endAt, payload.Invitation.EventSummary, payload.Invitation.Description, payload.Invitation.Location, payload.Invitation.OrganizerEmail, payload.Invitation.OrganizerFullName, attendees, ics.MethodRequest)

		icsOrganizer := createICS(startAt, endAt, payload.Invitation.EventSummary, payload.Invitation.Description, payload.Invitation.Location, payload.Invitation.OrganizerEmail, payload.Invitation.OrganizerFullName, attendees, ics.MethodPublish)

		icsFileAttendees, err := ioutil.TempFile(os.TempDir(), "*.ics")

		if err != nil {
			log.Printf("Failed to write to temporary file: %s", err)
			return c.JSON(http.StatusInternalServerError, API500Error)
		}

		icsFileOrganizer, err := ioutil.TempFile(os.TempDir(), "*.ics")

		if err != nil {
			log.Printf("Failed to write to temporary file: %s", err)
			return c.JSON(http.StatusInternalServerError, API500Error)
		}

		icsFileAttendees.Write([]byte(icsAttendees))
		icsFileOrganizer.Write([]byte(icsOrganizer))

		// Close the file
		err = icsFileAttendees.Close()
		if err != nil {
			log.Printf("Unable to close file %s", err)
			return c.JSON(http.StatusInternalServerError, API500Error)
		}

		err = icsFileOrganizer.Close()
		if err != nil {
			log.Printf("Unable to close file %s", err)
			return c.JSON(http.StatusInternalServerError, API500Error)
		}

		calendarICSFileAttendees = icsFileAttendees.Name()
		calendarICSFileOrganizer = icsFileOrganizer.Name()

		defer os.Remove(icsFileAttendees.Name())
		defer os.Remove(icsFileOrganizer.Name())
	}

	// send emails to users
	if settings.SendUsing == "MAILGUN" {
		_, err = sendEmailMailgun(payload.EmailSubject, payload.EmailBody, payload.EmailIsHTML, emailsAddress, calendarICSFileAttendees)
	} else {
		err = sendEmailSMTP(emailsAddress, payload.EmailSubject, payload.EmailBody, payload.EmailIsHTML, calendarICSFileAttendees)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, API500Error)
	}

	// send a separate email to the organizer
	// so the event gets created on his calendar
	// (Only if an invitation is created)
	if calendarICSFileOrganizer != "" {

		var organizer []string
		organizer = append(organizer, payload.Invitation.OrganizerEmail)

		// send email to organizer
		if settings.SendUsing == "MAILGUN" {
			_, err = sendEmailMailgun(payload.EmailSubject, payload.EmailBody, payload.EmailIsHTML, organizer, calendarICSFileOrganizer)
		} else {
			err = sendEmailSMTP(organizer, payload.EmailSubject, payload.EmailBody, payload.EmailIsHTML, calendarICSFileOrganizer)
		}
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, API500Error)
	}

	var res models.APIResponse

	res = models.APIResponse{
		Message:    "SENT_OK",
		StatusCode: http.StatusOK,
	}

	return c.JSON(res.StatusCode, res)

}

func sendEmailMailgun(subject string, body string, isHTML bool, recipients []string, attachment string) (string, error) {

	mg := mailgun.NewMailgun(settings.Mailgun.Domain, settings.Mailgun.SecretKey)

	sender := settings.SenderAddress

	message := mg.NewMessage(sender, subject, body, recipients[0])

	if isHTML {
		message.SetHtml(body)
	}

	if attachment != "" {
		message.AddAttachment(attachment)
	}

	for _, e := range recipients {
		message.AddCC(e)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	_, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return id, nil
}

func sendEmailSMTP(recipients []string, subject string, body string, isHTML bool, attachment string) error {

	m := gomail.NewMessage()

	m.SetHeader("From", settings.SenderAddress)
	m.SetHeader("Subject", subject)

	for _, email := range recipients {
		m.SetHeader("To", email)
	}

	if isHTML {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	m.Attach(attachment)

	// Settings for SMTP server
	port, err := strconv.Atoi(settings.SMTP.Port)

	if err != nil {
		fmt.Println(err)
		return err
	}

	d := gomail.NewDialer(settings.SMTP.Host, port, settings.SMTP.User, settings.SMTP.Password)

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}

	return err
}

func createICS(startAt time.Time, endAt time.Time, summary string, description string, location string, organizerEmail string, organizerFullName string, attendees map[string]string, methodRequest ics.Method) string {

	cal := ics.NewCalendar()
	cal.SetMethod(methodRequest)
	event := cal.AddEvent(fmt.Sprintf("%s%s", uuid.New(), organizerEmail))
	event.SetTimeTransparency(ics.TransparencyOpaque)
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetStartAt(startAt)
	event.SetEndAt(endAt)
	event.SetSummary(fmt.Sprintf("%s", summary))
	event.SetDescription(fmt.Sprintf("%s", description))
	event.SetLocation(fmt.Sprintf("%s", location))
	event.SetOrganizer(organizerEmail, ics.WithCN(organizerFullName))

	for email, fullName := range attendees {
		event.AddAttendee(email, ics.WithCN(fullName), ics.CalendarUserTypeIndividual, ics.ParticipationStatusNeedsAction, ics.ParticipationRoleReqParticipant, ics.WithRSVP(true))
	}

	return cal.Serialize()
}

func main() {

	// Echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.GET("/healthcheck", HealthcheckHandler)
	e.POST("/invite/", InviteHandler)

	// set default port
	if settings.ServicePort == "" {
		settings.ServicePort = "8000"
	}

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", settings.ServicePort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
