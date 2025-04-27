'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Elements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';
import StripeCheckoutForm from './StripeCheckoutForm';
import { useAppUser } from '@/components/AppUserContext';

// Initialize Stripe - replace with your publishable key
// In a real application, you would fetch this from an environment variable
const stripePromise = loadStripe('pk_test_your_stripe_key');

export default function StripePage() {
  const [clientSecret, setClientSecret] = useState<string | null>(null);
  const [checkoutData, setCheckoutData] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();
  const { cart } = useAppUser();

  useEffect(() => {
    // Get checkout data from session storage
    const storedData = sessionStorage.getItem('checkoutData');
    if (!storedData) {
      // Redirect back to checkout if no data is available
      router.push('/checkout');
      return;
    }

    const parsedData = JSON.parse(storedData);
    setCheckoutData(parsedData);

    const fetchPaymentIntent = async () => {
      try {
        // In a real application, you would call your backend to create a payment intent
        // Example API call:
        // const response = await fetch('/api/create-payment-intent', {
        //   method: 'POST',
        //   headers: { 'Content-Type': 'application/json' },
        //   body: JSON.stringify({
        //     amount: cart?.total_price || 0,
        //     email: parsedData.stripeEmail || parsedData.email,
        //   }),
        // });
        // const data = await response.json();
        // setClientSecret(data.clientSecret);

        // For demo purposes, we'll just simulate a successful response after a delay
        setTimeout(() => {
          setClientSecret('dummy_client_secret');
          setIsLoading(false);
        }, 1000);
      } catch (error) {
        console.error('Error creating payment intent:', error);
        setIsLoading(false);
      }
    };

    fetchPaymentIntent();
  }, [router, cart]);

  const handlePaymentSuccess = () => {
    // Clear checkout data from session storage
    sessionStorage.removeItem('checkoutData');

    // Redirect to a success page or order confirmation
    router.push('/checkout/success');
  };

  const handlePaymentError = (error: any) => {
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
                <span className='font-medium'>${cart?.total_price || 0}</span>
              </div>
              <div className='flex justify-between py-2 border-b border-gray-200'>
                <span className='text-gray-600'>Shipping</span>
                <span className='font-medium'>$0.00</span>
              </div>
              <div className='flex justify-between py-2 border-b border-gray-200'>
                <span className='text-gray-600'>Tax</span>
                <span className='font-medium'>$0.20</span>
              </div>
              <div className='flex justify-between py-3 font-bold'>
                <span>Total</span>
                <span>${cart?.total_price || 0}</span>
              </div>
            </div>

            <Elements
              stripe={stripePromise}
              options={{ clientSecret, appearance: { theme: 'stripe' } }}
            >
              <StripeCheckoutForm
                onSuccess={handlePaymentSuccess}
                onError={handlePaymentError}
              />
            </Elements>
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
