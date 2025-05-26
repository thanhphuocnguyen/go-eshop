'use client';

import React, { useState, useMemo, useEffect } from 'react';
import { useFormContext } from 'react-hook-form';
import { TabPanel, Button } from '@headlessui/react';
import {
  XMarkIcon,
  PlusIcon,
  MagnifyingGlassIcon,
} from '@heroicons/react/24/outline';
import { CategoryType, EditDiscountFormData } from '../_types';
import { AnimatePresence, motion } from 'framer-motion';
import { useCategories } from '../../hooks';
import useDiscountCategories from '../_lib/useDiscountCategores';
import { useParams } from 'next/navigation';

const mockAllCategories = [
  { id: 'cat1', name: 'Summer Collection' },
  { id: 'cat2', name: 'Beachwear' },
  { id: 'cat3', name: 'Accessories' },
  { id: 'cat4', name: 'Footwear' },
  { id: 'cat5', name: 'Seasonal' },
];

interface CategoriesPanelProps {}

export const CategoriesPanel: React.FC<CategoriesPanelProps> = ({}) => {
  const { id } = useParams<{ id: string }>();
  // Use form context to access form methods and state
  const { setValue, watch } = useFormContext<EditDiscountFormData>();
  const selectedCategories = watch('categories') || [];

  const [searchQuery, setSearchQuery] = useState('');
  const [availableCategories, setAvailableCategories] =
    useState<CategoryType[]>(mockAllCategories);

  const { categories, isLoading } = useCategories();
  const { categoryDiscounts, isLoading: discountLoading } =
    useDiscountCategories(id);

  const handleAddCategory = (category: CategoryType) => {
    const currentCategories = [...selectedCategories];
    currentCategories.push(category);
    setValue('categories', currentCategories, {
      shouldDirty: true,
    });

    // Remove from available categories
    setAvailableCategories(
      availableCategories.filter((c) => c.id !== category.id)
    );
  };

  const handleRemoveCategory = (categoryId: string) => {
    const removedCategory = selectedCategories.find((c) => c.id === categoryId);
    if (removedCategory) {
      setValue(
        'categories',
        selectedCategories.filter((c) => c.id !== categoryId),
        {
          shouldDirty: true,
        }
      );
      setAvailableCategories([...availableCategories, removedCategory]);
    }
  };

  // Filter available categories based on search query
  const filteredAvailableCategories = useMemo(() => {
    if (!searchQuery.trim()) return availableCategories;

    const query = searchQuery.toLowerCase().trim();
    return availableCategories.filter((category) =>
      category.name.toLowerCase().includes(query)
    );
  }, [availableCategories, searchQuery]);

  // Filter selected categories based on search query
  const filteredSelectedCategories = useMemo(() => {
    if (!searchQuery.trim()) return selectedCategories;

    const query = searchQuery.toLowerCase().trim();
    return selectedCategories.filter((category) =>
      category.name.toLowerCase().includes(query)
    );
  }, [selectedCategories, searchQuery]);

  useEffect(() => {
    if (categories) {
      setAvailableCategories(
        categories.map((cat) => ({ id: cat.id, name: cat.name }))
      );
    }
  }, [categories]);

  useEffect(() => {
    if (categoryDiscounts) {
      setValue('categories', categoryDiscounts);
    }
  }, [categoryDiscounts]);

  return (
    <TabPanel as={AnimatePresence} mode='wait'>
      <motion.div
        key='categories'
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        exit={{ opacity: 0, x: 20 }}
        transition={{ duration: 0.3 }}
      >
        <div className='pt-6 flex justify-between items-center mb-4'>
          <h3 className='text-lg font-medium'>Applied Categories</h3>
          <p className='text-sm text-gray-500'>
            If no categories are selected, this discount applies to all
            products.
          </p>
        </div>

        {/* Search Input */}
        <div className='relative mb-4'>
          <div className='absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none'>
            <MagnifyingGlassIcon className='h-5 w-5 text-gray-400' />
          </div>
          <input
            type='text'
            placeholder='Search categories by name...'
            className='w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent'
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>

        {/* Selected Categories */}
        <div className='border-b border-gray-200 pb-5 mb-5'>
          <h4 className='font-medium text-sm text-gray-700 mb-2'>
            Selected Categories
          </h4>
          
          {discountLoading ? (
            <div className='flex flex-col items-center justify-center py-8'>
              <div className='animate-spin rounded-full h-8 w-8 border-b-2 border-primary mb-2'></div>
              <p className='text-sm text-gray-500'>Loading discount categories...</p>
            </div>
          ) : selectedCategories.length > 0 ? (
            <div className='mb-6 grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3'>
              {filteredSelectedCategories.map((category) => (
                <div
                  key={category.id}
                  className='relative flex items-center border rounded-md p-3 shadow-sm hover:shadow-md transition-shadow duration-200'
                >
                  <div className='flex-1'>
                    <p className='font-medium'>{category.name}</p>
                  </div>
                  <Button
                    type='button'
                    onClick={() => handleRemoveCategory(category.id)}
                    className='text-gray-400 hover:text-red-500 transition-colors'
                  >
                    <XMarkIcon className='h-5 w-5' />
                  </Button>
                </div>
              ))}
            </div>
          ) : (
            <div className='mb-6 bg-gray-50 rounded-md p-6 text-center border border-gray-100'>
              <p className='text-gray-500'>
                No categories selected. This discount applies to all products.
              </p>
            </div>
          )}
        </div>

        {/* Available Categories with Loading State */}
        <h4 className='font-medium text-sm text-gray-700 mb-2'>
          Add Categories
        </h4>

        {isLoading ? (
          <div className='flex justify-center items-center py-10'>
            <div className='animate-spin rounded-full h-10 w-10 border-b-2 border-primary'></div>
            <span className='ml-3 text-sm text-gray-500'>
              Loading categories...
            </span>
          </div>
        ) : filteredAvailableCategories.length > 0 ? (
          <div className='grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3'>
            {filteredAvailableCategories.map((category) => (
              <div
                key={category.id}
                className='flex items-center border rounded-md p-3 cursor-pointer hover:border-primary hover:bg-gray-50 transition-all duration-200 shadow-sm'
                onClick={() => handleAddCategory(category)}
              >
                <div className='flex-1'>
                  <p className='font-medium'>{category.name}</p>
                </div>
                <Button
                  type='button'
                  className='text-gray-400 hover:text-green-500 transition-colors'
                >
                  <PlusIcon className='h-5 w-5' />
                </Button>
              </div>
            ))}
          </div>
        ) : (
          <div className='bg-gray-50 rounded-md p-6 text-center border border-gray-100'>
            <p className='text-gray-500'>
              {availableCategories.length === 0
                ? 'All available categories have been added to this discount.'
                : 'No categories match your search criteria.'}
            </p>
          </div>
        )}
      </motion.div>
    </TabPanel>
  );
};
