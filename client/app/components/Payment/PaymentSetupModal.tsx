'use client';

import { useState } from 'react';
import StyledModal from '../StyledModal';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { redirect, useRouter } from 'next/navigation';
import LoadingButton from '../Common/LoadingButton';
import { toast } from 'react-toastify';
import { Button } from '@headlessui/react';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { CheckoutDataResponse } from '@/app/(user)/(personal)/checkout/_lib/definitions';

interface PaymentSetupModalProps {
  orderId: string;
  total: number;
  isOpen: boolean;
  onClose: () => void;
}

type PaymentMethod = 'stripe' | 'cod';

const PaymentSetupModal: React.FC<PaymentSetupModalProps> = ({
  orderId,
  total,
  isOpen,
  onClose,
}) => {
  const [selectedMethod, setSelectedMethod] = useState<PaymentMethod | null>(
    null
  );
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleSetupPayment = async () => {
    if (!selectedMethod) {
      toast.error('Please select a payment method');
      return;
    }

    setIsLoading(true);
    try {
      const response = await clientSideFetch<Partial<CheckoutDataResponse>>(
        PUBLIC_API_PATHS.PAYMENTS,
        {
          method: 'POST',
          body: {
            orderId: orderId,
            payment_method: selectedMethod,
          },
        }
      );

      if (response.error || !response.data) {
        toast.error(response?.error?.details || 'Failed to setup payment');
        return;
      }

      if (selectedMethod === 'stripe') {
        // Redirect to stripe payment page or handle stripe checkout
        localStorage.setItem(
          'checkoutData',
          JSON.stringify({
            orderId: orderId,
            clientSecret: response.data.clientSecret,
            paymentIntentId: response.data.paymentId,
          } as CheckoutDataResponse)
        );
        redirect(`/checkout/payment/stripe`);
      } else {
        // For COD, just refresh the page to show updated payment info
        toast.success('Payment method set to Cash On Delivery');
      }

      router.refresh();
      onClose();
    } catch (error) {
      toast.error('Something went wrong while setting up payment');
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <StyledModal title='Setup Payment' open={isOpen} onClose={onClose}>
      <div className='p-4 w-full max-w-md mx-auto'>
        <p className='text-gray-600 mb-6 text-center'>
          Please select a payment method for your order
        </p>

        <div className='space-y-4 mb-8'>
          <div
            className={`border rounded-lg p-4 cursor-pointer ${
              selectedMethod === 'stripe'
                ? 'border-indigo-500 bg-indigo-50'
                : 'border-gray-200 hover:border-gray-300'
            }`}
            onClick={() => setSelectedMethod('stripe')}
          >
            <div className='flex items-center'>
              <div className='flex-shrink-0 h-6 w-6 flex items-center justify-center'>
                <input
                  type='radio'
                  checked={selectedMethod === 'stripe'}
                  onChange={() => setSelectedMethod('stripe')}
                  className='form-radio h-4 w-4 text-indigo-600'
                />
              </div>
              <div className='ml-3 flex justify-between w-full items-center'>
                <div>
                  <p className='font-medium text-gray-800'>Pay with Stripe</p>
                  <p className='text-sm text-gray-500'>
                    Pay securely with credit or debit card
                  </p>
                </div>
                <div className='flex items-center space-x-1'>
                  <svg
                    className='h-8 w-8'
                    fill='none'
                    xmlns='http://www.w3.org/2000/svg'
                    viewBox='0 0 24 24'
                  >
                    <path
                      d='M12 4c-4.418 0-8 3.582-8 8s3.582 8 8 8h8V12c0-4.418-3.582-8-8-8z'
                      fill='#6772E5'
                    />
                  </svg>
                </div>
              </div>
            </div>
          </div>

          <div
            className={`border rounded-lg p-4 cursor-pointer ${
              selectedMethod === 'cod'
                ? 'border-indigo-500 bg-indigo-50'
                : 'border-gray-200 hover:border-gray-300'
            }`}
            onClick={() => setSelectedMethod('cod')}
          >
            <div className='flex items-center'>
              <div className='flex-shrink-0 h-6 w-6 flex items-center justify-center'>
                <input
                  type='radio'
                  checked={selectedMethod === 'cod'}
                  onChange={() => setSelectedMethod('cod')}
                  className='form-radio h-4 w-4 text-indigo-600'
                />
              </div>
              <div className='ml-3 flex justify-between w-full items-center'>
                <div>
                  <p className='font-medium text-gray-800'>Cash on Delivery</p>
                  <p className='text-sm text-gray-500'>
                    Pay when your order is delivered
                  </p>
                </div>
                <div className='ml-2 text-green-600 font-medium'>
                  ${total.toFixed(2)}
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className='flex justify-end space-x-3'>
          <Button
            type='button'
            onClick={onClose}
            className='px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none'
          >
            Cancel
          </Button>
          <LoadingButton
            onClick={handleSetupPayment}
            isLoading={isLoading}
            className='px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none'
          >
            Confirm Payment Method
          </LoadingButton>
        </div>
      </div>
    </StyledModal>
  );
};

export default PaymentSetupModal;
