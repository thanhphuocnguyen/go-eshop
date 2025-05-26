'use client';

import LoadingInline from '@/app/components/Common/Loadings/LoadingInline';
import { Button } from '@headlessui/react';
import { CheckIcon, TrashIcon } from '@heroicons/react/16/solid';
import Image from 'next/image';
import Link from 'next/link';
import { useCartCtx } from '@/app/lib/contexts/CartContext';

export const CartItemSection = () => {
  const { cart, cartLoading, updateCartItemQuantity, removeFromCart } =
    useCartCtx();

  if (cartLoading && !cart) {
    return (
      <div className='flex justify-center items-center h-screen'>
        <LoadingInline />
      </div>
    );
  }
  if (!cart || !cart.cartItems.length) {
    return (
      <div className='flex justify-center my-20 items-center h-full'>
        <div className='text-center'>
          <h1 className='text-2xl font-bold'>Your cart is empty</h1>
          <Link href='/shop' className='mt-4 text-indigo-600 hover:underline'>
            Continue Shopping
          </Link>
        </div>
      </div>
    );
  }

  return (
    <>
      {cart.cartItems.map((e) => (
        <div className='mt-6 pb-4 border-b border-gray-300' key={e.id}>
          <div className='flex gap-4'>
            <div className='relative'>
              <Image
                src={e.imageUrl || '/images/placeholder.webp'}
                alt={e.name}
                height={100}
                width={100}
                className='rounded-md object-cover border border-lime-400'
              />
            </div>
            <div className='flex flex-col justify-between w-full'>
              <div>
                <div className='flex justify-between'>
                  <h2 className='text-lg font-semibold'>{e.name}</h2>
                  <div>${e.price.toFixed(2)}</div>
                </div>
                <div className='mt-2 flex flex-col gap-2'>
                  {e.attributes.map((e) => (
                    <div className='text-sm text-gray-400' key={e.name}>
                      <span key={e.name}>{e.name}: </span>
                      <span>{e.value}</span>
                    </div>
                  ))}
                </div>
              </div>
              <div className='flex justify-between w-full'>
                <div className='flex items-center gap-2'>
                  <span>
                    {e.stock > e.quantity ? (
                      <CheckIcon className='h-5 w-5 text-green-500' />
                    ) : null}
                  </span>
                  <span className='text-sm text-gray-500 mr-4'>In stock</span>
                  <div className='flex gap-2 items-center text-sm text-gray-500'>
                    <Button
                      onClick={() => {
                        if (e.quantity > 1) {
                          updateCartItemQuantity(e.id, -1);
                        } else {
                          removeFromCart(e.id);
                        }
                      }}
                      className='bg-gray-200 rounded-md px-2 py-1'
                    >
                      <span className='text-gray-500'>-</span>
                    </Button>
                    <span>{e.quantity}</span>
                    <Button
                      onClick={() => {
                        updateCartItemQuantity(e.id, 1);
                      }}
                      className='bg-gray-200 rounded-md px-2 py-1'
                    >
                      <span className='text-gray-500'>+</span>
                    </Button>
                  </div>
                </div>
                <Button
                  onClick={() => {
                    removeFromCart(e.id);
                  }}
                  className={
                    'text-indigo-600 font-medium flex gap-1 items-center hover:text-red-500'
                  }
                >
                  <TrashIcon className='size-5 text-red-200' />
                  Remove
                </Button>
              </div>
            </div>
          </div>
        </div>
      ))}

      <div className='mt-6'>
        <div className='flex justify-between'>
          <div>Subtotal</div>
          <div>
            $
            {cart.cartItems.reduce(
              (acc, curr) => (acc += curr.price * curr.quantity),
              0
            )}
          </div>
        </div>
      </div>

      <div className='text-gray-400'>
        Shipping and taxes will be calculated at checkout.
      </div>
      <Link
        href={'/checkout'}
        className='mt-4 w-full block text-center bg-indigo-600 text-white py-3 rounded-md'
      >
        Checkout
      </Link>
      <div className='mt-2'>
        <span className='mr-2'>or</span>
        <Link href='/shop' className='text-indigo-600'>
          Continue Shopping
        </Link>
      </div>
    </>
  );
};
