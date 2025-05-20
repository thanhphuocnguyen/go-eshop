'use client';

import React, { useState, useEffect } from 'react';
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

// Mock data for a specific discount
const mockDiscount = {
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
  createdAt: '2025-04-15T08:30:00Z',
  updatedAt: '2025-04-20T14:45:30Z',
  products: [
    { id: 'prod1', name: 'Summer T-Shirt', price: 29.99 },
    { id: 'prod2', name: 'Beach Shorts', price: 39.99 },
    { id: 'prod3', name: 'Sunglasses', price: 25.0 },
  ],
  categories: [
    { id: 'cat1', name: 'Summer Collection' },
    { id: 'cat2', name: 'Beachwear' },
  ],
  usageHistory: [
    {
      id: 'order1',
      orderId: 'ORD-12345',
      customerName: 'John Doe',
      amount: 89.97,
      discountAmount: 22.49,
      date: '2025-05-10T15:30:00Z',
    },
    {
      id: 'order2',
      orderId: 'ORD-12346',
      customerName: 'Jane Smith',
      amount: 129.99,
      discountAmount: 32.5,
      date: '2025-05-12T09:15:00Z',
    },
    {
      id: 'order3',
      orderId: 'ORD-12350',
      customerName: 'Robert Johnson',
      amount: 75.5,
      discountAmount: 18.88,
      date: '2025-05-14T17:22:00Z',
    },
  ],
};

