package null

import "reflect"

var (
	tString  = reflect.TypeOf(String{})
	tBoolean = reflect.TypeOf(Bool{})
	tInt     = reflect.TypeOf(Int{})
	tFloat   = reflect.TypeOf(Float{})
	tTime    = reflect.TypeOf(Time{})
)
