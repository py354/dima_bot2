package vkapi

import (
	"github.com/json-iterator/go"
	"log"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func CheckError(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}
