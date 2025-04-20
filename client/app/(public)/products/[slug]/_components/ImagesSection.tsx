'use client';
import { ProductImageModel } from '@/lib/definitions';
import { Button } from '@headlessui/react';
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
        <div className='grid grid-cols-4 gap-6' aria-orientation='horizontal'>
          {images.map((image) => (
            <Button
              key={image.id}
              className={`relative flex items-center justify-center h-72 rounded-md bg-white text-sm font-medium uppercase text-gray-900 cursor-pointer hover:bg-gray-50 focus:outline-none focus:ring focus:ring-offset-4 focus:ring-opacity-50 ${image.role === 'main' ? 'ring-2 ring-offset-2 ring-indigo-500' : ''}`}
              onClick={() => setMainImage(image)}
            >
              <span className='absolute inset-0 rounded-md overflow-hidden'>
                <Image
                  fill
                  src={image.url}
                  objectFit='cover'
                  objectPosition='center'
                  alt={'Thumbnail image'}
                />
              </span>
              {/* Selected ring */}
              <span
                className={`${image.role === 'main' ? 'ring-indigo-500' : 'ring-transparent'} absolute inset-0 rounded-md ring-2 ring-offset-2 pointer-events-none`}
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
            ></div>
          ))}
        </div>
      </div>

      {/* Main Image */}
      <div className='w-full aspect-w-1 aspect-h-1'>
        {mainImage && (
          <Image
            width={1000}
            height={400}
            src={mainImage.url}
            alt={'Model wearing Basic Tee in black.'}
            objectFit='cover'
            objectPosition='center'
            className='rounded-lg shadow-sm'
          />
        )}
      </div>
    </div>
  );
};
