package database

// Condition present database query condition
//
//	e.g:
//		Condition{"field", "123"}:                 field = "123"
//		Condition{"field", ">", 100}:              field > 100
//		Condition{"filed1 = 100 OR field2 = 200"}: filed1 = 100 OR field2 = 200
type Condition []any

type Conditions []Condition
