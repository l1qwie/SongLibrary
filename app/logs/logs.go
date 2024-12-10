package logs

import "log"

var (
	caption string
	DEBUG   bool
	INFO    bool
)

func SetDebug() {
	DEBUG = true
	caption = "DEBUG"
}

func StartPoint(point, method string) {
	log.Printf("[INFO] %s (%s) has been launched", point, method)
}

func FieldRequired(field string) {
	if DEBUG {
		log.Printf("[%s] parameter %s is required", caption, field)
	}
}

func FieldsRequired(field1, field2 string) {
	if DEBUG {
		log.Printf("[%s] parameters %s & %s are required", caption, field1, field2)
	}
}

func InputDataIsOK() {
	if DEBUG {
		log.Printf("[%s] input data has been fully checked. everything is OK", caption)
	}
}

func Nothing() {
	if DEBUG {
		log.Printf("[%s] at least something is required", caption)
	}
}
