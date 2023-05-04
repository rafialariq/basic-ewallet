package controller

// import (
// 	"bytes"
// 	"encoding/json"
// 	"final_project_easycash/model"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// type historyUsecaseMock struct {
// 	mock.Mock
// }

// func (h *historyUsecaseMock) HistoryByUser(user model.User) ([]model.Bill, error) {
// 	args := h.Called(&user)
// 	return args.Get(0).([]model.Bill), args.Error(1)
// }

// func (h *historyUsecaseMock) HistoryWithAccountFilter(user model.User, accountTypeId int) ([]model.Bill, error) {
// 	args := h.Called(&user, &accountTypeId)
// 	return args.Get(0).([]model.Bill), args.Error(1)
// }

// func (h *historyUsecaseMock) HistoryWithTypeFilter(user model.User, typeId string) ([]model.Bill, error) {
// 	args := h.Called(&user, &typeId)
// 	return args.Get(0).([]model.Bill), args.Error(1)
// }

// func (h *historyUsecaseMock) HistoryWithAmountFilter(user model.User, moreThan, lessThan float64) ([]model.Bill, error) {
// 	args := h.Called(&user, &moreThan, &lessThan)
// 	return args.Get(0).([]model.Bill), args.Error(1)
// }

// func TestHistoryController_FindAllByUser_Success(t *testing.T) {
// 	// create gin context
// 	w := httptest.NewRecorder()
// 	c, router := gin.CreateTestContext(w)

// 	// create mock historyUsecase
// 	mockHistoryUsecase := new(mocks.HistoryUsecase)

// 	// create expected response
// 	expectedRes := []model.Bill{
// 		{Id: 1, SenderTypeId: 1, SenderId: "082123456789", TypeId: "1", Amount: 80000, DestinationTypeId: 1, DestinationId: "085712345678"},
// 		{Id: 2, SenderTypeId: 1, SenderId: "082123456789", TypeId: "2", Amount: 45000, DestinationTypeId: 2, DestinationId: "7750821758759"},
// 		{Id: 3, SenderTypeId: 1, SenderId: "085712345678", TypeId: "1", Amount: 50000, DestinationTypeId: 1, DestinationId: "082123456789"},
// 	}

// 	// set mock behavior
// 	mockHistoryUsecase.On("HistoryByUser", mock.Anything).Return(expectedRes, nil)

// 	// set gin context request body
// 	requestBody := gin.H{
// 		"phoneNumber": "082123456789",
// 	}
// 	requestJson, _ := json.Marshal(requestBody)
// 	c.Request, _ = http.NewRequest("POST", "/history", bytes.NewBuffer(requestJson))

// 	// perform test
// 	h := &HistoryController{
// 		historyUsecase: mockHistoryUsecase,
// 	}
// 	router.POST("/history", h.FindAllByUser)
// 	router.ServeHTTP(w, c.Request)

// 	// assertions
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response []model.Bill
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedRes, response)

// 	// assert that mock was called with the correct parameter
// 	var user model.User
// 	mockHistoryUsecase.AssertCalled(t, "HistoryByUser", user)
// }
