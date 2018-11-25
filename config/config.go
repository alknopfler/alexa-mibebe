package config

import "os"

const(
	AppID					= "amzn1.ask.skill.d54a5af8-fa55-4236-8eb4-b1e39eb9668a"

	CardTitle 				= "Mi bebé"
	ImageSmall 				= "https://raw.githubusercontent.com/alknopfler/alexa-mibebe/master/images/icono108.jpg"
	ImageLong 				= "https://raw.githubusercontent.com/alknopfler/alexa-mibebe/master/images/icono500.jpg"

	AddBaby					= "addBaby"
	AddRecordPeso			= "addPeso"
	AddRecordToma			= "addToma"
)

var(
	AWS_Key 				= os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_Secret 				= os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_Region				= "eu-west-1"
	DynamoTablePeso			= "mibebe-peso"
	DynamoTableToma			= "mibebe-toma"
	DynamoTableName			= "mibebe-nombre"
)
