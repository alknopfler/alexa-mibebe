package function

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	cfg "github.com/alknopfler/alexa-mibebe/config"
)

func createRecord(data interface{}, dbTable string) error {
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
		TableName: aws.String(dbTable),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Println("Error inserting element: "+err.Error())
		return err
	}

	return nil
}


func existRecord(key, value, dbTable string) (bool,error) {
	sess, err := session.NewSession(&aws.Config{Region: &cfg.AWS_Region})
	if err != nil{
		log.Println("Error creating session with aws: " + err.Error())
		return false, err
	}
	svc := dynamodb.New(sess)

	filt := expression.Name(key).Equal(expression.Value(value))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		log.Println("Got error building expression: "+err.Error())
		return false, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(dbTable),
	}

	result, err := svc.Scan(params)

	if err != nil {
		log.Println("Query API call failed: "+err.Error())
		return false, err
	}

	if len(result.Items) > 0 {
		return true, nil
	}
	return false, nil
}

func getRecordsName(key, value string) ([]RecordName, error){

	sess, err := session.NewSession(&aws.Config{Region: &cfg.AWS_Region})
	if err != nil{
		log.Println("Error creating session with aws: " + err.Error())
		return nil, err
	}
	svc := dynamodb.New(sess)

	filt := expression.Name(key).Equal(expression.Value(value))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		log.Println("Got error building expression: "+err.Error())
		return nil, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(cfg.DynamoTableName),
	}

	result, err := svc.Scan(params)
	var item []RecordName
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &item)
	if err != nil {
		log.Println("Failed to unmarshal Record: "+ err.Error())
		return nil, err
	}
	log.Println(item)
	return item, nil

}


func getRecordsPeso(key, value, oldTime, newTime string) ([]RecordPeso, error){

	sess, err := session.NewSession(&aws.Config{Region: &cfg.AWS_Region})
	if err != nil{
		log.Println("Error creating session with aws: " + err.Error())
		return nil, err
	}
	svc := dynamodb.New(sess)

	/*filt := expression.Name(key).Equal(expression.Value(value))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		log.Println("Got error building expression: "+err.Error())
		return nil, err
	}*/

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]*string{
			"#F": 	aws.String("fecha"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":fo" : {
				S: 	aws.String(oldTime),
			},
			":fn" : {
				S: 	aws.String(newTime),
			},

		},
		FilterExpression:          aws.String("#F > :fo"),
		TableName:                 aws.String(cfg.DynamoTablePeso),
	}

	result, err := svc.Scan(params)
	if err!= nil{
		log.Println("error scanning")
	}
	log.Println(result)
	var item []RecordPeso
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &item)
	if err != nil {
		log.Println("Failed to unmarshal Record: "+ err.Error())
		return nil, err
	}
	log.Println("resultado: ")
	log.Println(item)
	return item, nil

}
