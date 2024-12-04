// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddProductToCart(ctx context.Context, arg AddProductToCartParams) (CartItem, error)
	ArchiveProduct(ctx context.Context, id int64) error
	CreateCart(ctx context.Context, userID int64) (Cart, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteProduct(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetCartByUserID(ctx context.Context, userID int64) ([]GetCartByUserIDRow, error)
	GetProduct(ctx context.Context, id int64) (Product, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	RemoveProductFromCart(ctx context.Context, arg RemoveProductFromCartParams) error
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateProductImage(ctx context.Context, arg UpdateProductImageParams) error
	UpdateProductQuantity(ctx context.Context, arg UpdateProductQuantityParams) error
	UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
