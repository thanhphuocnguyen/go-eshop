// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	AddProductToCart(ctx context.Context, arg AddProductToCartParams) (CartItem, error)
	ArchiveProduct(ctx context.Context, id int64) error
	ClearCart(ctx context.Context, cartID int32) error
	CountCartItem(ctx context.Context, cartID int32) (int64, error)
	CountProducts(ctx context.Context, arg CountProductsParams) (int64, error)
	CreateAddress(ctx context.Context, arg CreateAddressParams) (UserAddress, error)
	CreateCart(ctx context.Context, userID int64) (Cart, error)
	CreateCollection(ctx context.Context, arg CreateCollectionParams) (Category, error)
	CreateImage(ctx context.Context, arg CreateImageParams) (Image, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error)
	CreatePaymentTransaction(ctx context.Context, arg CreatePaymentTransactionParams) (Payment, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	DeleteAddress(ctx context.Context, arg DeleteAddressParams) error
	DeleteImage(ctx context.Context, id int32) error
	DeleteOrder(ctx context.Context, id int64) error
	DeletePaymentTransaction(ctx context.Context, id int32) error
	DeleteProduct(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetAddress(ctx context.Context, arg GetAddressParams) (UserAddress, error)
	GetAddresses(ctx context.Context, userID int64) ([]UserAddress, error)
	GetCart(ctx context.Context, userID int64) (Cart, error)
	GetCartItem(ctx context.Context, id int32) (CartItem, error)
	GetCartItemByProductID(ctx context.Context, productID int64) (CartItem, error)
	GetCartItems(ctx context.Context, cartID int32) ([]GetCartItemsRow, error)
	GetCollection(ctx context.Context, id int32) (Category, error)
	GetImageByExternalID(ctx context.Context, externalID pgtype.Text) (Image, error)
	GetImageByID(ctx context.Context, id int32) (Image, error)
	GetImagesByProductID(ctx context.Context, productID pgtype.Int8) ([]Image, error)
	GetImagesByVariantID(ctx context.Context, variantID pgtype.Int8) ([]Image, error)
	GetOrder(ctx context.Context, id int64) (Order, error)
	GetOrderDetails(ctx context.Context, id int64) ([]GetOrderDetailsRow, error)
	GetPaymentTransactionByID(ctx context.Context, id int32) (Payment, error)
	GetPaymentTransactionByOrderID(ctx context.Context, orderID int64) (Payment, error)
	GetPrimaryAddress(ctx context.Context, userID int64) (UserAddress, error)
	GetPrimaryImageByProductID(ctx context.Context, productID pgtype.Int8) (Image, error)
	GetPrimaryImageByVariantID(ctx context.Context, variantID pgtype.Int8) (Image, error)
	GetProduct(ctx context.Context, arg GetProductParams) (Product, error)
	GetProductDetail(ctx context.Context, arg GetProductDetailParams) ([]GetProductDetailRow, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (Session, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListOrderItems(ctx context.Context, arg ListOrderItemsParams) ([]OrderItem, error)
	ListOrders(ctx context.Context, arg ListOrdersParams) ([]ListOrdersRow, error)
	ListProducts(ctx context.Context, arg ListProductsParams) ([]ListProductsRow, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	RemoveProductFromCart(ctx context.Context, arg RemoveProductFromCartParams) error
	ResetPrimaryAddress(ctx context.Context, userID int64) error
	SetPrimaryAddress(ctx context.Context, arg SetPrimaryAddressParams) error
	SetPrimaryImage(ctx context.Context, id int32) error
	UnsetPrimaryImage(ctx context.Context, productID pgtype.Int8) error
	UpdateAddress(ctx context.Context, arg UpdateAddressParams) (UserAddress, error)
	UpdateCart(ctx context.Context, id int32) error
	UpdateCartItemQuantity(ctx context.Context, arg UpdateCartItemQuantityParams) error
	UpdateImage(ctx context.Context, arg UpdateImageParams) error
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	UpdatePaymentTransaction(ctx context.Context, arg UpdatePaymentTransactionParams) error
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error)
}

var _ Querier = (*Queries)(nil)
