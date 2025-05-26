'use client';

import { PaymentInfo } from '@/app/lib/definitions/order';
import { useState } from 'react';
import PaymentSetupModal from './PaymentSetupModal';
import { ShoppingBagIcon, CreditCardIcon } from '@heroicons/react/16/solid';
import { useRouter } from 'next/navigation';
import { CheckoutDataResponse } from '@/app/(user)/(personal)/checkout/_lib/definitions';

interface PaymentInfoSectionProps {
  paymentInfo: PaymentInfo | null;
  orderId: string;
  total: number;
}

const PaymentInfoSection: React.FC<PaymentInfoSectionProps> = ({
  paymentInfo,
  orderId,
  total,
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const router = useRouter();

  const handleContinuePayment = () => {
    if (paymentInfo && paymentInfo.intent_id && paymentInfo.client_secret) {
      // Store the checkout data in session storage
      const checkoutData: CheckoutDataResponse = {
        order_id: orderId,
        payment_id: paymentInfo.intent_id,
        client_secret: paymentInfo.client_secret,
        payment_intent_id: paymentInfo.intent_id,
        total_price: paymentInfo.amount,
      };
      sessionStorage.setItem('checkoutData', JSON.stringify(checkoutData));

      // Navigate to Stripe payment page
      router.push('/checkout/payment/stripe');
    }
  };

  return (
    <>
      <div className='font-semibold'>Payment information</div>
      {!paymentInfo ? (
        <div className='mt-2'>
          <div className='flex items-center text-sm text-red-500'>
            <span>No payment method set</span>
          </div>
          <button
            onClick={() => setIsModalOpen(true)}
            className='mt-2 flex items-center gap-1 px-3 py-1.5 text-sm text-white bg-indigo-600 rounded-md hover:bg-indigo-700 transition duration-150'
          >
            <ShoppingBagIcon className='w-4 h-4' />
            <span>Setup Payment</span>
          </button>
        </div>
      ) : (
        <div className='mt-2 text-sm text-gray-600'>
          <div>
            Method: <span className='font-medium'>{paymentInfo.method}</span>
          </div>
          <div>
            Status:{' '}
            <span
              className={`font-medium ${paymentInfo.status === 'pending' ? 'text-yellow-600' : paymentInfo.status === 'success' ? 'text-green-600' : ''}`}
            >
              {paymentInfo.status}
            </span>
          </div>
          {paymentInfo.intent_id && (
            <div>
              Transaction ID:{' '}
              <span className='font-medium'>{paymentInfo.intent_id}</span>
            </div>
          )}

          {paymentInfo.method === 'stripe' &&
            paymentInfo.status === 'pending' && (
              <button
                onClick={handleContinuePayment}
                className='mt-3 flex items-center gap-1 px-3 py-1.5 text-sm text-white bg-blue-600 rounded-md hover:bg-blue-700 transition duration-150'
              >
                <CreditCardIcon className='w-4 h-4' />
                <span>Continue Payment</span>
              </button>
            )}
        </div>
      )}

      <PaymentSetupModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        orderId={orderId}
        total={total}
      />
    </>
  );
};

export default PaymentInfoSection;
