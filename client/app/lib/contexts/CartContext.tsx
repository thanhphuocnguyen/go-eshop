'use client';
import React, { useCallback, useState } from 'react';
import { CartModel, GenericResponse } from '../definitions';
import { clientSideFetch } from '../api/apiClient';
import { PUBLIC_API_PATHS } from '../constants/api';
import { toast } from 'react-toastify';
import { useCart, useUser } from '@/app/hooks';

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
  children,
}: {
  children: React.ReactNode;
}) {
  const { user } = useUser();
  const [isLoading, setIsLoading] = useState(false);
  // Load user from cookies if available
  const { cart, mutateCart, cartLoading } = useCart(user?.id);

  // Remove item from cart
  const removeFromCart = useCallback(
    async (itemId: string) => {
      if (!user) return;

      try {
        setIsLoading(true);
        await clientSideFetch(`${PUBLIC_API_PATHS.CART_ITEM}/${itemId}`, {
          method: 'DELETE',
        });
        await mutateCart();
      } catch (error) {
        console.error('Error removing item from cart', error);
      } finally {
        setIsLoading(false);
      }
    },
    [mutateCart, user]
  );

  // Update cart item quantity
  const updateCartItemQuantity = async (itemId: string, quantity: number) => {
    if (!user || quantity < 1) return;
    setIsLoading(true);
    const cartItem = cart?.cartItems.find((item) => item.id === itemId);

    if (cartItem) {
      quantity = cartItem.quantity + quantity;
    }

    const { data, error } = await clientSideFetch<GenericResponse<boolean>>(
      PUBLIC_API_PATHS.UPDATE_CART_ITEM_QUANTITY.replace(':id', itemId),
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
                cartItems:
                  prev.cartItems.map((item) =>
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
      await clientSideFetch(`${PUBLIC_API_PATHS.CART}/clear`, {
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
        cartItemsCount: cart?.cartItems.length ?? 0,
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
