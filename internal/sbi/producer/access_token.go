package producer

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	nrf_context "github.com/free5gc/nrf/internal/context"
	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/httpwrapper"
)

func HandleAccessTokenRequest(request *httpwrapper.Request) *httpwrapper.Response {
	// Param of AccessTokenRsp
	logger.AccessTokenLog.Infoln("Handle AccessTokenRequest")

	accessTokenReq := request.Body.(models.AccessTokenReq)

	response, errResponse := AccessTokenProcedure(accessTokenReq)

	if response != nil {
		// status code is based on SPEC, and option headers
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else if errResponse != nil {
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, errResponse)
	}
	problemDetails := &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func AccessTokenProcedure(request models.AccessTokenReq) (response *models.AccessTokenRsp,
	errResponse *models.AccessTokenErr) {
	logger.AccessTokenLog.Infoln("In AccessTokenProcedure")

	var expiration int32 = 1000
	scope := request.Scope
	tokenType := "Bearer"
	now := int32(time.Now().Unix())

	// Create AccessToken
	accessTokenClaims := models.AccessTokenClaims{
		Iss:            nrf_context.Nrf_NfInstanceID, // TODO: NF instance id of the NRF
		Sub:            request.NfInstanceId,         // nfInstanceId of service consumer
		Aud:            request.TargetNfInstanceId,   // nfInstanceId of service producer
		Scope:          request.Scope,                // TODO: the name of the NF services for which the
		Exp:            now + expiration,             // access_token is authorized for use
		StandardClaims: jwt.StandardClaims{},
	}
	accessTokenClaims.IssuedAt = int64(now)

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
