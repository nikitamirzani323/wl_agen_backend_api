package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nikitamirzani323/wl_agen_backend_api/models"
)

const Fieldmember_home_redis = "LISTMEMBER_AGEN"

func Memberhome(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, client_idmasteragen, _, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	var obj entities.Model_member
	var arraobj []entities.Model_member
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldmember_home_redis + "_" + client_idmasteragen)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		member_id, _ := jsonparser.GetString(value, "member_id")
		member_username, _ := jsonparser.GetString(value, "member_username")
		member_timezone, _ := jsonparser.GetString(value, "member_timezone")
		member_ipaddress, _ := jsonparser.GetString(value, "member_ipaddress")
		member_lastlogin, _ := jsonparser.GetString(value, "member_lastlogin")
		member_name, _ := jsonparser.GetString(value, "member_name")
		member_phone, _ := jsonparser.GetString(value, "member_phone")
		member_email, _ := jsonparser.GetString(value, "member_email")
		member_status, _ := jsonparser.GetString(value, "member_status")
		member_status_css, _ := jsonparser.GetString(value, "member_status_css")
		member_create, _ := jsonparser.GetString(value, "member_create")
		member_update, _ := jsonparser.GetString(value, "member_update")

		obj.Member_id = member_id
		obj.Member_username = member_username
		obj.Member_timezone = member_timezone
		obj.Member_ipaddress = member_ipaddress
		obj.Member_lastlogin = member_lastlogin
		obj.Member_name = member_name
		obj.Member_phone = member_phone
		obj.Member_email = member_email
		obj.Member_status = member_status
		obj.Member_status_css = member_status_css
		obj.Member_create = member_create
		obj.Member_update = member_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_memberHome(client_idmasteragen)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldmember_home_redis+"_"+client_idmasteragen, result, 60*time.Minute)
		fmt.Println("MEMBER AGEN MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("MEMBER AGEN CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}

func MemberSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_membersave)
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
	client_idmaster, client_idmasteragen, client_admin, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	result, err := models.Save_member(
		client_admin,
		client_idmaster, client_idmasteragen, client.Member_username, client.Member_password,
		client.Member_name, client.Member_phone, client.Member_email, client.Member_status,
		client.Sdata, client.Member_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_member(client_idmasteragen)
	return c.JSON(result)
}
func _deleteredis_member(idmasteragen string) {
	val_master := helpers.DeleteRedis(Fieldmember_home_redis + "_" + idmasteragen)
	log.Printf("Redis Delete AGEN MEMBER : %d", val_master)

}