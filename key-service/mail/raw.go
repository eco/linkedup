package mail

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-gomail/gomail"
	"io"
)

func sendRaw(svc *ses.SES, toAddress string, data string) (result *ses.SendRawEmailOutput, err error) {
	source := aws.String(GmEmail)
	destinations := []*string{aws.String(toAddress)}

	msg := gomail.NewMessage()
	msg.SetHeader("From", GmEmail)
	msg.SetHeader("To", toAddress)
	msg.SetHeader("Subject", "Linked Up Shared Contacts")
	msg.SetBody("text/html", getInfoMessageBody())
	msg.Attach("contactInfo.csv", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err = w.Write([]byte(data))
		return err
	}))

	var emailRaw bytes.Buffer
	_, err = msg.WriteTo(&emailRaw)
	if err != nil {
		return
	}
	message := ses.RawMessage{Data: emailRaw.Bytes()}
	input := ses.SendRawEmailInput{Source: source, Destinations: destinations, RawMessage: &message}
	result, err = svc.SendRawEmail(&input)
	return
}

func getInfoMessageBody() string {
	return ExportedBody

}
