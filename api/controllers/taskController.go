package controllers

import (
	"assignment2/api/models"
	"assignment2/api/responses"
	"assignment2/api/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) AddTask(w http.ResponseWriter, request *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Task successfully created"}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	task := &models.Task{}
	err = json.Unmarshal(body, &task)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = task.ValidateTask()
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	savedTask, err := task.GetTask(a.DB)
	if savedTask != nil {
		resp["status"] = "failed"
		resp["message"] = "Task already exists, please choose another name"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	userID := request.Context().Value(utils.KEY_USER_ID).(float64)
	task.UserID = uint(userID)
	savedTask, err = task.SaveTask(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["task"] = savedTask
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) GetTasks(w http.ResponseWriter, request *http.Request) {
	userID := request.Context().Value(utils.KEY_USER_ID).(float64)
	tasks, err := models.GetUserTasks(int(userID), a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, tasks)
	return
}

func (a *App) DeleteTask(w http.ResponseWriter, request *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Task deleted successfully"}
	vars := mux.Vars(request)
	userID := int(request.Context().Value(utils.KEY_USER_ID).(float64))

	id, _ := strconv.Atoi(vars["id"])

	_, err := models.GetTaskByID(id, userID, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = models.DeleteTask(id, userID, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, resp)
	return
}

func (a *App) UpdateTask(w http.ResponseWriter, request *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Task updated successfully"}
	vars := mux.Vars(request)
	taskID, _ := strconv.Atoi(vars["id"])
	userID := request.Context().Value(utils.KEY_USER_ID).(float64)
	_, err := models.GetTaskByID(taskID, int(userID), a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	task := &models.Task{}
	err = json.Unmarshal(body, &task)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = task.UpdateTask(taskID, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, resp)
	return
}
