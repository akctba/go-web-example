package controllers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/starptech/go-web/models"
	"github.com/starptech/go-web/server"
	validator "gopkg.in/go-playground/validator.v9"
)

type Importer struct{}

type UserEntity struct {
	Name string `json:"name" validate:"required"`
}

func (ctrl Importer) ImportUser(e *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(UserEntity)

		if errB := c.Bind(u); errB != nil {
			b := server.Boom{Code: server.InvalidBindingModel, Message: "invalid user model", Details: errB}
			c.Logger().Error(errB)
			return c.JSON(http.StatusBadRequest, b)
		}

		errV := c.Validate(u)

		if errV != nil {
			err := errV.(validator.ValidationErrors)
			c.Logger().Error(err.Error())
			return errors.New(err.Error())
		}

		model := models.User{Name: u.Name}

		if errM := e.GetDB().Create(&model).Error; errM != nil {
			b := server.Boom{Code: server.EntityCreationError, Message: "user could not be created", Details: errM}
			c.Logger().Error(errM)
			return c.JSON(http.StatusBadRequest, b)
		}

		return c.JSON(http.StatusOK, u)
	}
}