package controllers

import (
    "github.com/revel/revel"
    "zodiac/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
    cipher := models.Cipher340
    return c.Render(cipher)
}

func (c App) Z408() revel.Result {
    cipher := models.Cipher408
    return c.Render(cipher)
}

func (c App) Z340() revel.Result {
    cipher := models.Cipher340
    return c.Render(cipher)
}
