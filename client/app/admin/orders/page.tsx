'use client';

import { useState, useEffect, Fragment } from 'react';
import { useRouter } from 'next/navigation';
import LoadingInline from '@/app/components/Common/Loadings/LoadingInline';
import { Breadcrumb } from '@/app/components/Common';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { Order } from '@/app/lib/definitions';
import { Button, Menu, Transition } from '@headlessui/react';

export default function AdminOrdersPage() {
  const router = useRouter();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [statusFilter, setStatusFilter] = useState('all');

  // Fetch orders on component mount and when page or filter changes
  useEffect(() => {
    fetchOrders();
  }, [currentPage, statusFilter]);

  const fetchOrders = async () => {
    setLoading(true);
    setError(null);

    try {
      // Adjust this API endpoint to your actual backend endpoint
      const { data, error, pagination } = await clientSideFetch<Order[]>(
        ADMIN_API_PATHS.ORDERS,
        {
          queryParams: {
            page: currentPage,
            pageSize: 10,
            status: statusFilter !== 'all' ? statusFilter : undefined,
          },
        }
      );

      if (error) {
        setError(error.details || 'Failed to fetch orders');
        setLoading(false);
        return;
      }

      if (data) {
        setOrders(data);
        if (pagination) {
          setTotalPages(pagination.totalPages || 1);
        }
      }
    } catch (err) {
      setError('An error occurred while fetching orders');
    } finally {
      setLoading(false);
    }
  };

  const navigateToOrderDetail = (orderId: string) => {
    router.push(`/admin/orders/${orderId}`);
  };

  // Format date to a more readable format
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  // Format currency
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount);
  };

  // Get color based on order status
  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'confirmed':
        return 'bg-blue-100 text-blue-800';
      case 'delivering':
        return 'bg-indigo-100 text-indigo-800';
      case 'delivered':
        return 'bg-green-100 text-green-800';
      case 'completed':
        return 'bg-green-500 text-white';
      case 'cancelled':
        return 'bg-gray-100 text-gray-800';
      case 'refunded':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  // Get color based on payment status
  const getPaymentStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'success':
        return 'bg-green-500 text-white';
      case 'failed':
        return 'bg-red-500 text-white';
      case 'refunded':
        return 'bg-purple-100 text-purple-800';
      case 'cancelled':
        return 'bg-gray-300 text-gray-800';
      case 'authorized':
        return 'bg-blue-100 text-blue-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  // Render the orders list
  const renderOrdersList = () => {
    if (loading && orders.length === 0) {
      return (
        <div className='flex flex-col items-center justify-center py-16'>
          <LoadingInline />
          <p className='mt-4 text-sm text-gray-500'>Loading orders...</p>
        </div>
      );
    }

    if (error) {
      return (
        <div className='bg-red-50 p-6 rounded-md text-red-700 mb-4 border border-red-200 flex items-start gap-4'>
          <div className='flex-shrink-0 mt-0.5'>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-6 w-6 text-red-600'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
              />
            </svg>
          </div>
          <div className='flex-1'>
            <h3 className='text-sm font-medium text-red-800 mb-1'>
              Error loading orders
            </h3>
            <p className='text-sm'>{error}</p>
            <Button
              className='mt-3 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 text-sm font-medium transition-colors duration-150 ease-in-out inline-flex items-center gap-2 shadow-sm'
              onClick={fetchOrders}
            >
              <svg
                xmlns='http://www.w3.org/2000/svg'
                className='h-4 w-4'
                fill='none'
                viewBox='0 0 24 24'
                stroke='currentColor'
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  strokeWidth={2}
                  d='M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15'
                />
              </svg>
              Try Again
            </Button>
          </div>
        </div>
      );
    }

    if (orders.length === 0) {
      return (
        <div className='p-12 text-center'>
          <div className='inline-block p-4 bg-gray-50 rounded-full mb-4'>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-12 w-12 text-gray-400'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z'
              />
            </svg>
          </div>
          <h3 className='text-lg font-medium text-gray-700 mb-1'>
            No orders found
          </h3>
          <p className='text-gray-500 mb-6'>
            Try changing your filter criteria or check back later.
          </p>
          <Button
            onClick={() => {
              setStatusFilter('all');
              fetchOrders();
            }}
            className='px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 text-sm font-medium inline-flex items-center gap-2 transition-colors duration-150 ease-in-out shadow-sm'
          >
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z'
              />
            </svg>
            Clear Filters
          </Button>
        </div>
      );
    }

    return (
      <div className='overflow-x-auto'>
        <table className='min-w-full divide-y divide-gray-200'>
          <thead className='bg-gray-50'>
            <tr>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Order ID
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Customer
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Total
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Items
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Status
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Payment Status
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Date
              </th>
              <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                Actions
              </th>
            </tr>
          </thead>
          <tbody className='bg-white divide-y divide-gray-200'>
            {orders.map((order) => (
              <tr key={order.id} className='hover:bg-gray-50'>
                <td className='px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600'>
                  {order.id.substring(0, 8)}...
                </td>
                <td className='px-6 py-4 whitespace-nowrap text-sm'>
                  <div>{order.customerName}</div>
                  <div className='text-xs text-gray-500'>
                    {order.customerEmail}
                  </div>
                </td>
                <td className='px-6 py-4 whitespace-nowrap text-sm'>
                  {formatCurrency(order.total)}
                </td>
                <td className='px-6 py-4 whitespace-nowrap text-sm text-center'>
                  {order.totalItems}
                </td>
                <td className='px-6 py-4 whitespace-nowrap'>
                  <span
                    className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(order.status)}`}
                  >
                    {order.status.charAt(0).toUpperCase() +
                      order.status.slice(1)}
                  </span>
                </td>
                <td className='px-6 py-4 whitespace-nowrap'>
                  <span
                    className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getPaymentStatusColor(order.paymentStatus)}`}
                  >
                    {order.paymentStatus.charAt(0).toUpperCase() +
                      order.paymentStatus.slice(1)}
                  </span>
                </td>
                <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                  {formatDate(order.createdAt)}
                </td>
                <td className='px-6 py-4 whitespace-nowrap text-sm text-right'>
                  <Button
                    onClick={() => navigateToOrderDetail(order.id)}
                    className='text-indigo-600 hover:text-indigo-900 font-medium inline-flex items-center gap-1 transition-colors duration-150'
                  >
                    <span>View Details</span>
                    <svg
                      xmlns='http://www.w3.org/2000/svg'
                      className='h-4 w-4'
                      fill='none'
                      viewBox='0 0 24 24'
                      stroke='currentColor'
                    >
                      <path
                        strokeLinecap='round'
                        strokeLinejoin='round'
                        strokeWidth={2}
                        d='M9 5l7 7-7 7'
                      />
                    </svg>
                  </Button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  };

  // Render pagination controls
  const renderPagination = () => {
    // Generate page numbers array for pagination
    const getPageNumbers = () => {
      const pages = [];
      const maxPagesToShow = 5;

      if (totalPages <= maxPagesToShow) {
        // Show all pages if total pages is less than or equal to max pages to show
        for (let i = 1; i <= totalPages; i++) {
          pages.push(i);
        }
      } else {
        // Always show first page
        pages.push(1);

        // Calculate start and end page numbers to show
        let startPage = Math.max(2, currentPage - 1);
        let endPage = Math.min(totalPages - 1, currentPage + 1);

        // Adjust if we're near the beginning
        if (currentPage <= 3) {
          endPage = 4;
        }

        // Adjust if we're near the end
        if (currentPage >= totalPages - 2) {
          startPage = totalPages - 3;
        }

        // Add ellipsis after first page if needed
        if (startPage > 2) {
          pages.push('...');
        }

        // Add middle pages
        for (let i = startPage; i <= endPage; i++) {
          pages.push(i);
        }

        // Add ellipsis before last page if needed
        if (endPage < totalPages - 1) {
          pages.push('...');
        }

        // Always show last page
        pages.push(totalPages);
      }

      return pages;
    };

    return (
      <div className='flex flex-col sm:flex-row justify-between items-center px-6 py-4 bg-white border-t border-gray-200 gap-4'>
        <div className='text-sm text-gray-500'>
          Showing page{' '}
          <span className='font-medium text-gray-700'>{currentPage}</span> of{' '}
          <span className='font-medium text-gray-700'>{totalPages}</span>
        </div>

        <nav className='flex items-center'>
          <Button
            onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
            disabled={currentPage === 1}
            className={`
              relative inline-flex items-center px-3 py-2 rounded-l-md border text-sm font-medium
              ${
                currentPage === 1
                  ? 'bg-gray-100 text-gray-400 cursor-not-allowed border-gray-200'
                  : 'bg-white text-gray-700 hover:bg-gray-50 border-gray-300'
              }
            `}
          >
            <span className='sr-only'>Previous</span>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-5 w-5'
              viewBox='0 0 20 20'
              fill='currentColor'
            >
              <path
                fillRule='evenodd'
                d='M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z'
                clipRule='evenodd'
              />
            </svg>
          </Button>

          <div className='hidden sm:flex'>
            {getPageNumbers().map((page, index) =>
              typeof page === 'number' ? (
                <Button
                  key={index}
                  onClick={() => setCurrentPage(page)}
                  className={`
                    relative inline-flex items-center px-4 py-2 border text-sm font-medium
                    ${
                      currentPage === page
                        ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600'
                        : 'bg-white border-gray-300 text-gray-700 hover:bg-gray-50'
                    }
                  `}
                >
                  {page}
                </Button>
              ) : (
                <span
                  key={index}
                  className='relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700'
                >
                  {page}
                </span>
              )
            )}
          </div>

          <Button
            onClick={() =>
              setCurrentPage((prev) => Math.min(prev + 1, totalPages))
            }
            disabled={currentPage === totalPages}
            className={`
              relative inline-flex items-center px-3 py-2 rounded-r-md border text-sm font-medium
              ${
                currentPage === totalPages
                  ? 'bg-gray-100 text-gray-400 cursor-not-allowed border-gray-200'
                  : 'bg-white text-gray-700 hover:bg-gray-50 border-gray-300'
              }
            `}
          >
            <span className='sr-only'>Next</span>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-5 w-5'
              viewBox='0 0 20 20'
              fill='currentColor'
            >
              <path
                fillRule='evenodd'
                d='M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z'
                clipRule='evenodd'
              />
            </svg>
          </Button>
        </nav>
      </div>
    );
  };

  // Render filter controls
  const renderFilterControls = () => {
    const statuses = [
      { value: 'all', label: 'All Statuses', icon: 'filter' },
      { value: 'pending', label: 'Pending', icon: 'clock' },
      { value: 'confirmed', label: 'Confirmed', icon: 'check' },
      { value: 'delivering', label: 'Delivering', icon: 'truck' },
      { value: 'delivered', label: 'Delivered', icon: 'package' },
      { value: 'completed', label: 'Completed', icon: 'check-circle' },
      { value: 'cancelled', label: 'Cancelled', icon: 'x-circle' },
      { value: 'refunded', label: 'Refunded', icon: 'currency-dollar' },
    ];

    // Function to render the icon based on status
    const renderStatusIcon = (iconName: string) => {
      switch (iconName) {
        case 'filter':
          return (
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z'
              />
            </svg>
          );
        case 'clock':
          return (
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z'
              />
            </svg>
          );
        case 'check':
          return (
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M5 13l4 4L19 7'
              />
            </svg>
          );
        case 'truck':
          return (
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0'
              />
            </svg>
          );
        case 'package':
          return (
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4'
              />
            </svg>
          );
        case 'check-circle':
          return (
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
              />
            </svg>
          );
        case 'x-circle':
          return (
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z'
              />
            </svg>
          );
        case 'currency-dollar':
          return (
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-4 w-4'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
              />
            </svg>
          );
        default:
          return null;
      }
    };

    return (
      <div className='mb-6 p-5 bg-white rounded-lg shadow-sm border border-gray-100'>
        <div className='flex flex-wrap items-center justify-between gap-4'>
          <div className='w-72'>
            <label className='block text-sm font-medium text-gray-700 mb-1'>
              Status Filter
            </label>

            <div className='relative'>
              <Menu as='div' className='relative inline-block text-left w-full'>
                {({ open }) => (
                  <>
                    <div>
                      <Menu.Button className='w-full flex justify-between items-center px-4 py-2.5 text-sm border border-gray-300 rounded-md bg-white text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition duration-150 ease-in-out'>
                        <div className='flex items-center gap-2'>
                          {renderStatusIcon(
                            statuses.find((s) => s.value === statusFilter)
                              ?.icon || 'filter'
                          )}
                          <span className='font-medium'>
                            {statuses.find((s) => s.value === statusFilter)
                              ?.label || 'All Statuses'}
                          </span>
                        </div>
                        <svg
                          className={`ml-2 h-5 w-5 text-gray-400 transition-transform duration-150 ${open ? 'transform rotate-180' : ''}`}
                          xmlns='http://www.w3.org/2000/svg'
                          viewBox='0 0 20 20'
                          fill='currentColor'
                          aria-hidden='true'
                        >
                          <path
                            fillRule='evenodd'
                            d='M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z'
                            clipRule='evenodd'
                          />
                        </svg>
                      </Menu.Button>
                    </div>

                    <Transition
                      as={Fragment}
                      enter='transition ease-out duration-100'
                      enterFrom='transform opacity-0 scale-95'
                      enterTo='transform opacity-100 scale-100'
                      leave='transition ease-in duration-75'
                      leaveFrom='transform opacity-100 scale-100'
                      leaveTo='transform opacity-0 scale-95'
                    >
                      <Menu.Items className='absolute z-10 w-full mt-1 bg-white border border-gray-200 rounded-md shadow-lg py-1 max-h-60 overflow-auto focus:outline-none ring-1 ring-black ring-opacity-5'>
                        {statuses.map((status) => (
                          <Menu.Item key={status.value}>
                            {({ active }) => (
                              <button
                                onClick={() => {
                                  setStatusFilter(status.value);
                                  setCurrentPage(1); // Reset to first page on filter change
                                }}
                                className={`
                                  ${active || status.value === statusFilter ? 'bg-indigo-50 text-indigo-700' : 'text-gray-700'}
                                  flex items-center gap-2 w-full text-left px-4 py-2.5 text-sm group
                                `}
                              >
                                <span
                                  className={`
                                  ${status.value === statusFilter ? 'text-indigo-600' : 'text-gray-400 group-hover:text-indigo-600'} 
                                  flex-shrink-0
                                `}
                                >
                                  {renderStatusIcon(status.icon)}
                                </span>
                                <span
                                  className={
                                    status.value === statusFilter
                                      ? 'font-medium'
                                      : ''
                                  }
                                >
                                  {status.label}
                                </span>
                                {status.value === statusFilter && (
                                  <span className='flex-shrink-0 ml-auto'>
                                    <svg
                                      className='h-4 w-4 text-indigo-600'
                                      xmlns='http://www.w3.org/2000/svg'
                                      viewBox='0 0 20 20'
                                      fill='currentColor'
                                    >
                                      <path
                                        fillRule='evenodd'
                                        d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z'
                                        clipRule='evenodd'
                                      />
                                    </svg>
                                  </span>
                                )}
                              </button>
                            )}
                          </Menu.Item>
                        ))}
                      </Menu.Items>
                    </Transition>
                  </>
                )}
              </Menu>
            </div>
          </div>

          <div className='flex items-center gap-3'>
            <Button
              onClick={() => fetchOrders()}
              className='px-4 py-2.5 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 text-sm font-medium inline-flex items-center transition duration-150 ease-in-out shadow-sm'
            >
              <svg
                xmlns='http://www.w3.org/2000/svg'
                className='h-4 w-4 mr-2'
                fill='none'
                viewBox='0 0 24 24'
                stroke='currentColor'
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  strokeWidth={2}
                  d='M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15'
                />
              </svg>
              Refresh
            </Button>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className='container mx-auto px-4 py-8'>
      <div className='mb-8'>
        <Breadcrumb
          items={[
            { label: 'Admin', href: '/admin' },
            { label: 'Orders', href: '/admin/orders' },
          ]}
        />
        <h1 className='text-3xl font-bold text-gray-900 mt-2'>
          Orders Management
        </h1>
      </div>

      {renderFilterControls()}
      <div className='bg-white rounded-lg shadow-sm overflow-hidden'>
        {renderOrdersList()}
        {orders.length > 0 && renderPagination()}
      </div>

      {error && (
        <Transition
          show={!!error}
          as={Fragment}
          enter='transform transition ease-in-out duration-300'
          enterFrom='translate-y-full opacity-0'
          enterTo='translate-y-0 opacity-100'
          leave='transform transition ease-in-out duration-300'
          leaveFrom='translate-y-0 opacity-100'
          leaveTo='translate-y-full opacity-0'
        >
          <div className='fixed bottom-4 right-4 max-w-sm bg-white border border-red-200 shadow-lg rounded-lg overflow-hidden'>
            <div className='p-4 flex'>
              <div className='flex-shrink-0'>
                <svg
                  xmlns='http://www.w3.org/2000/svg'
                  className='h-5 w-5 text-red-500'
                  viewBox='0 0 20 20'
                  fill='currentColor'
                >
                  <path
                    fillRule='evenodd'
                    d='M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z'
                    clipRule='evenodd'
                  />
                </svg>
              </div>
              <div className='ml-3 flex-1'>
                <p className='text-sm font-medium text-gray-900'>Error</p>
                <p className='mt-1 text-sm text-gray-500'>{error}</p>
              </div>
              <div className='ml-4 flex-shrink-0 flex'>
                <Button
                  onClick={() => setError(null)}
                  className='bg-white rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
                >
                  <span className='sr-only'>Close</span>
                  <svg
                    xmlns='http://www.w3.org/2000/svg'
                    className='h-5 w-5'
                    viewBox='0 0 20 20'
                    fill='currentColor'
                  >
                    <path
                      fillRule='evenodd'
                      d='M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z'
                      clipRule='evenodd'
                    />
                  </svg>
                </Button>
              </div>
            </div>
            <div className='bg-red-50 px-4 py-3'>
              <div className='flex justify-end'>
                <Button
                  onClick={() => {
                    setError(null);
                    fetchOrders();
                  }}
                  className='px-3 py-1.5 bg-red-600 text-white text-xs font-medium rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500'
                >
                  Retry
                </Button>
              </div>
            </div>
          </div>
        </Transition>
      )}
    </div>
  );
}
