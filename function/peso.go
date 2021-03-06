package function

import (
	"context"
	"fmt"
	cfg "github.com/alknopfler/alexa-mibebe/config"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"log"
	"strconv"
)

const parser = "2006-01-02T15:04:05"

type RecordPeso struct {
	Email  		string   `json:"email"`
	Timestamp	string 	 `json:"timestamp"`
	Fecha  		string	 `json:"fecha"`
	Nombre 		string   `json:"nombre"`
	Peso   		float64  `json:"peso"`
}


func (r *RecordPeso) newRecord(email, timestamp, fecha, nombre string, peso float64) RecordPeso{
	return RecordPeso{email, timestamp, fecha, nombre, peso}
}


func (r *RecordPeso) AddRecord(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response) {
	log.Println("register new peso")

	pEnt := request.Intent.Slots["kilos"].Value
	pDec := request.Intent.Slots["gramos"].Value
	log.Println("peso: "+pEnt+","+pDec)
	peso, _ :=  strconv.ParseFloat(pEnt + "." + pDec, 64)
	email := getUserId(aContext)

	if request.DialogState != "COMPLETED" {
		log.Println("Get into dialog to confirm name 'addPeso intent confirmation'")
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
				return
			}
			if pEnt != "" && exists{
				log.Println(email)
				listNames, err := getRecordsName("email", email)
				if err!=nil{
					log.Println("entra por error")
				}
				err = createRecord(r.newRecord(email, "\""+getTimestamp()+"\"", "\""+getTimeNow()+"\"",listNames[0].Nombre, peso), cfg.DynamoTablePeso)
				if err!= nil {
					response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorAddRecord, cfg.ImageSmall, cfg.ImageLong)
					response.SetOutputText(cfg.SpeechErrorAddRecord)
					return
				}
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechOnAddPeso, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechOnAddPeso)
				response.ShouldSessionEnd=true
				return

			}else{
				if pEnt == "" {
					response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoPeso, cfg.ImageSmall, cfg.ImageLong)
					response.SetOutputText(cfg.SpeechErrorNoPeso)
					return
				}
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNotExist, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechErrorNotExist)
				return

			}
		}else{
			response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoRegistered, cfg.ImageSmall, cfg.ImageLong)
			response.SetOutputText(cfg.SpeechErrorNoRegistered)
			return
		}

	}
}

func (r *RecordPeso) GetRecord(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response){
	log.Println("getting the peso")

	email := getUserId(aContext)

	oldTime := "2006-01-02"
	newTime := getTimeNow()
	log.Println(oldTime, newTime)

	listPesos, err := getRecordsPeso("email", email,"\""+oldTime+"\"","\""+newTime+"\"")
	if err != nil || listPesos == nil{
		response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoPeso, cfg.ImageSmall, cfg.ImageLong)
		response.SetOutputText(cfg.SpeechErrorNoPeso)
		response.ShouldSessionEnd = true

		return
	}
	peso := listPesos[len(listPesos)-1].Peso
	kilos, gramos := splitFloat(fmt.Sprintf("%.3f", peso))
	response.SetStandardCard(cfg.CardTitle, cfg.SpeechTotalPeso + listPesos[0].Nombre +" es "+ kilos + " kilogramos" + " con "+ gramos + " gramos", cfg.ImageSmall, cfg.ImageLong)
	response.SetOutputText(cfg.SpeechTotalPeso + listPesos[0].Nombre +" es " + kilos + " kilogramos" + " con "+ gramos + " gramos")
	response.ShouldSessionEnd = true

	return
}


