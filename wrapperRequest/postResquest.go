package wrapperRequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goledgerdev/gogreen-api/services/authAPI"
)

type FieldForm struct {
	Required bool
	FieldAPI string
	Type     string
	Default  interface{}
}

type PostData struct {
	FieldsForm []FieldForm
	Data       map[string]interface{}
	Auth       AuthRequest
}

func (p *PostData) parser(c *gin.Context, rawMap map[string]interface{}) error {

	p.Data = make(map[string]interface{}, 0)

	for _, f := range p.FieldsForm {

		// ** Check required
		if f.Required {
			_, has := rawMap[f.FieldAPI]
			if !has {
				return fmt.Errorf("%s is required but not found", f.FieldAPI)
			}
		}

		// ** Check type of field
		valueField, has := rawMap[f.FieldAPI]
		if has {
			typeReflect := reflect.TypeOf(valueField)

			if typeReflect == nil {
				return fmt.Errorf("field '%s' cannot null", f.FieldAPI)
			}

			if typeReflect.String() != f.Type {
				return fmt.Errorf("field '%s' not is type '%s' but is type '%s'", f.FieldAPI, f.Type, typeReflect)
			}
			p.Data[f.FieldAPI] = valueField
		} else {
			if f.Default != nil {
				p.Data[f.FieldAPI] = f.Default
			}
		}
	}
	return p.insertUserTypeRequest(c)
}

func (p *PostData) ParserWithMap(c *gin.Context, rawMap map[string]interface{}) error {
	return p.parser(c, rawMap)
}

func (p *PostData) ParserForm(c *gin.Context) error {
	dataForm := map[string]interface{}{}
	value := ""
	has := false

	for _, f := range p.FieldsForm {

		value = c.Request.FormValue(f.FieldAPI)
		if len(value) > 0 {
			has = true
		} else {
			has = false
		}

		// ** Check required
		if f.Required {
			if !has {
				return fmt.Errorf("%s is required but not found", f.FieldAPI)
			}
		}

		var v interface{}
		// ** Convert string to type in definition
		switch f.Type {
		case "float64":
			v, _ = strconv.ParseFloat(value, 64)

		case "string":
			v = value
		default:
			v = value
		}

		dataForm[f.FieldAPI] = v
	}

	return p.parser(c, dataForm)
}

func (p *PostData) Parser(c *gin.Context) error {

	// Get raw data from request
	rawData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("failed read body request")
	}

	// Convert rawData to map of interface
	var rawMap map[string]interface{}
	err = json.Unmarshal(rawData, &rawMap)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON")
	}

	return p.parser(c, rawMap)
}

func (p *PostData) insertUserTypeRequest(c *gin.Context) error {

	atorKey := c.Request.Header.Get("atorKey")
	userType := c.Request.Header.Get("userType")
	strAdminOrg := c.Request.Header.Get("adminOrg")
	p.AtorKey = atorKey

	if len(userType) == 0 {
		return errors.New("userType not defined")
	}

	u, err := authAPI.ParseUserType(userType)
	if err != nil {
		return nil
	}

	if u != authAPI.Admin && len(atorKey) == 0 {
		return errors.New("atorKey not defined")
	}

	if strAdminOrg == "true" {
		p.AdminOrganization = true
	}

	switch u {
	case authAPI.Admin:
		p.PathUser = AdminPath
		p.UserType = u

	case authAPI.Certificador:
		p.PathUser = CertificadorPath
		p.Data["certifierKey"] = atorKey
		p.UserType = u

	case authAPI.Participante:
		p.PathUser = ParticipantePath
		p.Data["participantKey"] = atorKey
		p.UserType = u

	case authAPI.Registrante:
		p.PathUser = RegistrantePath
		p.Data["registrantKey"] = atorKey
		p.UserType = u
	}
	return nil
}
