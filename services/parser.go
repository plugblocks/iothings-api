package services

import (
	"gitlab.com/plugblocks/iothings-api/models/sigfox"
	"gitlab.com/plugblocks/iothings-api/models"
)

func parse(message sigfox.Message, syntaxes sigfox.Syntax) []models.QuantitativeValue {
	/*rawData := message.Data
	decodedValues := []schema_org.QuantitativeValue{}

	for index, syntax := range syntaxes.Values {
		switch syntax.Type {
		case "byte":
		case "int":
		case "float":
		case "string":
		case "ssid":
		default:
			fmt.Println("Unknown Syntax")
		}
	}*/
	return nil
}
