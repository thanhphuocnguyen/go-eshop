'use client';

import React, { useState, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { TabPanel, Button } from '@headlessui/react';
import {
  XMarkIcon,
  PlusIcon,
  MagnifyingGlassIcon,
} from '@heroicons/react/24/outline';
import { DiscountFormData, ProductType } from '../_types';
import { AnimatePresence, motion } from 'framer-motion';

interface ProductsPanelProps {}

// Mock data for products and categories selection
const mockAllProducts = [
  { id: 'prod1', name: 'Summer T-Shirt', price: 29.99 },
  { id: 'prod2', name: 'Beach Shorts', price: 39.99 },
  { id: 'prod3', name: 'Sunglasses', price: 25.0 },
  { id: 'prod4', name: 'Beach Hat', price: 19.99 },
  { id: 'prod5', name: 'Flip Flops', price: 15.99 },
  { id: 'prod6', name: 'Beach Towel', price: 22.99 },
];

export const ProductsPanel: React.FC<ProductsPanelProps> = ({}) => {
  const { setValue, watch } = useFormContext<DiscountFormData>();
  const selectedProducts = watch('products') || [];
  const [availableProducts, setAvailableProducts] =
    useState<ProductType[]>(mockAllProducts);

  const [searchQuery, setSearchQuery] = useState('');

  const handleAddProduct = (product: ProductType) => {
    const currentProducts = [...selectedProducts];
    currentProducts.push(product);
    setValue('products', currentProducts);

    // Remove from available products
    setAvailableProducts(availableProducts.filter((p) => p.id !== product.id));
  };

  const handleRemoveProduct = (productId: string) => {
    const removedProduct = selectedProducts.find((p) => p.id === productId);
    if (removedProduct) {
      setValue(
        'products',
        selectedProducts.filter((p) => p.id !== productId)
      );
      setAvailableProducts([...availableProducts, removedProduct]);
    }
  };

  // Filter available products based on search query
  const filteredAvailableProducts = useMemo(() => {
    if (!searchQuery.trim()) return availableProducts;

    const query = searchQuery.toLowerCase().trim();
    return availableProducts.filter(
      (product) =>
        product.name.toLowerCase().includes(query) ||
        product.price.toString().includes(query)
    );
  }, [availableProducts, searchQuery]);

  // Filter selected products based on search query
  const filteredSelectedProducts = useMemo(() => {
    if (!searchQuery.trim()) return selectedProducts;

    const query = searchQuery.toLowerCase().trim();
    return selectedProducts.filter(
      (product) =>
        product.name.toLowerCase().includes(query) ||
        product.price.toString().includes(query)
    );
  }, [selectedProducts, searchQuery]);

  return (
    <TabPanel as={AnimatePresence} mode='wait'>
      <motion.div
        key='products'
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        exit={{ opacity: 0, x: 20 }}
        transition={{ duration: 0.3 }}
      >
        <div className='pt-6 flex justify-between items-center mb-4'>
          <h3 className='text-lg font-medium'>Applied Products</h3>
          <p className='text-sm text-gray-500'>
            If no products are selected, this discount applies to all products.
          </p>
        </div>

        {/* Search Input */}
        <div className='relative mb-4'>
          <div className='absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none'>
            <MagnifyingGlassIcon className='h-5 w-5 text-gray-400' />
          </div>
          <input
            type='text'
            placeholder='Search products by name or price...'
            className='w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent'
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>

        {/* Selected Products */}
        <div className='border-b border-gray-200 pb-5 mb-5'>
          <h4 className='font-medium text-sm text-gray-700 mb-2'>
            Selected Products
          </h4>
          {selectedProducts.length > 0 ? (
            <div className='mb-6 grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3'>
              {filteredSelectedProducts.map((product) => (
                <div
                  key={product.id}
                  className='relative flex items-center border rounded-md p-3 shadow-sm hover:shadow-md transition-shadow duration-200'
                >
                  <div className='flex-1'>
                    <p className='font-medium'>{product.name}</p>
                    <p className='text-gray-500 text-sm'>
                      ${product.price.toFixed(2)}
                    </p>
                  </div>
                  <Button
                    type='button'
                    onClick={() => handleRemoveProduct(product.id)}
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
                No products selected. This discount applies to all products.
              </p>
            </div>
          )}
        </div>

        {/* Available Products */}
        <h4 className='font-medium text-sm text-gray-700 mb-2'>Add Products</h4>
        {filteredAvailableProducts.length > 0 ? (
          <div className='grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3'>
            {filteredAvailableProducts.map((product) => (
              <div
                key={product.id}
                className='flex items-center border rounded-md p-3 cursor-pointer hover:border-primary hover:bg-gray-50 transition-all duration-200 shadow-sm'
                onClick={() => handleAddProduct(product)}
              >
                <div className='flex-1'>
                  <p className='font-medium'>{product.name}</p>
                  <p className='text-gray-500 text-sm'>
                    ${product.price.toFixed(2)}
                  </p>
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
              {availableProducts.length === 0
                ? 'All available products have been added to this discount.'
                : 'No products match your search criteria.'}
            </p>
          </div>
        )}
      </motion.div>
    </TabPanel>
  );
};
