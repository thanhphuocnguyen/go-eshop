'use client';
import React, { useCallback, useState } from 'react';
import { CartModel, GenericResponse, UserModel } from '../definitions';
import { useCart } from '../hooks';
import { getCookie } from 'cookies-next';
import { apiFetch } from '../apis/api';
import { PUBLIC_API_PATHS } from '../constants/api';
import { toast } from 'react-toastify';

interface CartContextType {
  cart: CartModel | undefined;
  cartLoading: boolean;
  cartItemsCount: number;
  removeFromCart: (itemId: string) => Promise<void>;
  updateCartItemQuantity: (itemId: string, quantity: number) => Promise<void>;
  clearCart: () => Promise<void>;
  refreshCart: () => void;
}

export const CartContext = React.createContext<CartContextType | undefined>(
  undefined
);

export function CartContextProvider({
  user,
  children,
}: {
  user: UserModel | null;
  children: React.ReactNode;
}) {
  const [isLoading, setIsLoading] = useState(false);

  // Load user from cookies if available
  const { cart, mutateCart, cartLoading } = useCart(user?.id);

  // Remove item from cart
  const removeFromCart = useCallback(
    async (itemId: string) => {
      const user = getCookie('user_id');
      if (!user) return;

      try {
        setIsLoading(true);
        await apiFetch(`${PUBLIC_API_PATHS.CART_ITEM}/${itemId}`, {
          method: 'DELETE',
        });
        await mutateCart();
      } catch (error) {
        console.error('Error removing item from cart', error);
      } finally {
        setIsLoading(false);
      }
    },
    [mutateCart]
  );

  // Update cart item quantity
  const updateCartItemQuantity = async (itemId: string, quantity: number) => {
    if (!user || quantity < 1) return;
    setIsLoading(true);
    const cartItem = cart?.cart_items.find((item) => item.id === itemId);

    if (cartItem) {
      quantity = cartItem.quantity + quantity;
    }

    const { data, error } = await apiFetch<GenericResponse<boolean>>(
      PUBLIC_API_PATHS.CART_ITEM_QUANTITY.replace(':id', itemId),
      {
        method: 'PUT',
        body: {
          quantity: quantity,
        },
      }
    );

    if (error || !data) {
      toast.error(
        <div>
          <p className='text-sm text-gray-700'>Error updating item quantity</p>
          <p className='text-sm text-gray-500'>{JSON.stringify(error)}</p>
        </div>
      );
    }

    if (data) {
      mutateCart(
        (prev) =>
          prev
            ? {
                ...prev,
                cart_items:
                  prev.cart_items.map((item) =>
                    item.id === itemId ? { ...item, quantity: quantity } : item
                  ) ?? [],
              }
            : undefined,
        {
          revalidate: false,
        }
      );
      toast.success(<div>Item quantity updated successfully</div>, {
        closeButton: true,
        autoClose: 2000,
      });
      if (!cartItem) {
        mutateCart();
      }
    }
    setIsLoading(false);
  };

  // Clear cart
  const clearCart = useCallback(async () => {
    try {
      setIsLoading(true);
      await apiFetch(`${PUBLIC_API_PATHS.CART}/clear`, {
        method: 'PUT',
      });
    } catch (error) {
      console.error('Error clearing cart', error);
    } finally {
      setIsLoading(false);
    }
  }, []);

  return (
    <CartContext.Provider
      value={{
        cart,
        cartLoading: cartLoading || isLoading,
        cartItemsCount: cart?.cart_items.length ?? 0,
        removeFromCart,
        updateCartItemQuantity,
        clearCart,
        refreshCart: mutateCart,
      }}
    >
      {children}
    </CartContext.Provider>
  );
}

export const useCartCtx = () => {
  const context = React.useContext(CartContext);
  if (context === undefined) {
    throw new Error('useCartCtx must be used within a CartContextProvider');
  }
  return context;
};
