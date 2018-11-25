package function

import (
	cfg "github.com/alknopfler/alexa-mibebe/config"
	"log"
	"context"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
	"strconv"
	"github.com/alknopfler/iso8601duration"
	"time"
)

const parser = "2006-01-02T15:04:05"

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

	pEnt := request.Intent.Slots["kilos"].Value
	pDec := request.Intent.Slots["gramos"].Value
	log.Println("peso: "+pEnt+","+pDec)
	peso, _ :=  strconv.ParseFloat(pEnt + "." + pDec, 64)
	email := getEmail(aContext)

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
				response.ShouldSessionEnd = true
				return
			}
			if pEnt != "" && exists{
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
				if pEnt == "" {
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

func (r *RecordPeso) GetRecord(context context.Context, request *alexa.Request, session *alexa.Session, aContext *alexa.Context, response *alexa.Response){
	log.Println("getting the peso")

	ultimo := request.Intent.Slots["ultimo"].Value
	tiempo := request.Intent.Slots["tiempo"].Value
	email := getEmail(aContext)


	log.Println(ultimo)
	log.Println(tiempo)

	d, err := duration.FromString(request.Intent.Slots["tiempo"].Value)
	if err != nil {
		//TODO return erro
		log.Println("error")
	}
	oldTime := formatNewTime(time.Now().Add(-d.ToDuration()))
	tOldTime := strconv.FormatInt(time.Now().Add(-d.ToDuration()).UTC().UnixNano(),10)
	newTime := getTimeNow()
	tNewTime := strconv.FormatInt(time.Now().UTC().UnixNano(),10)
	log.Println(oldTime, newTime)

	//result, err := getRecordsBetweenDate("fecha", "\""+formatNewTime(oldTime)+"\"", getTimeNow(),cfg.DynamoTablePeso)
	listPesos, err := getRecordsPeso("email", email)
	if err != nil{
		log.Println("ERror getintg records")
	}
	var peso float64 = 0
	for _,val := range listPesos{
		tVal := strconv.FormatInt(ShortDateFromString(val.Fecha).UTC().UnixNano(),10)
		log.Println(tVal)
		a, _ := strconv.Atoi(tVal)
		b, _ := strconv.Atoi(tOldTime)
		c, _ := strconv.Atoi(tVal)
		d, _ := strconv.Atoi(tNewTime)
		if a > b && c <= d {
			peso += val.Peso
		}else{
			log.Println("entra por else")
		}
	}
	log.Println(peso)
}

func formatNewTime(d time.Time) string{
	p := strconv.Itoa
	return p(d.Year()) + "-" + p(int(d.Month())) + "-" + p(d.Day()) + "T" + p(d.Hour()) + ":" + p(d.Minute()) + ":" + p(d.Second())
}

func ShortDateFromString(ds string) time.Time {
 	date, _ := time.Parse(parser, ds)
 	return date
}
