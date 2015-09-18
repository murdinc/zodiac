package controllers

import (
	"github.com/revel/revel"
	"zodiac/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	cipher := models.NewCipher(revel.Config.StringDefault("cipher.340.raw", ""), revel.Config.IntDefault("cipher.340.cols", 10))
	return c.Render(cipher)
}

func (c App) Z408() revel.Result {
	cipher := models.NewCipher(revel.Config.StringDefault("cipher.408.raw", ""), revel.Config.IntDefault("cipher.408.cols", 10))
	return c.Render(cipher)
}

func (c App) Z340() revel.Result {
	cipher := models.NewCipher(revel.Config.StringDefault("cipher.340.raw", ""), revel.Config.IntDefault("cipher.340.cols", 10))
	return c.Render(cipher)
}
