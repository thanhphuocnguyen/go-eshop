// Order types for the e-commerce application
import { OrderStatus } from './common';

/**
 * Represents an attribute snapshot in an order product
 */
export interface AttributeSnapshot {
  name: string;
  value: string;
}

/**
 * Represents a product in an order
 */
export interface OrderProduct {
  id: string;
  name: string;
  image_url?: string;
  quantity: number;
  line_total: number;
  attributes_snapshot: AttributeSnapshot[];
}

/**
 * Represents shipping information for an order
 */
export interface ShippingInfo {
  street: string;
  ward: string;
  district: string;
  city: string;
  phone: string;
}

/**
 * Represents payment information for an order
 */
export interface PaymentInfo {
  id: string;
  amount: number;
  method: string;
  status: string;
  gateway?: string;
  refund_id?: string;
  intent_id?: string;
  client_secret?: string;
  transaction_id?: string;
}

/**
 * Represents a basic order in list view
 */
export interface Order {
  id: string;
  customer_name: string;
  customer_email: string;
  total: number;
  total_items: number;
  status: string;
  payment_status: string;
  created_at: string;
  updated_at?: string;
}

/**
 * Represents a detailed order
 */
export interface OrderDetail {
  id: string;
  customer_name: string;
  customer_email: string;
  total: number;
  status: string;
  payment_status: string;
  created_at: string;
  products: OrderProduct[];
  shipping_info: ShippingInfo;
  payment_info?: PaymentInfo;
}

/**
 * Alternative order model with slightly different structure
 */
export type OrderModel = {
  id: string;
  total: number;
  status: OrderStatus;
  customer_name: string;
  customer_email: string;
  payment_info: PaymentInfo;
  shipping_info: ShippingInfoModel;
  products: OrderItemMode[];
  created_at: string;
};

/**
 * Represents a product order item with slightly different structure
 */
export type OrderItemMode = {
  id: string;
  name: string;
  image_url: string;
  attribute_snapshot: AttributeSnapshot[];
  line_total: number;
  quantity: number;
};

/**
 * Alternative shipping info structure with name field
 */
export type ShippingInfoModel = {
  name: string;
  phone: string;
  address: string;
  city: string;
  district: string;
  ward: string;
};

/**
 * Order model used in the order list page
 */
export type OrderListModel = {
  id: string;
  total: number;
  total_items: number;
  status: OrderStatus;
  payment_status: string;
  created_at: string;
  updated_at: string;
};

/**
 * Stats for orders summary
 */
export interface OrdersStats {
  total: number;
  pending: number;
  confirm: number;
  delivering: number;
  delivered: number;
  cancelled: number;
  refunded: number;
  completed: number;
  totalSpent: number;
}
