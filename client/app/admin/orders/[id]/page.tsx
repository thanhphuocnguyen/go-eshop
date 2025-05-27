'use client';

import { useState, useEffect, use } from 'react';
import Link from 'next/link';
import { Breadcrumb } from '@/app/components/Common';
import LoadingButton from '@/app/components/Common/LoadingButton';
import LoadingInline from '@/app/components/Common/Loadings/LoadingInline';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { OrderDetail, GenericResponse } from '@/app/lib/definitions';

export default function OrderDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const [orderDetail, setOrderDetail] = useState<OrderDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [changingStatus, setChangingStatus] = useState(false);
  console.log(orderDetail);
  useEffect(() => {
    fetchOrderDetail();
  }, []);

  const fetchOrderDetail = async () => {
    setLoading(true);
    setError(null);

    try {
      const { data, error } = await clientSideFetch<OrderDetail>(
        ADMIN_API_PATHS.ORDER_DETAIL.replace(':id', id),
        {}
      );

      if (error) {
        setError(error.details || 'Failed to fetch order details');
        setLoading(false);
        return;
      }

      if (data) {
        setOrderDetail(data);
      } else {
        setError('Order not found');
      }
    } catch (err) {
      setError('An error occurred while fetching the order details');
    } finally {
      setLoading(false);
    }
  };

  const changeOrderStatus = async (newStatus: string) => {
    setChangingStatus(true);
    setError(null);

    try {
      const { data, error } = await clientSideFetch<GenericResponse<boolean>>(
        ADMIN_API_PATHS.ORDER_DETAIL_STATUS.replace(':id', id),
        {
          method: 'PUT',
          body: { status: newStatus },
        }
      );

      if (error) {
        setError(error.details || 'Failed to change order status');
        setChangingStatus(false);
        return;
      }

      if (data) {
        // Update the order with the new status
        setOrderDetail((prev) =>
          prev ? { ...prev, status: newStatus } : null
        );
      } else {
        setError('Failed to change order status');
      }
    } catch (err) {
      setError('An error occurred while changing the order status');
    } finally {
      setChangingStatus(false);
    }
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

  // Get status options based on current status
  const getAvailableStatusOptions = (currentStatus: string) => {
    switch (currentStatus.toLowerCase()) {
      case 'pending':
        return ['confirmed', 'cancelled'];
      case 'confirmed':
        return ['delivering', 'cancelled'];
      case 'delivering':
        return ['delivered'];
      case 'delivered':
        return ['completed', 'refunded'];
      case 'completed':
        return ['refunded'];
      default:
        return [];
    }
  };

  if (loading) {
    return (
      <div className='container mx-auto px-4 py-8'>
        <div className='flex justify-center py-10'>
          <LoadingInline />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className='container mx-auto px-4 py-8'>
        <div className='bg-red-100 p-4 rounded-md text-red-700 mb-4'>
          <p>{error}</p>
          <button
            className='mt-2 text-red-700 font-medium hover:text-red-800'
            onClick={fetchOrderDetail}
          >
            Try Again
          </button>
        </div>
        <Link
          href='/admin/orders'
          className='text-blue-600 hover:text-blue-800'
        >
          &larr; Back to Orders
        </Link>
      </div>
    );
  }

  if (!orderDetail) {
    return (
      <div className='container mx-auto px-4 py-8'>
        <div className='p-4 text-center text-gray-500'>Order not found.</div>
        <Link
          href='/admin/orders'
          className='text-blue-600 hover:text-blue-800'
        >
          &larr; Back to Orders
        </Link>
      </div>
    );
  }

  return (
    <div className='container mx-auto px-4 py-8'>
      <div className='mb-8'>
        <Breadcrumb
          items={[
            { label: 'Admin', href: '/admin' },
            { label: 'Orders', href: '/admin/orders' },
            {
              label: `Order ${id.substring(0, 8)}...`,
              href: `/admin/orders/${id}`,
            },
          ]}
        />
        <div className='flex justify-between items-center mt-2'>
          <h1 className='text-3xl font-bold text-gray-900'>Order Details</h1>
          <Link
            href='/admin/orders'
            className='px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300'
          >
            Back to Orders
          </Link>
        </div>
      </div>

      <div className='bg-white p-6 rounded-lg shadow-lg'>
        <div className='grid grid-cols-1 md:grid-cols-2 gap-6 mb-6'>
          <div>
            <h3 className='font-semibold text-gray-700 mb-2'>
              Order Information
            </h3>
            <div className='bg-gray-50 p-4 rounded-md'>
              <p>
                <span className='font-medium'>Order ID:</span> {orderDetail.id}
              </p>
              <p>
                <span className='font-medium'>Date:</span>{' '}
                {formatDate(orderDetail.createdAt)}
              </p>
              <p>
                <span className='font-medium'>Total:</span>{' '}
                {formatCurrency(orderDetail.total)}
              </p>
              <p className='flex items-center mt-1'>
                <span className='font-medium mr-2'>Status:</span>
                <span
                  className={`px-2 py-1 text-xs leading-5 font-semibold rounded-full ${getStatusColor(orderDetail.status)}`}
                >
                  {orderDetail.status}
                </span>
              </p>
            </div>
          </div>

          <div>
            <h3 className='font-semibold text-gray-700 mb-2'>
              Customer Information
            </h3>
            <div className='bg-gray-50 p-4 rounded-md'>
              <p>
                <span className='font-medium'>Name:</span>{' '}
                {orderDetail.customerName}
              </p>
              <p>
                <span className='font-medium'>Email:</span>{' '}
                {orderDetail.customerEmail}
              </p>
              <p>
                <span className='font-medium'>Shipping Address:</span>
              </p>
              <p className='ml-4'>
                {orderDetail.shippingInfo.street},{' '}
                {orderDetail.shippingInfo.ward},
                <br />
                {orderDetail.shippingInfo.district},{' '}
                {orderDetail.shippingInfo.city}
              </p>
              <p>
                <span className='font-medium'>Phone:</span>{' '}
                {orderDetail.shippingInfo.phone}
              </p>
            </div>
          </div>
        </div>

        {orderDetail.paymentInfo && (
          <div className='mb-6'>
            <h3 className='font-semibold text-gray-700 mb-2'>
              Payment Information
            </h3>
            <div className='bg-gray-50 p-4 rounded-md'>
              <p>
                <span className='font-medium'>Payment ID:</span>{' '}
                {orderDetail.paymentInfo.id}
              </p>
              <p>
                <span className='font-medium'>Amount:</span>{' '}
                {formatCurrency(orderDetail.paymentInfo.amount)}
              </p>
              <p>
                <span className='font-medium'>Method:</span>{' '}
                {orderDetail.paymentInfo.method}
              </p>
              <p className='flex items-center mt-1'>
                <span className='font-medium mr-2'>Status:</span>
                <span
                  className={`px-2 py-1 text-xs leading-5 font-semibold rounded-full ${getStatusColor(orderDetail.paymentInfo.status)}`}
                >
                  {orderDetail.paymentInfo.status}
                </span>
              </p>
              {orderDetail.paymentInfo.gateway && (
                <p>
                  <span className='font-medium'>Gateway:</span>{' '}
                  {orderDetail.paymentInfo.gateway}
                </p>
              )}
              {orderDetail.paymentInfo.refundId && (
                <p>
                  <span className='font-medium'>Refund ID:</span>{' '}
                  {orderDetail.paymentInfo.refundId}
                </p>
              )}
            </div>
          </div>
        )}

        {/* Status change section */}
        {getAvailableStatusOptions(orderDetail.status).length > 0 && (
          <div className='mb-6'>
            <h3 className='font-semibold text-gray-700 mb-2'>
              Change Order Status
            </h3>
            <div className='flex flex-wrap gap-2'>
              {getAvailableStatusOptions(orderDetail.status).map((status) => (
                <LoadingButton
                  key={status}
                  className={`px-4 py-2 rounded-md text-sm font-medium ${
                    status === 'cancelled' || status === 'refunded'
                      ? 'bg-red-600 hover:bg-red-700 text-white'
                      : 'bg-blue-600 hover:bg-blue-700 text-white'
                  }`}
                  onClick={() => changeOrderStatus(status)}
                  isLoading={changingStatus}
                  disabled={changingStatus}
                >
                  Mark as {status}
                </LoadingButton>
              ))}
            </div>
          </div>
        )}

        <div>
          <h3 className='font-semibold text-gray-700 mb-4'>Order Items</h3>
          <div className='overflow-x-auto'>
            <table className='min-w-full divide-y divide-gray-200'>
              <thead className='bg-gray-50'>
                <tr>
                  <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                    Product
                  </th>
                  <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                    Attributes
                  </th>
                  <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                    Quantity
                  </th>
                  <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                    Price
                  </th>
                </tr>
              </thead>
              <tbody className='bg-white divide-y divide-gray-200'>
                {orderDetail.products.map((product) => (
                  <tr key={product.id}>
                    <td className='px-6 py-4 whitespace-nowrap'>
                      <div className='flex items-center'>
                        {product.imageUrl && (
                          <div className='flex-shrink-0 h-10 w-10 mr-4'>
                            <img
                              className='h-10 w-10 rounded-md object-cover'
                              src={product.imageUrl}
                              alt={product.name}
                            />
                          </div>
                        )}
                        <div>
                          <div className='text-sm font-medium text-gray-900'>
                            {product.name}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td className='px-6 py-4'>
                      <div className='text-sm text-gray-900'>
                        {product.attributesSnapshot?.map((attr, index) => (
                          <span
                            key={index}
                            className='px-2 py-1 mr-1 text-xs bg-gray-100 rounded'
                          >
                            {attr.name}: {attr.value}
                          </span>
                        ))}
                      </div>
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                      {product.quantity}
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                      {formatCurrency(product.lineTotal)}
                    </td>
                  </tr>
                ))}
              </tbody>
              <tfoot className='bg-gray-50'>
                <tr>
                  <td colSpan={3} className='px-6 py-4 text-right font-medium'>
                    Total:
                  </td>
                  <td className='px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900'>
                    {formatCurrency(orderDetail.total)}
                  </td>
                </tr>
              </tfoot>
            </table>
          </div>
        </div>
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
