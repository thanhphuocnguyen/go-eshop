'use client';

import { useState } from 'react';
import {
  PaymentElement,
  useStripe,
  useElements,
} from '@stripe/react-stripe-js';
import { StripeError } from '@stripe/stripe-js';

interface StripeCheckoutFormProps {
  onSuccess: () => void;
  onError: (error: StripeError | unknown) => void;
}

export default function StripeCheckoutForm({
  onSuccess,
  onError,
}: StripeCheckoutFormProps) {
  const stripe = useStripe();
  const elements = useElements();
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!stripe || !elements) {
      // Stripe.js hasn't yet loaded.
      return;
    }

    setIsProcessing(true);
    setErrorMessage(null);

    try {
      const { error } = await stripe.confirmPayment({
        elements,
        confirmParams: {
          // Make sure to change this to your payment completion page
          return_url: `${window.location.origin}/checkout/success`,
        },
        redirect: 'if_required',
      });

      if (error) {
        setErrorMessage(
          error.message || 'Something went wrong with your payment.'
        );
        onError(error);
      } else {
        // The payment was successful
        onSuccess();
      }
    } catch (e: unknown) {
      setErrorMessage('An unexpected error occurred.');
      onError(e);
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2 className='text-lg font-medium text-gray-900 mb-4'>
        Payment Details
      </h2>

      <div className='mb-6'>
        <PaymentElement />
      </div>

      {errorMessage && (
        <div className='p-3 rounded bg-red-50 text-red-700 mb-4'>
          {errorMessage}
        </div>
      )}

      <button
        disabled={isProcessing || !stripe || !elements}
        type='submit'
        className='w-full px-4 py-3 bg-indigo-600 text-white font-medium rounded-md 
                  hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 
                  focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed'
      >
        {isProcessing ? (
          <div className='flex items-center justify-center'>
            <div className='animate-spin rounded-full h-4 w-4 border-t-2 border-b-2 border-white mr-2'></div>
            Processing...
          </div>
        ) : (
          'Pay Now'
        )}
      </button>

      <p className='mt-4 text-sm text-gray-600 text-center'>
        Your payment information is processed securely through Stripe.
      </p>
    </form>
  );
}
