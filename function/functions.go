package function

import (
	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"log"
	cfg "github.com/alknopfler/alexa-mibebe/config"
	"context"
)

func AddBaby(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) {
	log.Println("register new baby")

	nombre := request.Intent.Slots["nombre"].Value

	if request.DialogState != "COMPLETED" {
		log.Println("Get into dialog to confirm name 'addBaby intent confirmation'")
		response.AddDialogDirective("Dialog.Delegate", "", "", &request.Intent)
		response.ShouldSessionEnd = false

	} else {

		log.Println(nombre)
		if nombre != "" {

			err := createRecord (newRecord(getEmail(aContext), getTimeNow(), nombre, 0,0))
			if err != nil{
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorAddRecord, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechErrorAddRecord)
			}

			response.SetStandardCard(cfg.CardTitle, cfg.SpeechOnAddRecord, cfg.ImageSmall, cfg.ImageLong)
			response.SetOutputText(cfg.SpeechOnAddRecord)
		}else{
			response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoName, cfg.ImageSmall, cfg.ImageLong)
			response.SetOutputText(cfg.SpeechErrorNoName)
		}
	}
	log.Printf("Set Output speech, value now: %s", response.OutputSpeech.Text)
}

func newRecord(email,fecha,nombre string, peso, toma float64) *cfg.Record{
	return &cfg.Record{email, fecha, nombre, peso, toma}
}

