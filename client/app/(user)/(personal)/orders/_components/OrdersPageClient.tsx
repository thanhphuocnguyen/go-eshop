'use client';

import { useMemo } from 'react';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { OrderStatus } from '@/app/lib/definitions/common';
import { OrderListModel, OrdersStats } from '@/app/lib/definitions/order';
import Link from 'next/link';
import useSWR from 'swr';
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  FunnelIcon,
} from '@heroicons/react/24/outline';
import dayjs from 'dayjs';
import clsx from 'clsx';
import { useSearchParams, useRouter, usePathname } from 'next/navigation';
import { formatCurrency } from '@/app/utils';

const getStatusBadgeColor = (status: OrderStatus) => {
  switch (status) {
    case OrderStatus.Pending:
      return 'bg-yellow-100 text-yellow-800';
    case OrderStatus.Confirmed:
      return 'bg-blue-100 text-blue-800';
    case OrderStatus.Delivering:
      return 'bg-purple-100 text-purple-800';
    case OrderStatus.Delivered:
      return 'bg-green-100 text-green-800';
    case OrderStatus.Cancelled:
      return 'bg-red-100 text-red-800';
    case OrderStatus.Refunded:
      return 'bg-orange-100 text-orange-800';
    case OrderStatus.Completed:
      return 'bg-green-100 text-green-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

const ITEMS_PER_PAGE = 10;
const ALL_STATUSES = 'all';

export default function OrdersPageClient() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  // Get query parameters
  const currentPage = parseInt(searchParams.get('page') || '1');
  const currentStatus = searchParams.get('status') || ALL_STATUSES;

  // Fetch orders data
  const { data, isLoading, error } = useSWR(
    PUBLIC_API_PATHS.ORDERS,
    clientSideFetch<OrderListModel[]>
  );

  const orders = data?.data || [];

  // Calculate order statistics
  const stats = useMemo(() => {
    const initialStats: OrdersStats = {
      total: 0,
      pending: 0,
      confirm: 0,
      delivering: 0,
      delivered: 0,
      cancelled: 0,
      refunded: 0,
      completed: 0,
      totalSpent: 0,
    };

    return orders.reduce((acc, order) => {
      acc.total++;
      acc.totalSpent += order.total;

      switch (order.status) {
        case OrderStatus.Pending:
          acc.pending++;
          break;
        case OrderStatus.Confirmed:
          acc.confirm++;
          break;
        case OrderStatus.Delivering:
          acc.delivering++;
          break;
        case OrderStatus.Delivered:
          acc.delivered++;
          break;
        case OrderStatus.Cancelled:
          acc.cancelled++;
          break;
        case OrderStatus.Refunded:
          acc.refunded++;
          break;
        case OrderStatus.Completed:
          acc.completed++;
          break;
      }

      return acc;
    }, initialStats);
  }, [orders]);

  // Filter orders by status
  const filteredOrders = useMemo(() => {
    if (currentStatus === ALL_STATUSES) {
      return orders;
    }
    return orders.filter((order) => order.status === currentStatus);
  }, [orders, currentStatus]);

  // Calculate pagination
  const totalPages = Math.ceil(filteredOrders.length / ITEMS_PER_PAGE);
  const paginatedOrders = useMemo(() => {
    const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
    return filteredOrders.slice(startIndex, startIndex + ITEMS_PER_PAGE);
  }, [filteredOrders, currentPage]);

  // Update URL with query parameters
  const setQueryParams = (params: { status?: string; page?: number }) => {
    const newParams = new URLSearchParams(searchParams.toString());

    if (params.status) {
      newParams.set('status', params.status);
      // Reset to page 1 when changing filter
      newParams.set('page', '1');
    }

    if (params.page) {
      newParams.set('page', params.page.toString());
    }

    router.push(`${pathname}?${newParams.toString()}`);
  };

  if (isLoading) {
    return (
      <div className='flex justify-center items-center p-8'>
        <div className='animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600'></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className='bg-red-50 border border-red-200 text-red-800 rounded-md p-4'>
        <p className='font-medium'>Error loading orders</p>
        <p className='text-sm'>Please try again later</p>
      </div>
    );
  }

  if (orders.length === 0) {
    return (
      <div className='text-center py-12 border-2 border-dashed border-gray-300 rounded-lg'>
        <h3 className='text-lg font-medium text-gray-900'>No orders yet</h3>
        <p className='mt-1 text-sm text-gray-500'>
          When you place orders, they will appear here.
        </p>
        <div className='mt-6'>
          <Link
            href='/shop'
            className='text-indigo-600 font-medium hover:text-indigo-500'
          >
            Browse products
            <span aria-hidden='true'> &rarr;</span>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className='max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8'>
      <h1 className='text-2xl font-bold text-gray-900 mb-6'>Your Orders</h1>

      {/* Order Statistics */}
      <div className='bg-white shadow rounded-lg mb-8 p-6'>
        <h2 className='text-lg font-medium text-gray-900 mb-4'>
          Order Statistics
        </h2>
        <div className='grid grid-cols-2 md:grid-cols-4 gap-4'>
          <div className='bg-indigo-50 p-4 rounded-lg'>
            <p className='text-sm text-indigo-600 font-medium'>Total Orders</p>
            <p className='text-2xl font-bold'>{stats.total}</p>
          </div>
          <div className='bg-green-50 p-4 rounded-lg'>
            <p className='text-sm text-green-600 font-medium'>
              Completed Orders
            </p>
            <p className='text-2xl font-bold'>{stats.completed}</p>
          </div>
          <div className='bg-yellow-50 p-4 rounded-lg'>
            <p className='text-sm text-yellow-600 font-medium'>In Progress</p>
            <p className='text-2xl font-bold'>
              {stats.pending + stats.confirm + stats.delivering}
            </p>
          </div>
          <div className='bg-gray-50 p-4 rounded-lg'>
            <p className='text-sm text-gray-600 font-medium'>Total Spent</p>
            <p className='text-2xl font-bold'>
              ${formatCurrency(stats.totalSpent)}
            </p>
          </div>
        </div>
      </div>

      {/* Filter Controls */}
      <div className='bg-white shadow rounded-lg p-4 mb-6'>
        <div className='flex flex-col sm:flex-row sm:items-center justify-between'>
          <div className='flex items-center mb-4 sm:mb-0'>
            <FunnelIcon className='h-5 w-5 text-gray-400 mr-2' />
            <span className='text-sm font-medium text-gray-700'>
              Filter by status:
            </span>
          </div>
          <div className='flex flex-wrap gap-2'>
            <button
              onClick={() => setQueryParams({ status: ALL_STATUSES })}
              className={clsx(
                'px-3 py-1 text-xs font-medium rounded-full',
                currentStatus === ALL_STATUSES
                  ? 'bg-indigo-100 text-indigo-800 ring-1 ring-indigo-500'
                  : 'bg-gray-100 text-gray-800 hover:bg-gray-200'
              )}
            >
              All
            </button>
            {Object.values(OrderStatus).map((status) => (
              <button
                key={status}
                onClick={() => setQueryParams({ status })}
                className={clsx(
                  'px-3 py-1 text-xs font-medium rounded-full',
                  currentStatus === status
                    ? 'ring-1 ring-indigo-500 ' + getStatusBadgeColor(status)
                    : getStatusBadgeColor(status).replace(
                        'bg-',
                        'bg-opacity-50 bg-'
                      )
                )}
              >
                {status.charAt(0).toUpperCase() + status.slice(1)}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Orders Table */}
      <div className='overflow-hidden'>
        <div className='border border-gray-200 rounded-lg shadow-sm divide-y divide-gray-200'>
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
                  Date
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
                  Items
                </th>
                <th
                  scope='col'
                  className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
                >
                  Total
                </th>
                <th scope='col' className='relative px-6 py-3'>
                  <span className='sr-only'>Details</span>
                </th>
              </tr>
            </thead>
            <tbody className='bg-white divide-y divide-gray-200'>
              {paginatedOrders.map((order) => (
                <tr key={order.id} className='hover:bg-gray-50'>
                  <td className='px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900'>
                    #{order.id.substring(0, 8)}
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                    {dayjs(order.createdAt).format('MMM DD, YYYY')}
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap'>
                    <span
                      className={clsx(
                        'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                        getStatusBadgeColor(order.status)
                      )}
                    >
                      {order.status}
                    </span>
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                    {order.totalItems}
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900 font-medium'>
                    ${formatCurrency(order.total)}
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap text-right text-sm font-medium'>
                    <Link
                      href={`/orders/${order.id}`}
                      className='text-indigo-600 hover:text-indigo-900'
                    >
                      View details
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className='bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6 mt-4 rounded-lg shadow-sm'>
            <div className='flex-1 flex justify-between sm:hidden'>
              <button
                onClick={() =>
                  setQueryParams({ page: Math.max(currentPage - 1, 1) })
                }
                disabled={currentPage === 1}
                className={`${
                  currentPage === 1
                    ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                    : 'bg-white text-gray-700 hover:bg-gray-50'
                } relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md`}
              >
                Previous
              </button>
              <button
                onClick={() =>
                  setQueryParams({
                    page: Math.min(currentPage + 1, totalPages),
                  })
                }
                disabled={currentPage === totalPages}
                className={`${
                  currentPage === totalPages
                    ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                    : 'bg-white text-gray-700 hover:bg-gray-50'
                } ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md`}
              >
                Next
              </button>
            </div>
            <div className='hidden sm:flex-1 sm:flex sm:items-center sm:justify-between'>
              <div>
                <p className='text-sm text-gray-700'>
                  Showing{' '}
                  <span className='font-medium'>
                    {(currentPage - 1) * ITEMS_PER_PAGE + 1}
                  </span>{' '}
                  to{' '}
                  <span className='font-medium'>
                    {Math.min(
                      currentPage * ITEMS_PER_PAGE,
                      filteredOrders.length
                    )}
                  </span>{' '}
                  of{' '}
                  <span className='font-medium'>{filteredOrders.length}</span>{' '}
                  orders
                </p>
              </div>
              <div>
                <nav
                  className='relative z-0 inline-flex rounded-md shadow-sm -space-x-px'
                  aria-label='Pagination'
                >
                  <button
                    onClick={() =>
                      setQueryParams({ page: Math.max(currentPage - 1, 1) })
                    }
                    disabled={currentPage === 1}
                    className={`${
                      currentPage === 1
                        ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                        : 'bg-white text-gray-500 hover:bg-gray-50'
                    } relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 text-sm font-medium`}
                  >
                    <span className='sr-only'>Previous</span>
                    <ChevronLeftIcon className='h-5 w-5' aria-hidden='true' />
                  </button>

                  {Array.from({ length: totalPages }, (_, i) => i + 1).map(
                    (page) => (
                      <button
                        key={page}
                        onClick={() => setQueryParams({ page })}
                        className={clsx(
                          page === currentPage
                            ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600'
                            : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50',
                          'relative inline-flex items-center px-4 py-2 border text-sm font-medium'
                        )}
                      >
                        {page}
                      </button>
                    )
                  )}

                  <button
                    onClick={() =>
                      setQueryParams({
                        page: Math.min(currentPage + 1, totalPages),
                      })
                    }
                    disabled={currentPage === totalPages}
                    className={`${
                      currentPage === totalPages
                        ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                        : 'bg-white text-gray-500 hover:bg-gray-50'
                    } relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 text-sm font-medium`}
                  >
                    <span className='sr-only'>Next</span>
                    <ChevronRightIcon className='h-5 w-5' aria-hidden='true' />
                  </button>
                </nav>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
