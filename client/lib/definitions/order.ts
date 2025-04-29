import { OrderStatus } from './common';

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

type OrderAttributeSnapshot = {
  name: string;
  value: string;
};
export type OrderItemMode = {
  id: string;
  name: string;
  image_url: string;
  attribute_snapshot: OrderAttributeSnapshot[];
  line_total: number;
  quantity: number;
};

export type ShippingInfoModel = {
  name: string;
  phone: string;
  address: string;
  city: string;
  district: string;
  ward: string;
};

export type PaymentInfo = {
  method: string;
  status: string;
  transaction_id: string;
  amount: number;
};
