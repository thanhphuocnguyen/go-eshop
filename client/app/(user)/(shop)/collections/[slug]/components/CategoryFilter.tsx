'use client';

import {
  Checkbox,
  Disclosure,
  DisclosureButton,
  DisclosurePanel,
} from '@headlessui/react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import { FilterOption } from '../../lib/definitions';

export function CheckboxGroup({
  categories,
  attributes,
  brands,
}: {
  categories: FilterOption[];
  brands: FilterOption[];
  attributes: Record<string, FilterOption[]>;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  // Get current filters from URL
  const currentCategory = searchParams.get('category');
  const currentBrand = searchParams.get('brand');
  const currentAttribute = searchParams.get('attribute');

  // State to track checked filters
  const [checkedCategories, setCheckedCategories] = useState<string[]>(
    currentCategory ? [currentCategory] : []
  );
  const [checkedBrands, setCheckedBrands] = useState<string[]>(
    currentBrand ? [currentBrand] : []
  );
  const [checkedAttributes, setCheckedAttributes] = useState<string[]>(
    currentAttribute ? [currentAttribute] : []
  );

  // Update URL when filters change
  useEffect(() => {
    const params = new URLSearchParams(searchParams);

    // Handle category filter
    params.delete('category');
    if (checkedCategories.length > 0) {
      params.set('category', checkedCategories[0]);
    }

    // Handle brand filter
    params.delete('brand');
    if (checkedBrands.length > 0) {
      params.set('brand', checkedBrands[0]);
    }

    // Handle attribute filter
    params.delete('attribute');
    if (checkedAttributes.length > 0) {
      params.set('attribute', checkedAttributes[0]);
    }

    // Reset to page 1 when changing filters
    params.delete('page');

    // Update URL - maintaining the collection slug path
    router.replace(`${pathname}?${params.toString()}`, { scroll: false });
  }, [
    checkedCategories,
    checkedBrands,
    checkedAttributes,
    pathname,
    router,
    searchParams,
  ]);

  // Handle category checkbox change
  const handleCategoryChange = (category: FilterOption) => {
    if (checkedCategories.includes(category.id)) {
      // If already checked, uncheck it
      setCheckedCategories([]);
    } else {
      // If not checked, make it the only checked category (single selection mode)
      setCheckedCategories([category.id]);
    }
  };

  // Handle brand checkbox change
  const handleBrandChange = (brand: FilterOption) => {
    if (checkedBrands.includes(brand.id)) {
      // If already checked, uncheck it
      setCheckedBrands([]);
    } else {
      // If not checked, make it the only checked brand (single selection mode)
      setCheckedBrands([brand.id]);
    }
  };

  // Handle attribute checkbox change
  const handleAttributeChange = (attributeKey: string, valueId: string) => {
    const attributeValue = `${attributeKey}:${valueId}`;
    if (checkedAttributes.includes(attributeValue)) {
      // If already checked, uncheck it
      setCheckedAttributes([]);
    } else {
      // If not checked, make it the only checked attribute (single selection mode)
      setCheckedAttributes([attributeValue]);
    }
  };

  // Check if any filters are active
  const hasActiveFilters =
    checkedCategories.length > 0 ||
    checkedBrands.length > 0 ||
    checkedAttributes.length > 0;

  // Handle clearing all filters
  const clearAllFilters = () => {
    setCheckedCategories([]);
    setCheckedBrands([]);
    setCheckedAttributes([]);
  };

  return (
    <div className='bg-white rounded-lg shadow-sm p-5 sticky top-24 border border-gray-100'>
      {/* CATEGORIES SECTION */}
      <Disclosure defaultOpen>
        {({ open }) => (
          <div>
            <DisclosureButton className='font-medium text-gray-900 mb-4 pb-2 border-b flex items-center w-full justify-between'>
              <div className='flex items-center'>
                <svg
                  xmlns='http://www.w3.org/2000/svg'
                  className='h-4 w-4 mr-2 text-indigo-600'
                  fill='none'
                  viewBox='0 0 24 24'
                  stroke='currentColor'
                >
                  <path
                    strokeLinecap='round'
                    strokeLinejoin='round'
                    strokeWidth={2}
                    d='M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z'
                  />
                </svg>
                Filter by Category
              </div>
              <svg
                className={`w-4 h-4 transform transition-transform ${
                  open ? 'rotate-180' : ''
                }`}
                fill='none'
                stroke='currentColor'
                viewBox='0 0 24 24'
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  strokeWidth={2}
                  d='M19 9l-7 7-7-7'
                />
              </svg>
            </DisclosureButton>
            <DisclosurePanel className='mb-4'>
              {categories.length > 0 ? (
                <div className='space-y-2'>
                  {categories.map((category) => {
                    const isChecked = checkedCategories.includes(category.id);

                    return (
                      <div key={category.id} className='flex items-center'>
                        <Checkbox
                          checked={isChecked}
                          onChange={() => handleCategoryChange(category)}
                          className={`flex items-center group text-sm py-1.5 px-2 w-full rounded-md transition-colors ${
                            isChecked
                              ? 'bg-indigo-50 text-indigo-700 font-medium'
                              : 'text-gray-700 hover:text-indigo-600 hover:bg-gray-50'
                          }`}
                        >
                          <div className='flex items-center w-full cursor-pointer'>
                            <span
                              className={`h-4 w-4 border rounded flex items-center justify-center mr-3 transition-colors ${
                                isChecked
                                  ? 'bg-indigo-600 border-indigo-600'
                                  : 'bg-white border-gray-300 group-hover:border-indigo-400'
                              }`}
                            >
                              {isChecked && (
                                <svg
                                  width='8'
                                  height='6'
                                  viewBox='0 0 8 6'
                                  fill='none'
                                  xmlns='http://www.w3.org/2000/svg'
                                >
                                  <path
                                    d='M7 1L3 5L1 3'
                                    stroke='white'
                                    strokeWidth='1.5'
                                    strokeLinecap='round'
                                    strokeLinejoin='round'
                                  />
                                </svg>
                              )}
                            </span>
                            <span>{category.name}</span>
                            {isChecked && (
                              <span className='ml-auto bg-indigo-100 text-indigo-700 text-xs py-0.5 px-1.5 rounded-full'>
                                Active
                              </span>
                            )}
                          </div>
                        </Checkbox>
                      </div>
                    );
                  })}
                </div>
              ) : (
                <div className='text-center py-4'>
                  <p className='text-sm text-gray-500'>
                    No categories available
                  </p>
                </div>
              )}
            </DisclosurePanel>
          </div>
        )}
      </Disclosure>

      {/* BRANDS SECTION */}
      <Disclosure defaultOpen>
        {({ open }) => (
          <div>
            <DisclosureButton className='font-medium text-gray-900 mb-4 pb-2 border-b flex items-center w-full justify-between'>
              <div className='flex items-center'>
                <svg
                  xmlns='http://www.w3.org/2000/svg'
                  className='h-4 w-4 mr-2 text-indigo-600'
                  fill='none'
                  viewBox='0 0 24 24'
                  stroke='currentColor'
                >
                  <path
                    strokeLinecap='round'
                    strokeLinejoin='round'
                    strokeWidth={2}
                    d='M16 4v12l-4-2-4 2V4M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z'
                  />
                </svg>
                Filter by Brand
              </div>
              <svg
                className={`w-4 h-4 transform transition-transform ${
                  open ? 'rotate-180' : ''
                }`}
                fill='none'
                stroke='currentColor'
                viewBox='0 0 24 24'
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  strokeWidth={2}
                  d='M19 9l-7 7-7-7'
                />
              </svg>
            </DisclosureButton>
            <DisclosurePanel className='mb-4'>
              {brands.length > 0 ? (
                <div className='space-y-2'>
                  {brands.map((brand) => {
                    const isChecked = checkedBrands.includes(brand.id);

                    return (
                      <div key={brand.id} className='flex items-center'>
                        <Checkbox
                          checked={isChecked}
                          onChange={() => handleBrandChange(brand)}
                          className={`flex items-center group text-sm py-1.5 px-2 w-full rounded-md transition-colors ${
                            isChecked
                              ? 'bg-indigo-50 text-indigo-700 font-medium'
                              : 'text-gray-700 hover:text-indigo-600 hover:bg-gray-50'
                          }`}
                        >
                          <div className='flex items-center w-full cursor-pointer'>
                            <span
                              className={`h-4 w-4 border rounded flex items-center justify-center mr-3 transition-colors ${
                                isChecked
                                  ? 'bg-indigo-600 border-indigo-600'
                                  : 'bg-white border-gray-300 group-hover:border-indigo-400'
                              }`}
                            >
                              {isChecked && (
                                <svg
                                  width='8'
                                  height='6'
                                  viewBox='0 0 8 6'
                                  fill='none'
                                  xmlns='http://www.w3.org/2000/svg'
                                >
                                  <path
                                    d='M7 1L3 5L1 3'
                                    stroke='white'
                                    strokeWidth='1.5'
                                    strokeLinecap='round'
                                    strokeLinejoin='round'
                                  />
                                </svg>
                              )}
                            </span>
                            <span>{brand.name}</span>
                            {isChecked && (
                              <span className='ml-auto bg-indigo-100 text-indigo-700 text-xs py-0.5 px-1.5 rounded-full'>
                                Active
                              </span>
                            )}
                          </div>
                        </Checkbox>
                      </div>
                    );
                  })}
                </div>
              ) : (
                <div className='text-center py-4'>
                  <p className='text-sm text-gray-500'>No brands available</p>
                </div>
              )}
            </DisclosurePanel>
          </div>
        )}
      </Disclosure>

      {/* ATTRIBUTES SECTION */}
      <Disclosure defaultOpen={Object.keys(attributes).length > 0}>
        {({ open }) => (
          <div>
            <DisclosureButton className='font-medium text-gray-900 mb-4 pb-2 border-b flex items-center w-full justify-between'>
              <div className='flex items-center'>
                <svg
                  xmlns='http://www.w3.org/2000/svg'
                  className='h-4 w-4 mr-2 text-indigo-600'
                  fill='none'
                  viewBox='0 0 24 24'
                  stroke='currentColor'
                >
                  <path
                    strokeLinecap='round'
                    strokeLinejoin='round'
                    strokeWidth={2}
                    d='M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z'
                  />
                </svg>
                Filter by Attributes
              </div>
              <svg
                className={`w-4 h-4 transform transition-transform ${
                  open ? 'rotate-180' : ''
                }`}
                fill='none'
                stroke='currentColor'
                viewBox='0 0 24 24'
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  strokeWidth={2}
                  d='M19 9l-7 7-7-7'
                />
              </svg>
            </DisclosureButton>
            <DisclosurePanel className='mb-4'>
              {Object.keys(attributes).length > 0 ? (
                <div className='space-y-4'>
                  {Object.entries(attributes).map(([attributeKey, values]) => (
                    <div key={attributeKey} className="mb-2">
                      <div className="font-medium text-sm text-gray-700 mb-1">{attributeKey}</div>
                      <div className="ml-2 space-y-2">
                        {values.map(value => {
                          const attributeValue = `${attributeKey}:${value.id}`;
                          const isChecked = checkedAttributes.includes(attributeValue);
                          
                          return (
                            <div key={value.id} className='flex items-center'>
                              <Checkbox
                                checked={isChecked}
                                onChange={() => handleAttributeChange(attributeKey, value.id)}
                                className={`flex items-center group text-sm py-1.5 px-2 w-full rounded-md transition-colors ${
                                  isChecked
                                    ? 'bg-indigo-50 text-indigo-700 font-medium'
                                    : 'text-gray-700 hover:text-indigo-600 hover:bg-gray-50'
                                }`}
                              >
                                <div className='flex items-center w-full cursor-pointer'>
                                  <span
                                    className={`h-4 w-4 border rounded flex items-center justify-center mr-3 transition-colors ${
                                      isChecked
                                        ? 'bg-indigo-600 border-indigo-600'
                                        : 'bg-white border-gray-300 group-hover:border-indigo-400'
                                    }`}
                                  >
                                    {isChecked && (
                                      <svg
                                        width='8'
                                        height='6'
                                        viewBox='0 0 8 6'
                                        fill='none'
                                        xmlns='http://www.w3.org/2000/svg'
                                      >
                                        <path
                                          d='M7 1L3 5L1 3'
                                          stroke='white'
                                          strokeWidth='1.5'
                                          strokeLinecap='round'
                                          strokeLinejoin='round'
                                        />
                                      </svg>
                                    )}
                                  </span>
                                  <span>{value.name}</span>
                                  {isChecked && (
                                    <span className='ml-auto bg-indigo-100 text-indigo-700 text-xs py-0.5 px-1.5 rounded-full'>
                                      Active
                                    </span>
                                  )}
                                </div>
                              </Checkbox>
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className='text-center py-4'>
                  <p className='text-sm text-gray-500'>
                    No attributes available
                  </p>
                </div>
              )}
            </DisclosurePanel>
          </div>
        )}
      </Disclosure>

      {/* Clear filters button - only show when a filter is active */}
      {hasActiveFilters && (
        <div className='mt-6 pt-3 border-t'>
          <button
            onClick={clearAllFilters}
            className='text-xs font-medium text-indigo-600 hover:text-indigo-800 flex items-center justify-center bg-indigo-50 rounded-md py-2 px-3 transition-colors w-full'
          >
            <svg
              className='h-3 w-3 mr-1'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M6 18L18 6M6 6l12 12'
              />
            </svg>
            Clear all filters
          </button>
        </div>
      )}
    </div>
  );
}
