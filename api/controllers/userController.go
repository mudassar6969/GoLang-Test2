package controllers

import (
	"assignment2/api/models"
	"assignment2/api/responses"
	"assignment2/api/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (a *App) RegisterUser(w http.ResponseWriter, request *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Registered successfully"}

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	user := &models.User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	err = user.ValidateUser("")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	savedUser, _ := user.GetUser(a.DB)

	if savedUser != nil {
		resp["status"] = "failed"
		resp["message"] = "User already registered, please login"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = user.BeforeSave()
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userCreated, err := user.SaveUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["user"] = userCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) LoginUser(w http.ResponseWriter, request *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Logged in"}
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	user := &models.User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	err = user.ValidateUser("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	savedUser, _ := user.GetUser(a.DB)

	if savedUser == nil {
		resp["status"] = "failed"
		resp["message"] = "Login failed, please signup"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = models.CheckPasswordHash(user.Password, savedUser.Password)

	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "Password incorrect"
		responses.JSON(w, http.StatusForbidden, resp)
		return
	}

	token, err := utils.EncodeAuthToken(savedUser.ID)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)
	return
}
