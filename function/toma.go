package function

import (
	"context"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
	cfg "github.com/alknopfler/alexa-mibebe/config"
	"log"
	"strconv"
)

type RecordToma struct {
	Email  string    `json:"email"`
	Fecha  string	 `json:"fecha"`
	Nombre string    `json:"nombre"`
	Toma   int   	 `json:"toma"`
}

func (r *RecordToma) newRecord(email,fecha,nombre string, toma int) RecordToma{
	return RecordToma{email, fecha, nombre, toma}
}


func (r *RecordToma) AddRecord(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) {
	log.Println("register toma")

	ml, _ :=  strconv.Atoi(request.Intent.Slots["mililitros"].Value)
	email := getEmail(aContext)

	if request.DialogState != "COMPLETED" {
		log.Println("Get into dialog to confirm name 'addPeso intent confirmation'")
		response.AddDialogDirective("Dialog.Delegate", "", "", &request.Intent)
		response.ShouldSessionEnd = false

	} else {
		log.Println(request.Intent.ConfirmationStatus)
		log.Println(ml)
		if	request.Intent.ConfirmationStatus == "CONFIRMED"{
			exists,err := existRecord("email",email,cfg.DynamoTableName)
			if err!= nil {
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorExist, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechErrorExist)
				response.ShouldSessionEnd = true
				return
			}
			if ml != 0 && exists{
				log.Println(email)
				listNames, err := getRecordsName("email", email)
				if err!=nil{
					log.Println("entra por error")
				}
				err = createRecord(r.newRecord(email, getTimeNow(),listNames[0].Nombre, peso), cfg.DynamoTablePeso)
				if err!= nil {
					response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorAddRecord, cfg.ImageSmall, cfg.ImageLong)
					response.SetOutputText(cfg.SpeechErrorAddRecord)
					response.ShouldSessionEnd = true
					return
				}
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechOnAddPeso, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechOnAddPeso)
				response.ShouldSessionEnd = true
				return

			}else{
				if ml == 0 {
					response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoPeso, cfg.ImageSmall, cfg.ImageLong)
					response.SetOutputText(cfg.SpeechErrorNoPeso)
					response.ShouldSessionEnd = true
					return
				}
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNotExist, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechErrorNotExist)
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
