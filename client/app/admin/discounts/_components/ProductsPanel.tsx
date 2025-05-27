'use client';

import React, { useState, useMemo, useEffect } from 'react';
import { useFormContext } from 'react-hook-form';
import { TabPanel, Button } from '@headlessui/react';
import {
  XMarkIcon,
  PlusIcon,
  MagnifyingGlassIcon,
  ArrowLeftIcon,
  ArrowRightIcon,
} from '@heroicons/react/24/outline';
import { EditDiscountFormData, ProductType } from '../_types';
import { AnimatePresence, motion } from 'framer-motion';
import useDiscountProducts from '../_lib/useDiscountProducts';
import { useParams } from 'next/navigation';
import { useProducts } from '@/app/hooks';

interface ProductsPanelProps {}

export const ProductsPanel: React.FC<ProductsPanelProps> = ({}) => {
  const { id } = useParams<{ id: string }>();
  const { setValue, watch } = useFormContext<EditDiscountFormData>();
  const selectedProducts = watch('products', []);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [availableProducts, setAvailableProducts] = useState<ProductType[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');

  // Implement debounce for search
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(searchQuery);
      setPage(1); // Reset to first page when search changes
    }, 500);

    return () => clearTimeout(timer);
  }, [searchQuery]);

  const { products, isLoading, pagination } = useProducts({
    page,
    limit: 20,
    debouncedSearch,
  });
  const { productDiscounts, isLoading: isDiscountsLoading } =
    useDiscountProducts(id as string);

  // Calculate total pages whenever totalItems changes
  useEffect(() => {
    if (pagination?.pageSize) {
      setTotalPages(Math.ceil(pagination.pageSize / 20));
    }
  }, [pagination]);

  const handleAddProduct = (product: ProductType) => {
    const currentProducts = [...selectedProducts!];
    currentProducts.push(product);
    setValue('products', currentProducts, {
      shouldDirty: true,
    });

    // Remove from available products
    setAvailableProducts((prev) => prev.filter((p) => p.id !== product.id));
  };

  const handleRemoveProduct = (productId: string) => {
    const removedProduct = selectedProducts!.find((p) => p.id === productId);
    if (removedProduct) {
      setValue(
        'products',
        selectedProducts!.filter((p) => p.id !== productId),
        {
          shouldDirty: true,
        }
      );
      // Add back to available products if it should be visible on current page
      setAvailableProducts((prev) => {
        // Check if the product is already in available products
        if (!prev.some((p) => p.id === removedProduct.id)) {
          return [...prev, removedProduct];
        }
        return prev;
      });
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
    return selectedProducts!.filter(
      (product) =>
        product.name.toLowerCase().includes(query) ||
        product.price.toString().includes(query)
    );
  }, [selectedProducts, searchQuery]);

  // Update available products when products data changes
  useEffect(() => {
    if (products) {
      // Filter out products that are already selected
      setAvailableProducts(
        products.map((e) => ({
          id: e.id,
          name: e.name,
          price: e.minPrice,
        }))
      );
    }
  }, [products]);

  // Handle page changes
  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= totalPages) {
      setPage(newPage);
    }
  };

  useEffect(() => {
    if (productDiscounts) {
      setValue('products', productDiscounts);
    }
  }, [productDiscounts]);

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

        {/* Selected Products with Loading UI for discount products */}
        <div className='border-b border-gray-200 pb-5 mb-5'>
          <h4 className='font-medium text-sm text-gray-700 mb-2'>
            Selected Products
          </h4>

          {isDiscountsLoading ? (
            <div className='flex flex-col items-center justify-center py-8'>
              <div className='animate-spin rounded-full h-8 w-8 border-b-2 border-primary mb-2'></div>
              <p className='text-sm text-gray-500'>
                Loading discount products...
              </p>
            </div>
          ) : selectedProducts!.length > 0 ? (
            <div className='mb-6 grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3'>
              {filteredSelectedProducts!.map((product) => (
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

        {/* Available Products with Loading State */}
        <h4 className='font-medium text-sm text-gray-700 mb-2'>Add Products</h4>

        {isLoading ? (
          <div className='flex justify-center items-center py-10'>
            <div className='animate-spin rounded-full h-10 w-10 border-b-2 border-primary'></div>
          </div>
        ) : Boolean(filteredAvailableProducts?.length) ? (
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
              {availableProducts.length === 0 && !isLoading && !debouncedSearch
                ? 'All available products have been added to this discount.'
                : 'No products match your search criteria.'}
            </p>
          </div>
        )}

        {/* Pagination Controls */}
        {totalPages > 1 && (
          <div className='flex items-center justify-between border-t border-gray-200 px-4 py-3 sm:px-6 mt-6'>
            <div className='flex flex-1 justify-between sm:hidden'>
              <Button
                onClick={() => handlePageChange(page - 1)}
                disabled={page === 1}
                className={`relative inline-flex items-center px-4 py-2 text-sm font-medium rounded-md ${
                  page === 1
                    ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                    : 'bg-white text-gray-700 hover:bg-gray-50'
                } border border-gray-300`}
              >
                Previous
              </Button>
              <Button
                onClick={() => handlePageChange(page + 1)}
                disabled={page === totalPages}
                className={`ml-3 relative inline-flex items-center px-4 py-2 text-sm font-medium rounded-md ${
                  page === totalPages
                    ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                    : 'bg-white text-gray-700 hover:bg-gray-50'
                } border border-gray-300`}
              >
                Next
              </Button>
            </div>
            <div className='hidden sm:flex sm:flex-1 sm:items-center sm:justify-between'>
              <div>
                <p className='text-sm text-gray-700'>
                  Showing page <span className='font-medium'>{page}</span> of{' '}
                  <span className='font-medium'>{totalPages}</span> pages
                </p>
              </div>
              <div>
                <nav
                  className='isolate inline-flex -space-x-px rounded-md shadow-sm'
                  aria-label='Pagination'
                >
                  <Button
                    onClick={() => handlePageChange(page - 1)}
                    disabled={page === 1}
                    className={`relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 ${
                      page === 1 ? 'cursor-not-allowed' : 'hover:bg-gray-50'
                    }`}
                  >
                    <span className='sr-only'>Previous</span>
                    <ArrowLeftIcon className='h-5 w-5' aria-hidden='true' />
                  </Button>

                  {/* Page numbers */}
                  {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                    // Show current page and pages around it
                    let pageNum: number;
                    if (totalPages <= 5) {
                      // If 5 or fewer pages, show all
                      pageNum = i + 1;
                    } else if (page <= 3) {
                      // Near the start
                      pageNum = i + 1;
                    } else if (page >= totalPages - 2) {
                      // Near the end
                      pageNum = totalPages - 4 + i;
                    } else {
                      // In the middle
                      pageNum = page - 2 + i;
                    }

                    return (
                      <Button
                        key={pageNum}
                        onClick={() => handlePageChange(pageNum)}
                        className={`relative inline-flex items-center px-4 py-2 text-sm font-semibold ${
                          page === pageNum
                            ? 'bg-primary text-white focus:z-20'
                            : 'text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50'
                        }`}
                      >
                        {pageNum}
                      </Button>
                    );
                  })}

                  <Button
                    onClick={() => handlePageChange(page + 1)}
                    disabled={page === totalPages}
                    className={`relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 ${
                      page === totalPages
                        ? 'cursor-not-allowed'
                        : 'hover:bg-gray-50'
                    }`}
                  >
                    <span className='sr-only'>Next</span>
                    <ArrowRightIcon className='h-5 w-5' aria-hidden='true' />
                  </Button>
                </nav>
              </div>
            </div>
          </div>
        )}
      </motion.div>
    </TabPanel>
  );
};
