'use client';

import { TextField } from '@/components/FormFields';
import { Button, Fieldset } from '@headlessui/react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { CheckoutFormSchema, CheckoutFormValues } from './_lib/definitions';
import Image from 'next/image';
import { TrashIcon } from '@heroicons/react/24/solid';

export default function CheckoutPage() {
  const { register, watch } = useForm<CheckoutFormValues>({
    resolver: zodResolver(CheckoutFormSchema),
    defaultValues: {
      address: '',
      cardCvc: '',
      cardExpiry: '',
      cardNumber: '',
      city: '',
      country: '',
      email: '',
      firstName: '',
      lastName: '',
      phone: '',
      paypalEmail: '',
      state: '',
      zip: '',
      paymentMethod: 'credit_card',
      termsAccepted: false,
    },
  });

  const paymentMethod = watch('paymentMethod');

  return (
    <div className='bg-gray-50 m-auto h-full p-10'>
      <div className='flex gap-20 container mx-auto'>
        <div className='w-1/2'>
          <h3 className='text-lg font-semibold text-gray-600 mb-4'>
            Contact Information
          </h3>
          <Fieldset>
            <TextField
              {...register('email')}
              type='email'
              label='Email address'
            />
            <hr className='my-8' />
            <h4 className='text-lg font-semibold text-gray-600'>
              Shipping Information
            </h4>
            <div className='grid grid-cols-2 gap-4 mt-4'>
              <TextField
                {...register('firstName')}
                type='text'
                label='First name'
              />
              <TextField
                {...register('lastName')}
                type='text'
                label='Last name'
              />
              <TextField
                className='col-span-2'
                {...register('address')}
                type='text'
                label='Address'
                placeholder='Street address'
              />
              <TextField
                {...register('city')}
                type='text'
                label='City'
                placeholder='City'
              />
              <TextField
                {...register('country')}
                type='text'
                label='Country'
                placeholder='Country'
              />
              <TextField
                {...register('state')}
                type='text'
                label='State/Province'
                placeholder='State/Province'
              />
              <TextField
                {...register('zip')}
                type='text'
                label='Zip/Postal code'
                placeholder='Zip/Postal code'
              />
              <TextField
                {...register('phone')}
                type='text'
                className='col-span-2'
                label='Phone number'
                placeholder='Phone number'
              />
            </div>

            <hr className='my-8' />
            <h4 className='text-lg font-semibold text-gray-600 mb-4'>
              Payment Method
            </h4>

            <div className='flex flex-col space-y-4'>
              <div className='grid grid-cols-2 gap-4'>
                <div 
                  className={`border rounded-lg p-4 flex items-center cursor-pointer transition-all ${
                    paymentMethod === 'credit_card' 
                      ? 'border-indigo-600 bg-indigo-50' 
                      : 'border-gray-200 hover:border-indigo-300'
                  }`}
                  onClick={() => {
                    const creditCardRadio = document.getElementById('credit_card') as HTMLInputElement;
                    if (creditCardRadio) {
                      creditCardRadio.checked = true;
                      creditCardRadio.dispatchEvent(new Event('change', { bubbles: true }));
                    }
                  }}
                >
                  <input
                    type='radio'
                    id='credit_card'
                    value='credit_card'
                    className='h-5 w-5 text-indigo-600 focus:ring-indigo-500 border-gray-300'
                    {...register('paymentMethod')}
                  />
                  <label htmlFor='credit_card' className='ml-3 flex-1 cursor-pointer'>
                    <div className='font-medium text-gray-800'>Credit Card</div>
                    <div className='text-sm text-gray-500'>Pay with Visa, Mastercard, etc.</div>
                  </label>
                  <div className='flex gap-2 ml-2'>
                    <div className='h-8 w-12 bg-gray-100 rounded-md flex items-center justify-center text-xs font-medium text-gray-800'>Visa</div>
                    <div className='h-8 w-12 bg-gray-100 rounded-md flex items-center justify-center text-xs font-medium text-gray-800'>MC</div>
                  </div>
                </div>

                <div 
                  className={`border rounded-lg p-4 flex items-center cursor-pointer transition-all ${
                    paymentMethod === 'paypal' 
                      ? 'border-indigo-600 bg-indigo-50' 
                      : 'border-gray-200 hover:border-indigo-300'
                  }`}
                  onClick={() => {
                    const paypalRadio = document.getElementById('paypal') as HTMLInputElement;
                    if (paypalRadio) {
                      paypalRadio.checked = true;
                      paypalRadio.dispatchEvent(new Event('change', { bubbles: true }));
                    }
                  }}
                >
                  <input
                    type='radio'
                    id='paypal'
                    value='paypal'
                    className='h-5 w-5 text-indigo-600 focus:ring-indigo-500 border-gray-300'
                    {...register('paymentMethod')}
                  />
                  <label htmlFor='paypal' className='ml-3 flex-1 cursor-pointer'>
                    <div className='font-medium text-gray-800'>PayPal</div>
                    <div className='text-sm text-gray-500'>Pay with your PayPal account</div>
                  </label>
                  <div className='h-8 w-20 bg-blue-100 rounded-md flex items-center justify-center text-sm font-bold text-blue-700'>PayPal</div>
                </div>
              </div>

              {paymentMethod === 'credit_card' && (
                <div className='mt-6 p-5 border border-gray-200 rounded-lg bg-white space-y-4'>
                  <TextField
                    {...register('cardNumber')}
                    type='text'
                    label='Card number'
                    placeholder='0000 0000 0000 0000'
                  />
                  <div className='grid grid-cols-2 gap-4'>
                    <TextField
                      {...register('cardExpiry')}
                      type='text'
                      label='Expiration date'
                      placeholder='MM/YY'
                    />
                    <TextField
                      {...register('cardCvc')}
                      type='text'
                      label='CVC'
                      placeholder='000'
                    />
                  </div>
                </div>
              )}

              {paymentMethod === 'paypal' && (
                <div className='mt-6 p-5 border border-gray-200 rounded-lg bg-white'>
                  <TextField
                    {...register('paypalEmail')}
                    type='email'
                    label='PayPal email'
                    placeholder='Enter your PayPal email'
                  />
                </div>
              )}

              <div className='mt-6 flex items-start'>
                <input
                  type='checkbox'
                  id='terms'
                  className='h-5 w-5 mt-0.5 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded'
                  {...register('termsAccepted')}
                />
                <label htmlFor='terms' className='ml-3 text-sm text-gray-600'>
                  I agree to the terms and conditions and the privacy policy
                </label>
              </div>
            </div>
          </Fieldset>
        </div>
        <div className='w-1/2'>
          <h3 className='text-lg font-semibold text-gray-600 mb-4'>
            Order summary
          </h3>
          <div className='border border-gray-200 bg-white rounded-md'>
            <div className='flex gap-4 p-6 border-b border-gray-200'>
              <div className='h-40 w-32 relative'>
                <Image
                  fill
                  objectFit='conditions'
                  src={'/images/product-placeholder.webp'}
                  alt='Product Image'
                  className='rounded-md border border-lime-300'
                />
              </div>
              <div className='flex-1 flex flex-col justify-between'>
                <div>
                  <div className='flex justify-between'>
                    <span className='font-medium'>Basic tee</span>
                    <Button>
                      <TrashIcon className='size-6 text-gray-500' />
                    </Button>
                  </div>
                  <div className='text-gray-500'>Black</div>
                  <div className='text-gray-500'>Large</div>
                </div>
                <div className='flex justify-between'>
                  <span>$32.00</span>
                  <span className='text-gray-500'>Qty: 1</span>
                </div>
              </div>
            </div>
            <div className='px-6 pt-6'>
              <div className='flex flex-col gap-6'>
                <div className='flex justify-between'>
                  <span>Subtotal</span>
                  <span>$64.00</span>
                </div>
                <div className='flex justify-between'>
                  <span>Shipping</span>
                  <span>$5.00</span>
                </div>
                <div className='flex justify-between'>
                  <span>Taxes</span>
                  <span>$5.20</span>
                </div>
              </div>
              <hr className='my-6' />
              <div className='flex justify-between font-semibold'>
                <span>Total</span>
                <span>$74.20</span>
              </div>
            </div>
            <hr className='my-6' />
            <div className='px-6 pb-6'>
              <Button
                type='submit'
                className='w-full bg-indigo-600 h-12 text-white py-2 rounded-md hover:bg-indigo-700'
              >
                Confirm Order
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
