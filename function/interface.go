package function

import (
	"context"
)

type Record interface {
	AddRecord(context.Context, *alexa.Request, *alexa.Session, *alexa.Context, *alexa.Response)
	GetRecord(context.Context, *alexa.Request, *alexa.Session, *alexa.Context, *alexa.Response)
}
