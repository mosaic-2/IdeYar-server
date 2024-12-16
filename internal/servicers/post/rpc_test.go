package postImpl

import (
	"context"
	"database/sql"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"

	"github.com/mosaic-2/IdeYar-server/internal/dbutil"
	"github.com/mosaic-2/IdeYar-server/internal/model"
	"github.com/mosaic-2/IdeYar-server/pkg/postServicePb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	args := m.Called(fc)
	return args.Error(0)
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Raw(sql string, values ...interface{}) *gorm.DB {
	args := m.Called(sql, values)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestCreate(t *testing.T) {
	mockDB := new(MockDB)
	ctx := dbutil.WithGormDB(context.Background(), mockDB)

	server := &Server{}
	req := &postServicePb.CreateRequest{}

	mockDB.On("Transaction", mock.Anything).Return(nil)

	resp, err := server.Create(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestGetPost(t *testing.T) {
	mockDB := new(MockDB)
	ctx := dbutil.WithGormDB(context.Background(), mockDB)

	server := &Server{}
	req := &postServicePb.GetPostRequest{Id: 1}

	mockDB.On("Model", mock.Anything).Return(&gorm.DB{})

	post := model.Post{
		ID:          1,
		Title:       "Test Post",
		MinimumFund: decimal.NewFromInt(100),
		UserID:      1,
		FundRaised:  decimal.NewFromInt(20),
	}
	mockDB.On("Take", &post).Return(nil)

	resp, err := server.GetPost(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "Test Post", resp.Title)
}

func TestSearchPost(t *testing.T) {
	mockDB := new(MockDB)
	ctx := dbutil.WithGormDB(context.Background(), mockDB)

	server := &Server{}
	req := &postServicePb.SearchPostRequest{Title: "Search Title", Page: 0}

	mockDB.On("Raw", mock.Anything, mock.Anything, mock.Anything).Return(&gorm.DB{})

	resp, err := server.SearchPost(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestLandingPosts(t *testing.T) {
	mockDB := new(MockDB)
	ctx := dbutil.WithGormDB(context.Background(), mockDB)

	server := &Server{}
	req := &emptypb.Empty{}

	mockDB.On("Raw", mock.Anything).Return(&gorm.DB{})

	resp, err := server.LandingPosts(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
