'use client';

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from 'react';
import { apiFetch } from '@/lib/api/api';
import { API_PATHS } from '@/lib/constants/api';
import { getCookie } from 'cookies-next';
import { toast } from 'react-toastify';
import { GenericResponse } from '@/lib/definitions';

// Define types for cart items and cart
interface CartItem {
  id: string;
  product_id: string;
  variant_id: string;
  name: string;
  quantity: number;
  price: number;
  discount: number;
  stock: number;
  sku: string;
  image_url?: string;
  attributes: Array<{
    name: string;
    value: string;
  }>;
}

interface Cart {
  id: string;
  user_id: string;
  total_price: number;
  cart_items: CartItem[];
  updated_at: string;
  created_at: string;
}

interface AppUserContextType {
  isLoading: boolean;
  cart: Cart | null;
  cartItemsCount: number;
  addToCart: (variantID: string, quantity: number) => Promise<void>;
  removeFromCart: (itemId: string) => Promise<void>;
  updateCartItemQuantity: (itemId: string, quantity: number) => Promise<void>;
  clearCart: () => Promise<void>;
  refreshCart: () => Promise<void>;
}

const AppUserContext = createContext<AppUserContextType | undefined>(undefined);

export function AppUserProvider({ children }: { children: React.ReactNode }) {
  const [cart, setCart] = useState<Cart | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [cartItemsCount, setCartItemsCount] = useState(0);

  // Load user from cookies if available

  // Fetch cart data when user changes or when required
  const fetchCart = useCallback(async () => {
    const user = getCookie('user_id');
    if (!user) {
      setCart(null);
      setCartItemsCount(0);
      return;
    }

    try {
      setIsLoading(true);
      const response = await apiFetch<{ data: Cart }>(API_PATHS.CART);
      if (response && response.data) {
        setCart(response.data);
        setCartItemsCount(response.data.cart_items?.length || 0);
      }
    } catch (error) {
      console.error('Error fetching cart', error);
    } finally {
      setIsLoading(false);
    }
  }, []);

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
        await fetchCart();
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
    [fetchCart]
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
        await fetchCart();
      } catch (error) {
        console.error('Error removing item from cart', error);
      } finally {
        setIsLoading(false);
      }
    },
    [fetchCart]
  );

  // Update cart item quantity
  const updateCartItemQuantity = useCallback(
    async (itemId: string, quantity: number) => {
      const user = getCookie('user_id');
      if (!user) return;

      try {
        setIsLoading(true);
        await apiFetch(`${API_PATHS.CART_ITEM}/${itemId}/quantity`, {
          method: 'PUT',
          body: {
            quantity: quantity,
          },
        });
        await fetchCart();
      } catch (error) {
        console.error('Error updating cart item quantity', error);
      } finally {
        setIsLoading(false);
      }
    },
    [fetchCart]
  );

  // Clear cart
  const clearCart = useCallback(async () => {
    try {
      setIsLoading(true);
      await apiFetch(`${API_PATHS.CART}/clear`, {
        method: 'PUT',
      });
      setCart(null);
      setCartItemsCount(0);
    } catch (error) {
      console.error('Error clearing cart', error);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Refresh cart data
  const refreshCart = useCallback(async () => {
    await fetchCart();
  }, [fetchCart]);

  // Load cart data when user is available
  useEffect(() => {
    const user = getCookie('user_id');
    if (user) {
      fetchCart();
    } else {
      setCart(null);
      setCartItemsCount(0);
    }
  }, [fetchCart]);

  const value = {
    isLoading,
    cart,
    cartItemsCount,
    addToCart,
    removeFromCart,
    updateCartItemQuantity,
    clearCart,
    refreshCart,
  };

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
