package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

type AppConfig struct {
	S3Region     string `default:"us-west-2"`
	SesRegion    string `default:"us-west-2"`
	SesSourceArn string
	MailTo       string        `required:"true"`
	MailFrom     string        `required:"true"`
	Template     string        `required:"true"`
	S3PresignTTL time.Duration `default:"160h"`
}

type TemplateData struct {
	Files   []File
	Subject string
}

type File struct {
	Url      string
	FileName string
}

var (
	appConfig AppConfig
)

func init() {
	err := envconfig.Process("app", &appConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func s3Presign(s3Svc *s3.S3, bucket string, key string) (string, error) {
	req, _ := s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(appConfig.S3PresignTTL)
}

func sendEmail(sesSvc *ses.SES, s3Svc *s3.S3, event events.S3Event) (*ses.SendTemplatedEmailOutput, error) {

	bucketMap := map[string]bool{}

	fileList := TemplateData{
		Subject: "New in S3:",
	}

	for _, record := range event.Records {
		s3 := record.S3
		if url, err := s3Presign(s3Svc, s3.Bucket.Name, s3.Object.Key); err == nil {
			fileList.Files = append(fileList.Files, File{FileName: s3.Object.Key, Url: url})
			bucketMap[s3.Bucket.Name] = true
		} else {
			log.Printf("ERROR: %+v", err)
		}
	}

	for bucket, _ := range bucketMap {
		fileList.Subject += " " + bucket
	}

	templateJson, err := json.Marshal(fileList)
	if err != nil {
		return nil, err
	}

	templateJsonString := string(templateJson)

	toAddressList := []*string{&appConfig.MailTo}

	var sesSourceArn *string
	if appConfig.SesSourceArn != "" {
		sesSourceArn = &appConfig.SesSourceArn
	}

	input := ses.SendTemplatedEmailInput{
		Destination:  &ses.Destination{ToAddresses: toAddressList},
		Source:       &appConfig.MailFrom,
		SourceArn:    sesSourceArn,
		Template:     &appConfig.Template,
		TemplateData: &templateJsonString,
	}

	return sesSvc.SendTemplatedEmail(&input)
}

func handleRequest(ctx context.Context, event events.S3Event) error {

	log.Printf("CONFIG: %+v", appConfig)
	log.Printf("EVENT: %+v", event)

	awsSession := session.Must(session.NewSession())
	sesSvc := ses.New(awsSession, aws.NewConfig().WithRegion(appConfig.SesRegion))
	s3Svc := s3.New(awsSession, aws.NewConfig().WithRegion(appConfig.S3Region))

	emailOutput, err := sendEmail(sesSvc, s3Svc, event)
	if err == nil {
		log.Printf("EMAIL: %+v", emailOutput)
	}
	return err
}

func main() {
	runtime.Start(handleRequest)
}
