import Image from 'next/image';
import React from 'react';
import Link from 'next/link';

interface ProductCardProps {
  ID: number;
  slug: string;
  name: string;
  image?: string | null;
  priceFrom?: number | null;
  priceTo?: number | null;
  rating: number;
}
const ProductCard: React.FC<ProductCardProps> = (props) => {
  const { name, image, priceFrom, priceTo, rating } = props;
  return (
    <div className='w-full h-full p-3 flex flex-col max-w-sm bg-white border border-gray-200 rounded-lg shadow-sm dark:bg-gray-800 dark:border-gray-700'>
      <Link
        className='h-60 w-full relative'
        href='/products/[id]'
        as={`/products/${props.slug}`}
      >
        <Image
          src={image || '/images/product-placeholder.webp'}
          alt='product image'
          className='rounded-t-lg object-cover'
          sizes='(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw'
          fill
        />
      </Link>
      <div className='p-3 flex flex-col justify-between flex-grow'>
        <div className='flex flex-col w-full'>
          <Link href='/products/[id]' as={`/products/${props.slug}`}>
            <h5 className='text-base whitespace-nowrap text-nowrap mt-4 font-semibold tracking-tight text-gray-900 dark:text-white text-ellipsis'>
              {name}
            </h5>
          </Link>
          <div className='flex items-center mt-2.5 mb-2'>
            <div className='flex items-center space-x-1 rtl:space-x-reverse'>
            <svg
              className='w-4 h-4 text-yellow-700'
              aria-hidden='true'
              xmlns='http://www.w3.org/2000/svg'
              fill='currentColor'
              viewBox='0 0 22 20'
            >
              <path d='M20.924 7.625a1.523 1.523 0 0 0-1.238-1.044l-5.051-.734-2.259-4.577a1.534 1.534 0 0 0-2.752 0L7.365 5.847l-5.051.734A1.535 1.535 0 0 0 1.463 9.2l3.656 3.563-.863 5.031a1.532 1.532 0 0 0 2.226 1.616L11 17.033l4.518 2.375a1.534 1.534 0 0 0 2.226-1.617l-.863-5.03L20.537 9.2a1.523 1.523 0 0 0 .387-1.575Z' />
            </svg>
            <svg
              className='w-4 h-4 text-yellow-700'
              aria-hidden='true'
              xmlns='http://www.w3.org/2000/svg'
              fill='currentColor'
              viewBox='0 0 22 20'
            >
              <path d='M20.924 7.625a1.523 1.523 0 0 0-1.238-1.044l-5.051-.734-2.259-4.577a1.534 1.534 0 0 0-2.752 0L7.365 5.847l-5.051.734A1.535 1.535 0 0 0 1.463 9.2l3.656 3.563-.863 5.031a1.532 1.532 0 0 0 2.226 1.616L11 17.033l4.518 2.375a1.534 1.534 0 0 0 2.226-1.617l-.863-5.03L20.537 9.2a1.523 1.523 0 0 0 .387-1.575Z' />
            </svg>
            <svg
              className='w-4 h-4 text-yellow-700'
              aria-hidden='true'
              xmlns='http://www.w3.org/2000/svg'
              fill='currentColor'
              viewBox='0 0 22 20'
            >
              <path d='M20.924 7.625a1.523 1.523 0 0 0-1.238-1.044l-5.051-.734-2.259-4.577a1.534 1.534 0 0 0-2.752 0L7.365 5.847l-5.051.734A1.535 1.535 0 0 0 1.463 9.2l3.656 3.563-.863 5.031a1.532 1.532 0 0 0 2.226 1.616L11 17.033l4.518 2.375a1.534 1.534 0 0 0 2.226-1.617l-.863-5.03L20.537 9.2a1.523 1.523 0 0 0 .387-1.575Z' />
            </svg>
            <svg
              className='w-4 h-4 text-yellow-700'
              aria-hidden='true'
              xmlns='http://www.w3.org/2000/svg'
              fill='currentColor'
              viewBox='0 0 22 20'
            >
              <path d='M20.924 7.625a1.523 1.523 0 0 0-1.238-1.044l-5.051-.734-2.259-4.577a1.534 1.534 0 0 0-2.752 0L7.365 5.847l-5.051.734A1.535 1.535 0 0 0 1.463 9.2l3.656 3.563-.863 5.031a1.532 1.532 0 0 0 2.226 1.616L11 17.033l4.518 2.375a1.534 1.534 0 0 0 2.226-1.617l-.863-5.03L20.537 9.2a1.523 1.523 0 0 0 .387-1.575Z' />
            </svg>
            <svg
              className='w-4 h-4 text-gray-200 dark:text-gray-600'
              aria-hidden='true'
              xmlns='http://www.w3.org/2000/svg'
              fill='currentColor'
              viewBox='0 0 22 20'
            >
              <path d='M20.924 7.625a1.523 1.523 0 0 0-1.238-1.044l-5.051-.734-2.259-4.577a1.534 1.534 0 0 0-2.752 0L7.365 5.847l-5.051.734A1.535 1.535 0 0 0 1.463 9.2l3.656 3.563-.863 5.031a1.532 1.532 0 0 0 2.226 1.616L11 17.033l4.518 2.375a1.534 1.534 0 0 0 2.226-1.617l-.863-5.03L20.537 9.2a1.523 1.523 0 0 0 .387-1.575Z' />
            </svg>
          </div>
          <span className='bg-blue-100 rounded-md text-blue-800 text-xs font-semibold px-2.5 py-0.5 dark:bg-blue-200 dark:text-blue-800 ms-3'>
            {rating.toFixed(1)}
          </span>
        </div>
        <div className='flex items-center justify-between mb-4'>
          <span className='text-xl font-bold text-gray-900 dark:text-white'>
            ${priceFrom} - ${priceTo}
          </span>
        </div>
        </div>
        <div className='mt-auto'>
          <Link
            href='/products/[id]'
            as={`/products/${props.ID}`}
            className='w-full text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800'
          >
            Add to cart
          </Link>
        </div>
      </div>
    </div>
  );
};

export default ProductCard;
