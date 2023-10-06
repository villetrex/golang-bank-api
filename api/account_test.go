package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"tutorial.sqlc.dev/app/util"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		acccountID    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			acccountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.Expect().GetAccount(gmock.Any(), gmock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code) // store status code in the code field of the recorder
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:       "NotFound",
			acccountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.Expect().GetAccount(gmock.Any(), gmock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code) // store status code in the code field of the recorder
			},
		},
		{
			name:       "InternalError",
			acccountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.Expect().GetAccount(gmock.Any(), gmock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code) // store status code in the code field of the recorder
			},
		},
		{
			name:       "InvalidId",
			acccountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.Expect().GetAccount(gmock.Any(), gmock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code) // store status code in the code field of the recorder
			},
		},
	}

	// make test to work for multiple scenerios
	for i := range testcases {

		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gmock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			// build stubs for the mock store
			// store.Expect().GetAccount(gmock.Any(), gmock.Eq(account.ID)).Times(1).Return(account, nil)

			//start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/accounts/%d", tc.acccountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
			// check response
			// require.Equal(t, http.StatusOK, recorder.Code) // store status code in the code field of the recorder
			// requireBodyMatchAccount(t, recorder.Body, account)
		})

	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing, body *bytes.Buffer, account db.Account) {
	data, err := uoutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)

}
