package function

import (
	cfg "github.com/alknopfler/alexa-mibebe/config"
	"log"
	"context"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
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
	log.Println("register peso")

	peso := request.Intent.Slots["peso"].Value

	if peso != "" {

		err := createRecord(r.newRecord(getEmail(aContext), getTimeNow(), "", 0), cfg.DynamoTablePeso)
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

	}
}






