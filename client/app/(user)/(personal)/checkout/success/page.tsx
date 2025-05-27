'use client';

import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { CheckCircleIcon } from '@heroicons/react/24/solid';
import { Button } from '@headlessui/react';

export default function PaymentSuccessPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [orderDetails, setOrderDetails] = useState<{
    orderId: string;
    paymentId?: string;
  } | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Get paymentIntent and paymentIntentClientSecret from URL if available
    const paymentIntent = searchParams.get('paymentIntent');

    // Try to get order details from session storage
    const storedData = sessionStorage.getItem('checkoutData');

    if (storedData) {
      try {
        const parsedData = JSON.parse(storedData);
        setOrderDetails({
          orderId: parsedData.orderId,
          paymentId: paymentIntent || undefined,
        });
        // Clear the checkout data from session storage
        sessionStorage.removeItem('checkoutData');
      } catch (error) {
        console.error('Error parsing stored checkout data:', error);
      }
    }

    setIsLoading(false);
  }, [searchParams]);

  const handleViewOrder = () => {
    if (orderDetails?.orderId) {
      router.push(`/orders/${orderDetails.orderId}`);
    } else {
      router.push('/orders');
    }
  };

  const handleContinueShopping = () => {
    router.push('/');
  };

  if (isLoading) {
    return (
      <div className='h-screen flex justify-center items-center'>
        <div className='animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500'></div>
      </div>
    );
  }

  return (
    <div className='bg-gray-50 min-h-screen py-12'>
      <div className='max-w-3xl mx-auto px-4 text-center'>
        <CheckCircleIcon className='h-24 w-24 text-green-500 mx-auto' />

        <h1 className='mt-6 text-3xl font-bold text-gray-900'>
          Payment Successful!
        </h1>

        <p className='mt-4 text-lg text-gray-600'>
          Thank you for your purchase. Your order has been confirmed and will be
          shipped soon.
        </p>

        {orderDetails?.orderId && (
          <p className='mt-2 text-md text-gray-500'>
            Order ID:{' '}
            <span className='font-medium'>{orderDetails.orderId}</span>
          </p>
        )}

        <div className='mt-12 bg-white p-8 rounded-lg shadow-md'>
          <div className='mb-8 pb-8 border-b border-gray-200'>
            <h2 className='text-xl font-semibold text-gray-800 mb-4'>
              What happens next?
            </h2>

            <div className='grid grid-cols-1 md:grid-cols-3 gap-6'>
              <div className='flex flex-col items-center'>
                <div className='bg-indigo-100 p-3 rounded-full'>
                  <svg
                    xmlns='http://www.w3.org/2000/svg'
                    className='h-6 w-6 text-indigo-600'
                    fill='none'
                    viewBox='0 0 24 24'
                    stroke='currentColor'
                  >
                    <path
                      strokeLinecap='round'
                      strokeLinejoin='round'
                      strokeWidth={2}
                      d='M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z'
                    />
                  </svg>
                </div>
                <h3 className='font-medium text-gray-800 mt-3'>
                  Order Confirmation
                </h3>
                <p className='text-sm text-gray-600 text-center mt-1'>
                  You&apos;ll receive an email with your order details
                </p>
              </div>

              <div className='flex flex-col items-center'>
                <div className='bg-indigo-100 p-3 rounded-full'>
                  <svg
                    xmlns='http://www.w3.org/2000/svg'
                    className='h-6 w-6 text-indigo-600'
                    fill='none'
                    viewBox='0 0 24 24'
                    stroke='currentColor'
                  >
                    <path
                      strokeLinecap='round'
                      strokeLinejoin='round'
                      strokeWidth={2}
                      d='M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4'
                    />
                  </svg>
                </div>
                <h3 className='font-medium text-gray-800 mt-3'>
                  Order Processing
                </h3>
                <p className='text-sm text-gray-600 text-center mt-1'>
                  We&apos;ll prepare your items for shipping
                </p>
              </div>

              <div className='flex flex-col items-center'>
                <div className='bg-indigo-100 p-3 rounded-full'>
                  <svg
                    xmlns='http://www.w3.org/2000/svg'
                    className='h-6 w-6 text-indigo-600'
                    fill='none'
                    viewBox='0 0 24 24'
                    stroke='currentColor'
                  >
                    <path
                      strokeLinecap='round'
                      strokeLinejoin='round'
                      strokeWidth={2}
                      d='M8 7h12m0 0l-4-4m4 4l-4 4m-8 6H4m0 0l4 4m-4-4l4-4'
                    />
                  </svg>
                </div>
                <h3 className='font-medium text-gray-800 mt-3'>Shipping</h3>
                <p className='text-sm text-gray-600 text-center mt-1'>
                  You&apos;ll receive tracking info when your order ships
                </p>
              </div>
            </div>
          </div>

          <div className='flex flex-col sm:flex-row justify-center gap-4'>
            <Button
              onClick={handleViewOrder}
              className='bg-indigo-600 text-white px-8 py-3 rounded-md hover:bg-indigo-700 transition'
            >
              View Order Details
            </Button>

            <Button
              onClick={handleContinueShopping}
              className='bg-white text-indigo-600 border border-indigo-600 px-8 py-3 rounded-md hover:bg-indigo-50 transition'
            >
              Continue Shopping
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
