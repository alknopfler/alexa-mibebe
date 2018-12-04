package function

import (
	"context"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"log"
	cfg "github.com/alknopfler/alexa-mibebe/config"
)

type RecordName struct {
	Email  string    `json:"email"`
	Nombre string    `json:"nombre"`
}

func (r *RecordName) newRecord(email,nombre string) RecordName{
	return RecordName{email,  nombre}
}


func (r *RecordName) AddRecord(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) {
	log.Println("register new baby")

	nombre := request.Intent.Slots["nombre"].Value

	if request.DialogState != "COMPLETED" {
		log.Println("Get into dialog to confirm name 'addBaby intent confirmation'")
		response.AddDialogDirective("Dialog.Delegate", "", "", &request.Intent)
		response.ShouldSessionEnd = false

	} else {
		log.Println(request.Intent.ConfirmationStatus)
		log.Println(nombre)
		if	request.Intent.ConfirmationStatus == "CONFIRMED"{
			exists,err := existRecord("email",getEmail(aContext),cfg.DynamoTableName)
			if err!= nil {
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorExist, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechErrorExist)
				response.ShouldSessionEnd = true
				return
			}
			if nombre != "" && !exists{

				err := createRecord(r.newRecord(getEmail(aContext), "\""+nombre+"\""),cfg.DynamoTableName)
				if err!= nil {
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
				if nombre == "" {
					response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoName, cfg.ImageSmall, cfg.ImageLong)
					response.SetOutputText(cfg.SpeechErrorNoName)
					response.ShouldSessionEnd = true
					return
				}
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorExist, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechErrorExist)
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

func (r *RecordName) GetRecord(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) {
	log.Println("Get baby name")

	listNames, err := getRecordsName("email", getEmail(aContext))
	if err!=nil{
		log.Println("entra por error")
	}
	if len(listNames)==1 {
		response.SetStandardCard(cfg.CardTitle, cfg.SpeechNameis, cfg.ImageSmall, cfg.ImageLong)
		response.SetOutputText(cfg.SpeechNameis + " " + listNames[0].Nombre)
		response.ShouldSessionEnd = true

		return
	}else{
		response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNotExist, cfg.ImageSmall, cfg.ImageLong)
		response.SetOutputText(cfg.SpeechErrorNotExist)
		response.ShouldSessionEnd = true

		return
	}
}