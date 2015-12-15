/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 Jane Lee <jane@webconn.me>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"github.com/webconnme/go-webconn"
	"github.com/webconnme/go-webconn-gpio"
	"log"
	"encoding/json"
	"fmt"
	"strconv"
)

var client webconn.Webconn
var g [8](*gpio.Gpio)

type channel struct {
	Channel string
	Data string
}

func D8_OUT(buf []byte) error{

	data := string(buf)
	log.Println(">>>out data : ",data)

	var gout channel
	err := json.Unmarshal(buf,&gout)
	if err != nil {
		log.Println("json.Unmarshal():", err)
		return err
	}


	ndx ,_ := strconv.Atoi(gout.Channel)
	if ndx >= 8 {
		return nil
	}

	if gout.Data == "1" {
		if err := g[ndx].Out(gpio.HIGH); err != nil {
			log.Println(err)
				return err
		}

	} else if gout.Data == "0" {
		if err := g[ndx].Out(gpio.LOW); err != nil {
			log.Println(err)
			return err
		}
	}
	fmt.Println(">>>data out : ",gout)

	return nil
}

func main() {

	ch := 248
	for i:=0;i<8;i++ {

		g[i] = &gpio.Gpio{ch, gpio.OUT}
		err := g[i].Open()
		if err != nil {
			log.Println(err)
		}
		defer g[i].Close()
		ch = ch+1
	}

	client = webconn.NewClient("http://192.168.4.180:3006/v01/do8ch/80")
	client.AddHandler("do",D8_OUT)

	client.Run()
}
