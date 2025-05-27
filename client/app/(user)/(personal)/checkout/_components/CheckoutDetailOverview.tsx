'use client';

import React from 'react';
import { Button, Fieldset } from '@headlessui/react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm, useWatch } from 'react-hook-form';
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
import { useUser } from '@/app/hooks';
import useSWR from 'swr';

// Import smaller components
import { ContactInformationSection } from './ContactInformationSection';
import { ShippingInformationSection } from './ShippingInformationSection';
import { PaymentMethodSection } from './PaymentMethodSection';
import { OrderSummarySection } from './OrderSummarySection';
import { DiscountSection } from './DiscountSection';
import { OrderTotalSection } from './OrderTotalSection';

// Import custom hooks
import { useDiscountLogic } from './hooks/useDiscountLogic';
import { useAddressLogic } from './hooks/useAddressLogic';

type DiscountModel = {
  id: string;
  code: string;
  description: string;
  discountType: 'percentage' | 'fixed';
  discountValue: number;
  startsAt: string;
  expiresAt: string;
  categoryId?: string;
  productId?: string;
};
const CheckoutDetailOverview: React.FC = () => {
  const router = useRouter();
  const { user, isLoading } = useUser();
  const { cart, cartLoading } = useCartCtx();

  // Fetch discounts
  const { data: discounts, isLoading: discountsLoading } = useSWR(
    PUBLIC_API_PATHS.CART_DISCOUNT_LIST,
    (url) => clientSideFetch<DiscountModel[]>(url).then((data) => data.data),
    {
      revalidateOnFocus: false,
      revalidateIfStale: false,
      revalidateOnReconnect: false,
      onError: (err) => {
        toast.error(
          <div>
            <h3 className='text-lg font-semibold text-red-600 mb-2'>
              Error fetching discounts
            </h3>
            <p className='text-sm text-gray-500'>{JSON.stringify(err)}</p>
          </div>
        );
      },
    }
  );

  // Form setup
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
        paymentMethod: 'cod',
        discountCode: '',
        termsAccepted: false,
      },
    });

  const paymentMethod = useWatch({ control, name: 'paymentMethod' });
  const selectedDiscountCode = useWatch({ control, name: 'discountCode' });

  // Custom hooks for address and discount logic
  const { isNewAddress, selectedAddressId, handleAddressChange, handleAddNewAddress } = 
    useAddressLogic({ user, setValue, reset });

  const { appliedDiscount, discountAmount, handleApplyDiscount, handleRemoveDiscount } = 
    useDiscountLogic({
      selectedDiscountCode,
      discounts,
      cart,
      setValue,
    });

  // Handle discount apply with toast
  const onApplyDiscount = () => {
    const result = handleApplyDiscount();
    if (result) {
      if (result.success) {
        toast.success(result.message);
      } else {
        toast.error(result.message);
      }
    }
  };

  // Handle discount remove with toast
  const onRemoveDiscount = () => {
    const result = handleRemoveDiscount();
    if (result) {
      toast.info(result.message);
    }
  };

  // Submit handler
  const onSubmit = async (body: CheckoutFormValues) => {
    const { data, error } = await clientSideFetch<CheckoutDataResponse>(
      PUBLIC_API_PATHS.CHECKOUT,
      {
        method: 'POST',
        body: {
          ...body,
          address: body.addressId ? undefined : body.address,
          discountCode: appliedDiscount ? appliedDiscount.code : undefined,
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
        redirect(`orders/${data.orderId}`, RedirectType.replace);
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
      if (body.paymentMethod === 'stripe') {
        if (!data.clientSecret || !data.paymentId) {
          toast.error(
            <div>
              <h3 className='text-lg font-semibold text-red-600 mb-2'>
                Error checkout
              </h3>
              <p className='text-sm text-gray-500'>Invalid payment data</p>
            </div>
          );
          redirect(`orders/${data.orderId}`, RedirectType.replace);
        }
        router.push('/checkout/payment/stripe');
      } else {
        // Handle COD checkout
        console.log('Processing COD order', body);
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
        <Fieldset>
          <ContactInformationSection register={register} />

          <ShippingInformationSection
            register={register}
            user={user}
            isNewAddress={isNewAddress}
            selectedAddressId={selectedAddressId}
            onAddressChange={handleAddressChange}
            onAddNewAddress={handleAddNewAddress}
          />

          <PaymentMethodSection
            setValue={setValue}
            watch={watch}
            paymentMethod={paymentMethod}
          />
        </Fieldset>
      </div>

      <div className='w-1/2'>
        <h3 className='text-lg font-semibold text-gray-600 mb-4'>
          Order summary
        </h3>
        
        <OrderSummarySection cart={cart} />

        <div className='px-6 pt-6'>
          <DiscountSection
            register={register}
            setValue={setValue}
            watch={watch}
            discounts={discounts}
            discountsLoading={discountsLoading}
            appliedDiscount={appliedDiscount}
            discountAmount={discountAmount}
            onApplyDiscount={onApplyDiscount}
            onRemoveDiscount={onRemoveDiscount}
          />

          <OrderTotalSection
            cart={cart}
            appliedDiscount={appliedDiscount}
            discountAmount={discountAmount}
          />
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
  );
};

export default CheckoutDetailOverview;
