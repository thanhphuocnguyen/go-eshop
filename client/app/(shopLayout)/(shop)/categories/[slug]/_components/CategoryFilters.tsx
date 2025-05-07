'use client';

import { FunnelIcon, StarIcon } from '@heroicons/react/24/solid';

interface CategoryFiltersProps {
  priceRange: { min: number; max: number };
  minPrice: number;
  maxPrice: number;
  selectedRating: number | null;
  setMinPrice: (price: number) => void;
  setMaxPrice: (price: number) => void;
  setSelectedRating: (rating: number | null) => void;
  resetFilters: () => void;
  filterOpen: boolean;
  toggleFilters: () => void;
}

export default function CategoryFilters({
  priceRange,
  minPrice,
  maxPrice,
  selectedRating,
  setMinPrice,
  setMaxPrice,
  setSelectedRating,
  resetFilters,
  filterOpen,
  toggleFilters,
}: CategoryFiltersProps) {
  return (
    <>
      {/* Mobile Filters Button */}
      <div className='block lg:hidden mb-4'>
        <button
          onClick={toggleFilters}
          className='flex items-center justify-center w-full py-2 px-4 border border-gray-300 bg-white text-gray-700 rounded-md shadow-sm hover:bg-gray-50'
        >
          <FunnelIcon className='h-5 w-5 mr-2' />
          {filterOpen ? 'Hide Filters' : 'Show Filters'}
        </button>
      </div>

      {/* Filter sidebar - hidden on mobile by default */}
      <div className={`lg:block lg:w-1/4 ${filterOpen ? 'block' : 'hidden'}`}>
        <div className='sticky top-4 bg-white p-4 border border-gray-200 rounded-lg shadow-sm'>
          <div className='flex justify-between items-center mb-4'>
            <h2 className='text-lg font-semibold text-gray-800'>Filters</h2>
            <button
              onClick={resetFilters}
              className='text-sm text-indigo-600 hover:text-indigo-800'
            >
              Reset All
            </button>
          </div>

          {/* Price Range Filter */}
          <div className='mb-6'>
            <h3 className='font-medium text-gray-700 mb-2'>Price Range</h3>
            <div className='flex items-center gap-2 mb-4'>
              <div className='w-1/2'>
                <label className='block text-sm text-gray-600 mb-1'>Min</label>
                <input
                  type='number'
                  value={minPrice}
                  onChange={(e) => setMinPrice(Number(e.target.value))}
                  className='w-full p-2 border border-gray-300 rounded-md'
                />
              </div>
              <div className='w-1/2'>
                <label className='block text-sm text-gray-600 mb-1'>Max</label>
                <input
                  type='number'
                  value={maxPrice}
                  onChange={(e) => setMaxPrice(Number(e.target.value))}
                  className='w-full p-2 border border-gray-300 rounded-md'
                />
              </div>
            </div>
            <input
              type='range'
              min={priceRange.min}
              max={priceRange.max}
              value={minPrice}
              onChange={(e) => setMinPrice(Number(e.target.value))}
              className='w-full mb-2'
            />
            <input
              type='range'
              min={priceRange.min}
              max={priceRange.max}
              value={maxPrice}
              onChange={(e) => setMaxPrice(Number(e.target.value))}
              className='w-full'
            />
          </div>

          {/* Rating Filter */}
          <div className='mb-6'>
            <h3 className='font-medium text-gray-700 mb-2'>Rating</h3>
            {[5, 4, 3, 2, 1].map((rating) => (
              <div key={rating} className='flex items-center mb-2'>
                <input
                  type='radio'
                  id={`rating-${rating}`}
                  name='rating'
                  checked={selectedRating === rating}
                  onChange={() => setSelectedRating(rating)}
                  className='mr-2'
                />
                <label
                  htmlFor={`rating-${rating}`}
                  className='flex items-center cursor-pointer'
                >
                  {Array.from({ length: 5 }).map((_, index) => (
                    <StarIcon
                      key={index}
                      className={`h-4 w-4 ${
                        index < rating ? 'text-yellow-400' : 'text-gray-300'
                      }`}
                    />
                  ))}
                  <span className='ml-2 text-sm text-gray-600'>& Up</span>
                </label>
              </div>
            ))}
            {selectedRating && (
              <button
                onClick={() => setSelectedRating(null)}
                className='text-xs text-indigo-600 hover:text-indigo-800 mt-1'
              >
                Clear Rating Filter
              </button>
            )}
          </div>

          {/* Availability Filter */}
          <div className='mb-6'>
            <h3 className='font-medium text-gray-700 mb-2'>Availability</h3>
            <div className='flex items-center mb-2'>
              <input type='checkbox' id='in-stock' className='mr-2' />
              <label
                htmlFor='in-stock'
                className='text-sm text-gray-600 cursor-pointer'
              >
                In Stock
              </label>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
