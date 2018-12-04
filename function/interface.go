package function

import (
	"context"
	"github.com/alknopfler/alexa-skills-kit-golang"
)

type Record interface {
	AddRecord(context.Context, *alexa.Request, *alexa.Session, *alexa.Context, *alexa.Response)
	GetRecord(context.Context, *alexa.Request, *alexa.Session, *alexa.Context, *alexa.Response)
}
