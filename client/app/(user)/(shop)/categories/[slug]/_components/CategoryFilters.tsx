'use client';

import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import {
  Disclosure,
  DisclosureButton,
  DisclosurePanel,
} from '@headlessui/react';
import { StarIcon } from '@heroicons/react/24/solid';
import { StarIcon as StarOutlineIcon } from '@heroicons/react/24/outline';
import { useState } from 'react';

interface CategoryFiltersProps {
  priceRange: {
    min: number;
    max: number;
  };
  minPrice?: number;
  maxPrice?: number;
  selectedRating?: number | null;
  categoryId: string;
}

export default function CategoryFilters({
  priceRange,
  minPrice: initialMinPrice,
  maxPrice: initialMaxPrice,
  selectedRating: initialRating,
}: CategoryFiltersProps) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  // Local state for UI interactions
  const [minPrice, setMinPrice] = useState<number>(
    initialMinPrice || priceRange.min
  );
  const [maxPrice, setMaxPrice] = useState<number>(
    initialMaxPrice || priceRange.max
  );
  const [selectedRating, setSelectedRating] = useState<number | null>(
    initialRating || null
  );
  const [filterOpen, setFilterOpen] = useState(false);

  // Apply filters when user clicks apply
  const applyFilters = () => {
    const params = new URLSearchParams(searchParams.toString());

    // Update or remove price filters
    if (minPrice !== priceRange.min) {
      params.set('minPrice', minPrice.toString());
    } else {
      params.delete('minPrice');
    }

    if (maxPrice !== priceRange.max) {
      params.set('maxPrice', maxPrice.toString());
    } else {
      params.delete('maxPrice');
    }

    // Update or remove rating filter
    if (selectedRating) {
      params.set('rating', selectedRating.toString());
    } else {
      params.delete('rating');
    }

    // Reset to page 1 when changing filters
    params.delete('page');

    // Update URL with new params
    router.push(`${pathname}?${params.toString()}`, { scroll: false });
  };

  // Reset all filters
  const resetFilters = () => {
    setMinPrice(priceRange.min);
    setMaxPrice(priceRange.max);
    setSelectedRating(null);

    router.push(pathname);
  };

  // For mobile toggle
  const toggleFilters = () => {
    setFilterOpen(!filterOpen);
  };

  // Check if any filters are active
  const hasActiveFilters =
    minPrice !== priceRange.min ||
    maxPrice !== priceRange.max ||
    selectedRating !== null;

  return (
    <>
      {/* Mobile filter button */}
      <div className='lg:hidden mb-4'>
        <button
          onClick={toggleFilters}
          className='w-full flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50'
        >
          <svg
            className='w-5 h-5 mr-2'
            fill='none'
            stroke='currentColor'
            viewBox='0 0 24 24'
          >
            <path
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth={2}
              d='M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4'
            />
          </svg>
          {filterOpen ? 'Hide Filters' : 'Show Filters'}
        </button>
      </div>

      {/* Filter sidebar - hidden on mobile unless toggled */}
      <div
        className={`lg:block ${filterOpen ? 'block' : 'hidden'} lg:w-64 w-full`}
      >
        <div className='bg-white rounded-lg shadow-sm p-5 border border-gray-100 sticky top-24'>
          <h3 className='text-lg font-medium text-gray-900 mb-4'>Filters</h3>

          {/* Price Range Filter */}
          <Disclosure defaultOpen>
            {({ open }) => (
              <div>
                <DisclosureButton className='flex justify-between w-full px-4 py-2 text-sm font-medium text-left text-gray-900 bg-gray-50 rounded-lg hover:bg-gray-100 focus:outline-none focus-visible:ring'>
                  <span>Price Range</span>
                  <svg
                    className={`${open ? 'transform rotate-180' : ''} w-4 h-4 text-gray-500`}
                    xmlns='http://www.w3.org/2000/svg'
                    viewBox='0 0 20 20'
                    fill='currentColor'
                  >
                    <path
                      fillRule='evenodd'
                      d='M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z'
                      clipRule='evenodd'
                    />
                  </svg>
                </DisclosureButton>
                <DisclosurePanel className='px-4 pt-4 pb-2'>
                  <div className='grid grid-cols-2 gap-4 mb-4'>
                    <div>
                      <label className='block text-sm font-medium text-gray-700'>
                        Min ($)
                      </label>
                      <input
                        type='number'
                        min={priceRange.min}
                        max={priceRange.max}
                        value={minPrice}
                        onChange={(e) =>
                          setMinPrice(
                            parseInt(e.target.value) || priceRange.min
                          )
                        }
                        className='mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm'
                      />
                    </div>
                    <div>
                      <label className='block text-sm font-medium text-gray-700'>
                        Max ($)
                      </label>
                      <input
                        type='number'
                        min={priceRange.min}
                        max={priceRange.max}
                        value={maxPrice}
                        onChange={(e) =>
                          setMaxPrice(
                            parseInt(e.target.value) || priceRange.max
                          )
                        }
                        className='mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm'
                      />
                    </div>
                  </div>
                  <input
                    type='range'
                    min={priceRange.min}
                    max={priceRange.max}
                    value={minPrice}
                    onChange={(e) => setMinPrice(parseInt(e.target.value))}
                    className='w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer'
                  />
                  <input
                    type='range'
                    min={priceRange.min}
                    max={priceRange.max}
                    value={maxPrice}
                    onChange={(e) => setMaxPrice(parseInt(e.target.value))}
                    className='w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer mt-2'
                  />
                </DisclosurePanel>
              </div>
            )}
          </Disclosure>

          {/* Rating Filter */}
          <div className='mt-6'>
            <h4 className='text-sm font-medium text-gray-900 mb-3'>Rating</h4>
            {[5, 4, 3, 2, 1].map((rating) => (
              <div key={rating} className='flex items-center mb-2'>
                <button
                  onClick={() =>
                    setSelectedRating(selectedRating === rating ? null : rating)
                  }
                  className={`flex items-center p-2 w-full rounded-md text-sm ${
                    selectedRating === rating
                      ? 'bg-indigo-50 text-indigo-700'
                      : 'text-gray-700 hover:bg-gray-50'
                  }`}
                >
                  {[1, 2, 3, 4, 5].map((star) =>
                    star <= rating ? (
                      <StarIcon
                        key={star}
                        className='h-5 w-5 text-yellow-400'
                      />
                    ) : (
                      <StarOutlineIcon
                        key={star}
                        className='h-5 w-5 text-gray-300'
                      />
                    )
                  )}
                  <span className='ml-2'>{rating === 1 ? '& Up' : `& Up`}</span>
                </button>
              </div>
            ))}
          </div>

          {/* Action Buttons */}
          <div className='mt-6 space-y-3'>
            <button
              onClick={applyFilters}
              className='w-full bg-indigo-600 text-white py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
            >
              Apply Filters
            </button>

            {hasActiveFilters && (
              <button
                onClick={resetFilters}
                className='w-full bg-white text-gray-700 py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
              >
                Reset Filters
              </button>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
