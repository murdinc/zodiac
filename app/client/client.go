package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"encoding/json"

	"github.com/murdinc/zodiac/app/controllers"
)

func main() {
	fmt.Print("Generating Keys...\n")

	minKill := 1
	minWord := 35
	maxSymbols := 7

	wordList, err := controllers.GetWordList("../words/")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	jsonKeys := make(chan []byte)

	var bro sync.WaitGroup
	bro.Add(1)

	for p := 0; p < 60; p++ {
		go func(wordList controllers.WordList, minKill int, minWord int, maxSymbols int) {

			for {

				cipher := Cipher340() // unsolved one
				//cipher := Cipher408() // solved one

				key, err := cipher.RandomKey(cipher.Symbols, maxSymbols)
				if err == nil {

					cipher.SetKeyFromKeyDoc(key)

					// Kill Count over Min, so scan for words!
					if cipher.KillCount >= minKill {

						cipher.ScanForWords(&wordList)

						// Word count over Min, so store in ES
						if cipher.FoundWordsTotal >= minWord {
							//fmt.Printf("Key: [%s] - Kill Count: [%d] - Word Count: [%d]\n", cipher.KeyID, cipher.KillCount, cipher.FoundWordsTotal)

							key.FoundWordsTotal = cipher.FoundWordsTotal
							key.FoundWords = cipher.FoundWords
							key.KeyID = controllers.HashKey(cipher.Key)

							json, err := json.Marshal(key)
							if err == nil {
								jsonKeys <- json
							}

						}
					}
				}

			}

		}(*wordList, minKill, minWord, maxSymbols)
	}

	go func() {
		for key := range jsonKeys {
			//fmt.Printf("%s", response)
			sendKey(key)
		}
	}()

	bro.Wait() // lol good luck bro

}

func sendKey(key []byte) {
	url := "http://zodiac.sudoba.sh/cipher/key"

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(key))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("Status: [%s] Response: [%s]\n", resp.Status, body)

}

func Cipher408() *controllers.Cipher {
	return controllers.NewCipher("Z408", "@%P/Z/UB%kOR~pX~BWV+eGYF5@HP*K]qYeMJYvU[k#qTtNQYD06S7/@$BPORAU%fRlqEkvLMZJdr\\pFHVWe!Y*+qGD@K[65qX!09S7RNt]YElO!qGBTQS$BLd/P$B*XqEHMUvRRkcZKqp[6Wq]!0LMr@$BPDR+j~5\\N7eEUHkFZcpOVW[0+tL6lvR5H[@DR&TYr\\de/*XJQAP0M!RUt%L6NVEKH~Gr[]Jk0@!LMlNA6Z7P9UpkA@$BVW\\+VTtOPv~SrlfUe5#D9G%%[MNk6ScE/@%%ZfAP$BVpeXqWq&F$!c+*@A@B%OT0RUc+&dYq&vSqWVZeGYKE&TYA@%$Lt&H]FBX@9XADd\\#L]~q&ed$$5e0PORXQF%GcZ*JTtq&!J[+rBPQW5VEXr@W[5qEHM6~U[k")
}

func Cipher340() *controllers.Cipher {
	return controllers.NewCipher("Z340", "HER>plvVPk|1LTG3dNp+B7$O%DWY.<^Kf6ByIcM+UZGW76L#$HJSpp#vl!^V4pO++RK3&@M+9tjd|0FP+P2k/p!RvFlO-^dCkF>2D7$0+Kq%i3UcXGV.9L|7G3Jfj$O+&NY9+*L@d<M+b+ZR3FBcyA52K-9lUV+vJ+Op#<FBy-U+R/0tE|DYBpbTMKO2<clRJ|^0T0M.+PBF95@Sy$+N|0FBc7i!RlGFNvf030b.cV0t++yBX1^I2@CE>VUZ0-+|c.49BK7Opv.fMqG3RcT+L03C<+FlWB|6L++6WC9WcPOSHT/76p|FkdW<#tB&YOB^-Cc>MDHNpkS9ZO!A|Ki+")
}
