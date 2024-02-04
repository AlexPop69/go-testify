package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandlerWhenRequestIsOkAndBodyIsNotEmpty(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", http.NoBody)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	actualStatus := responseRecorder.Code
	require.Equal(t, 200, actualStatus)
	require.NotEmpty(t, responseRecorder.Body)
}

// Город, передаваемый, в параметре city, не поддерживается.
// Сервис возвращает код ответа 400 и ошибку "wrong city value" в теле ответа.
func TestMainHandlerWhenCityIsNotOk(t *testing.T) {
	wrongCity := "UnExistsCity"
	req := httptest.NewRequest("GET", "/cafe?count=2&city="+wrongCity, http.NoBody)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	actualStatus := responseRecorder.Code
	assert.Equal(t, 400, actualStatus)

	expectedMsg := "wrong city value"
	assert.Equal(t, expectedMsg, responseRecorder.Body.String())
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=7&city=moscow", http.NoBody)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	city := "moscow"
	allCafes := strings.Join(cafeList[city], ",")
	assert.Equal(t, allCafes, responseRecorder.Body.String())

	totalCount := len(cafeList[city])
	cafesFromResponse := strings.Split(responseRecorder.Body.String(), ",")
	actualCount := len(cafesFromResponse)
	assert.Equal(t, totalCount, actualCount)

}
