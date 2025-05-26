'use client';

import { useState } from 'react';
import Link from 'next/link';
import {
  TagIcon,
  PencilIcon,
  TrashIcon,
  PlusCircleIcon,
  CheckCircleIcon,
  XCircleIcon,
} from '@heroicons/react/24/outline';
import dayjs from 'dayjs';
import useSWR from 'swr';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { toast } from 'react-toastify';

type DiscountModel = {
  id: string;
  code: string;
  description: string;
  discountType: 'percentage' | 'fixed_amount';
  discountValue: number;
  minPurchaseAmount: number | null;
  maxDiscountAmount: number | null;
  usageLimit: number | null;
  usedCount: number;
  isActive: boolean;
  startsAt: string;
  expiresAt: string;
};

// Mock data for discounts
const mockDiscounts: DiscountModel[] = [
  {
    id: '1',
    code: 'SUMMER25',
    description: 'Summer sale discount',
    discountType: 'percentage',
    discountValue: 25,
    minPurchaseAmount: 50,
    maxDiscountAmount: 100,
    usageLimit: 1000,
    usedCount: 342,
    isActive: true,
    startsAt: '2025-05-01T00:00:00Z',
    expiresAt: '2025-08-31T23:59:59Z',
  },
  {
    id: '2',
    code: 'WELCOME10',
    description: 'Welcome discount for new customers',
    discountType: 'fixed_amount',
    discountValue: 10,
    minPurchaseAmount: 30,
    maxDiscountAmount: null,
    usageLimit: 1,
    usedCount: 546,
    isActive: true,
    startsAt: '2025-01-01T00:00:00Z',
    expiresAt: '2025-12-31T23:59:59Z',
  },
  {
    id: '3',
    code: 'FLASH50',
    description: 'Flash sale 50% off',
    discountType: 'percentage',
    discountValue: 50,
    minPurchaseAmount: 100,
    maxDiscountAmount: 200,
    usageLimit: 500,
    usedCount: 500,
    isActive: false,
    startsAt: '2025-04-01T00:00:00Z',
    expiresAt: '2025-04-03T23:59:59Z',
  },
  {
    id: '4',
    code: 'FREESHIP',
    description: 'Free shipping on all orders',
    discountType: 'fixed_amount',
    discountValue: 15,
    minPurchaseAmount: 0,
    maxDiscountAmount: 15,
    usageLimit: 2000,
    usedCount: 1203,
    isActive: true,
    startsAt: '2025-01-01T00:00:00Z',
    expiresAt: '2025-12-31T23:59:59Z',
  },
  {
    id: '5',
    code: 'VIP20',
    description: 'VIP customer discount',
    discountType: 'percentage',
    discountValue: 20,
    minPurchaseAmount: null,
    maxDiscountAmount: null,
    usageLimit: null,
    usedCount: 89,
    isActive: true,
    startsAt: '2025-01-01T00:00:00Z',
    expiresAt: '2025-12-31T23:59:59Z',
  },
];

