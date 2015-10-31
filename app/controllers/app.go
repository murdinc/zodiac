package controllers

import (
	"errors"
	"strings"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

type CipherLink struct {
	Name string
	URL  string
}

func (c App) Index() revel.Result {

	cipherLinks := make([]CipherLink, 0)

	config := revel.Config.Options("cipher.")

	for _, str := range config {

		if strings.HasSuffix(str, "Name") {
			revel.INFO.Print(str)
			name := revel.Config.StringDefault(str, "")
			url := "/cipher/" + name
			cipherLinks = append(cipherLinks, CipherLink{Name: name, URL: url})
		}

	}

	return c.Render(cipherLinks)
}

func (c App) Display(sort string) revel.Result {

	cipher, err := c.buildCipher()
	if err != nil {
		return c.RenderError(err)
	}

	key := KeyDoc{}

	switch sort {
	case "newest":
		revel.INFO.Print("Request to Display Latest Key for Cipher: " + cipher.Name)

		key, err = GetNewestKey(cipher)
		if err != nil {
			return c.RenderError(err)
		}
		cipher.SetKeyFromKeyDoc(key)
	case "best":
		revel.INFO.Print("Request to Display Best Key for Cipher: " + cipher.Name)
		key, err = GetBestKey(cipher)
		if err != nil {
			return c.RenderError(err)
		}
		cipher.SetKeyFromKeyDoc(key)
	case "generate":
		revel.INFO.Print("Request to Generate New Key and Display for Cipher: " + cipher.Name)
		cipher.SetKeyFromKeyMap(cipher.RandomKey(cipher.Symbols))

		CreateIndex()
		IndexKey(cipher)
	case "hash":
		hash := c.Params.Get("hash")
		revel.INFO.Print("Request to Display Key by hash for Cipher: " + cipher.Name + " Key: " + hash)

		key, err = GetKeyByHash(hash)
		if err != nil {
			return c.RenderError(err)
		}

		cipher.SetKeyFromKeyDoc(key)
	default:
		// default TODO
	}

	return c.Render(cipher)

}

func (c App) buildCipher() (*Cipher, error) {
	cipherName := c.Params.Get("cipherName")

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

	indexName := c.Params.Get("indexName")

	revel.INFO.Print("Request delete Index: " + indexName)

	resp, err := DeleteIndex(indexName)
	if err != nil {
		return c.RenderJson(err)
	}

	return c.RenderJson(resp)
}
