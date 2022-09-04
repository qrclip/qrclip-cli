package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/golang-jwt/jwt"
	"github.com/mdp/qrterminal"
	"golang.org/x/term"
)

var gDisplayUserData = true // TO DISPLAY THE USER DATA JUST ONCE

// Login ///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Login(pUsername string, pPassword string) {
	ShowInfo("LOGIN, PLEASE WAIT...")
	if pUsername == "" {
		loginQRCode()
	} else {
		loginWithUsernamePassword(pUsername, pPassword)
	}
}

// Logout //////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Logout() {
	ShowSuccess("LOGOUT")
	ShowSuccess("Cleaning credentials stored at:")
	ShowSuccess(getQRClipConfigFilePath())

	// CLEAR THE CREDENTIALS AND SAVE - FILE IS NEVER DELETED
	var tLogInResponseDto LogInResponseDto
	tLogInResponseDto.AccessToken = ""
	tLogInResponseDto.RefreshToken = ""

	handleLogInResponse(tLogInResponseDto)
}

// loginQRCode /////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func loginQRCode() {
	tLoginApprovalDto := getLoginApproval()

	displayLoginQRCode(tLoginApprovalDto)

	// ASK USER TO APPROVE IN THE OTHER DEVICE
userInputGOTO: // SECOND GOTO I EVER USED :D WHY NOT ?! THE LANGUAGE STARTS WITH GO
	ShowInfoYellow("After approved in the other device,")
	ShowInfoYellow("press any key to continue, or \"a\" to abort...")
	GetUserInput() // HERE IF USER WANTS ABORTS AND PROGRAM EXITS

	// CHECK IF LOGIN APPROVED
	tLogInResponseDto := qrLogin(tLoginApprovalDto)

	// IF LOGIN NOT APPROVED GOTO THE USER INPUT OR IN CASE ITS APPROVED HANDLE IT
	if tLogInResponseDto.Error != "" {
		ShowError("Login not approved yet!")
		goto userInputGOTO // GOTO THE TOP AGAIN UNTIL USER ABORTS
	} else {
		handleLogInResponse(tLogInResponseDto)
	}
}

// getLoginApproval ////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getLoginApproval() LoginApprovalDto {
	tResponse, tErr := http.Get(gApiUrl + "/users/login-approval")
	if tErr != nil {
		ExitWithError("Checking if login approved")
	}
	defer tResponse.Body.Close()

	var tLoginApprovalDto LoginApprovalDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tLoginApprovalDto)
	if tErr != nil {
		ExitWithError("Decoding login approved response")
	}

	return tLoginApprovalDto
}

// displayLoginQRCode ////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func displayLoginQRCode(tLoginApprovalDto LoginApprovalDto) {
	tConfig := GetQRCodeTerminalConfig()

	tData, tErr := json.Marshal(tLoginApprovalDto)
	if tErr != nil {
		ExitWithError("Display QR Login, " + tErr.Error())
		return
	}

	qrterminal.GenerateWithConfig(string(tData), tConfig)
}

// qrLogin /////////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func qrLogin(pLoginApprovalDto LoginApprovalDto) LogInResponseDto {
	tErrorPrefix := "QR code Login, "

	tJson, tErr := json.Marshal(pLoginApprovalDto)
	if tErr != nil {
		ExitWithError("Transforming LoginApprovalDto to JSON")
	}

	tUrl := gApiUrl + "/users/qr-login"

	// CREATE REQUEST
	tRequest, tErr := http.NewRequest(http.MethodPut, tUrl, bytes.NewBuffer(tJson))
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	tRequest.Header.Set("Content-Type", "application/json")

	// SEND REQUEST
	tClient := &http.Client{}
	tResponse, tErr := tClient.Do(tRequest)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer tResponse.Body.Close()

	// PARSE RESPONSE
	var tLogInResponseDto LogInResponseDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tLogInResponseDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	return tLogInResponseDto
}

// handleLogInResponse /////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func handleLogInResponse(pLogInResponseDto LogInResponseDto) {
	tConfig, _ := GetQRClipConfig()

	tConfig.AccessToken = pLogInResponseDto.AccessToken
	tConfig.RefreshToken = pLogInResponseDto.RefreshToken

	if tConfig.AccessToken != "" {
		ShowSuccess("LOGIN OK")
	}

	SaveQRClipConfigFile(tConfig)
}

// CheckJwtToken ///////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func CheckJwtToken() string {
	tLogInDto, tErr := getStoredLogin()
	if tErr != nil {
		ShowWarning("No Credentials found, please login to be able to send larger files!")
		return ""
	} else {
		tUserId := getUserIdFromJWT(tLogInDto.AccessToken)
		tAccountType := getAccountTypeFromJWT(tLogInDto.AccessToken)

		// TO DISPLAY ONLY ONE TIME
		if gDisplayUserData {
			ShowSuccess("USING USER ACCOUNT")
			ShowSuccess(" ID: " + tUserId)
			ShowSuccess(" TYPE: " + tAccountType)
			ShowSuccess(" ")
			gDisplayUserData = false
		}

		// CHECK IF JWT IS STILL VALID
		if checkJWTValidity(tLogInDto.AccessToken) {
			return tLogInDto.AccessToken
		}

		// IN CASE IT'S NOT REFRESH IT
		tJwt, tErr := refreshCredentials(tLogInDto)
		if tErr != nil {
			return ""
		} else {
			return tJwt
		}
	}
}

