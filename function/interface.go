package function

import (
	"context"
	"github.com/ericdaugherty/alexa-skills-kit-golang"
)

type Record interface {
	AddRecord(context.Context, *alexa.Request, *alexa.Session, *alexa.Context, *alexa.Response)
}
