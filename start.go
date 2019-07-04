package main

import (
	"fmt"
	//"github.com/influxdata/influxdb/client/v2"
	client "github.com/influxdata/influxdb1-client/v2"
	"io/ioutil"
	"net/http"
	"time"

	"strconv"
)

const (
	MyDB     = "YANGWUQI" //数据库名
	username = "yangwuqi" //用户名
	password = "19951129" //密码
)

var dataIndex int

func main() {
	url1 := "http://192.168.213.136:9100/metrics"
	httpGet(dataIndex, url1)
}

func httpGet(dataIndex int, url1 string) {

	dataNumber := 0
	for {
		fmt.Println("每过5秒Get抓取数据！")
		resp, err := http.Get(url1)
		if err != nil {
			fmt.Println("http get error!")
			return
		}
		defer resp.Body.Close()

		body, err1 := ioutil.ReadAll(resp.Body)

		if err1 != nil {
			fmt.Println("http read error!")
			return
		}

		cli, err3 := client.NewHTTPClient(client.HTTPConfig{ //set the connection with the database
			Addr: "http://10.15.165.3:8086",
		})
		if err3 != nil {
			panic(err3)
		}

		batchPoints, _ := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  "HEYHEY" + strconv.Itoa(dataIndex),
			Precision: "us",
		})
		defer cli.Close()

		if dataNumber == 0 {
			query0 := client.NewQuery("CREATE DATABASE HEYHEY"+strconv.Itoa(dataIndex), "", "")
			if response, err4 := cli.Query(query0); err4 == nil && response.Error() == nil {
				fmt.Println("create database success!")
			}
		}

		changeLineNumber := 0
		lineStart := 0
		lineEnd := 0
		for i := 0; i < len(body); i++ {
			if body[i] == '\n' || body[i] == '\r' || i == len(body)-1 {
				changeLineNumber++
				lineEnd = i
				writeLine1, writeline2 := ParseLine(lineStart, lineEnd, body) //get the line and parse, return a string can be added to the database
				lineStart = i + 1
				if len(writeLine1) == 0 || len(writeline2) == 0 {
					continue
				}

				tags := map[string]string{"": ""}
				fields := map[string]interface{}{
					writeLine1: writeline2,
				}
				point, err4 := client.NewPoint(url1, tags, fields, time.Now())
				if err4 != nil {
					panic(err4)
				}
				batchPoints.AddPoint(point)
			}
		}
		err5 := cli.Write(batchPoints)
		if err5 != nil {
			panic(err5)
		}
		fmt.Println(dataIndex, "   ", dataNumber, " data written ok!")
		dataNumber++
		time.Sleep(5000000000)
	}
}

func ParseLine(lineStart, lineEnd int, body []byte) (line1, line2 string) {

	if lineEnd-lineStart == 0 {
		return "", ""
	}

	line := body[lineStart:lineEnd]
	if line[0] == '#' {

	} else {
		for i := 0; i < len(line); i++ {
			if line[i] == ' ' {
				line1 = string(line[0:i])
				line2 = string(line[i+1:])
				//fmt.Println(line1,"   ",line2)
			}
		}
	}
	return line1, line2
}
