package emails

import (
	"log"
	"strconv"
    "fmt"
	"example.com/weather-report/config"
    "example.com/weather-report/weather-api"
	"gopkg.in/mail.v2"
)


func sendEmail(
    toEmail string, 
    subject string, 
    message string, 
    conf *config.Config,
) error {
    m := mail.NewMessage()
    m.SetHeader("From", conf.AdminEmailConfig.Email)
    m.SetHeader("To", toEmail)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", message)

    tlsPort, err := strconv.Atoi(conf.SmtpServerConfig.TlsPort)
    if err != nil {
        log.Printf("Error converting TLS port to int: %s", err)
        return err
    }
    d := mail.NewDialer(
        conf.SmtpServerConfig.Host, 
        tlsPort, 
        conf.AdminEmailConfig.Email, 
        conf.AdminEmailConfig.AppPassword,
    )
    
    if err := d.DialAndSend(m); err != nil {
        log.Printf("Error sending email: %s", err)
        return err
    }
    return nil
}

func SendConfirmationEmail(
    toEmail string, 
    token string,
    conf *config.Config,
) error {
    subject := "Weather Report Confirmation"
    message := fmt.Sprintf("Your token: %s", token)
    
    return sendEmail(
        toEmail,
        subject,
        message,
        conf,
    )
}

func SendWeatherReportEmail(
    toEmail string,
    weatherReport *weatherApi.WeatherResponse,
    conf *config.Config,
) error {
    subject := "Weather Report"
    message := fmt.Sprintf("Temperature: %f\nHumidity: %d\nDescription: %s", 
        weatherReport.TempC, weatherReport.Humidity, weatherReport.Description)
    
    return sendEmail(toEmail, subject, message, conf)
}
