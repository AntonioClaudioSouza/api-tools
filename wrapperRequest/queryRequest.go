package wrapperRequest

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/goledgerdev/gogreen-api/services/authAPI"
)

type ConfigPage struct {
	BookMark string
	Limit    int32
}

type QueryRequest struct {
	Page ConfigPage
	Data map[string]interface{}
	Auth AuthRequest
}

// ** Replace field in query by new value
func (q *QueryRequest) ReplaceFieldQuery(nameField string, value interface{}) {
	q.Data[nameField] = value
}

// ** Add field in query if not exist
func (q *QueryRequest) AddFieldQuery(nameField string, value interface{}) {

	if len(nameField) > 0 {

		_, hasExist := q.Data[nameField]
		if !hasExist {
			q.Data[nameField] = value
		}
	}
}

// ** Add field with value of atorKey into query
func (q *QueryRequest) AddFieldQueryAtorKey(c *gin.Context, nameField string) error {

	dataField := GetAtorKey(c)
	if len(dataField) == 0 {
		return errors.New("atorKey not defined")
	}

	q.Data[nameField] = dataField
	return nil
}

func (q *QueryRequest) GetConfigPageJson() ([]byte, error) {

	pMap := map[string]interface{}{
		"bookmark": q.Page.BookMark,
		"pageSize": q.Page.Limit,
	}

	pJson, err := json.Marshal(pMap)
	if err != nil {
		return nil, err
	}

	return pJson, nil
}

func ParserResquestData(c *gin.Context) (QueryRequest, error) {

	queryRes := QueryRequest{
		Page: ConfigPage{
			BookMark: "",
			Limit:    0,
		},
		Data: map[string]interface{}{},
	}

	//** Decode if exist
	b64 := c.Query("request")

	if len(b64) == 0 {

		er := queryRes.insertUserTypeRequest(c)
		if er != nil {
			return queryRes, er
		}

		return queryRes, nil
	}

	// ** Decode
	jsonArray, err := base64.RawStdEncoding.DecodeString(b64)
	if err != nil {
		return queryRes, err
	}

	var dataRes map[string]interface{}
	err = json.Unmarshal(jsonArray, &dataRes)
	if err != nil {
		return queryRes, err
	}

	//Clean parameters empty from front
	for k, v := range dataRes {

		switch p := v.(type) {
		case string:
			{
				if len(p) == 0 {
					delete(dataRes, k)
				}
			}
		}
	}

	limit, _ := dataRes["limit"].(float64)
	bookmark := c.Query("bookmark")

	res := QueryRequest{
		Page: ConfigPage{
			BookMark: bookmark,
			Limit:    int32(limit),
		},
		Data: dataRes,
	}

	err = res.insertUserTypeRequest(c)
	if err != nil {
		return queryRes, err
	}

	return res, nil
}

func (q *QueryRequest) insertUserTypeRequest(c *gin.Context) error {

	atorKey := c.Request.Header.Get("atorKey")
	userType := c.Request.Header.Get("userType")
	strAdminOrg := c.Request.Header.Get("adminOrg")
	q.AtorKey = atorKey

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
		q.AdminOrganization = true
	}

	switch u {
	case authAPI.Admin:
		q.PathUser = AdminPath
		q.UserType = u

	case authAPI.Certificador:
		q.PathUser = CertificadorPath
		q.Data["certifierKey"] = atorKey
		q.UserType = u

	case authAPI.Participante:
		q.PathUser = ParticipantePath
		q.Data["participantKey"] = atorKey
		q.UserType = u

	case authAPI.Registrante:
		q.PathUser = RegistrantePath
		q.Data["registrantKey"] = atorKey
		q.UserType = u
	}
	return nil
}
