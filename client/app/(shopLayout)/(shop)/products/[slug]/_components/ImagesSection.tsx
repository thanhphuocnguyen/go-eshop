'use client';
import { ProductImageModel } from '@/app/lib/definitions';
import { Button } from '@headlessui/react';
import clsx from 'clsx';
import Image from 'next/image';
import React from 'react';

interface ImagesSectionProps {
  images: ProductImageModel[];
}
export const ImagesSection: React.FC<ImagesSectionProps> = ({ images }) => {
  const [mainImage, setMainImage] = React.useState<ProductImageModel | null>(
    images[0] || null
  );

  return (
    <div className='col-span-2 flex flex-col-reverse'>
      {/* Image selector */}
      <div className='hidden mt-6 w-full max-w-2xl mx-auto sm:block lg:max-w-none'>
        <div className='grid grid-cols-5 gap-6' aria-orientation='horizontal'>
          {images.map((image) => (
            <Button
              key={image.id}
              className={clsx(
                'relative flex items-center justify-center h-56 rounded-md bg-white text-sm font-medium uppercase text-gray-900 cursor-pointer',
                'hover:bg-gray-50 focus:outline-none focus:ring focus:ring-offset-4 focus:ring-opacity-50',
                image.role === 'main'
                  ? 'ring-2 ring-offset-2 ring-indigo-500'
                  : ''
              )}
              onClick={() => setMainImage(image)}
            >
              <span className='absolute inset-0 rounded-md overflow-hidden'>
                <Image
                  fill
                  src={image.url}
                  alt={'Thumbnail image'}
                  className='object-cover w-full h-full'
                  sizes='(max-width: 768px) 96px, 96px'
                />
              </span>
              {/* Selected ring */}
              <span
                className={clsx(
                  'absolute inset-0 rounded-md ring-2 ring-offset-2 pointer-events-none',
                  image.role === 'main' ? 'ring-indigo-500' : 'ring-transparent'
                )}
                aria-hidden='true'
              />
            </Button>
          ))}
          {/* Add placeholder boxes if fewer than 4 images */}
          {Array.from({
            length: Math.max(0, 4 - images.length),
          }).map((_, idx) => (
            <div
              key={`placeholder-${idx}`}
              className='relative h-full rounded-md bg-gray-100'
            />
          ))}
        </div>
      </div>

      {/* Main Image */}
      <div className='flex relative justify-center h-[400px] sm:h-[500px] lg:h-[1000px] w-full'>
        {mainImage && (
          <Image
            fill
            priority
            src={mainImage.url}
            alt={'Model wearing Basic Tee in black.'}
            className='object-cover rounded-md shadow-lg'
            sizes='(max-width: 768px) 100vw, 100vw'
          />
        )}
      </div>
    </div>
  );
};
