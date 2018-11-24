package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ericdaugherty/alexa-skills-kit-golang"

	f "github.com/alknopfler/alexa-mibebe/function"
	cfg "github.com/alknopfler/alexa-mibebe/config"
)


var a = &alexa.Alexa{ApplicationID: cfg.AppID, RequestHandler: &MiBebe{}, IgnoreTimestamp: true}


// Mibebe struct for request from the mibebe skill.
type MiBebe struct{}

// Handle processes calls from Lambda
func Handle(ctx context.Context, requestEnv *alexa.RequestEnvelope) (interface{}, error) {
	return a.ProcessRequest(ctx, requestEnv)
}

// OnSessionStarted called when a new session is created.
func (h *MiBebe) OnSessionStarted(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) error {
	log.Printf("OnSessionStarted requestId=%s, sessionId=%s", request.RequestID, session.SessionID)
		//Can be usefull to login internally with the end service

	return nil
}

// OnLaunch called with a reqeust is received of type LaunchRequest
func (h *MiBebe) OnLaunch(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) error {
	log.Printf("OnLaunch requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	response.SetStandardCard(cfg.CardTitle, cfg.SpeechOnLaunch, cfg.ImageSmall, cfg.ImageLong)
	response.SetOutputText(cfg.SpeechOnLaunch)
	response.SetRepromptSSML(cfg.SpeechOnLaunch)

	return nil
}

// OnIntent called with a reqeust is received of type IntentRequest
func (h *MiBebe) OnIntent(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) error {
	log.Printf("OnIntent requestId=%s, sessionId=%s, intent=%s", request.RequestID, session.SessionID, request.Intent.Name)

	switch request.Intent.Name {
	case cfg.AddBaby:
		f.AddBaby(context, request, session, aContext, response)
	default:
		return errors.New("Invalid Intent")
	}

	return nil
}

// OnSessionEnded called with a reqeust is received of type SessionEndedRequest
func (h *MiBebe) OnSessionEnded(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) error {
	log.Printf("OnSessionEnded requestId=%s, sessionId=%s", request.RequestID, session.SessionID)

	return nil
}


func main() {
	lambda.Start(Handle)
}