package handler

//
//func TestHandler_UpdateListJSONHandler(t *testing.T) {
//	t.Run("1. Should return error when json is invalid", func(t *testing.T) {
//		h := Handler{Mux: chi.NewMux()}
//		h.Post("/updates/", h.UpdateListJSONHandler())
//
//		requestBody := `[{"ID":"TestMetric","MType":"Counter","Delta":10}]`
//		request := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBufferString(requestBody))
//
//		w := httptest.NewRecorder()
//
//		h.ServeHTTP(w, request)
//		res := w.Result()
//		defer res.Body.Close()
//
//		require.Equal(t, res.StatusCode, http.StatusBadRequest)
//	})
//}
