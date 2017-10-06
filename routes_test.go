package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	store, err := NewStore("test.db")
	require.NoError(t, err)

	err = store.Clean()
	require.NoError(t, err)

	defer store.Close()

	server := NewServer(false, 8000, store)

	t.Run("GET /api/debts", func(t *testing.T) {
		t.Run("GET non existent debt", func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/api/debts/abc", nil)
			response := httptest.NewRecorder()

			server.routes().ServeHTTP(response, request)
			require.Equal(t, http.StatusNotFound, response.Code)
		})

		t.Run("GET an existent debt", func(t *testing.T) {
			debt := store.NewDebt(nil)
			debt.Data = &DebtData{"a": "abc"}
			err := store.SaveDebt(debt)
			require.NoError(t, err)

			request, _ := http.NewRequest("GET", "/api/debts/"+debt.ID, nil)
			response := httptest.NewRecorder()

			server.routes().ServeHTTP(response, request)
			require.Equal(t, http.StatusOK, response.Code)

			returnedDebtData := &DebtData{}
			content, err := ioutil.ReadAll(response.Body)
			err = json.Unmarshal(content, returnedDebtData)
			require.NoError(t, err)
			require.Equal(t, debt.Data, returnedDebtData)
		})
	})

	t.Run("POST /api/debts", func(t *testing.T) {
		t.Run("create a new debt", func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/api/debts", nil)
			response := httptest.NewRecorder()

			server.routes().ServeHTTP(response, request)
			require.Equal(t, http.StatusCreated, response.Code)

			debt, err := store.GetDebt(response.Header().Get("X-Created-Id"))
			require.NoError(t, err)
			require.NotNil(t, debt)
		})
	})

	t.Run("PUT /api/debts/:id/:key", func(t *testing.T) {
		t.Run("update non existent", func(t *testing.T) {
			request, _ := http.NewRequest("PUT", "/api/debts/123/abc", nil)
			response := httptest.NewRecorder()

			server.routes().ServeHTTP(response, request)
			require.Equal(t, http.StatusNotFound, response.Code)
		})

		t.Run("update with invalid key", func(t *testing.T) {
			debt := store.NewDebt(nil)
			err := store.SaveDebt(debt)
			require.NoError(t, err)

			request, _ := http.NewRequest("PUT", "/api/debts/"+debt.ID+"/invalid", nil)
			response := httptest.NewRecorder()

			server.routes().ServeHTTP(response, request)
			require.Equal(t, http.StatusForbidden, response.Code)
		})

		t.Run("update with invalid payload", func(t *testing.T) {

		})

		t.Run("update valid", func(t *testing.T) {
			debt := store.NewDebt(nil)
			err := store.SaveDebt(debt)
			require.NoError(t, err)

			updatedData := &DebtData{"a": float64(5)}
			updatedSONContent, err := json.Marshal(updatedData)
			require.NoError(t, err)
			contentReader := bytes.NewReader(updatedSONContent)

			url := "/api/debts/" + debt.ID + "/" + debt.Secret

			request, _ := http.NewRequest("PUT", url, contentReader)
			response := httptest.NewRecorder()

			server.routes().ServeHTTP(response, request)
			require.Equal(t, http.StatusOK, response.Code)

			updatedDebt, err := store.GetDebt(debt.ID)
			require.NoError(t, err)
			require.Equal(t, updatedData, updatedDebt.Data)
		})
	})

}
