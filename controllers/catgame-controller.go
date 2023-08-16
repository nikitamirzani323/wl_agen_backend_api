package controllers

import (
	"fmt"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nikitamirzani323/wl_agen_backend_api/models"
)

const Fieldcategame_home_redis = "LISTCATEGAME_BACKEND"
const Fieldcategame_home_client_redis = "LISTCATEGAME_FRONTEND"

func CateGamehome(c *fiber.Ctx) error {
	var obj entities.Model_categame
	var arraobj []entities.Model_categame
	var objprovider entities.Model_providershare
	var arraobjprovider []entities.Model_providershare
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcategame_home_redis)
	jsonredis := []byte(resultredis)
	listprovider_RD, _, _, _ := jsonparser.Get(jsonredis, "listprovider")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		categame_id, _ := jsonparser.GetString(value, "categame_id")
		categame_name, _ := jsonparser.GetString(value, "categame_name")
		categame_status, _ := jsonparser.GetString(value, "categame_status")
		categame_status_css, _ := jsonparser.GetString(value, "categame_status_css")
		categame_create, _ := jsonparser.GetString(value, "categame_create")
		categame_update, _ := jsonparser.GetString(value, "categame_update")

		var objgame entities.Model_game
		var arraobjgame []entities.Model_game
		record_game_RD, _, _, _ := jsonparser.Get(value, "categame_list")
		jsonparser.ArrayEach(record_game_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			game_id, _ := jsonparser.GetString(value, "game_id")
			game_idcategame, _ := jsonparser.GetString(value, "game_idcategame")
			game_idprovider, _ := jsonparser.GetInt(value, "game_idprovider")
			game_nmprovider, _ := jsonparser.GetString(value, "game_nmprovider")
			game_name, _ := jsonparser.GetString(value, "game_name")
			game_img, _ := jsonparser.GetString(value, "game_img")
			game_multiplier, _ := jsonparser.GetInt(value, "game_multiplier")
			game_urlstaging, _ := jsonparser.GetString(value, "game_urlstaging")
			game_urlproduction, _ := jsonparser.GetString(value, "game_urlproduction")
			game_status, _ := jsonparser.GetString(value, "game_status")
			game_status_css, _ := jsonparser.GetString(value, "game_status_css")
			game_create, _ := jsonparser.GetString(value, "game_create")
			game_update, _ := jsonparser.GetString(value, "game_update")

			objgame.Game_id = game_id
			objgame.Game_idcategame = game_idcategame
			objgame.Game_idprovider = int(game_idprovider)
			objgame.Game_nmprovider = game_nmprovider
			objgame.Game_name = game_name
			objgame.Game_img = game_img
			objgame.Game_multiplier = int(game_multiplier)
			objgame.Game_urlstaging = game_urlstaging
			objgame.Game_urlproduction = game_urlproduction
			objgame.Game_status = game_status
			objgame.Game_status_css = game_status_css
			objgame.Game_create = game_create
			objgame.Game_update = game_update
			arraobjgame = append(arraobjgame, objgame)
		})

		obj.Categame_id = categame_id
		obj.Categame_name = categame_name
		obj.Categame_list = arraobjgame
		obj.Categame_status = categame_status
		obj.Categame_status_css = categame_status_css
		obj.Categame_create = categame_create
		obj.Categame_update = categame_update
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listprovider_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		provider_id, _ := jsonparser.GetInt(value, "provider_id")
		provider_name, _ := jsonparser.GetString(value, "provider_name")

		objprovider.Provider_id = int(provider_id)
		objprovider.Provider_name = provider_name
		arraobjprovider = append(arraobjprovider, objprovider)
	})
	if !flag {
		result, err := models.Fetch_categameHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcategame_home_redis, result, 60*time.Minute)
		fmt.Println("CATEGAME MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("CATEGAME CACHE")
		return c.JSON(fiber.Map{
			"status":       fiber.StatusOK,
			"message":      "Success",
			"record":       arraobj,
			"listprovider": arraobjprovider,
			"time":         time.Since(render_page).String(),
		})
	}
}
func CateGameSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_categamesave)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, _, client_admin, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	// admin, idrecord, name, status, sData string
	result, err := models.Save_categame(
		client_admin,
		client.Categame_id, client.Categame_name, client.Categame_status, client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_categame()
	return c.JSON(result)
}
func GameSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_gamesave)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, _, client_admin, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	result, err := models.Save_game(
		client_admin,
		client.Game_id, client.Game_idcategame, client.Game_name, client.Game_img,
		client.Game_urlstaging, client.Game_urlproduction, client.Game_status,
		client.Sdata, client.Game_idprovider, client.Game_multiplier)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_categame()
	return c.JSON(result)
}
func _deleteredis_categame() {
	val_master := helpers.DeleteRedis(Fieldcategame_home_redis)
	fmt.Printf("Redis Delete BACKEND CATEGAME : %d", val_master)

	val_client := helpers.DeleteRedis(Fieldcategame_home_client_redis)
	fmt.Printf("Redis Delete CLIENT CATEGAME : %d", val_client)

}
