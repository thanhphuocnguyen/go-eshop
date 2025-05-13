import { Metadata } from 'next';
import CheckoutDetailOverview from './_components/CheckoutDetailOverview';

export const metadata: Metadata = {
  title: 'Checkout | E-Shop',
  description: 'Secure checkout page for your purchases.',
};

export default function CheckoutPage() {
  return (
    <div className='bg-gray-50 m-auto h-full p-10'>
      <CheckoutDetailOverview />
    </div>
  );
}