export default function DiscountsPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [activeFilter, setActiveFilter] = useState<string | null>(null);

  // Filter discounts based on search term and active filter
  const { data: discounts, isLoading } = useSWR(
    [ADMIN_API_PATHS.DISCOUNTS, searchTerm, activeFilter],
    ([url, searchTerm, activeFilter]) =>
      clientSideFetch<DiscountModel[]>(url, {
        queryParams: {
          search: searchTerm ? searchTerm : undefined,
          status: activeFilter,
        },
      }).then((resp) => {
        return resp.data;
      }),
    {
      onError: (error) => {
        toast.error(
          <div>
            Failed to fetch discounts:
            <div>{JSON.stringify(error)}</div>
          </div>
        );
      },
      fallbackData: mockDiscounts,
    }
  );

  const handleDelete = (id: string) => {
    // In a real implementation, this would make an API call to delete the discount
    alert(`Delete discount ${id}`);
  };

  return (
    <section className='flex flex-col gap-4 p-4'>
      <div className='flex justify-between items-center'>
        <div>
          <h1 className='text-2xl font-semibold'>Discounts</h1>
          <p className='text-gray-500'>
            Create and manage discounts for your products.
          </p>
        </div>
        <Link href='/admin/discounts/new'>
          <button className='flex items-center gap-2 bg-primary text-white py-2 px-4 rounded-md hover:bg-primary/80'>
            <PlusCircleIcon className='w-5 h-5' />
            Create Discount
          </button>
        </Link>
      </div>

      <div className='flex flex-col gap-3 md:flex-row md:items-center md:justify-between'>
        <div className='flex gap-2'>
          <button
            onClick={() => setActiveFilter(null)}
            className={`px-3 py-1 rounded-md ${activeFilter === null ? 'bg-primary text-white' : 'bg-gray-100'}`}
          >
            All
          </button>
          <button
            onClick={() => setActiveFilter('active')}
            className={`px-3 py-1 rounded-md ${activeFilter === 'active' ? 'bg-primary text-white' : 'bg-gray-100'}`}
          >
            Active
          </button>
          <button
            onClick={() => setActiveFilter('inactive')}
            className={`px-3 py-1 rounded-md ${activeFilter === 'inactive' ? 'bg-primary text-white' : 'bg-gray-100'}`}
          >
            Inactive
          </button>
        </div>

        <div className='relative'>
          <input
            type='text'
            placeholder='Search discounts...'
            className='border rounded-md pl-10 pr-4 py-2 w-full md:w-64'
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
          <div className='absolute inset-y-0 left-0 flex items-center pl-3'>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-5 w-5 text-gray-400'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z'
              />
            </svg>
          </div>
        </div>
      </div>

      <div className='overflow-x-auto bg-white rounded-lg shadow'>
        <table className='min-w-full divide-y divide-gray-200'>
          <thead className='bg-gray-50'>
            <tr>
              <th
                scope='col'
                className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
              >
                Code
              </th>
              <th
                scope='col'
                className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
              >
                Type
              </th>
              <th
                scope='col'
                className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
              >
                Value
              </th>
              <th
                scope='col'
                className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
              >
                Date Range
              </th>
              <th
                scope='col'
                className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
              >
                Status
              </th>
              <th
                scope='col'
                className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
              >
                Usage
              </th>
              <th
                scope='col'
                className='px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider'
              >
                Actions
              </th>
            </tr>
          </thead>
          <tbody className='bg-white divide-y divide-gray-200'>
            {discounts.map((discount) => (
              <tr key={discount.id} className='hover:bg-gray-50'>
                <td className='px-6 py-4 whitespace-nowrap'>
                  <div className='flex items-center'>
                    <div className='flex-shrink-0 h-10 w-10 flex items-center justify-center text-primary'>
                      <TagIcon className='h-6 w-6' />
                    </div>
                    <div className='ml-4'>
                      <div className='text-sm font-medium text-gray-900'>
                        {discount.code}
                      </div>
                      <div className='text-sm text-gray-500'>
                        {discount.description}
                      </div>
                    </div>
                  </div>
                </td>
                <td className='px-6 py-4 whitespace-nowrap'>
                  <span className='px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800'>
                    {discount.discountType === 'percentage'
                      ? 'Percentage'
                      : 'Fixed Amount'}
                  </span>
                </td>
                <td className='px-6 py-4 whitespace-nowrap'>
                  {discount.discountType === 'percentage'
                    ? `${discount.discountValue}%`
                    : `$${discount.discountValue.toFixed(2)}`}
                  <div className='text-xs text-gray-500'>
                    {discount.minPurchaseAmount
                      ? `Min: $${discount.minPurchaseAmount}`
                      : 'No minimum'}
                  </div>
                  <div className='text-xs text-gray-500'>
                    {discount.maxDiscountAmount
                      ? `Max: $${discount.maxDiscountAmount}`
                      : 'No maximum'}
                  </div>
                </td>
                <td className='px-6 py-4 whitespace-nowrap'>
                  <div className='text-sm text-gray-900'>
                    {dayjs(discount.startsAt).format('MMM D, YYYY')}
                  </div>
                  <div className='text-sm text-gray-500'>
                    to {dayjs(discount.expiresAt).format('MMM D, YYYY')}
                  </div>
                </td>
                <td className='px-6 py-4 whitespace-nowrap'>
                  {discount.isActive ? (
                    <span className='px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800'>
                      <CheckCircleIcon className='h-4 w-4 mr-1' />
                      Active
                    </span>
                  ) : (
                    <span className='px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-100 text-red-800'>
                      <XCircleIcon className='h-4 w-4 mr-1' />
                      Inactive
                    </span>
                  )}
                </td>
                <td className='px-6 py-4 whitespace-nowrap'>
                  <div className='text-sm text-gray-900'>
                    {discount.usedCount} / {discount.usageLimit || '∞'}
                  </div>
                  {discount.usageLimit && (
                    <div className='w-full bg-gray-200 rounded-full h-2.5'>
                      <div
                        className='bg-primary h-2.5 rounded-full'
                        style={{
                          width: `${Math.min(100, (discount.usedCount / discount.usageLimit) * 100)}%`,
                        }}
                      ></div>
                    </div>
                  )}
                </td>
                <td className='px-6 py-4 whitespace-nowrap text-center text-sm font-medium'>
                  <div className='flex justify-center space-x-2'>
                    <Link
                      href={`/admin/discounts/${discount.id}`}
                      className='text-gray-500 hover:text-gray-900'
                    >
                      <span className='sr-only'>View</span>
                      <svg
                        xmlns='http://www.w3.org/2000/svg'
                        className='h-5 w-5'
                        fill='none'
                        viewBox='0 0 24 24'
                        stroke='currentColor'
                      >
                        <path
                          strokeLinecap='round'
                          strokeLinejoin='round'
                          strokeWidth={2}
                          d='M15 12a3 3 0 11-6 0 3 3 0 016 0z'
                        />
                        <path
                          strokeLinecap='round'
                          strokeLinejoin='round'
                          strokeWidth={2}
                          d='M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z'
                        />
                      </svg>
                    </Link>
                    <Link
                      href={`/admin/discounts/${discount.id}/edit`}
                      className='text-indigo-600 hover:text-indigo-900'
                    >
                      <span className='sr-only'>Edit</span>
                      <PencilIcon className='h-5 w-5' />
                    </Link>
                    <button
                      onClick={() => handleDelete(discount.id)}
                      className='text-red-600 hover:text-red-900'
                    >
                      <span className='sr-only'>Delete</span>
                      <TrashIcon className='h-5 w-5' />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  );
}
