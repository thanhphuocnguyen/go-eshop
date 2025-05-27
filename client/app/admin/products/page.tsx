'use client';

import Loading from '@/app/loading';
import { Button } from '@headlessui/react';
import Link from 'next/link';
import { useState, useEffect } from 'react';

import dayjs from 'dayjs';
import Image from 'next/image';
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  MagnifyingGlassIcon,
} from '@heroicons/react/24/outline';
import { useProducts, useDebounce } from '@/app/hooks';

export default function Page() {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [searchInput, setSearchInput] = useState('');
  const debouncedSearch = useDebounce(searchInput, 500); // 500ms delay
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const { products, pagination, isLoading } = useProducts({
    page,
    limit,
    debouncedSearch,
  });

  const handlePreviousPage = () => {
    if (page > 1) {
      setPage(page - 1);
    }
  };

  const handleNextPage = () => {
    if (page < totalPages) {
      setPage(page + 1);
    }
  };

  useEffect(() => {
    setPage(1);
  }, [debouncedSearch]);

  const handleLimitChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setLimit(Number(e.target.value));
    setPage(1); // Reset to first page when changing limit
  };

  useEffect(() => {
    if (pagination) {
      const { total, totalPages } = pagination;
      setTotal(total);
      setTotalPages(totalPages);
    }
  }, [pagination, limit]);

  if (isLoading) return <Loading />;

  return (
    <div className='h-full overflow-auto'>
      <div className='flex justify-between items-center pt-4 pb-8'>
        <h2 className='text-2xl font-semibold text-primary'>Product List</h2>
        <Button
          as={Link}
          href={'/admin/products/new'}
          className='btn btn-lg btn-primary'
        >
          Add new
        </Button>
      </div>

      {/* Search Filter */}
      <div className='mb-6'>
        <div className='flex gap-2'>
          <div className='relative flex-grow'>
            <div className='absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none'>
              <MagnifyingGlassIcon className='w-5 h-5 text-gray-500' />
            </div>
            <input
              type='text'
              className='bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 p-2.5'
              placeholder='Search products...'
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
            />
          </div>
        </div>
      </div>

      <div className='relative overflow-x-auto shadow-md sm:rounded-lg'>
        <table className='w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400'>
          <thead className='text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400'>
            <tr>
              <th scope='col' className='px-6 py-3'>
                Image
              </th>
              <th scope='col' className='px-6 py-3'>
                Product name
              </th>
              <th scope='col' className='px-6 py-3'>
                Price
              </th>
              <th scope='col' className='px-6 py-3'>
                SKU
              </th>
              <th scope='col' className='px-6 py-3'>
                Created At
              </th>
              <th scope='col' className='px-6 py-3'>
                Action
              </th>
            </tr>
          </thead>
          <tbody>
            {products?.map((product) => (
              <tr
                key={product.id}
                className='odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700 border-gray-200'
              >
                <td className='px-6 py-4'>
                  {product?.imageUrl && (
                    <div className='h-10 w-10 relative'>
                      <Image
                        src={product.imageUrl}
                        alt={product.name}
                        fill
                        className='object-cover rounded'
                      />
                    </div>
                  )}
                </td>
                <th
                  scope='row'
                  className='px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white'
                >
                  {product.name}
                </th>
                <td className='px-6 py-4'>
                  ${product.minPrice.toFixed(2)} - $
                  {product.maxPrice.toFixed(2)}
                </td>
                <td className='px-6 py-4'>{product.sku}</td>
                <td className='px-6 py-4'>
                  {product.createdAt &&
                    dayjs(product.createdAt).format('MMM D, YYYY')}
                </td>
                <td className='px-6 py-4'>
                  <Link
                    href={`/admin/products/${product.id}`}
                    className='font-medium text-blue-600 dark:text-blue-500 hover:underline'
                  >
                    Edit
                  </Link>
                </td>
              </tr>
            ))}
            {(!products || products.length === 0) && (
              <tr className='odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700 border-gray-200'>
                <td colSpan={6} className='px-6 py-4 text-center'>
                  No products found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination Component */}
      <div className='flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6 mt-4 rounded-md'>
        <div className='flex flex-1 justify-between sm:hidden'>
          <Button
            onClick={handlePreviousPage}
            disabled={page <= 1}
            className={`relative inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 ${page <= 1 ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50'}`}
          >
            Previous
          </Button>
          <Button
            onClick={handleNextPage}
            disabled={page >= totalPages}
            className={`relative ml-3 inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 ${page >= totalPages ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50'}`}
          >
            Next
          </Button>
        </div>
        <div className='hidden sm:flex sm:flex-1 sm:items-center sm:justify-between'>
          <div>
            <p className='text-sm text-gray-700'>
              Showing{' '}
              <span className='font-medium'>{(page - 1) * limit + 1}</span> to{' '}
              <span className='font-medium'>
                {Math.min(page * limit, total)}
              </span>{' '}
              of <span className='font-medium'>{total}</span> products
            </p>
          </div>
          <div className='flex items-center gap-4'>
            {/* Page limit selector */}
            <div className='flex items-center'>
              <label htmlFor='pageLimit' className='mr-2 text-sm text-gray-700'>
                Show:
              </label>
              <select
                id='pageLimit'
                value={limit}
                onChange={handleLimitChange}
                className='block w-20 rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6'
              >
                <option value={5}>5</option>
                <option value={10}>10</option>
                <option value={25}>25</option>
                <option value={50}>50</option>
                <option value={100}>100</option>
              </select>
            </div>

            <nav
              className='isolate inline-flex -space-x-px rounded-md shadow-sm'
              aria-label='Pagination'
            >
              <Button
                onClick={handlePreviousPage}
                disabled={page <= 1}
                className={`relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ${page <= 1 ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50'}`}
              >
                <span className='sr-only'>Previous</span>
                <ChevronLeftIcon className='h-5 w-5' aria-hidden='true' />
              </Button>
              {/* Page number display */}
              <span className='relative inline-flex items-center px-4 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300'>
                {page} / {totalPages}
              </span>
              <Button
                onClick={handleNextPage}
                disabled={page >= totalPages}
                className={`relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ${page >= totalPages ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50'}`}
              >
                <span className='sr-only'>Next</span>
                <ChevronRightIcon className='h-5 w-5' aria-hidden='true' />
              </Button>
            </nav>
          </div>
        </div>
      </div>
    </div>
  );
}
