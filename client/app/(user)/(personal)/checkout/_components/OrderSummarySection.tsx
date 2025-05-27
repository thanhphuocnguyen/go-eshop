'use client';

import { Button } from '@headlessui/react';
import { TrashIcon } from '@heroicons/react/24/solid';
import Image from 'next/image';
import { CartModel } from '@/app/lib/definitions';

interface OrderSummarySectionProps {
  cart: CartModel | undefined;
}

export const OrderSummarySection: React.FC<OrderSummarySectionProps> = ({
  cart,
}) => {
  if (!cart?.cartItems?.length) {
    return (
      <div className='border border-gray-200 bg-white rounded-md p-6'>
        <p className='text-gray-500 text-center'>Your cart is empty</p>
      </div>
    );
  }

  return (
    <div className='border border-gray-200 bg-white rounded-md'>
      {cart.cartItems.map((item) => (
        <div
          key={item.variantId}
          className='flex gap-4 p-6 border-b border-gray-200'
        >
          <div className='h-28 w-24 relative'>
            <Image
              fill
              objectFit='cover'
              src={item.imageUrl ?? '/images/logos/logo.webp'}
              alt='Product Image'
              className='rounded-md border border-lime-300'
            />
          </div>
          <div
            key={item.variantId}
            className='flex-1 flex flex-col justify-between'
          >
            <div>
              <div className='flex justify-between'>
                <span className='font-medium'>{item.name}</span>
                <Button>
                  <TrashIcon className='size-6 text-red-200' />
                </Button>
              </div>
              <div className='flex flex-col gap-1 mt-2'>
                {item.attributes.map((attribute) => (
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
              <span>${item.price}</span>
              <span className='text-gray-500'>Qty: {item.quantity}</span>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
