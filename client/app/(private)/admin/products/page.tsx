'use client';

import Loading from '@/app/loading';
import { API_PATHS } from '@/lib/constants/api';
import { ProductListModel, GenericResponse } from '@/lib/definitions';
import { Button } from '@headlessui/react';
import Link from 'next/link';
import { useState } from 'react';
import useSWR from 'swr';

import dayjs from 'dayjs';
import Image from 'next/image';
import { apiFetch } from '@/lib/api/api';

export default function Page() {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [search, setSearch] = useState('');
  const [total, setTotal] = useState(0);

  const { data, isLoading } = useSWR(
    [API_PATHS.PRODUCTS, page, limit, search],
    ([url, page, limit, search]) =>
      apiFetch<GenericResponse<ProductListModel[]>>(
        `${url}?page=${page}&limit=${limit}&search=${search}`,
        {}
      ).then((data) => data.data),
    {
      refreshInterval: 0,
      revalidateOnFocus: false,
      onError: (err) => {
        console.error(err);
      },
    }
  );

  if (isLoading) return <Loading />;

  return (
    <div className='h-full'>
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
            {data?.map((product) => (
              <tr
                key={product.id}
                className='odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700 border-gray-200'
              >
                <td className='px-6 py-4'>
                  {product.image_url && (
                    <div className='h-10 w-10 relative'>
                      <Image
                        src={product.image_url}
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
                <td className='px-6 py-4'>${product.price.toFixed(2)}</td>
                <td className='px-6 py-4'>{product.sku}</td>
                <td className='px-6 py-4'>
                  {product.created_at &&
                    dayjs(product.created_at).format('MMM D, YYYY')}
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
            {(!data || data.length === 0) && (
              <tr className='odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700 border-gray-200'>
                <td colSpan={6} className='px-6 py-4 text-center'>
                  No products found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
