// 'use client';
import { apiFetch } from '@/lib/api/api';
import { API_PATHS } from '@/lib/constants/api';
import { GenericResponse, ProductDetailModel } from '@/lib/definitions';
import {
  CurrencyDollarIcon,
  GlobeAsiaAustraliaIcon,
  StarIcon as SolidStarIcon,
} from '@heroicons/react/16/solid';
import { StarIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';
import React from 'react';
import AddToCart from './_components/AddToCart';
// Example using react-icons, adjust if you use a different library

// Helper to render stars based on rating
// const renderStars = (rating) => {
//   const stars = [];
//   for (let i = 1; i <= 5; i++) {
//     if (i <= Math.floor(rating)) {
//       stars.push(<SolidStarIcon key={i} className='text-yellow-400' />);
//     } else if (i === Math.ceil(rating) && !Number.isInteger(rating)) {
//       stars.push(<SolidStarIcon key={i} className='text-yellow-400' />);
//     } else {
//       stars.push(<StarIcon key={i} className='text-yellow-400' />);
//     }
//   }
//   return stars;
// };

async function ProductDetailPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;

  const { data, error } = await apiFetch<GenericResponse<ProductDetailModel>>(
    API_PATHS.PRODUCT_DETAIL.replace(':id', slug),
    {
      nextOptions: {
        next: {
          tags: ['product'],
        },
      },
    }
  );
  if (error) {
    console.error(error);
    return <div>Error loading product details</div>;
  }
  if (!data) {
    return <div>Loading...</div>;
  }

  return (
    <div className='container mx-auto px-4 py-8'>
      <div className='lg:grid lg:grid-cols-2 lg:gap-x-8 lg:items-start'>
        {/* Image gallery */}
        <div className='flex flex-col-reverse'>
          {/* Image selector */}
          <div className='hidden mt-6 w-full max-w-2xl mx-auto sm:block lg:max-w-none'>
            <div
              className='grid grid-cols-4 gap-6'
              aria-orientation='horizontal'
            >
              {data.images.map((image) => (
                <button
                  key={image.id}
                  className={`relative flex items-center justify-center h-24 rounded-md bg-white text-sm font-medium uppercase text-gray-900 cursor-pointer hover:bg-gray-50 focus:outline-none focus:ring focus:ring-offset-4 focus:ring-opacity-50 ${image.role === 'main' ? 'ring-2 ring-offset-2 ring-indigo-500' : ''}`}
                  // onClick={() => setMainImage(image)}
                >
                  <span className='absolute inset-0 rounded-md overflow-hidden'>
                    <Image
                      width={100}
                      height={100}
                      src={image.url}
                      alt={'Thumbnail image'}
                      className='w-full h-full object-center object-cover'
                    />
                  </span>
                  {/* Selected ring */}
                  <span
                    className={`${image.role === 'main' ? 'ring-indigo-500' : 'ring-transparent'} absolute inset-0 rounded-md ring-2 ring-offset-2 pointer-events-none`}
                    aria-hidden='true'
                  />
                </button>
              ))}
              {/* Add placeholder boxes if fewer than 4 images */}
              {Array.from({
                length: Math.max(0, 4 - data.images.length),
              }).map((_, idx) => (
                <div
                  key={`placeholder-${idx}`}
                  className='relative h-24 rounded-md bg-gray-100'
                ></div>
              ))}
            </div>
          </div>

          {/* Main Image */}
          <div className='w-full aspect-w-1 aspect-h-1'>
            {data.images[0] && (
              <Image
                width={1000}
                height={1000}
                src={data.images[0].url}
                alt={'Model wearing Basic Tee in black.'}
                className='w-full h-full object-center object-cover rounded-lg shadow-sm'
              />
            )}
          </div>
        </div>

        {/* Product info */}
        <div className='mt-10 px-4 sm:px-0 sm:mt-16 lg:mt-0'>
          <h1 className='text-3xl font-extrabold tracking-tight text-gray-900'>
            {data.name}
          </h1>

          <div className='mt-3'>
            <h2 className='sr-only'>Product information</h2>
            <p className='text-3xl text-gray-900'>${data.price}</p>
          </div>

          {/* Reviews */}
          <div className='mt-3'>
            <h3 className='sr-only'>Reviews</h3>
            <div className='flex items-center'>
              <div className='flex items-center'>
                {/* {renderStars(data.rating)} */}
              </div>
              <p className='sr-only'>{5} out of 5 stars</p>
              <a
                // href={data.href}
                className='ml-3 text-sm font-medium text-indigo-600 hover:text-indigo-500'
              >
                See all {0} reviews
              </a>
            </div>
          </div>

          <div className='mt-6'>
            <h3 className='sr-only'>Description</h3>
            <div
              className='text-base text-gray-700 space-y-6'
              dangerouslySetInnerHTML={{ __html: data.description }} // Use only if description is trusted HTML
              // Or just: <p className="text-base text-gray-700">{data.description}</p>
            />
          </div>

          <form className='mt-6'>
            {/* Colors */}
            <div>
              <h3 className='text-sm text-gray-900 font-medium'>Color</h3>
              <div className='flex items-center space-x-3 mt-2'>
                {/* {data.colors.map((color) => (
                  <button
                    key={color.name}
                    type='button' // Important: prevent form submission
                    className={`relative -m-0.5 flex items-center justify-center rounded-full p-0.5 focus:outline-none ${
                      selectedColor.name === color.name
                        ? `ring-2 ring-offset-1 ${color.selectedClass}`
                        : ''
                    }`}
                    onClick={() => setSelectedColor(color)}
                    aria-label={color.name}
                  >
                    <span
                      aria-hidden='true'
                      className={`h-8 w-8 ${color.class} border border-black border-opacity-10 rounded-full`}
                    />
                  </button>
                ))} */}
              </div>
            </div>

            {/* Sizes */}
            <div className='mt-10'>
              <div className='flex items-center justify-between'>
                <h3 className='text-sm text-gray-900 font-medium'>Size</h3>
                <a
                  href='#'
                  className='text-sm font-medium text-indigo-600 hover:text-indigo-500'
                >
                  See sizing chart
                </a>
              </div>

              <div className='grid grid-cols-4 gap-4 sm:grid-cols-8 lg:grid-cols-4 mt-4'></div>
            </div>

            <AddToCart />
          </form>
          <div
            className='mt-6 text-base text-gray-700 space-y-4'
            dangerouslySetInnerHTML={{
              __html: data.description,
            }}
          ></div>
          {/* Details/Fabric & Care */}
          <div className='mt-10 pt-10 border-t border-gray-200'>
            <h3 className='text-sm font-medium text-gray-900'>Fabric & Care</h3>
            <div className='mt-4 prose prose-sm text-gray-500'>
              <ul role='list'>
                {/* {data.details.map((item) => (
                  <li key={item}>{item}</li>
                ))} */}
              </ul>
            </div>
          </div>

          {/* Info Boxes */}
          <div className='mt-8 grid grid-cols-1 gap-y-8 sm:grid-cols-2 sm:gap-x-6 lg:grid-cols-1 xl:grid-cols-2'>
            {/* International Delivery */}
            <div className='border border-gray-200 rounded-lg p-6 text-center'>
              <div className='flex items-center justify-center text-gray-400 mb-2'>
                <GlobeAsiaAustraliaIcon className='siz-6' />
              </div>
              <p className='text-sm font-medium text-gray-900'>
                International delivery
              </p>
              <p className='mt-1 text-sm text-gray-500'>
                Get your order in 2 years
              </p>{' '}
              {/* Update text as needed */}
            </div>
            {/* Loyalty Rewards */}
            <div className='border border-gray-200 rounded-lg p-6 text-center'>
              <div className='flex items-center justify-center text-gray-400 mb-2'>
                <CurrencyDollarIcon className='size-6' />
              </div>
              <p className='text-sm font-medium text-gray-900'>
                Loyalty rewards
              </p>
              <p className='mt-1 text-sm text-gray-500'>
                Don't look at other tees
              </p>{' '}
              {/* Update text as needed */}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default ProductDetailPage;
