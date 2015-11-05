package app

import (
	"time"

	"github.com/murdinc/zodiac/app/controllers"
	"github.com/revel/revel"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	// revel.OnAppStart(InitDB())
	// revel.OnAppStart(FillCache())
	revel.OnAppStart(func() {
		//jobs.In(time.Second, GenerateKeys{})
	})
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

type GenerateKeys struct{}

func (g GenerateKeys) Run() {

	revel.INFO.Print("Generating Random Keys...")

	minKill := revel.Config.IntDefault("app.killcount.min", 1)
	minWord := revel.Config.IntDefault("app.wordcount.min", 40)
	maxSymbols := revel.Config.IntDefault("cipher.maxSymbols", 4)

	wordList, err := controllers.GetWordList(revel.AppPath + "/words/")
	if err != nil {
		return
	}

	for p := 0; p < 60; p++ {
		go func(wordList controllers.WordList, minKill int, minWord int, maxSymbols int) {

			es := controllers.NewIndex(revel.Config.StringDefault("elasticsearch.host", "localhost"))

			for {

				cipher, err := controllers.BuildCipher("Z408")
				if err != nil {
					return
				}

				key, err := cipher.RandomKey(cipher.Symbols, maxSymbols)

				if err != nil {

					cipher.SetKeyFromKeyDoc(key)

					// Kill Count over Min, so scan for words!
					if cipher.KillCount >= minKill {
						cipher.ScanForWords(&wordList)

						// Word count over Min, so store in ES
						if cipher.FoundWordsTotal >= minWord {
							revel.INFO.Printf("Key: [%s] - Kill Count: [%d] - Word Count: [%d]", cipher.KeyID, cipher.KillCount, cipher.FoundWordsTotal)
							es.IndexKey(cipher)
						}
					}
				}

			}

		}(*wordList, minKill, minWord, maxSymbols)
	}

	// Key Count Avg
	cipher, err := controllers.BuildCipher("Z408")
	if err != nil {
		return
	}

	startCount := cipher.KeyCount()
	currCount := 0

	for {
		time.Sleep(time.Minute)
		currCount = cipher.KeyCount()
		kpm := (currCount - startCount)

		startCount = currCount
		revel.INFO.Printf("[%d] Keys per minute", kpm)
	}

}