// getStoredLogin //////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getStoredLogin() (LogInDto, error) {
	var tLogInDto LogInDto
	tConfig, tError := GetQRClipConfig()
	if tError != nil {
		return tLogInDto, errors.New("NO LOGIN SAVED")
	}
	tLogInDto.AccessToken = tConfig.AccessToken
	tLogInDto.RefreshToken = tConfig.RefreshToken
	if tLogInDto.RefreshToken == "" {
		return tLogInDto, errors.New("NO REFRESH TOKEN AVAILABLE")
	} else {
		return tLogInDto, nil
	}
}

// refreshCredentials //////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func refreshCredentials(pLogInDto LogInDto) (string, error) {
	tErrorPrefix := "Refresh Token, "
	var tRefreshTokenRequestDto RefreshTokenRequestDto
	tRefreshTokenRequestDto.Token = pLogInDto.RefreshToken
	tRefreshTokenRequestDto.UserId = getUserIdFromJWT(pLogInDto.AccessToken)

	tJson, tErr := json.Marshal(tRefreshTokenRequestDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tResponse, tErr := http.Post(gApiUrl+"/users/refresh-token", "application/json", bytes.NewBuffer(tJson))
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer tResponse.Body.Close()

	if tResponse.StatusCode != 201 {
		return "", errors.New("LOGIN ERROR")
	}

	var tLogInResponseDto LogInResponseDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tLogInResponseDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	if tLogInResponseDto.Error == "" {
		handleLogInResponse(tLogInResponseDto)
		return tLogInResponseDto.AccessToken, nil
	} else {
		return "", errors.New("LOGIN ERROR")
	}

}

// getJWTClaim /////////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getJWTClaim(pJwt string, pKeyName string) (interface{}, error) {
	tToken, _, tErr := new(jwt.Parser).ParseUnverified(pJwt, jwt.MapClaims{})
	if tErr != nil {
		ShowError(tErr.Error())
	}
	if tClaims, ok := tToken.Claims.(jwt.MapClaims); ok {
		if tClaims[pKeyName] != nil {
			tClaim := tClaims[pKeyName]
			return tClaim, nil
		}
	}
	return nil, nil
}

// getUserIdFromJWT ////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getUserIdFromJWT(pJwt string) string {
	tClaim, tErr := getJWTClaim(pJwt, "id")
	if tErr == nil {
		return fmt.Sprint(tClaim)
	}
	return ""
}

// checkJWTValidity ////////////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func checkJWTValidity(pJwt string) bool {
	tClaim, tErr := getJWTClaim(pJwt, "exp")
	if tErr != nil {
		return false
	}

	var tExp time.Time
	switch iat := tClaim.(type) {
	case float64:
		tExp = time.Unix(int64(iat), 0)
	case json.Number:
		tV, _ := iat.Int64()
		tExp = time.Unix(tV, 0)
	}
	tRemainder := time.Until(tExp)

	// IN CASE THE TIME OF VALIDITY IS SMALLER THAN 30 SECONDS RETURN FALSE
	return tRemainder.Seconds() > 30
}

// getAccountTypeFromJWT ///////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getAccountTypeFromJWT(pJwt string) string {
	tClaim, tErr := getJWTClaim(pJwt, "type")
	if tErr != nil {
		return ""
	}

	var tType = fmt.Sprint(tClaim)
	if tType == "1" {
		return "FREE"
	}

	if tType == "2" {
		return "PREMIUM"
	}

	return ""
}

// loginWithUsernamePassword ///////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func loginWithUsernamePassword(pUsername string, pPassword string) {
	if pPassword == "" {
		// GET PASSWORD FROM USER
		color.Set(color.FgYellow)
		fmt.Print("Enter Password: ")
		color.Unset()
		tBytePassword, tErr := term.ReadPassword(int(syscall.Stdin))
		if tErr != nil {
			ExitWithError("Error getting password")
		}
		pPassword = string(tBytePassword)
	}

	tErrorPrefix := "Login username and password, "
	var tLogInRequestDto LogInRequestDto
	tLogInRequestDto.Email = pUsername
	tLogInRequestDto.Password = pPassword

	tJson, tErr := json.Marshal(tLogInRequestDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	tResponse, tErr := http.Post(gApiUrl+"/users/sign-in", "application/json", bytes.NewBuffer(tJson))
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}
	defer tResponse.Body.Close()

	var tLogInResponseDto LogInResponseDto
	tErr = json.NewDecoder(tResponse.Body).Decode(&tLogInResponseDto)
	if tErr != nil {
		ExitWithError(tErrorPrefix + tErr.Error())
	}

	if tLogInResponseDto.Error == "" {
		handleLogInResponse(tLogInResponseDto)
	} else {
		ExitWithError("LOGIN FAILED!")
	}
}
