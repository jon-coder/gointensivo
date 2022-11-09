package database

import (
	"database/sql"
	"testing"

	"github.com/jon-coder/gointensivo/internal/order/entity"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuit struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuit) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuit) TearDownTest() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuit))
}

func (suite *OrderRepositoryTestSuit) TestGivenAnOrder_WhenSave_ThenShouldReturnSaveOrder() {
	order, err := entity.NewOrder("123", 12.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}
