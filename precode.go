package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

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

	// city := req.URL.Query().Get("city")
	// assert.NotEqual(t, "moscow", city)

	actualStatus := responseRecorder.Code
	assert.Equal(t, 400, actualStatus)

	expectedMsg := "wrong city value"
	assert.Equal(t, expectedMsg, responseRecorder.Body.String())
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", http.NoBody)

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
