package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var r *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	workflowClient := SetupCadence()
	r = SetupRouter(workflowClient)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func TestApiStartWorkflow(t *testing.T) {
	req, _ := http.NewRequest("POST", "/fibonacci/100", nil)
	response := executeRequest(req)

	assert.Equal(t, http.StatusCreated, response.Code, "Status should be 201 created.")

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Contains(t, m["address"], "fibonacci/polling", "Address should contain fibonacci/polling")
	assert.NotNil(t, m["id"], "Id should not be nil")
}

func TestApiStartWorkflowWithNegativeNumber(t *testing.T) {
	req, _ := http.NewRequest("POST", "/fibonacci/-1", nil)
	response := executeRequest(req)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Status should be 400 Bad Request.")
}
