package wrapperRequest

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	PathUser          string
	UserType          string
	AtorKey           string
	AdminOrganization bool
}

func getUserDataLogged(c *gin.Context) (*AuthRequest,"", error) {
	result := AuthRequest{
		AtorKey: c.Request.Header.Get("atorKey"),
	}

	userType := c.Request.Header.Get("userType")
	if len(userType) == 0 {
		return nil, "",errors.New("userType not defined")
	}

	strAdminOrg := c.Request.Header.Get("adminOrg")
	if strAdminOrg == "true" {
		result.AdminOrganization = true
	}

	/*
		//atorKey := c.Request.Header.Get("atorKey")
		userType := c.Request.Header.Get("userType")
		//strAdminOrg := c.Request.Header.Get("adminOrg")
		//q.AtorKey = atorKey

		//if len(userType) == 0 {
		//	return errors.New("userType not defined")
		//}

		u, err := authAPI.ParseUserType(userType)
		if err != nil {
			return nil
		}

		if u != authAPI.Admin && len(atorKey) == 0 {
			return errors.New("atorKey not defined")
		}

		//if strAdminOrg == "true" {
		//	q.AdminOrganization = true
		//}

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
	*/
}
