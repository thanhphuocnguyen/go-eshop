'use client';
import { useAppUser } from '@/components/AppUserContext';
import LoadingInline from '@/components/Common/Loadings/LoadingInline';
import { Button } from '@headlessui/react';
import { CheckIcon, TrashIcon } from '@heroicons/react/16/solid';
import Image from 'next/image';
import Link from 'next/link';

export default function CartPage() {
  const { cart, cartLoading, updateCartItemQuantity, removeFromCart } =
    useAppUser();

  if (cartLoading && !cart) {
    return (
      <div className='flex justify-center items-center h-screen'>
        <LoadingInline />
      </div>
    );
  }

  if (!cart) {
    return (
      <div className='flex justify-center items-center h-screen'>
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
    <div className='mx-auto flex flex-col'>
      <div className='w-[768px] m-auto py-14'>
        <div className='w-full'>
          <h1 className='text-3xl font-bold text-center mb-12'>
            Shopping Cart
          </h1>
          <hr />

          {cart.cart_items.map((e) => (
            <div className='mt-6 pb-4 border-b border-gray-300' key={e.id}>
              <div className='flex gap-4'>
                <div className='relative'>
                  <Image
                    src={e.image_url || '/images/placeholder.webp'}
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
                      <span className='text-sm text-gray-500 mr-4'>
                        In stock
                      </span>
                      <div className='flex gap-2 items-center text-sm text-gray-500'>
                        <Button
                          onClick={() => {
                            if (e.quantity > 1) {
                              updateCartItemQuantity(e.id, e.quantity - 1);
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
                            updateCartItemQuantity(e.id, e.quantity + 1);
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
                {cart.cart_items.reduce(
                  (acc, curr) => (acc += curr.price * curr.quantity),
                  0
                )}
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
          </div>
        </div>
      </div>
      <div className='border-t border-gray-300 m-auto w-full py-32 bg-gray-100 shadow-sm px-20 flex gap-16'>
        <div className='flex flex-col items-center'>
          {/* Image */}
          <Image
            src={'/images/logos/icon-returns-light.svg'}
            alt='Returns'
            width={100}
            height={100}
          />
          <div className='text-lg mb-3 font-bold'>Free returns</div>
          <div className='text-center'>
            Not what you expected? Place it back in the parcel and attach the
            pre-paid postage stamp.
          </div>
        </div>
        <div className='flex flex-col items-center'>
          {/* Image */}
          <Image
            src={'/images/logos/icon-calendar-light.svg'}
            alt='Returns'
            width={100}
            height={100}
          />
          <div className='text-lg mb-3 font-bold'>Same day delivery</div>
          <div className='text-center'>
            We offer a delivery service that has never been done before.
            Checkout today and receive your products within hours.
          </div>
        </div>
        <div className='flex flex-col items-center'>
          {/* Image */}
          <Image
            src={'/images/logos/icon-gift-card-light.svg'}
            alt='Returns'
            width={100}
            height={100}
          />
          <div className='text-lg mb-3 font-bold'>All year discount</div>
          <div className='text-center'>
            Looking for a deal? You can use the code &quot;ALLYEAR &quot; at
            checkout and get money off all year round.
          </div>
        </div>
        <div className='flex mb-3 flex-col items-center'>
          {/* Image */}
          <Image
            src={'/images/logos/icon-planet-light.svg'}
            alt='Returns'
            width={100}
            height={100}
          />
          <div className='text-lg font-bold'>For the planet</div>
          <div className='text-center'>
            We’ve pledged 1% of sales to the preservation and restoration of the
            natural environment.
          </div>
        </div>
      </div>
    </div>
  );
}

//
//

//
//

//
//

//
//
