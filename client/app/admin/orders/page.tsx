'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import LoadingInline from '@/components/Common/Loadings/LoadingInline';
import { Breadcrumb } from '@/components/Common';
import { apiFetch } from '@/lib/apis/api';
import { ADMIN_API_PATHS } from '@/lib/constants/api';
import { GenericResponse, Order } from '@/lib/definitions';

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
      const { data, error, pagination } = await apiFetch<
        GenericResponse<Order[]>
      >(ADMIN_API_PATHS.ORDERS, {
        queryParams: {
          page: currentPage,
          page_size: 10,
          status: statusFilter !== 'all' ? statusFilter : undefined,
        },
      });
      
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
        <div className='flex justify-center py-10'>
          <LoadingInline />
        </div>
      );
    }

    if (error) {
      return (
        <div className='bg-red-100 p-4 rounded-md text-red-700 mb-4'>
          <p>{error}</p>
          <button
            className='mt-2 text-red-700 font-medium hover:text-red-800'
            onClick={fetchOrders}
          >
            Try Again
          </button>
        </div>
      );
    }

    if (orders.length === 0) {
      return (
        <div className='p-4 text-center text-gray-500'>No orders found.</div>
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
                  <div>{order.customer_name}</div>
                  <div className='text-xs text-gray-500'>
                    {order.customer_email}
                  </div>
                </td>
                <td className='px-6 py-4 whitespace-nowrap text-sm'>
                  {formatCurrency(order.total)}
                </td>
                <td className='px-6 py-4 whitespace-nowrap text-sm text-center'>
                  {order.total_items}
                </td>
                <td className='px-6 py-4 whitespace-nowrap'>
                  <span
                    className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(order.status)}`}
                  >
                    {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                  </span>
                </td>
                <td className='px-6 py-4 whitespace-nowrap'>
                  <span
                    className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getPaymentStatusColor(order.payment_status)}`}
                  >
                    {order.payment_status.charAt(0).toUpperCase() + order.payment_status.slice(1)}
                  </span>
                </td>
                <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                  {formatDate(order.created_at)}
                </td>
                <td className='px-6 py-4 whitespace-nowrap text-sm text-right'>
                  <button
                    onClick={() => navigateToOrderDetail(order.id)}
                    className='text-indigo-600 hover:text-indigo-900 font-medium'
                  >
                    View Details
                  </button>
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
    return (
      <div className='flex justify-between items-center px-6 py-4 bg-white border-t border-gray-200'>
        <div className='text-sm text-gray-500'>
          Showing page {currentPage} of {totalPages}
        </div>
        <div className='flex space-x-2'>
          <button
            onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
            disabled={currentPage === 1}
            className={`px-3 py-1 rounded ${currentPage === 1 ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'}`}
          >
            Previous
          </button>
          <button
            onClick={() =>
              setCurrentPage((prev) => Math.min(prev + 1, totalPages))
            }
            disabled={currentPage === totalPages}
            className={`px-3 py-1 rounded ${currentPage === totalPages ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'}`}
          >
            Next
          </button>
        </div>
      </div>
    );
  };

  // Render filter controls
  const renderFilterControls = () => {
    return (
      <div className='mb-4 p-4 bg-white rounded-lg shadow-sm'>
        <div className='flex flex-wrap items-center justify-center gap-4'>
          <div>
            <label
              htmlFor='statusFilter'
              className='block text-sm font-medium text-gray-700 mb-1'
            >
              Status Filter
            </label>
            <select
              id='statusFilter'
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value);
                setCurrentPage(1); // Reset to first page on filter change
              }}
              className='border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500'
            >
              <option value='all'>All Statuses</option>
              <option value='pending'>Pending</option>
              <option value='confirmed'>Confirmed</option>
              <option value='delivering'>Delivering</option>
              <option value='delivered'>Delivered</option>
              <option value='completed'>Completed</option>
              <option value='cancelled'>Cancelled</option>
              <option value='refunded'>Refunded</option>
            </select>
          </div>

          <button
            onClick={() => fetchOrders()}
            className='px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm font-medium'
          >
            Refresh
          </button>
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
        <div className='fixed bottom-4 right-4 bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded shadow-md'>
          <p className='font-bold'>Error</p>
          <p>{error}</p>
          <button
            onClick={() => setError(null)}
            className='absolute top-2 right-2 text-red-500 hover:text-red-700'
          >
            &times;
          </button>
        </div>
      )}
    </div>
  );
}
