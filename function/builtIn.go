package function

import (
	"context"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
	cfg "github.com/alknopfler/alexa-mibebe/config"
)

func Cancel(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) {
	response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoRegistered, cfg.ImageSmall, cfg.ImageLong)
	response.SetOutputText(cfg.SpeechErrorNoRegistered)
	response.ShouldSessionEnd = true
	return

}

func Navigate(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) {
	response.SetStandardCard(cfg.CardTitle, cfg.SpeechNavigate, cfg.ImageSmall, cfg.ImageLong)
	response.SetOutputText(cfg.SpeechNavigate)
	response.ShouldSessionEnd = true
	return
}

func Help(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) {
	response.SetStandardCard(cfg.CardTitle, cfg.SpeechHelp, cfg.ImageSmall, cfg.ImageLong)
	response.SetOutputText(cfg.SpeechHelp)
	response.ShouldSessionEnd = false
	return
}

