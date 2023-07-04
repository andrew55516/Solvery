package api

import (
	mockdb "Solvery/db/mock"
	db "Solvery/db/sqlc"
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListUserEntriesApi(t *testing.T) {
	testCases := []struct {
		name          string
		emai          string
		pageId        int
		pageSize      int
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			emai:     "test@example.com",
			pageId:   1,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListUserEntriesParams{
					UserEmail: "test@example.com",
					Limit:     5,
					Offset:    0,
				}
				store.EXPECT().ListUserEntries(gomock.Any(), arg).Times(1).Return([]db.Entry{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			emai:     "test@example.com",
			pageId:   1,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListUserEntriesParams{
					UserEmail: "test@example.com",
					Limit:     5,
					Offset:    0,
				}
				store.EXPECT().ListUserEntries(gomock.Any(), arg).Times(1).Return([]db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidPageId",
			emai:     "test@example.com",
			pageId:   0,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListUserEntries(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "InvalidPageSize",
			emai:     "test@example.com",
			pageId:   1,
			pageSize: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListUserEntries(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/entries?email=%s&page_id=%d&page_size=%d", tc.emai, tc.pageId, tc.pageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListAllEntriesApi(t *testing.T) {
	testCases := []struct {
		name          string
		pageId        int
		pageSize      int
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageId:   1,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAllEntriesParams{
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().ListAllEntries(gomock.Any(), arg).Times(1).Return([]db.Entry{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			pageId:   1,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAllEntriesParams{
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().ListAllEntries(gomock.Any(), arg).Times(1).Return([]db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidPageId",
			pageId:   0,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAllEntries(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "InvalidPageSize",
			pageId:   1,
			pageSize: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAllEntries(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/entries?page_id=%d&page_size=%d", tc.pageId, tc.pageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
