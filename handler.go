package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	cfg "github.com/alknopfler/alexa-mibebe/config"
	f "github.com/alknopfler/alexa-mibebe/function"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
)


var a = &alexa.Alexa{ApplicationID: cfg.AppID, RequestHandler: &MiBebe{}, IgnoreTimestamp: true}

// Mibebe struct for request from the mibebe skill.
type MiBebe struct{}

// Handle processes calls from Lambda
func Handle(ctx context.Context, requestEnv *alexa.RequestEnvelope) (interface{}, error) {
	return a.ProcessRequest(ctx, requestEnv)
}

// OnSessionStarted called when a new session is created.
//
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

	response.ShouldSessionEnd = false
	//response.SetRepromptText(cfg.SpeechNavigate)


	return nil
}

// OnIntent called with a reqeust is received of type IntentRequest
func (h *MiBebe) OnIntent(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) error {
	log.Printf("OnIntent requestId=%s, sessionId=%s, intent=%s", request.RequestID, session.SessionID, request.Intent.Name)

	switch request.Intent.Name {
	case cfg.AddBaby:
		var f f.RecordName
		f.AddRecord(context, request, session, aContext, response)

	case cfg.GetBaby:
		var f f.RecordName
		f.GetRecord(context, request, session, aContext, response)

	case cfg.AddRecordPeso:
		var f f.RecordPeso
		f.AddRecord(context, request, session, aContext, response)

	case cfg.GetRecordPeso:
		var f f.RecordPeso
		f.GetRecord(context, request, session, aContext, response)

	case cfg.AddRecordToma:
		var f f.RecordToma
		f.AddRecord(context, request, session, aContext, response)

	case cfg.GetRecordToma:
		var f f.RecordToma
		f.GetRecord(context, request, session, aContext, response)

	case cfg.GetRecordTomaHoy:
		var f f.RecordToma
		f.GetRecord(context, request, session, aContext, response)

	case cfg.Cancel,cfg.Stop:
		f.Cancel(context, request, session, aContext, response)
		response.ShouldSessionEnd = true
		return nil

	case cfg.Navigate:
		f.Navigate(context, request, session, aContext, response)

	case cfg.Help:
		f.Help(context, request, session, aContext, response)

	default:
		log.Println("Entra por default Intent")
		f.Navigate(context, request, session, aContext, response)
	}
	log.Println("antes del nil onlaunch")
	return nil
}

// OnSessionEnded called with a reqeust is received of type SessionEndedRequest
func (h *MiBebe) OnSessionEnded(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) error {
	log.Printf("OnSessionEnded requestId=%s, sessionId=%s", request.RequestID, session.SessionID)
	response.ShouldSessionEnd = true
	return nil
}


func main() {
	lambda.Start(Handle)
}
