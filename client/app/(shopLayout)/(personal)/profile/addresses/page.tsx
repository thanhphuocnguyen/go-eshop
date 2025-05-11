import { Metadata } from 'next';
import AddressesClient from './AddressesClient';
import { getUserCache } from '@/lib/cache/user';

export const metadata: Metadata = {
  title: 'Shipping Addresses | E-Shop',
  description: 'Manage your shipping addresses for faster checkout.',
};

export default async function AddressesPage() {
  const user = await getUserCache();
  if (!user) {
    return <div>Unauthorized</div>;
  }
  return <AddressesClient user={user} />;
}
