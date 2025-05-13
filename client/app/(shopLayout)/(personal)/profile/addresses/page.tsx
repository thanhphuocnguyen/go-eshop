import { Metadata } from 'next';
import AddressesClient from './AddressesClient';

export const metadata: Metadata = {
  title: 'Shipping Addresses | E-Shop',
  description: 'Manage your shipping addresses for faster checkout.',
};

export default async function AddressesPage() {
  return <AddressesClient />;
}
