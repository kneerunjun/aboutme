package mail

/* ==================
Very purpose of the package is to encompass sending emails from servers.
- Dialling smpt servers
- Logging on to servers
- Sending simple and mails with attachments
We extend the dialer from go-mail/mail and then have functions implemented on it from over an interface
=====================*/
import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-mail/mail"
)

// This here helps us send emails as notifications and files
// MailNotify : interface to operate on simple mail notifier
type MailNotify interface {
	SendMessage(from, to, sub, body string) error
	SendErrNotification(from, to string) error
	SendFileAttach(from, to, sub, body, filepath string) error
	DialConfig(cfg MailConfig)
}

// MailConfig : flywheel object for creating a new mail notifier
type MailConfig struct {
	Host   string
	Port   int
	UName  string
	Passwd string // AppPassword if you are using 2FA
}

// NewMailNotify : generic constructor for creating a notifier
// Dials the connection with smtp server using the configuration
/*
	notifier, err   := NewMailNotify(MailConfig{], reflect.TypeOf(&GmailNotify{}))
	if err != nil {
		fmt.Errorf("failed to initiate mail notifier %s",err)
	}
*/
func NewMailNotify(cfg MailConfig, typ reflect.Type) (MailNotify, error) {
	itf := reflect.New(typ.Elem()).Interface()
	notify, ok := itf.(MailNotify)
	if !ok || notify == nil {
		return nil, fmt.Errorf("failed to create using reflect type %s", typ.String())
	}
	notify.DialConfig(cfg) //dials the smtp server
	return notify, nil
}

/*
	====================

=======================
*/
type GmailNotify struct {
	*mail.Dialer
}

// SendMessage : sends a simple text/html message to recepient no attachment, but body of the message is customizable
//
// from, to , sub, body  : string literals to set headers on the email
//
// Error when the email isnt send, nil when everything goes as expected.
//
// For slower internet connection since the timeout is 10 seconds, can lead to errors.
// Typically emails with larger size do take a while before they are sent.
/*
	notifier, err   := NewMailNotify(MailConfig{], reflect.TypeOf(&GmailNotify{}))
	if err != nil {
		fmt.Errorf("failed to initiate mail notifier %s",err)
		return
	}
	err :=notifier.SendMessage("michael.scott@dundermifflin.com", "packaging@dundermifflin.com", "Ice queen", "Check this out in bahamas")
	if err !=nil{
		fmt.Errorf("failed to send email %s", err)
		return
	}
	return
*/
func (gmn *GmailNotify) SendMessage(from, to, sub, body string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", sub)
	msg.SetBody("text/html", body)
	gmn.Timeout = 10 * time.Second // kind of ok for a medium sized message, less when sending with attachement
	if err := gmn.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send message %s", err)
	}
	return nil
}

// SendErrNotification : sends an error notification instead of the email. When the intended email isnt sent this can be used to notifiy the user of the problem.
//
// When using SendEmail() since it resides on a separate coroutine, the user would not know until the email send attempt was made. API server would confirm as soon as the email coroutine is fired.
//
// from,to	: Only the sender and receivers email ids, subject is standard with error
//
// Error when sending this fails,
/*
	notifier, err   := NewMailNotify(MailConfig{], reflect.TypeOf(&GmailNotify{}))
	if err != nil {
		fmt.Errorf("failed to initiate mail notifier %s",err)
		return
	}
	go func(){
		err :=notifier.SendMessage("michael.scott@dundermifflin.com", "packaging@dundermifflin.com", "Ice queen", "Check this out in bahamas")
		if err !=nil{
			notifier.SendErrNotification("michael.scott@dundermifflin.com", "packaging@dundermifflin.com")
			return
		}
	}() // co rotuine would continue in the backgroud
	return // API handler returns instantly
*/
func (gmn *GmailNotify) SendErrNotification(from, to string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Error in sending email message")
	body := "We were trying to send a message that failed. Kindly check with the admin for missed notifications"
	msg.SetBody("text/html", body)
	gmn.Timeout = 10 * time.Second // kind of ok for a small sized message, less when sending with attachement
	if err := gmn.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send message %s", err)
	}
	return nil
}

// SendFileAttach : Just regular SendMessage() but with file attachment
//
// from, to , sub, body filepath : string literals to set headers on the email
//
// Error when the email isnt send, nil when everything goes as expected.
//
// Timeout is bumped up to 40 seconds since file attachments generally take more time to upload
/*
	notifier, err   := NewMailNotify(MailConfig{], reflect.TypeOf(&GmailNotify{}))
	if err != nil {
		fmt.Errorf("failed to initiate mail notifier %s",err)
		return
	}
	err :=notifier.SendFileAttach("michael.scott@dundermifflin.com", "packaging@dundermifflin.com", "Ice queen", "Check this out in bahamas", "/images/jan_and_me.jpg")
	if err !=nil{
		fmt.Errorf("failed to send email %s", err)
		return
	}
	return
*/
func (gmn *GmailNotify) SendFileAttach(from, to, sub, body, filepath string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", sub)
	msg.SetBody("text/html", body)
	msg.Attach(filepath)
	gmn.Timeout = 300 * time.Second
	if err := gmn.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send message %s", err)
	}
	return nil
}

func (gmn *GmailNotify) DialConfig(cfg MailConfig) {
	gmn.Dialer = mail.NewDialer(cfg.Host, cfg.Port, cfg.UName, cfg.Passwd)
}
