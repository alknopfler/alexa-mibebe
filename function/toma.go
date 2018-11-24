package function

import (
	"context"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
	cfg "github.com/alknopfler/alexa-mibebe/config"
	"log"
)

type RecordToma struct {
	Email  string    `json:"email"`
	Fecha  string	 `json:"fecha"`
	Nombre string    `json:"nombre"`
	Toma   float64   `json:"toma"`
}

func (r *RecordToma) newRecord(email,fecha,nombre string, toma float64) RecordToma{
	return RecordToma{email, fecha, nombre, toma}
}


func (r *RecordToma) AddRecord(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) {
	log.Println("register toma")

	nombre := request.Intent.Slots["nombre"].Value

	if request.DialogState != "COMPLETED" {
		log.Println("Get into dialog to confirm name 'addBaby intent confirmation'")
		response.AddDialogDirective("Dialog.Delegate", "", "", &request.Intent)
		response.ShouldSessionEnd = false

	} else {
		log.Println(request.Intent.ConfirmationStatus)
		log.Println(nombre)
		if	request.Intent.ConfirmationStatus == "CONFIRMED"{

			if nombre != "" {

				err := createRecord(r.newRecord(getEmail(aContext), getTimeNow(), nombre, 0), cfg.DynamoTablePeso)
				if err != nil{
					response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorAddRecord, cfg.ImageSmall, cfg.ImageLong)
					response.SetOutputText(cfg.SpeechErrorAddRecord)
					response.ShouldSessionEnd = true
					return
				}

				response.SetStandardCard(cfg.CardTitle, cfg.SpeechOnAddRecord, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechOnAddRecord)
				response.ShouldSessionEnd = true
				return

			}else{
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoName, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechErrorNoName)
				response.ShouldSessionEnd = true
				return
			}
		}else{
			response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoRegistered, cfg.ImageSmall, cfg.ImageLong)
			response.SetOutputText(cfg.SpeechErrorNoRegistered)
			response.ShouldSessionEnd = true
			return
		}

	}
}
