'use client';

import React, { use, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import {
  ArrowLeftIcon,
  PencilIcon,
  TagIcon,
  ShoppingBagIcon,
  ClipboardDocumentListIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';
import dayjs from 'dayjs';
import { ProductsTable, CategoriesTable, UsageTable } from '../_components';
import useSWR from 'swr';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';

// Mock data for a specific discount
export type Discount = {
  id: string;
  code: string;
  description: string;
  discountType: 'percentage' | 'fixed_amount';
  discountValue: number;
  minPurchase?: number;
  maxDiscount?: number;
  usageLimit?: number;
  usedCount?: number;
  isActive: boolean;
  startsAt: string;
  expiresAt: string;
  createdAt: string;
  updatedAt: string;
  products: Array<{ id: string; name: string; price: number }>;
  categories: Array<{ id: string; name: string }>;
  usageHistory: Array<{
    id: string;
    orderId: string;
    customerName: string;
    amount: number;
    discountAmount: number;
    date: string;
  }>;
};

export default function DiscountDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const router = useRouter();
  const [activeTab, setActiveTab] = useState('details');

  // Use mock data for demonstration
  // const discount = mockDiscount;

  // /* Uncomment to use real API
  const { data: discount, isLoading } = useSWR(
    ADMIN_API_PATHS.DISCOUNT.replace(':id', id),
    (url) =>
      clientSideFetch<Discount>(url).then((res) => {
        if (res.error) {
          throw new Error(res.error.details);
        }
        return res.data;
      }),
    {
      onError: (err) => {
        toast.error(
          <div>
            Failed to fetch discount details.
            <div>{JSON.stringify(err)}</div>
          </div>
        );
      },
    }
  );

  if (isLoading) {
    return (
      <div className='p-4'>
        <div className='flex justify-center items-center h-64'>
          <div className='animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary'></div>
        </div>
      </div>
    );
  }

  if (!discount) {
    return (
      <div className='p-4'>
        <div className='flex flex-col items-center justify-center h-64'>
          <h2 className='text-xl font-semibold'>Discount not found</h2>
          <p className='text-gray-500 mt-2'>
            The discount you're looking for doesn't exist or has been removed.
          </p>
          <button
            onClick={() => router.push('/admin/discounts')}
            className='mt-4 px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80'
          >
            Back to Discounts
          </button>
        </div>
      </div>
    );
  }

  const isActive =
    discount.isActive &&
    dayjs(discount.startsAt).isBefore(dayjs()) &&
    dayjs(discount.expiresAt).isAfter(dayjs());
  const isExpired = dayjs(discount.expiresAt).isBefore(dayjs());
  const isUpcoming = dayjs(discount.startsAt).isAfter(dayjs());

  let statusClasses =
    'px-2 py-1 inline-flex items-center text-xs font-medium rounded-full ';
  if (isActive) {
    statusClasses += 'bg-green-100 text-green-800';
  } else if (isExpired) {
    statusClasses += 'bg-red-100 text-red-800';
  } else if (isUpcoming) {
    statusClasses += 'bg-yellow-100 text-yellow-800';
  } else {
    statusClasses += 'bg-gray-100 text-gray-800';
  }

  const getStatusText = () => {
    if (isActive) return 'Active';
    if (isExpired) return 'Expired';
    if (isUpcoming) return 'Upcoming';
    return 'Inactive';
  };

  const usagePercentage = discount.usageLimit
    ? Math.min(100, ((discount.usedCount ?? 0) / discount.usageLimit) * 100)
    : 0;

  return (
    <div className='p-4'>
      <div className='mb-4 flex flex-col sm:flex-row sm:items-center sm:justify-between'>
        <div className='flex items-center mb-4 sm:mb-0'>
          <button
            onClick={() => router.back()}
            className='mr-4 p-2 rounded-full hover:bg-gray-100'
          >
            <ArrowLeftIcon className='h-5 w-5' />
          </button>
          <div>
            <div className='flex items-center'>
              <h1 className='text-2xl font-semibold mr-3'>{discount.code}</h1>
              <span className={statusClasses}>{getStatusText()}</span>
            </div>
            <p className='text-gray-500'>{discount.description}</p>
          </div>
        </div>

        <div className='flex space-x-2'>
          <Link href={`/admin/discounts/${id}/edit`}>
            <button className='flex items-center gap-1 px-3 py-1.5 border border-gray-300 rounded bg-white hover:bg-gray-50'>
              <PencilIcon className='h-4 w-4' />
              Edit
            </button>
          </Link>
        </div>
      </div>

      <div className='bg-white rounded-lg shadow overflow-hidden'>
        <div className='border-b'>
          <nav className='flex -mb-px'>
            <button
              onClick={() => setActiveTab('details')}
              className={`px-4 py-3 text-sm font-medium ${
                activeTab === 'details'
                  ? 'border-b-2 border-primary text-primary'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <span className='flex items-center'>
                <TagIcon className='h-4 w-4 mr-2' />
                Details
              </span>
            </button>
            <button
              onClick={() => setActiveTab('products')}
              className={`px-4 py-3 text-sm font-medium ${
                activeTab === 'products'
                  ? 'border-b-2 border-primary text-primary'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <span className='flex items-center'>
                <ShoppingBagIcon className='h-4 w-4 mr-2' />
                Products
              </span>
            </button>
            <button
              onClick={() => setActiveTab('categories')}
              className={`px-4 py-3 text-sm font-medium ${
                activeTab === 'categories'
                  ? 'border-b-2 border-primary text-primary'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <span className='flex items-center'>
                <ClipboardDocumentListIcon className='h-4 w-4 mr-2' />
                Categories
              </span>
            </button>
            <button
              onClick={() => setActiveTab('usage')}
              className={`px-4 py-3 text-sm font-medium ${
                activeTab === 'usage'
                  ? 'border-b-2 border-primary text-primary'
                  : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <span className='flex items-center'>
                <ChartBarIcon className='h-4 w-4 mr-2' />
                Usage
              </span>
            </button>
          </nav>
        </div>

        <div className='p-6'>
          {activeTab === 'details' && (
            <div className='grid grid-cols-1 md:grid-cols-2 gap-6'>
              <div className='col-span-1'>
                <h3 className='text-lg font-medium mb-4'>
                  Discount Information
                </h3>
                <div className='space-y-3'>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Code
                    </span>
                    <span className='mt-1 block'>{discount.code}</span>
                  </div>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Description
                    </span>
                    <span className='mt-1 block'>
                      {discount.description || 'No description provided'}
                    </span>
                  </div>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Type
                    </span>
                    <span className='mt-1 block capitalize'>
                      {discount.discountType.replace('_', ' ')}
                    </span>
                  </div>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Value
                    </span>
                    <span className='mt-1 block'>
                      {discount.discountType === 'percentage'
                        ? `${discount.discountValue}%`
                        : `$${discount.discountValue.toFixed(2)}`}
                    </span>
                  </div>
                </div>
              </div>

              <div className='col-span-1'>
                <h3 className='text-lg font-medium mb-4'>
                  Usage & Restrictions
                </h3>
                <div className='space-y-3'>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Minimum Purchase
                    </span>
                    <span className='mt-1 block'>
                      {discount.minPurchase
                        ? `$${discount.minPurchase.toFixed(2)}`
                        : 'No minimum'}
                    </span>
                  </div>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Maximum Discount
                    </span>
                    <span className='mt-1 block'>
                      {discount.maxDiscount
                        ? `$${discount.maxDiscount.toFixed(2)}`
                        : 'No maximum'}
                    </span>
                  </div>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Usage Limit
                    </span>
                    <span className='mt-1 block'>
                      {discount.usageLimit
                        ? `${discount.usedCount} of ${discount.usageLimit} used`
                        : 'Unlimited'}
                    </span>
                    {discount.usageLimit && (
                      <div className='w-full bg-gray-200 rounded-full h-2.5 mt-2'>
                        <div
                          className='bg-primary h-2.5 rounded-full'
                          style={{ width: `${usagePercentage}%` }}
                        ></div>
                      </div>
                    )}
                  </div>
                </div>
              </div>

              <div className='col-span-1'>
                <h3 className='text-lg font-medium mb-4'>Dates</h3>
                <div className='space-y-3'>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Valid From
                    </span>
                    <span className='mt-1 block'>
                      {dayjs(discount.startsAt).format('MMMM D, YYYY h:mm A')}
                    </span>
                  </div>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Expires On
                    </span>
                    <span className='mt-1 block'>
                      {dayjs(discount.expiresAt).format('MMMM D, YYYY h:mm A')}
                    </span>
                  </div>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Created
                    </span>
                    <span className='mt-1 block'>
                      {dayjs(discount.createdAt).format('MMMM D, YYYY h:mm A')}
                    </span>
                  </div>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Last Updated
                    </span>
                    <span className='mt-1 block'>
                      {dayjs(discount.updatedAt).format('MMMM D, YYYY h:mm A')}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'products' && (
            <div>
              <div className='flex justify-between items-center mb-4'>
                <h3 className='text-lg font-medium'>Applied Products</h3>
                <button className='px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50'>
                  Manage Products
                </button>
              </div>

              <ProductsTable
                id={discount.id}
                discountType={discount.discountType}
                discountValue={discount.discountValue}
              />
            </div>
          )}

          {activeTab === 'categories' && (
            <div>
              <div className='flex justify-between items-center mb-4'>
                <h3 className='text-lg font-medium'>Applied Categories</h3>
                <button className='px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50'>
                  Manage Categories
                </button>
              </div>

              <CategoriesTable id={discount.id} />
            </div>
          )}

          {activeTab === 'usage' && (
            <UsageTable
              usageHistory={discount.usageHistory}
              overview={{
                usedCount: discount.usedCount || 0,
                usageLimit: discount.usageLimit,
              }}
            />
          )}
        </div>
      </div>
    </div>
  );
}
