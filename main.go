package main

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"net/url"
	"strings"
	"time"
	"fmt"
)

func GetQ() string {

	g, _ := url.Parse("http://127.0.0.1:8018/")
	gq := g.Query()
	gq.Set("auth", "ebkhttpsqs2016")
	gq.Set("name", "local")
	gq.Set("opt", "get")
	g.RawQuery = gq.Encode()
	//	res, err := http.Get(u.String())
//	fmt.Println(g.String())
	r := gorequest.New()
	_, body, errs := r.Get(g.String()).
		Type("multipart").
		Retry(5, 10 * time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		End()
	if errs != nil{
		return fmt.Sprint(errs)
	}

	return body

}

func main() {

	for{
	//	text := `d,2019-01-23 02:49:15-0500,ekaterina.vaulina@eventbank.com,<EI_EDM:711:14151:b4fb79f5-ef1b-47a9-8f58-be851399051f:JavaMail.ebprod@ebcproweb01>,no-reply@eventbankmail.com,eb-share1,com,,"smtp;250 2.6.0 <EI_EDM:711:14151:b4fb79f5-ef1b-47a9-8f58-be851399051f:JavaMail.ebprod@ebcproweb01> [InternalId=11686606012594, Hostname=BN7PR08MB5522.namprd08.prod.outlook.com] 70483 bytes in 1.009, 68.151 KB/sec Queued mail for delivery",,,,,,,`
		text := GetQ()
		//fmt.Println(text)

		if strings.Contains(text,"GET_END"){
			time.Sleep(6 * time.Second)
			continue
		}
		Consume(text)
	}

}


func Consume(text string) {
	if strings.Contains(text, ",") {
		segmenttext := strings.Split(text, ",")
		if len(segmenttext) > 4 {
			var rbody []map[string]interface{}
			t := make(map[string]interface{})
			if segmenttext[0] == "d" {
				t["event"] = "delivered"
			} else if  strings.Contains(segmenttext[7],"hard") {

				t["event"] = "dropped"
			}else {

				t["event"] = "bounce"
			}

			t["email"] = segmenttext[2]

			t["Tracking-ID"] = segmenttext[3]

			rbody = append(rbody, t)
			b, _ := json.Marshal(rbody)
			jssyr := string(b)

			var burl string

			if strings.Contains(segmenttext[6],"cn"){
				burl = "https://api.eventbank.cn/v1/sendgrid/callback/edm"

			}else if strings.Contains(segmenttext[6],"com"){
				burl = "https://api.eventbank.com/v1/sendgrid/callback/edm"

			}else {

				burl = "https://api.eventbank.com/v1/sendgrid/callback/edm"
			}

			request := gorequest.New()
			_, body, _ := request.Post(burl).
				Retry(5, 10 * time.Second, http.StatusBadRequest, http.StatusInternalServerError).
				Send(jssyr).
				End()

	//		var jsonStr = []byte(jssyr)
			fmt.Println(body)







			//----- process - end


		}

	}






}
