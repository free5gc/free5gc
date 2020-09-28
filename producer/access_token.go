package producer

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/logger"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func HandleAccessTokenRequest(request *http_wrapper.Request) *http_wrapper.Response {
	// Param of AccessTokenRsp
	logger.AccessTokenLog.Infoln("Handle AccessTokenRequest")

	accessTokenReq := request.Body.(models.AccessTokenReq)

	response, errResponse := AccessTokenProcedure(accessTokenReq)

	if response != nil {
		// status code is based on SPEC, and option headers
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if errResponse != nil {
		return http_wrapper.NewResponse(http.StatusBadRequest, nil, errResponse)
	}
	problemDetails := &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func AccessTokenProcedure(request models.AccessTokenReq) (response *models.AccessTokenRsp,
	errResponse *models.AccessTokenErr) {
	logger.AccessTokenLog.Infoln("In AccessTokenProcedure")

	var expiration int32 = 1000
	var scope = request.Scope
	var tokenType = "Bearer"

	// Create AccessToken
	var accessTokenClaims = models.AccessTokenClaims{
		Iss:            "1234567",                  // TODO: NF instance id of the NRF
		Sub:            request.NfInstanceId,       // nfInstanceId of service consumer
		Aud:            request.TargetNfInstanceId, // nfInstanceId of service producer
		Scope:          request.Scope,              // TODO: the name of the NF services for which the
		Exp:            expiration,                 //       access_token is authorized for use
		StandardClaims: jwt.StandardClaims{},
	}

	mySigningKey := []byte("NRF") // AllYourBase
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err := token.SignedString(mySigningKey)

	if err != nil {
		logger.AccessTokenLog.Warnln("Signed string error: ", err)
		errResponse = &models.AccessTokenErr{
			Error: "invalid_request",
		}

		return nil, errResponse
	}

	response = &models.AccessTokenRsp{
		AccessToken: accessToken,
		TokenType:   tokenType,
		ExpiresIn:   expiration,
		Scope:       scope,
	}

	return response, nil

}
