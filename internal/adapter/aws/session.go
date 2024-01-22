package aws

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	doOnceAwsSession sync.Once
	awsSession       *session.Session
)

// Session returns a singleton AWS session
func Session() *session.Session {
	doOnceAwsSession.Do(func() {
		awsSession = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	})
	return awsSession
}
