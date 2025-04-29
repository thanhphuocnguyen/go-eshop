'use client';

import { PaymentInfo } from '@/lib/definitions/order';
import { useState } from 'react';
import PaymentSetupModal from './PaymentSetupModal';
import { ShoppingBagIcon } from '@heroicons/react/16/solid';

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

  return (
    <>
      <div className="font-semibold">Payment information</div>
      {!paymentInfo ? (
        <div className="mt-2">
          <div className="flex items-center text-sm text-red-500">
            <span>No payment method set</span>
          </div>
          <button
            onClick={() => setIsModalOpen(true)}
            className="mt-2 flex items-center gap-1 px-3 py-1.5 text-sm text-white bg-indigo-600 rounded-md hover:bg-indigo-700 transition duration-150"
          >
            <ShoppingBagIcon className="w-4 h-4" />
            <span>Setup Payment</span>
          </button>
        </div>
      ) : (
        <div className="mt-2 text-sm text-gray-600">
          <div>Method: <span className="font-medium">{paymentInfo.method}</span></div>
          <div>Status: <span className="font-medium">{paymentInfo.status}</span></div>
          {paymentInfo.transaction_id && (
            <div>Transaction ID: <span className="font-medium">{paymentInfo.transaction_id}</span></div>
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