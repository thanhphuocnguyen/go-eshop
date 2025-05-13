export interface CartItem {
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

export interface CartModel {
  id: string;
  user_id: string;
  total_price: number;
  shipping_fee?: number;
  tax?: number;
  discount?: number;
  cart_items: CartItem[];
  updated_at: string;
  created_at: string;
}
