'use client';

import React, { useState } from 'react';
import { ProductImageModel } from '@/app/lib/definitions';
import { Button } from '@headlessui/react';
import clsx from 'clsx';
import Image from 'next/image';

import {
  MagnifyingGlassIcon,
  ArrowsPointingOutIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from '@heroicons/react/24/outline';

interface ImagesSectionProps {
  images: ProductImageModel[];
}
export const ImagesSection: React.FC<ImagesSectionProps> = ({ images }) => {
  const [mainImage, setMainImage] = useState<ProductImageModel | null>(
    images[0] || null
  );
  const [isZoomed, setIsZoomed] = useState(false);
  const [currentIndex, setCurrentIndex] = useState(0);

  const handleZoom = () => {
    setIsZoomed(!isZoomed);
  };

  const nextImage = () => {
    const nextIndex = (currentIndex + 1) % images.length;
    setCurrentIndex(nextIndex);
    setMainImage(images[nextIndex]);
  };

  const prevImage = () => {
    const prevIndex = (currentIndex - 1 + images.length) % images.length;
    setCurrentIndex(prevIndex);
    setMainImage(images[prevIndex]);
  };

  const handleImageClick = (image: ProductImageModel, index: number) => {
    setMainImage(image);
    setCurrentIndex(index);
  };

  return (
    <div className='col-span-2 flex flex-col-reverse'>
      {/* Image selector */}
      <div className='hidden mt-6 w-full max-w-2xl mx-auto sm:block lg:max-w-none'>
        <div className='grid grid-cols-5 gap-6' aria-orientation='horizontal'>
          {images.map((image, index) => (
            <Button
              key={image.id}
              className={clsx(
                'relative flex items-center justify-center h-56 rounded-md bg-white text-sm font-medium uppercase text-gray-900 cursor-pointer',
                'hover:bg-gray-50 focus:outline-none focus:ring focus:ring-offset-4 focus:ring-opacity-50',
                mainImage?.id === image.id
                  ? 'ring-2 ring-offset-2 ring-indigo-500'
                  : ''
              )}
              onClick={() => handleImageClick(image, index)}
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
                  mainImage?.id === image.id
                    ? 'ring-indigo-500'
                    : 'ring-transparent'
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
      <div className='group relative'>
        <div
          className={`relative overflow-hidden rounded-lg ${isZoomed ? 'cursor-zoom-out' : 'cursor-zoom-in'} h-[400px] sm:h-[500px] lg:h-[650px] w-full`}
        >
          {mainImage && (
            <Image
              fill
              priority
              src={mainImage.url}
              alt={'Product image'}
              className={`object-contain rounded-md shadow-lg transition-transform duration-500 ease-out ${isZoomed ? 'scale-150' : 'scale-100'}`}
              sizes='(max-width: 768px) 100vw, 100vw'
              onClick={handleZoom}
            />
          )}

          {/* Zoom button */}
          <button
            className='absolute top-4 right-4 bg-white bg-opacity-75 p-2 rounded-full shadow-md z-10 opacity-0 group-hover:opacity-100 transition-opacity duration-200'
            onClick={handleZoom}
          >
            {isZoomed ? (
              <MagnifyingGlassIcon className='h-5 w-5 text-gray-700' />
            ) : (
              <ArrowsPointingOutIcon className='h-5 w-5 text-gray-700' />
            )}
          </button>

          {/* Navigation arrows - only visible when hovering and if there are multiple images */}
          {images.length > 1 && (
            <>
              <button
                className='absolute left-4 top-1/2 transform -translate-y-1/2 bg-white bg-opacity-75 p-2 rounded-full shadow-md z-10 opacity-0 group-hover:opacity-100 transition-opacity duration-200'
                onClick={prevImage}
              >
                <ChevronLeftIcon className='h-5 w-5 text-gray-700' />
              </button>
              <button
                className='absolute right-4 top-1/2 transform -translate-y-1/2 bg-white bg-opacity-75 p-2 rounded-full shadow-md z-10 opacity-0 group-hover:opacity-100 transition-opacity duration-200'
                onClick={nextImage}
              >
                <ChevronRightIcon className='h-5 w-5 text-gray-700' />
              </button>
            </>
          )}
        </div>

        {/* Mobile thumbnails slider */}
        <div className='flex overflow-x-auto gap-3 mt-4 pb-2 sm:hidden'>
          {images.map((image, index) => (
            <button
              key={image.id}
              className={clsx(
                'flex-shrink-0 w-20 h-20 rounded-md overflow-hidden border-2',
                mainImage?.id === image.id
                  ? 'border-indigo-500'
                  : 'border-transparent'
              )}
              onClick={() => handleImageClick(image, index)}
            >
              <Image
                width={80}
                height={80}
                src={image.url}
                alt={'Thumbnail'}
                className='object-cover w-full h-full'
              />
            </button>
          ))}
        </div>
      </div>
    </div>
  );
};

interface ImagesSectionProps {
  images: ProductImageModel[];
}
