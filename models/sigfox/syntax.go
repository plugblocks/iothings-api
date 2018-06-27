package sigfox

type Syntax struct {
	Values []Value
}

type Value struct {
	Key  string //Name of the data
	Size uint8  //Size to parse
	Type string //Ugly: string for the data type: int, float ...
	Unit string
}

//TODO: finish Sensit Syntax Description
func getSensitSyntax() Syntax {
	var SensitSyntax Syntax

	batteryMsb := Value{"batteryMsb", 1, "byte", ""}
	eventType := Value{"eventType", 2, "string", ""}
	timeframe := Value{"timeframe", 2, "int", "min"}

	SensitSyntax.Values = []Value{batteryMsb, eventType, timeframe}
	return SensitSyntax
}

func getWiFiSyntax() Syntax {
	var WiFiSyntax Syntax

	ssid1 := Value{"ssid1", 12, "ssid", ""}
	ssid2 := Value{"ssid2", 12, "ssid", ""}

	WiFiSyntax.Values = []Value{ssid1, ssid2}
	return WiFiSyntax
}

const SigfoxSyntaxesCollection = "sigfoxSyntaxes"
