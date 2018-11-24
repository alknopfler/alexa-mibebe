package function

import (
	cfg "github.com/alknopfler/alexa-mibebe/config"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

)
func createRecord(data *cfg.Record) error {
	sess, err := session.NewSession(&aws.Config{Region: &cfg.AWS_Region})
	if err != nil{
		log.Println("Error creating session with aws: " + err.Error())
		return err
	}
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(&data)
	if err != nil {
		log.Println("Got error marshalling map: "+err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(cfg.DynamoTable),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Println("Error inserting element: "+err.Error())
		return err
	}

	return nil
}
