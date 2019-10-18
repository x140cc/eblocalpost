package main

import (
	"bytes"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"time"
	"os"
	"bufio"
)



var senders map[string]bool = make(map[string]bool)
func readsenders(){
	file, err := os.Open("/etc/pmta/senders2")
	if err != nil {
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strp := scanner.Text()
		senders[strp] = true
	}
}
func main() {

		readsenders()
		file, err := os.Open(os.Args[1])
		if err != nil {
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		i := 0
		for scanner.Scan() {
			i++
			strp := scanner.Text()
			text := strp
			if text == "HTTPSQS_GET_END" {
				println("queue empty")
				time.Sleep(20 * time.Second)
			}
			println(i)
			if strings.Contains(text, ",") {
				segmenttext := strings.Split(text, ",")
				if len(segmenttext) > 4 && senders[segmenttext[4]] == true{
					var rbody []map[string]interface{}
					t := make(map[string]interface{})
					if segmenttext[0] == "d" {
						t["event"] = "delivered"
					} else {

						t["event"] = "bounce"
					}

					t["email"] = segmenttext[5]
					t["Tracking-ID"] = segmenttext[len(segmenttext)-2]
					println(segmenttext[4])
					println(segmenttext[5])


					rbody = append(rbody, t)
					b, _ := json.Marshal(rbody)
					jssyr := string(b)

					var url string = ""
					println(segmenttext[len(segmenttext)-2])
					if strings.Contains(segmenttext[len(segmenttext)-2],"usweb"){

						url = "https://example.com/v1/x/callback/edm"
					}else {

						url = "https://example.com/v1/x/callback/edm"
					}


					fmt.Println("URL:>", url)

					var jsonStr = []byte(jssyr)
					req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
					req.Header.Set("X-Custom-Header", "myvalue")
					req.Header.Set("Content-Type", "application/json")

					client := &http.Client{
					}
					resp, err := client.Do(req)
					if err != nil {

						fmt.Println("post fail")
					}
					defer resp.Body.Close()

					fmt.Println("response Status:", resp.Status)
					fmt.Println("response Headers:", resp.Header)
					body, _ := ioutil.ReadAll(resp.Body)
					fmt.Println("response Body:", string(body))
				}

			}





		}


	}

