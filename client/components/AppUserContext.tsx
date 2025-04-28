'use client';

import { createContext, useCallback, useContext, useState } from 'react';
import { apiFetch } from '@/lib/apis/api';
import { API_PATHS } from '@/lib/constants/api';
import { getCookie } from 'cookies-next/client';
import { toast } from 'react-toastify';
import { CartModel, GenericResponse, UserModel } from '@/lib/definitions';
import { useCart } from '@/lib/hooks';
import { useUser } from '@/lib/hooks/useUser';
import { KeyedMutator } from 'swr';

// Define types for cart items and cart

interface AppUserContextType {
  cartLoading: boolean;
  isLoading: boolean;
  cart: CartModel | undefined;
  user: UserModel | undefined;
  mutateUser?: KeyedMutator<UserModel>;
  cartItemsCount: number;
  addToCart: (variantID: string, quantity: number) => Promise<void>;
  removeFromCart: (itemId: string) => Promise<void>;
  updateCartItemQuantity: (itemId: string, quantity: number) => Promise<void>;
  clearCart: () => Promise<void>;
  refreshCart: () => void;
}

const AppUserContext = createContext<AppUserContextType | undefined>(undefined);

export function AppUserProvider({ children }: { children: React.ReactNode }) {
  const [isLoading, setIsLoading] = useState(false);
  const { user, mutateUser } = useUser(getCookie('access_token'));
  // Load user from cookies if available
  const { cart, mutateCart, cartLoading } = useCart(user?.id);

  // Add item to cart
  const addToCart = useCallback(
    async (variantID: string, quantity: number) => {
      const user = getCookie('user_id');
      if (!user) return;

      setIsLoading(true);
      const { data, error } = await apiFetch<GenericResponse<string>>(
        API_PATHS.CART_ITEM.replaceAll(':variant_id', variantID.toString()),
        {
          method: 'PUT',
          body: {
            variant_id: variantID,
            quantity: quantity,
          },
        }
      );

      if (data) {
        await mutateCart();
        toast.success(
          <div>
            <p className='text-sm text-gray-700'>Item added to cart</p>
            <p className='text-sm text-gray-500'>{JSON.stringify(data)}</p>
          </div>
        );
      }
      if (error) {
        toast.error(
          <div>
            <p className='text-sm text-gray-700'>Error adding item to cart</p>
            <p className='text-sm text-gray-500'>{JSON.stringify(error)}</p>
          </div>
        );
      }
    },
    [mutateCart]
  );

  // Remove item from cart
  const removeFromCart = useCallback(
    async (itemId: string) => {
      const user = getCookie('user_id');
      if (!user) return;

      try {
        setIsLoading(true);
        await apiFetch(`${API_PATHS.CART_ITEM}/${itemId}`, {
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
    const { data, error } = await apiFetch<GenericResponse<boolean>>(
      `${API_PATHS.CART_ITEM}/${itemId}/quantity`,
      {
        method: 'PUT',
        body: {
          quantity: quantity,
        },
      }
    );

    if (error) {
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
    }
    setIsLoading(false);
  };

  // Clear cart
  const clearCart = useCallback(async () => {
    try {
      setIsLoading(true);
      await apiFetch(`${API_PATHS.CART}/clear`, {
        method: 'PUT',
      });
    } catch (error) {
      console.error('Error clearing cart', error);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Load cart data when user is available

  const value = {
    cartLoading,
    cart,
    user,
    isLoading,
    cartItemsCount: cart ? cart.cart_items.length : 0,
    addToCart,
    removeFromCart,
    updateCartItemQuantity,
    clearCart,
    refreshCart: () => mutateCart(),
    mutateUser,
  } satisfies AppUserContextType;

  return (
    <AppUserContext.Provider value={value}>{children}</AppUserContext.Provider>
  );
}

export const useAppUser = (): AppUserContextType => {
  const context = useContext(AppUserContext);
  if (context === undefined) {
    throw new Error('useAppUser must be used within an AppUserProvider');
  }
  return context;
};
