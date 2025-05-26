'use client';

import React from 'react';
import { TextField } from '@/app/components/FormFields';
import {
  Button,
  Checkbox,
  Fieldset,
  Radio,
  RadioGroup,
} from '@headlessui/react';
import Image from 'next/image';
import { TrashIcon } from '@heroicons/react/24/solid';
import { CheckIcon } from '@heroicons/react/24/outline';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm, useWatch } from 'react-hook-form';

import { useEffect, useState } from 'react';
import { redirect, RedirectType, useRouter } from 'next/navigation';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';
import { useCartCtx } from '@/app/lib/contexts/CartContext';
import {
  CheckoutDataResponse,
  CheckoutFormSchema,
  CheckoutFormValues,
} from '../_lib/definitions';
import { useUser } from '@/app/lib/hooks/useUser';

const CheckoutDetailOverview: React.FC = () => {
  const router = useRouter();
  const { user, isLoading } = useUser();
  const { cart, cartLoading } = useCartCtx();
  const [isNewAddress, setIsNewAddress] = useState(true);
  const [selectedAddressId, setSelectedAddressId] = useState<string | null>(
    null
  );

  const { register, control, watch, reset, setValue, handleSubmit } =
    useForm<CheckoutFormValues>({
      resolver: zodResolver(CheckoutFormSchema),
      defaultValues: {
        address: {
          city: '',
          street: '',
          district: '',
          phone: '',
        },
        fullname: '',
        email: '',
        payment_method: 'cod',
        terms_accepted: false,
      },
    });

  const paymentMethod = useWatch({ control, name: 'payment_method' });

  const handleAddressChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const addressId = e.target.value;

    if (addressId === null) {
      setIsNewAddress(true);
      setSelectedAddressId(null);

      setValue('address', {
        city: '',
        district: '',
        phone: '',
        street: '',
        ward: '',
      });
    } else {
      setIsNewAddress(false);
      setSelectedAddressId(addressId);

      const selectedAddress = user?.addresses?.find(
        (addr) => addr.id === addressId
      );
      if (selectedAddress) {
        setValue('address_id', selectedAddress.id);
        setValue('address', {
          street: selectedAddress.street,
          city: selectedAddress.city,
          district: selectedAddress.district,
          ward: selectedAddress.ward,
          phone: selectedAddress.phone,
        });
      }
    }
  };

  useEffect(() => {
    if (user) {
      const defaultValue: Partial<CheckoutFormValues> = {
        email: user.email,
        fullname: user.fullname,
        payment_method: 'cod',
        terms_accepted: false,
      };

      if (user.addresses && user.addresses.length > 0) {
        const defaultAddress =
          user.addresses.find((addr) => addr.default) || user.addresses[0];
        defaultValue.address_id = defaultAddress.id;
        defaultValue.address = {
          street: defaultAddress.street,
          city: defaultAddress.city,
          district: defaultAddress.district,
          phone: defaultAddress.phone,
          ward: defaultAddress.ward,
        };
        setIsNewAddress(false);
        setSelectedAddressId(defaultAddress.id);
      } else {
        setIsNewAddress(true);
        setSelectedAddressId(null);
      }

      reset(defaultValue);
    }
  }, [user, reset, setValue]);

  const onSubmit = async (body: CheckoutFormValues) => {
    // Save form data to session storage for the next step
    const { data, error } = await clientSideFetch<CheckoutDataResponse>(
      PUBLIC_API_PATHS.CHECKOUT,
      {
        method: 'POST',
        body: {
          ...body,
          address: body.address_id ? undefined : body.address,
          payment_receipt_email: body.payment_receipt_email
            ? body.payment_receipt_email
            : undefined,
        },
      }
    );

    if (error) {
      if (error.code === 'payment_gateway_error') {
        toast.error(
          <div>
            <h3 className='text-lg font-semibold text-red-600 mb-2'>
              Payment gateway error
            </h3>
            <div>{JSON.stringify(error)}</div>
          </div>
        );
        redirect(`orders/${data.order_id}`, RedirectType.replace);
      }

      toast.error(
        <div>
          <h3 className='text-lg font-semibold text-red-600 mb-2'>
            Error checkout
          </h3>
          <p className='text-sm text-gray-500'>{JSON.stringify(error)}</p>
        </div>
      );
      return;
    }

    if (data) {
      sessionStorage.setItem('checkoutData', JSON.stringify(data));
      // If Stripe is selected, redirect to the Stripe payment page
      if (body.payment_method === 'stripe') {
        if (!data.client_secret || !data.payment_id) {
          toast.error(
            <div>
              <h3 className='text-lg font-semibold text-red-600 mb-2'>
                Error checkout
              </h3>
              <p className='text-sm text-gray-500'>Invalid payment data</p>
            </div>
          );
          redirect(`orders/${data.order_id}`, RedirectType.replace);
        }
        router.push('/checkout/payment/stripe');
      } else {
        // Handle COD checkout
        // You would typically call your API to create the order here
        console.log('Processing COD order', body);
        // Implement your COD order processing logic
      }
    } else {
      toast.error(
        <div>
          <h3 className='text-lg font-semibold text-red-600 mb-2'>
            Error create order
          </h3>
          <p className='text-sm text-gray-500'>Unknown error</p>
        </div>
      );
    }
  };
  if (isLoading || cartLoading) {
    return (
      <div className='flex items-center justify-center h-screen'>
        <div className='loader'></div>
      </div>
    );
  }

  return (
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

          {user?.addresses && user.addresses.length > 0 && (
            <div className='mt-4 mb-6'>
              <label className='block text-sm font-medium text-gray-700 mb-1'>
                Select Address
              </label>
              <div className='flex gap-4'>
                <select
                  className='flex-1 py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-sm'
                  value={selectedAddressId === null ? -1 : selectedAddressId}
                  onChange={handleAddressChange}
                >
                  {user.addresses.map((address) => (
                    <option key={address.id} value={address.id}>
                      {address.street}, {address.district}, {address.city}
                      {address.default ? ' (Default)' : ''}
                    </option>
                  ))}
                  <option value='-1'>+ Add new address</option>
                </select>

                {!isNewAddress && (
                  <Button
                    type='button'
                    onClick={() => {
                      setIsNewAddress(true);
                      setSelectedAddressId(null);
                    }}
                    className='bg-white border border-gray-300 rounded-md shadow-sm px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50'
                  >
                    Add New
                  </Button>
                )}
              </div>
            </div>
          )}

          <div className='grid grid-cols-2 gap-6 mt-4'>
            <TextField
              {...register('fullname')}
              type='text'
              label='Full name'
            />
            <TextField
              className=''
              {...register('address.street')}
              type='text'
              label='Street address'
              placeholder='Street address'
              disabled={!isNewAddress}
            />
            <TextField
              {...register('address.city')}
              type='text'
              label='City'
              placeholder='City'
              disabled={!isNewAddress}
            />
            <TextField
              {...register('address.district')}
              type='text'
              label='District'
              placeholder='District'
              disabled={!isNewAddress}
            />
            <TextField
              {...register('address.ward')}
              type='text'
              label='Ward'
              placeholder=''
              disabled={!isNewAddress}
            />
            <TextField
              {...register('address.phone')}
              type='phone'
              className=''
              label='Phone number'
              placeholder='Phone number'
              disabled={!isNewAddress}
            />
          </div>

          <hr className='my-8' />
          <h4 className='text-lg font-semibold text-gray-600 mb-4'>
            Payment Method
          </h4>

          <div className='flex flex-col space-y-6'>
            <RadioGroup
              value={paymentMethod}
              onChange={(value) =>
                setValue('payment_method', value as 'stripe' | 'cod')
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

            {paymentMethod === 'stripe' && (
              <div className='mt-6 p-5 border border-gray-200 rounded-lg bg-white'>
                <TextField
                  {...register('payment_receipt_email')}
                  type='email'
                  label='Email for payment receipt'
                  placeholder='Enter your email for payment receipt'
                />
                <p className='mt-2 text-sm text-gray-600'>
                  You will be redirected to our secure payment page after
                  confirming your order.
                </p>
              </div>
            )}

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
                checked={!!watch('terms_accepted')}
                onChange={(checked) => setValue('terms_accepted', checked)}
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
        </Fieldset>
      </div>
      <div className='w-1/2'>
        <h3 className='text-lg font-semibold text-gray-600 mb-4'>
          Order summary
        </h3>
        <div className='border border-gray-200 bg-white rounded-md'>
          {cart?.cartItems.map((e) => (
            <div
              key={e.variantId}
              className='flex gap-4 p-6 border-b border-gray-200'
            >
              <div className='h-28 w-24 relative'>
                <Image
                  fill
                  objectFit='contains'
                  src={e.imageUrl ?? '/images/logos/logo.webp'}
                  alt='Product Image'
                  className='rounded-md border border-lime-300'
                />
              </div>
              <div
                key={e.variantId}
                className='flex-1 flex flex-col justify-between'
              >
                <div>
                  <div className='flex justify-between'>
                    <span className='font-medium'>{e.name}</span>
                    <Button>
                      <TrashIcon className='size-6 text-red-200' />
                    </Button>
                  </div>
                  <div className='flex flex-col gap-1 mt-2'>
                    {e.attributes.map((attribute) => (
                      <div
                        key={attribute.name}
                        className='text-md text-gray-500'
                      >
                        {attribute.name}: {attribute.value}
                      </div>
                    ))}
                  </div>
                </div>
                <div className='flex justify-between'>
                  <span>${e.price}</span>
                  <span className='text-gray-500'>Qty: {e.quantity}</span>
                </div>
              </div>
            </div>
          ))}

          <div className='px-6 pt-6'>
            <div className='flex flex-col gap-6'>
              <div className='flex justify-between'>
                <span>Subtotal</span>
                <span>${cart?.totalPrice}</span>
              </div>
              <div className='flex justify-between'>
                <span>Shipping</span>
                <span>$0.00</span>
              </div>
              <div className='flex justify-between'>
                <span>Taxes</span>
                <span>$0.20</span>
              </div>
            </div>
            <hr className='my-6' />
            <div className='flex justify-between font-semibold'>
              <span>Total</span>
              <span>${cart?.totalPrice ?? 0}</span>
            </div>
          </div>
          <hr className='my-6' />
          <div className='px-6 pb-6'>
            <Button
              onClick={handleSubmit(onSubmit, (err) => {
                console.log(err);
              })}
              className='w-full bg-indigo-600 h-12 text-white py-2 rounded-md hover:bg-indigo-700'
            >
              Confirm Order
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CheckoutDetailOverview;
