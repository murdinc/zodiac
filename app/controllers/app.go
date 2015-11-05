package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

type CipherLink struct {
	Name       string
	ByDateURL  string
	ByCountURL string
}

type Nav struct {
	Visible     bool
	Index       int
	Previous    int
	Next        int
	Total       int
	Description string
}

func (c App) Index() revel.Result {

	cipherLinks := make([]CipherLink, 0)

	config := revel.Config.Options("cipher.")

	for _, str := range config {

		if strings.HasSuffix(str, "name") {
			revel.INFO.Print(str)
			name := revel.Config.StringDefault(str, "")
			revel.INFO.Print(name)
			dateurl := "/cipher/" + name + "/key/date/"
			counturl := "/cipher/" + name + "/key/wordcount/"
			cipherLinks = append(cipherLinks, CipherLink{Name: name, ByDateURL: dateurl, ByCountURL: counturl})
		}

	}

	return c.Render(cipherLinks)
}

func (c App) PutKey() revel.Result {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return c.RenderJson(err)
	}

	key := KeyDoc{}
	json.Unmarshal(body, &key)

	if HashKey(key.Key) == key.KeyID {
		es := NewIndex(revel.Config.StringDefault("elasticsearch.host", "localhost"))
		es.IndexKeyDoc(key)

		revel.INFO.Printf("Got Key [%s] for Cipher [%s] from Host [%s]", key.KeyID, key.CipherName, c.Request.Host)

		return c.RenderJson(true)
	}

	revel.ERROR.Printf("Got Invalid Key [%s] for Cipher [%s] from Host [%s]", key.KeyID, key.CipherName, c.Request.Host)
	return c.RenderJson("Invalid Key Hash!")

}

func (c App) Display(sort string) revel.Result {

	var err error

	es := NewIndex(revel.Config.StringDefault("elasticsearch.host", "localhost"))

	cipherName := c.Params.Get("cipherName")
	maxSymbols := revel.Config.IntDefault("cipher.maxSymbols", 4)

	cipher, err := BuildCipher(cipherName)
	if err != nil {
		return c.RenderError(err)
	}

	// Navigation
	index := 1
	previous, next := 0, 0
	total := cipher.KeyCount()

	c.Params.Bind(&index, "index")
	if index < 1 {
		index++
	}

	offset := index - 1

	if index < total {
		next = index + 1
	}

	if index > 1 {
		previous = index - 1
	}

	nav := Nav{
		Index:    index,
		Previous: previous,
		Next:     next,
		Total:    total,
	}

	key := KeyDoc{}

	switch sort {
	case "date":
		revel.INFO.Printf("Request to Display Latest Key for Cipher: [%s] starting from offset [%d]", cipher.Name, offset)
		key, err = es.GetKeyByDate(cipher, offset)
		revel.INFO.Printf("KEY ID: %s", key.KeyID)
		if err != nil {
			return c.RenderError(err)
		}

		cipher.SetKeyFromKeyDoc(key)
		nav.Visible = true
		nav.Description = "Sorted By Date"

	case "wordcount":
		revel.INFO.Printf("Request to Display Best Key for Cipher: [%s] starting from offset [%d]", cipher.Name, offset)
		key, err = es.GetKeyByWordcount(cipher, offset)
		if err != nil {
			return c.RenderError(err)
		}

		cipher.SetKeyFromKeyDoc(key)

		nav.Visible = true
		nav.Description = "Sorted By Word Count"

	case "generate":
		revel.INFO.Print("Request to Generate New Key and Display for Cipher: " + cipher.Name)

		key, err = cipher.RandomKey(cipher.Symbols, maxSymbols)

		if err == nil {
			cipher.SetKeyFromKeyDoc(key)

			wordList, err := GetWordList(revel.AppPath + "/words/")
			if err == nil {
				// scan for words!
				cipher.ScanForWords(wordList)
			}

			//es.IndexKey(cipher)
		}

		key.FoundWords = cipher.FoundWords
		key.FoundWordsTotal = cipher.FoundWordsTotal

	case "solution":
		revel.INFO.Print("Request to Display Cipher from Z408 Solution: " + cipher.Name)

		key, err = cipher.Z408Solution()

		if err == nil {
			cipher.SetKeyFromKeyDoc(key)

			wordList, err := GetWordList(revel.AppPath + "/words/")
			if err == nil {
				// scan for words!
				cipher.ScanForWords(wordList)
			}
		}

	case "hash":
		hash := c.Params.Get("hash")
		revel.INFO.Print("Request to Display Key by hash for Cipher: " + cipher.Name + " Key: " + hash)

		key, err = es.GetKeyByHash(hash)

		cipher.SetKeyFromKeyDoc(key)
	default:
		// default TODO
	}

	// Raw JSON responses
	if c.Params.Values.Encode() == "application/json" {
		if err != nil {
			return c.RenderJson(err)
		}

		return c.RenderJson(key)
	}

	// Default View
	if err != nil {
		revel.INFO.Print("ERROR")
		return c.RenderError(err)
	}

	return c.Render(cipher, key, nav)

}

func BuildCipher(cipherName string) (*Cipher, error) {

	cipherString := revel.Config.StringDefault("cipher."+cipherName+".raw", "")
	if cipherString == "" {
		return &Cipher{}, errors.New("Invalid Request for Cipher Page: " + cipherName)
	}

	cipher := NewCipher(cipherName, cipherString)

	cipher.SetCols(revel.Config.IntDefault("cipher."+cipherName+".cols", 10))
	cipher.SetSolved(revel.Config.BoolDefault("cipher."+cipherName+".solved", false))

	return cipher, nil

}

func (c App) DeleteCiphersIndex() revel.Result {

	es := NewIndex(revel.Config.StringDefault("elasticsearch.host", "localhost"))

	indexName := c.Params.Get("indexName")

	revel.INFO.Print("Request delete Index: " + indexName)

	resp, err := es.DeleteIndex(indexName)
	if err != nil {
		return c.RenderJson(err)
	}

	return c.RenderJson(resp)
}
