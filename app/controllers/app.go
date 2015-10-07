package controllers

import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {

	cipher := NewCipher(revel.Config.StringDefault("cipher.408.raw", ""))

	cipher.SetCols(revel.Config.IntDefault("cipher.408.cols", 10))
	cipher.SetSolved(revel.Config.BoolDefault("cipher.408.solved", false))
	cipher.SetKey(cipher.Z408Key())

	return c.Render(cipher)
}

func (c App) Z408() revel.Result {

	cipher := NewCipher(revel.Config.StringDefault("cipher.408.raw", ""))

	cipher.SetCols(revel.Config.IntDefault("cipher.408.cols", 10))
	cipher.SetSolved(revel.Config.BoolDefault("cipher.408.solved", false))
	cipher.SetKey(cipher.Z408Key())

	return c.Render(cipher)
}

func (c App) Z340(random bool) revel.Result {

	cipherString := revel.Config.StringDefault("cipher.340.raw", "")

	cipher := NewCipher(cipherString)

	cipher.SetCols(revel.Config.IntDefault("cipher.340.cols", 10))
	cipher.SetSolved(revel.Config.BoolDefault("cipher.340.solved", false))
	cipher.SetKey(cipher.RandomKey(cipher.Symbols))

	return c.Render(cipher)
}
