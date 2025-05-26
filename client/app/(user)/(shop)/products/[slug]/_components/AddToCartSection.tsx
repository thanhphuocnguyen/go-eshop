'use client';
import { useState } from 'react';

interface AddToCartProps {
  price: number;
  inStock: boolean;
  onAddToCart?: (quantity: number) => void;
}

export const AddToCartSection = ({
  price,
  inStock = true,
  onAddToCart,
}: AddToCartProps) => {
  const [quantity, setQuantity] = useState(1);
  const [isWishlisted, setIsWishlisted] = useState(false);

  const increaseQuantity = () => {
    if (quantity < 10) setQuantity(quantity + 1);
  };

  const decreaseQuantity = () => {
    if (quantity > 1) setQuantity(quantity - 1);
  };

  const handleAddToCart = () => {
    if (inStock && onAddToCart) {
      onAddToCart(quantity);
    }
  };

  const toggleWishlist = () => {
    setIsWishlisted(!isWishlisted);
  };

  return (
    <div className='mt-8'>
      <div className='mb-5'>
        <label
          htmlFor='quantity'
          className='block text-sm font-medium text-gray-700 mb-1'
        >
          Quantity
        </label>
        <div className='flex items-center'>
          <button
            type='button'
            className={`rounded-l-md p-2 border border-gray-300 ${
              quantity <= 1
                ? 'bg-gray-100 text-gray-400'
                : 'bg-white text-gray-700 hover:bg-gray-50'
            }`}
            onClick={decreaseQuantity}
            disabled={quantity <= 1}
          >
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M20 12H4'
              />
            </svg>
          </button>

          <div className='w-12 text-center py-2 border-t border-b border-gray-300'>
            {quantity}
          </div>

          <button
            type='button'
            className={`rounded-r-md p-2 border border-gray-300 ${
              quantity >= 10
                ? 'bg-gray-100 text-gray-400'
                : 'bg-white text-gray-700 hover:bg-gray-50'
            }`}
            onClick={increaseQuantity}
            disabled={quantity >= 10}
          >
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M12 4v16m8-8H4'
              />
            </svg>
          </button>
        </div>
      </div>

      <div className='flex space-x-4'>
        <button
          onClick={handleAddToCart}
          disabled={!inStock}
          className={`flex-1 items-center justify-center rounded-md border border-transparent py-3 px-8 flex text-base font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 ${
            inStock
              ? 'bg-indigo-600 hover:bg-indigo-700'
              : 'bg-gray-400 cursor-not-allowed'
          }`}
        >
          <svg
            xmlns='http://www.w3.org/2000/svg'
            className='h-5 w-5 mr-2'
            fill='none'
            viewBox='0 0 24 24'
            stroke='currentColor'
          >
            <path
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth={2}
              d='M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z'
            />
          </svg>
          {inStock ? 'Add to cart' : 'Out of stock'}
        </button>

        <button
          onClick={toggleWishlist}
          className='rounded-md border border-gray-300 p-3 flex items-center justify-center text-gray-700 hover:bg-gray-50'
        >
          {isWishlisted ? (
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-6 w-6 text-red-500'
              viewBox='0 0 20 20'
              fill='currentColor'
            >
              <path
                fillRule='evenodd'
                d='M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z'
                clipRule='evenodd'
              />
            </svg>
          ) : (
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-6 w-6'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z'
              />
            </svg>
          )}
        </button>
      </div>
    </div>
  );
};
