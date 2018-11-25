package function

import (
	cfg "github.com/alknopfler/alexa-mibebe/config"
	"log"
	"context"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"reflect"
	"strconv"
)

type RecordPeso struct {
	Email  string    `json:"email"`
	Fecha  string	 `json:"fecha"`
	Nombre string    `json:"nombre"`
	Peso   float64   `json:"peso"`
}


func (r *RecordPeso) newRecord(email,fecha,nombre string, peso float64) RecordPeso{
	return RecordPeso{email, fecha, nombre, peso}
}


func (r *RecordPeso) AddRecord(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) {
	log.Println("register new peso")

	p := request.Intent.Slots["peso"].Value
	peso, _ :=  strconv.ParseFloat(p, 64)
	email := getEmail(aContext)

	if request.DialogState != "COMPLETED" {
		log.Println("Get into dialog to confirm name 'addPeo intent confirmation'")
		response.AddDialogDirective("Dialog.Delegate", "", "", &request.Intent)
		response.ShouldSessionEnd = false

	} else {
		log.Println(request.Intent.ConfirmationStatus)
		log.Println(peso)
		if	request.Intent.ConfirmationStatus == "CONFIRMED"{
			exists,err := existRecord("email",email,cfg.DynamoTableName)
			if err!= nil {
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorExist, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechErrorExist)
				response.ShouldSessionEnd = true
				return
			}
			if p != "" && exists{
				nombre, err := getRecord("email", email, cfg.DynamoTableName)
				if err!=nil{
					log.Println("entra por error")
				}
				log.Println(nombre)
				log.Println(getTimeNow())
				log.Println(reflect.ValueOf(nombre).Elem().FieldByName("nombre").String())
				log.Println(peso)
				err = createRecord(r.newRecord(email, getTimeNow(),reflect.ValueOf(nombre).Elem().FieldByName("nombre").String(), peso), cfg.DynamoTableName)
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
				if p == "" {
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






