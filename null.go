package null

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

func RegisterNullStruct(rb *bsoncodec.RegistryBuilder) *bsoncodec.RegistryBuilder {
	rb.RegisterDecoder(tTime, Time{})
	rb.RegisterEncoder(tTime, Time{})
	rb.RegisterDecoder(tString, String{})
	rb.RegisterEncoder(tString, String{})
	rb.RegisterDecoder(tBoolean, Bool{})
	rb.RegisterEncoder(tBoolean, Bool{})
	rb.RegisterDecoder(tInt, Int{})
	rb.RegisterEncoder(tInt, Int{})
	rb.RegisterDecoder(tFloat, Float{})
	rb.RegisterEncoder(tFloat, Float{})
	return rb
}
