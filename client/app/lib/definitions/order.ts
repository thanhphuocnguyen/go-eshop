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
  imageUrl?: string;
  quantity: number;
  lineTotal: number;
  attributesSnapshot: AttributeSnapshot[];
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
  refundId?: string;
  intentId?: string;
  clientSecret?: string;
  transactionId?: string;
}

/**
 * Represents a basic order in list view
 */
export interface Order {
  id: string;
  customerName: string;
  customerEmail: string;
  total: number;
  totalItems: number;
  status: string;
  paymentStatus: string;
  createdAt: string;
  updatedAt?: string;
}

/**
 * Represents a detailed order
 */
export interface OrderDetail {
  id: string;
  customerName: string;
  customerEmail: string;
  total: number;
  status: string;
  paymentStatus: string;
  createdAt: string;
  products: OrderProduct[];
  shippingInfo: ShippingInfo;
  paymentInfo?: PaymentInfo;
}

/**
 * Alternative order model with slightly different structure
 */
export type OrderModel = {
  id: string;
  total: number;
  status: OrderStatus;
  customerName: string;
  customerEmail: string;
  paymentInfo: PaymentInfo;
  shippingInfo: ShippingInfoModel;
  products: OrderItemModel[];
  createdAt: string;
};

/**
 * Represents a product order item with slightly different structure
 */
export type OrderItemModel = {
  id: string;
  variantId: string;
  name: string;
  imageUrl: string;
  lineTotal: number;
  quantity: number;
  rating?: {
    id: string;
    rating: number;
    title: string;
    content: string;
    imageUrl?: string;
    createdAt: string;
  };
  attributeSnapshot: AttributeSnapshot[];
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
  totalItems: number;
  status: OrderStatus;
  paymentStatus: string;
  createdAt: string;
  updatedAt: string;
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
