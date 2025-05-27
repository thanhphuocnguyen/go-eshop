'use client';

import { useEffect, useState } from 'react';
import { redirect, useRouter } from 'next/navigation';
import { Elements } from '@stripe/react-stripe-js';
import { loadStripe, StripeError } from '@stripe/stripe-js';
import StripeCheckoutForm from './StripeCheckoutForm';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';
import { CheckoutDataResponse, PaymentResponse } from '../../_lib/definitions';

// Initialize Stripe - replace with your publishable key
// In a real application, you would fetch this from an environment variable
const stripePromise = loadStripe(
  process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || ''
);

export default function StripePage() {
  const [totalPrice, setTotalPrice] = useState(0);
  const [clientSecret, setClientSecret] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  async function getPaymentById(paymentId: string) {
    setIsLoading(true);
    const { data, error } = await clientSideFetch<PaymentResponse>(
      PUBLIC_API_PATHS.PAYMENT.replace(':id', paymentId)
    );

    if (error) {
      console.error('Error fetching payment:', error);
      toast.error(
        error.details || 'Failed to fetch payment details. Please try again.'
      );
      return null;
    }

    if (data.details) {
      setClientSecret(data.details.clientSecret);
      setTotalPrice(data.details.totalPrice);
    }
    setIsLoading(false);
  }

  useEffect(() => {
    // Get checkout data from session storage
    const storedData = sessionStorage.getItem('checkoutData');

    if (!storedData) {
      redirect('/orders');
    } else {
      const parsedData = JSON.parse(storedData) as CheckoutDataResponse;
      if (parsedData.clientSecret) {
        setTotalPrice(parsedData.totalPrice);
        setClientSecret(parsedData.clientSecret);
      } else if (parsedData.paymentId) {
        getPaymentById(parsedData.paymentId);
      } else {
        redirect('/profiles/orders');
      }
    }
  }, []);

  const handlePaymentSuccess = () => {
    // Clear checkout data from session storage
    sessionStorage.removeItem('checkoutData');

    // Redirect to a success page or order confirmation
    router.push('/checkout/success');
  };

  const handlePaymentError = (error: StripeError | unknown) => {
    console.error('Payment error:', error);
    // You could redirect to an error page or display an error message
  };

  return (
    <div className='bg-gray-50 min-h-screen py-12'>
      <div className='max-w-3xl mx-auto px-4'>
        <h1 className='text-2xl font-bold text-gray-900 mb-8'>
          Complete Your Payment
        </h1>

        {isLoading ? (
          <div className='bg-white p-8 rounded-lg shadow-md'>
            <div className='flex justify-center'>
              <div className='animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500'></div>
            </div>
            <p className='text-center mt-4 text-gray-600'>
              Preparing your payment...
            </p>
          </div>
        ) : clientSecret ? (
          <div className='bg-white p-8 rounded-lg shadow-md'>
            <div className='mb-6'>
              <h2 className='text-lg font-medium text-gray-900 mb-2'>
                Order Summary
              </h2>
              <div className='flex justify-between py-2 border-b border-gray-200'>
                <span className='text-gray-600'>Subtotal</span>
                <span className='font-medium'>${totalPrice}</span>
              </div>
              <div className='flex justify-between py-2 border-b border-gray-200'>
                <span className='text-gray-600'>Shipping</span>
                <span className='font-medium'>$0.00</span>
              </div>
              <div className='flex justify-between py-2 border-b border-gray-200'>
                <span className='text-gray-600'>Tax</span>
                <span className='font-medium'>$0.00</span>
              </div>
              <div className='flex justify-between py-3 font-bold'>
                <span>Total</span>
                <span>${totalPrice}</span>
              </div>
            </div>

            {clientSecret && (
              <Elements
                stripe={stripePromise}
                options={{ clientSecret, appearance: { theme: 'stripe' } }}
              >
                <StripeCheckoutForm
                  onSuccess={handlePaymentSuccess}
                  onError={handlePaymentError}
                />
              </Elements>
            )}
          </div>
        ) : (
          <div className='bg-white p-8 rounded-lg shadow-md'>
            <p className='text-center text-red-500'>
              There was an issue preparing your payment. Please try again.
            </p>
            <button
              onClick={() => router.push('/checkout')}
              className='w-full mt-4 px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none'
            >
              Return to Checkout
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
