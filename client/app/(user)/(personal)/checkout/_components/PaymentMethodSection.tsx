'use client';

import { Checkbox, Radio, RadioGroup } from '@headlessui/react';
import { CheckIcon } from '@heroicons/react/24/outline';
import { UseFormSetValue, UseFormWatch } from 'react-hook-form';
import { CheckoutFormValues } from '../_lib/definitions';

interface PaymentMethodSectionProps {
  setValue: UseFormSetValue<CheckoutFormValues>;
  watch: UseFormWatch<CheckoutFormValues>;
  paymentMethod: 'stripe' | 'cod';
}

export const PaymentMethodSection: React.FC<PaymentMethodSectionProps> = ({
  setValue,
  watch,
  paymentMethod,
}) => {
  return (
    <div>
      <hr className='my-8' />
      <h4 className='text-lg font-semibold text-gray-600 mb-4'>
        Payment Method
      </h4>

      <div className='flex flex-col space-y-6'>
        <RadioGroup
          value={paymentMethod}
          onChange={(value) =>
            setValue('paymentMethod', value as 'stripe' | 'cod')
          }
          className='grid grid-cols-2 gap-4'
        >
          <Radio
            value='stripe'
            className={({ checked }) =>
              `border rounded-lg p-4 flex items-center cursor-pointer transition-all 
                ${checked ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-indigo-300'}`
            }
          >
            {({ checked }) => (
              <>
                <div className='h-5 w-5 mr-3 flex items-center justify-center'>
                  <div
                    className={`h-3 w-3 rounded-full ${checked ? 'bg-indigo-600' : 'bg-gray-300'}`}
                  />
                </div>
                <div className='flex-1'>
                  <div className='font-medium text-gray-800'>Stripe</div>
                  <div className='text-sm text-gray-500'>
                    Pay securely with Stripe
                  </div>
                </div>
                <div className='h-8 w-20 bg-purple-100 rounded-md flex items-center justify-center text-sm font-bold text-purple-700'>
                  Stripe
                </div>
              </>
            )}
          </Radio>

          <Radio
            value='cod'
            className={({ checked }) =>
              `border rounded-lg p-4 flex items-center cursor-pointer transition-all 
                ${checked ? 'border-indigo-600 bg-indigo-50' : 'border-gray-200 hover:border-indigo-300'}`
            }
          >
            {({ checked }) => (
              <>
                <div className='h-5 w-5 mr-3 flex items-center justify-center'>
                  <div
                    className={`h-3 w-3 rounded-full ${checked ? 'bg-indigo-600' : 'bg-gray-300'}`}
                  />
                </div>
                <div className='flex-1'>
                  <div className='font-medium text-gray-800'>
                    Cash on Delivery
                  </div>
                  <div className='text-sm text-gray-500'>
                    Pay with cash upon delivery
                  </div>
                </div>
                <div className='h-8 w-20 bg-green-100 rounded-md flex items-center justify-center text-sm font-bold text-green-700'>
                  COD
                </div>
              </>
            )}
          </Radio>
        </RadioGroup>

        {paymentMethod === 'cod' && (
          <div className='mt-6 p-5 border border-gray-200 rounded-lg bg-white'>
            <p className='text-sm text-gray-600'>
              You will pay in cash when your order is delivered. No
              additional information is required.
            </p>
          </div>
        )}

        <div className='mt-6'>
          <Checkbox
            checked={!!watch('termsAccepted')}
            onChange={(checked) => setValue('termsAccepted', checked)}
            className='flex items-center'
          >
            {({ checked }) => (
              <>
                <div className='flex h-5 w-5 items-center justify-center rounded border border-gray-300 bg-white'>
                  {checked && (
                    <CheckIcon className='h-4 w-4 text-indigo-600' />
                  )}
                </div>
                <span className='ml-3 text-sm text-gray-600'>
                  I agree to the terms and conditions and the privacy policy
                </span>
              </>
            )}
          </Checkbox>
        </div>
      </div>
    </div>
  );
};
