package mail

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMailNotify(t *testing.T) {
	notify, err := NewMailNotify(MailConfig{Host: "smtp.gmail.com", Port: 587, UName: "awatiniranjan@gmail.com", Passwd: "imbilafrkzilxvwv"}, reflect.TypeOf(&GmailNotify{}))
	assert.Nil(t, err, "Unexpected error creating new notify object")
	assert.NotNil(t, notify, "Unexpected nil value on notify object")
}

func TestSendEmail(t *testing.T) {
	notify, _ := NewMailNotify(MailConfig{Host: "smtp.gmail.com", Port: 587, UName: "awatiniranjan@gmail.com", Passwd: ""}, reflect.TypeOf(&GmailNotify{}))
	body := "Hi there <br>This is a test email message, if you can see this smtp automation seems to be working well<br><br>Best regards,<br>Niranjan"
	err := notify.SendMessage("awatiniranjan@gmail.com", "kneerunjun@gmail.com", "Test email", body)
	assert.Nil(t, err, "Unexpected error when sending email")

	// We need to see if I can change the sender to kneerunjun@gmail.com despite being logged on the smtp server with awatiniranjan@gmail.com
	// IMP: this will not have any effect whatsoever, will send a email from awatiniranjan@gmail.com to  awatiniranjan@gmail.com
	// NOTE: changing the email id here does not change anything, it always follows what you have used when logging into smtp server
	err = notify.SendMessage("kneerunjun@gmail.com", "awatiniranjan@gmail.com", "Test email", body)
	assert.Nil(t, err, "Unexpected error when sending email")

	// TEST: now lets try sending email with a name in the To field
	err = notify.SendMessage("NiranjanAwati", "awatiniranjan@gmail.com", "Test email", body)
	assert.NotNil(t, err, "To field always expects the same address as used when logging into the smtp server, this should have failed")
}

func TestSendErrEmail(t *testing.T) {
	notify, _ := NewMailNotify(MailConfig{Host: "smtp.gmail.com", Port: 587, UName: "awatiniranjan@gmail.com", Passwd: ""}, reflect.TypeOf(&GmailNotify{}))
	err := notify.SendErrNotification("awatiniranjan@gmail.com", "kneerunjun@gmail.com")
	assert.Nil(t, err, "Unexpected error when sending error email")
}

func TestSendEmailWithAttach(t *testing.T) {
	notify, _ := NewMailNotify(MailConfig{Host: "smtp.gmail.com", Port: 587, UName: "awatiniranjan@gmail.com", Passwd: ""}, reflect.TypeOf(&GmailNotify{}))
	body := "Hi there <br>This is a test email message, if you can see this message with attachment alongside,I guess everything is working fine<br><br>Best regards,<br>Niranjan"
	err := notify.SendFileAttach("awatiniranjan@gmail.com", "kneerunjun@gmail.com", "Test email", body, "/home/niranjan/Pictures/louis-reed-53jnUK5LqEY-unsplash.jpg")
	assert.Nil(t, err, "Unexpected error when sending email")
}
