package main

import (
	"crypto/sha256"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"log"
	"fmt"
)

type Header struct {
	Name string `json:"name"`
}

func main()  {
	key := []byte("test")
	h := hmac.New(sha256.New, key)
	json_header, err := json.Marshal(&Header{
		Name:"zhouhao",
	})
	if err != nil{
		log.Printf(err.Error())
	}
	encode_string := base64.URLEncoding.EncodeToString(json_header)
	h.Write([]byte(encode_string))
	secret := base64.StdEncoding.EncodeToString(h.Sum(nil))

	//received header info
	//encode_string2 := strings.Replace(encode_string, "e" , "a", 2)
	//fmt.Println(encode_string)
	//fmt.Println(encode_string2)
	encode_string2 := encode_string
	key2 := []byte("test")
	h2 := hmac.New(sha256.New, key2)
	h2.Write([]byte(encode_string2))
	secret2 := base64.StdEncoding.EncodeToString(h2.Sum(nil))
	if secret == secret2 {
		payloadByte, _ := base64.URLEncoding.DecodeString(encode_string2)
		header := &Header{}
		json.Unmarshal(payloadByte, header)
		fmt.Println(header.Name)
	}

	//json_payload, err3 := base64.URLEncoding.DecodeString(string(h.Write(payload)))
	//if err3 != nil {
	//	log.Printf(err3.Error())
	//}
	//fmt.Println(json_payload)
}
