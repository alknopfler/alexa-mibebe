package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ericdaugherty/alexa-skills-kit-golang"

	f "github.com/alknopfler/alkalarm-alexa-skills/function"
	cfg "github.com/alknopfler/alexa-mibebe/config"
)


var a = &alexa.Alexa{ApplicationID: cfg.AppID, RequestHandler: &AlkAlarm{}, IgnoreTimestamp: true}


// Alkalarm struct for request from the alkalarm skill.
type AlkAlarm struct{}

// Handle processes calls from Lambda
func Handle(ctx context.Context, requestEnv *alexa.RequestEnvelope) (interface{}, error) {
	return a.ProcessRequest(ctx, requestEnv)
}

// OnSessionStarted called when a new session is created.
func (h *AlkAlarm) OnSessionStarted(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) error {
	log.Printf("OnSessionStarted requestId=%s, sessionId=%s", request.RequestID, session.SessionID)
		//Can be usefull to login internally with the end service

	return nil
}

// OnLaunch called with a reqeust is received of type LaunchRequest
func (h *AlkAlarm) OnLaunch(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) error {
	log.Printf("OnLaunch requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	response.SetStandardCard(cfg.CardTitle, cfg.SpeechOnLaunch, cfg.ImageSmall, cfg.ImageLong)
	response.SetOutputText(cfg.SpeechOnLaunch)
	response.SetRepromptSSML(cfg.SpeechOnLaunch)

	response.ShouldSessionEnd = true

	return nil
}

// OnIntent called with a reqeust is received of type IntentRequest
func (h *AlkAlarm) OnIntent(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) error {
	log.Printf("OnIntent requestId=%s, sessionId=%s, intent=%s", request.RequestID, session.SessionID, request.Intent.Name)

	switch request.Intent.Name {
	case cfg.ActiveFullIntent:
		f.ActivateAlarmFull(request,response)
	case cfg.ActivePartialIntent:
		f.ActivateAlarmPartial(request,response)
	case cfg.DeactiveIntent:
		f.DeactivateAlarm(request,response)
	case cfg.StatusIntent:
		f.StatusAlarm(request,response)
	default:
		return errors.New("Invalid Intent")
	}

	return nil
}

// OnSessionEnded called with a reqeust is received of type SessionEndedRequest
func (h *AlkAlarm) OnSessionEnded(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) error {
	log.Printf("OnSessionEnded requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	return nil
}


func main() {
	lambda.Start(Handle)
}