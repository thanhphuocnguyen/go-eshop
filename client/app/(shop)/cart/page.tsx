import { apiFetch } from '@/lib/apis/api';
import { API_PATHS } from '@/lib/constants/api';
import { GenericResponse } from '@/lib/definitions';
import { CartModel } from '@/lib/definitions/cart';
import { Button } from '@headlessui/react';
import { CheckIcon } from '@heroicons/react/16/solid';
import Image from 'next/image';
import Link from 'next/link';

export default async function CartPage() {
  const { data, error } = await apiFetch<GenericResponse<CartModel>>(
    API_PATHS.CART,
    {
      method: 'GET',
    }
  );

  if (error) {
    return (
      <div>
        <h1>Cart</h1>
        <p>Error loading cart: {JSON.stringify(error)}</p>
      </div>
    );
  }

  return (
    <div className='container h-full flex flex-col items-center justify-center'>
      <h1 className='text-2xl font-bold text-center my-6'>Shopping Cart</h1>
      <div>
        <hr />
        {data.cart_items.map((e) => (
          <div className='pb-4 border-b border-gray-300' key={e.id}>
            <div className='flex items-center gap-4'>
              <div className='h-52 w-40 relative'>
                <Image
                  src={e.image_url || '/images/placeholder.webp'}
                  alt={e.name}
                  fill
                  objectFit='cover'
                  className='rounded-md border-lime-400'
                />
              </div>
              <div>
                <div className='flex flex-col justify-between'>
                  <div>
                    <div className='flex justify-between'>
                      <h2 className='text-lg font-semibold'>{e.name}</h2>
                      <div>${e.price}</div>
                    </div>
                    <div className='flex gap-2'>
                      {e.attributes.map((e) => (
                        <div key={e.name}>
                          <span key={e.name}>{e.name}: </span>
                          <span>{e.value}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                  <div className='flex'>
                    <div>
                      {e.stock > e.quantity ? (
                        <CheckIcon className='h-5 w-5 text-green-500' />
                      ) : null}
                      <span className='text-sm text-gray-500'>In stock</span>
                    </div>
                    <Button className={'text-indigo-600 font-medium'}>
                      Remove
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        ))}
        <div className='mt-6'>
          <div className='flex justify-between'>
            <div>Subtotal</div>
            <div>
              ${data.cart_items.reduce((acc, curr) => (curr.price += acc), 0)}
            </div>
          </div>
          <div className='text-gray-400'>
            Shipping and taxes will be calculated at checkout.
          </div>
          <Button className='mt-4 w-full bg-indigo-600 text-white py-2 rounded-md'>
            Checkout
          </Button>
          or <Link href='/shop'>Continue Shopping</Link>
        </div>
      </div>
    </div>
  );
}