export default function DiscountDetailPage({
  params,
}: {
  params: { id: string };
}) {
  const router = useRouter();
  const [discount, setDiscount] = useState<any | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('details');

  useEffect(() => {
    // In a real implementation, this would fetch the discount data from an API
    const fetchDiscount = async () => {
      try {
        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 500));
        setDiscount(mockDiscount);
      } catch (error) {
        console.error('Error fetching discount:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchDiscount();
  }, [params.id]);

  if (loading) {
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
    discount.isActive && new Date(discount.expiresAt) > new Date();
  const isExpired = new Date(discount.expiresAt) < new Date();
  const isUpcoming = new Date(discount.startsAt) > new Date();

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
    ? Math.min(100, (discount.usedCount / discount.usageLimit) * 100)
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
          <Link href={`/admin/discounts/${params.id}/edit`}>
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
                      {discount.minPurchaseAmount
                        ? `$${discount.minPurchaseAmount.toFixed(2)}`
                        : 'No minimum'}
                    </span>
                  </div>
                  <div>
                    <span className='block text-sm font-medium text-gray-500'>
                      Maximum Discount
                    </span>
                    <span className='mt-1 block'>
                      {discount.maxDiscountAmount
                        ? `$${discount.maxDiscountAmount.toFixed(2)}`
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

              {discount.products && discount.products.length > 0 ? (
                <div className='overflow-x-auto'>
                  <table className='min-w-full divide-y divide-gray-200'>
                    <thead className='bg-gray-50'>
                      <tr>
                        <th
                          scope='col'
                          className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
                        >
                          Product
                        </th>
                        <th
                          scope='col'
                          className='px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider'
                        >
                          Price
                        </th>
                        <th
                          scope='col'
                          className='px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider'
                        >
                          Discounted Price
                        </th>
                      </tr>
                    </thead>
                    <tbody className='bg-white divide-y divide-gray-200'>
                      {discount.products.map((product: any) => (
                        <tr key={product.id}>
                          <td className='px-6 py-4 whitespace-nowrap'>
                            <div className='text-sm font-medium text-gray-900'>
                              {product.name}
                            </div>
                            <div className='text-xs text-gray-500'>
                              {product.id}
                            </div>
                          </td>
                          <td className='px-6 py-4 whitespace-nowrap text-right text-sm text-gray-500'>
                            ${product.price.toFixed(2)}
                          </td>
                          <td className='px-6 py-4 whitespace-nowrap text-right text-sm text-gray-900'>
                            {discount.discountType === 'percentage'
                              ? `$${(product.price * (1 - discount.discountValue / 100)).toFixed(2)}`
                              : `$${Math.max(0, product.price - discount.discountValue).toFixed(2)}`}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className='bg-gray-50 rounded-md p-8 text-center'>
                  <p className='text-gray-500'>
                    This discount applies to all products.
                  </p>
                </div>
              )}
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

              {discount.categories && discount.categories.length > 0 ? (
                <div className='overflow-x-auto'>
                  <table className='min-w-full divide-y divide-gray-200'>
                    <thead className='bg-gray-50'>
                      <tr>
                        <th
                          scope='col'
                          className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
                        >
                          Category Name
                        </th>
                        <th
                          scope='col'
                          className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
                        >
                          Category ID
                        </th>
                      </tr>
                    </thead>
                    <tbody className='bg-white divide-y divide-gray-200'>
                      {discount.categories.map((category: any) => (
                        <tr key={category.id}>
                          <td className='px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900'>
                            {category.name}
                          </td>
                          <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                            {category.id}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className='bg-gray-50 rounded-md p-8 text-center'>
                  <p className='text-gray-500'>
                    No specific categories selected for this discount.
                  </p>
                </div>
              )}
            </div>
          )}

          {activeTab === 'usage' && (
            <div>
              <div className='mb-6'>
                <h3 className='text-lg font-medium mb-3'>Usage Overview</h3>
                <div className='grid grid-cols-1 md:grid-cols-3 gap-4'>
                  <div className='p-4 border rounded-md'>
                    <div className='text-sm font-medium text-gray-500'>
                      Total Uses
                    </div>
                    <div className='mt-1 text-3xl font-semibold'>
                      {discount.usedCount}
                    </div>
                  </div>
                  <div className='p-4 border rounded-md'>
                    <div className='text-sm font-medium text-gray-500'>
                      Remaining
                    </div>
                    <div className='mt-1 text-3xl font-semibold'>
                      {discount.usageLimit
                        ? discount.usageLimit - discount.usedCount
                        : '∞'}
                    </div>
                  </div>
                  <div className='p-4 border rounded-md'>
                    <div className='text-sm font-medium text-gray-500'>
                      Usage Limit
                    </div>
                    <div className='mt-1 text-3xl font-semibold'>
                      {discount.usageLimit || '∞'}
                    </div>
                  </div>
                </div>
              </div>

              <div>
                <h3 className='text-lg font-medium mb-3'>Recent Usage</h3>

                {discount.usageHistory && discount.usageHistory.length > 0 ? (
                  <div className='overflow-x-auto'>
                    <table className='min-w-full divide-y divide-gray-200'>
                      <thead className='bg-gray-50'>
                        <tr>
                          <th
                            scope='col'
                            className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
                          >
                            Order ID
                          </th>
                          <th
                            scope='col'
                            className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
                          >
                            Customer
                          </th>
                          <th
                            scope='col'
                            className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
                          >
                            Date
                          </th>
                          <th
                            scope='col'
                            className='px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider'
                          >
                            Order Total
                          </th>
                          <th
                            scope='col'
                            className='px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider'
                          >
                            Discount Applied
                          </th>
                        </tr>
                      </thead>
                      <tbody className='bg-white divide-y divide-gray-200'>
                        {discount.usageHistory.map((usage: any) => (
                          <tr key={usage.id}>
                            <td className='px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900'>
                              {usage.orderId}
                            </td>
                            <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                              {usage.customerName}
                            </td>
                            <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                              {dayjs(usage.date).format('MMM D, YYYY h:mm A')}
                            </td>
                            <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-right'>
                              ${usage.amount.toFixed(2)}
                            </td>
                            <td className='px-6 py-4 whitespace-nowrap text-sm text-red-600 font-medium text-right'>
                              -${usage.discountAmount.toFixed(2)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <div className='bg-gray-50 rounded-md p-8 text-center'>
                    <p className='text-gray-500'>
                      This discount hasn't been used yet.
                    </p>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
