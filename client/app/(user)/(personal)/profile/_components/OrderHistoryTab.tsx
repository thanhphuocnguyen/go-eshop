'use client';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { OrderStatus } from '@/app/lib/definitions';
import Link from 'next/link';
import useSWR from 'swr';
import { ChevronRightIcon } from '@heroicons/react/24/outline';
import dayjs from 'dayjs';
import clsx from 'clsx';
import { formatCurrency } from '@/app/utils';

type OrderModel = {
  id: string;
  total: number;
  totalItems: number;
  status: OrderStatus;
  paymentStatus: string;
  createdAt: string;
  updatedAt: string;
};

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

export default function OrderHistoryTab() {
  const { data, isLoading, error } = useSWR(
    PUBLIC_API_PATHS.ORDERS,
    clientSideFetch<OrderModel[]>
  );

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

  const orders = data?.data || [];

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
            {orders.map((order) => (
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
                    className='text-indigo-600 hover:text-indigo-900 inline-flex items-center'
                  >
                    Details
                    <ChevronRightIcon className='ml-1 w-4 h-4' />
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className='mt-4 text-right'>
        <Link
          href='/orders'
          className='text-indigo-600 hover:text-indigo-900 text-sm font-medium'
        >
          View all orders
        </Link>
      </div>
    </div>
  );
}
