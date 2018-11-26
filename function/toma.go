package function

import (
	"context"
	cfg "github.com/alknopfler/alexa-mibebe/config"
	"github.com/alknopfler/iso8601duration"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"log"
	"strconv"
	"time"
)

type RecordToma struct {
	Email  		string    	`json:"email"`
	Timestamp	string		`json:"timestamp"`
	Fecha  		string	    `json:"fecha"`
	Nombre 		string      `json:"nombre"`
	Toma   		int     	`json:"toma"`
}

func (r *RecordToma) newRecord(email, timestamp, fecha, nombre string, toma int) RecordToma{
	return RecordToma{email, timestamp, fecha, nombre, toma}
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
				err = createRecord(r.newRecord(email,"\""+getTimestamp()+"\"", "\""+getTimeNow()+"\"","\""+listNames[0].Nombre+"\"", ml), cfg.DynamoTableToma)
				if err!= nil {
					response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorAddRecord, cfg.ImageSmall, cfg.ImageLong)
					response.SetOutputText(cfg.SpeechErrorAddRecord)
					response.ShouldSessionEnd = true
					return
				}
				response.SetStandardCard(cfg.CardTitle, cfg.SpeechOnAddToma, cfg.ImageSmall, cfg.ImageLong)
				response.SetOutputText(cfg.SpeechOnAddToma)
				response.ShouldSessionEnd = true
				return

			}else{
				if ml == 0 {
					response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoPeso, cfg.ImageSmall, cfg.ImageLong)
					response.SetOutputText(cfg.SpeechErrorNoToma)
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


func (r *RecordToma) GetRecord(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response){
	log.Println("getting the toma")

	email := getEmail(aContext)
 	var oldTime, newTime string
	if request.Intent.Slots["tiempo"].Value == ""{
 		oldTime = getTimeNow()
 		newTime = oldTime
	}else {
		d, err := duration.FromString(request.Intent.Slots["tiempo"].Value)
		if err != nil {
			log.Println("error formatting string")
			response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoToma, cfg.ImageSmall, cfg.ImageLong)
			response.SetOutputText(cfg.SpeechErrorNoToma)
			response.ShouldSessionEnd = true
			return
		}
		oldTime = formatNewTime(time.Now().Add(-d.ToDuration()))
		newTime = getTimeNow()
	}
	log.Println(oldTime, newTime)

	listTomas, err := getRecordsToma("email", email,"\""+oldTime+"\"","\""+newTime+"\"")
	if err != nil {
		response.SetStandardCard(cfg.CardTitle, cfg.SpeechErrorNoToma, cfg.ImageSmall, cfg.ImageLong)
		response.SetOutputText(cfg.SpeechErrorNoToma)
		response.ShouldSessionEnd = true
	}
	var toma int
	for _, val := range listTomas{
		toma += val.Toma
	}
	response.SetStandardCard(cfg.CardTitle, cfg.SpeechTotalToma + strconv.Itoa(toma) + " mililitros", cfg.ImageSmall, cfg.ImageLong)
	response.SetOutputText(cfg.SpeechTotalPeso + strconv.Itoa(toma) + " mililitros")
	response.ShouldSessionEnd = true
	return

}

