package mail

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-gomail/gomail"
	"io"
)

const (
	//LinkedUpSender is the account we send the info request email from
	LinkedUpSender = "linkedup@sfbw.io"
)

func sendRaw(svc *ses.SES, toAddress string, data string) (result *ses.SendRawEmailOutput, err error) {
	source := aws.String(LinkedUpSender)
	destinations := []*string{aws.String(toAddress)}

	//https://linkedup.sfbw.io/s/export/index.html?id=1284763463&token=584353

	msg := gomail.NewMessage()
	msg.SetHeader("From", LinkedUpSender)
	msg.SetHeader("To", toAddress)
	msg.SetHeader("Subject", "Hello!")
	msg.SetBody("text/html", "Hello <b>this is ur data</b>!")
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
