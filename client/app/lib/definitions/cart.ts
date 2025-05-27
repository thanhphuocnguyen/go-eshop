export interface CartItem {
  id: string;
  productId: string;
  variantId: string;
  name: string;
  quantity: number;
  price: number;
  discount: number;
  stock: number;
  sku: string;
  imageUrl?: string;
  categoryId?: string; // Add categoryId for category-based discounts
  attributes: Array<{
    name: string;
    value: string;
  }>;
}

export interface CartModel {
  id: string;
  userId: string;
  totalPrice: number;
  shippingFee?: number;
  tax?: number;
  discount?: number;
  cartItems: CartItem[];
  updatedAt: string;
  createdAt: string;
}
