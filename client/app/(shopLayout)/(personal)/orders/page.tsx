import { Metadata } from 'next';
import OrdersPageClient from './_components/OrdersPageClient';

export const metadata: Metadata = {
  title: 'Your Orders | eShop',
  description: 'View and manage your order history',
};

export default function OrdersPage() {
  return <OrdersPageClient />;
}
