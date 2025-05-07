'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import LoadingButton from '@/components/Common/LoadingButton';
import LoadingInline from '@/components/Common/Loadings/LoadingInline';
import { Breadcrumb } from '@/components/Common';

// Define TypeScript interfaces for our data models
interface Order {
  id: string;
  total: number;
  total_items: number;
  status: string;
  payment_status: string;
  customer_name: string;
  customer_email: string;
  created_at: string;
  updated_at: string;
}

interface OrderDetail {
  id: string;
  total: number;
  status: string;
  customer_name: string;
  customer_email: string;
  payment_info: {
    id: string;
    refund_id?: string;
    amount: number;
    intent_id?: string;
    client_secret?: string;
    gateway?: string;
    method: string;
    status: string;
  } | null;
  shipping_info: {
    street: string;
    ward: string;
    district: string;
    city: string;
    phone: string;
  };
  products: Array<{
    id: string;
    name: string;
    image_url?: string;
    attributes_snapshot: Array<{
      name: string;
      value: string;
    }>;
    line_total: number;
    quantity: number;
  }>;
  created_at: string;
}

interface ApiResponse<T> {
  success: boolean;
  message?: string;
  data?: T;
  error?: {
    code: string;
    details: string;
    stack?: string;
  };
  pagination?: {
    total: number;
    page: number;
    pageSize: number;
    totalPages: number;
    hasNextPage: boolean;
    hasPreviousPage: boolean;
  };
  meta?: {
    timestamp: string;
    requestId: string;
    path: string;
    method: string;
  };
}

export default function AdminOrdersPage() {
  const router = useRouter();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [selectedOrder, setSelectedOrder] = useState<OrderDetail | null>(null);
  const [viewingOrderDetail, setViewingOrderDetail] = useState(false);
  const [changingStatus, setChangingStatus] = useState(false);
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
      const response = await fetch(`/api/v1/admin/orders?page=${currentPage}&page_size=10${statusFilter !== 'all' ? `&status=${statusFilter}` : ''}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          // Add your authentication headers here
        },
      });
      
      if (!response.ok) {
        throw new Error(`Error fetching orders: ${response.status}`);
      }
      
      const data: ApiResponse<Order[]> = await response.json();
      
      if (data.success && data.data) {
        setOrders(data.data);
        if (data.pagination) {
          setTotalPages(data.pagination.totalPages);
        }
      } else {
        setError(data.error?.details || 'Failed to fetch orders');
      }
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const fetchOrderDetail = async (orderId: string) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/v1/admin/orders/${orderId}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          // Add your authentication headers here
        },
      });
      
      if (!response.ok) {
        throw new Error(`Error fetching order details: ${response.status}`);
      }
      
      const data: ApiResponse<OrderDetail> = await response.json();
      
      if (data.success && data.data) {
        setSelectedOrder(data.data);
        setViewingOrderDetail(true);
      } else {
        setError(data.error?.details || 'Failed to fetch order details');
      }
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const changeOrderStatus = async (orderId: string, newStatus: string) => {
    setChangingStatus(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/v1/admin/orders/${orderId}/status`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          // Add your authentication headers here
        },
        body: JSON.stringify({ status: newStatus }),
      });
      
      if (!response.ok) {
        throw new Error(`Error changing order status: ${response.status}`);
      }
      
      const data: ApiResponse<Order> = await response.json();
      
      if (data.success && data.data) {
        // Update the selected order with the new status
        setSelectedOrder(prev => prev ? { ...prev, status: newStatus } : null);
        
        // Also update the order in the list
        setOrders(prevOrders => 
          prevOrders.map(order => 
            order.id === orderId ? { ...order, status: newStatus } : order
          )
        );
      } else {
        setError(data.error?.details || 'Failed to change order status');
      }
    } catch (err) {
      setError((err as Error).message);
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
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount);
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

  // Render the orders list
  const renderOrdersList = () => {
    if (loading && orders.length === 0) {
      return (
        <div className="flex justify-center py-10">
          <LoadingInline />
        </div>
      );
    }

    if (error) {
      return (
        <div className="bg-red-100 p-4 rounded-md text-red-700 mb-4">
          <p>{error}</p>
          <button 
            className="mt-2 text-red-700 font-medium hover:text-red-800"
            onClick={fetchOrders}
          >
            Try Again
          </button>
        </div>
      );
    }

    if (orders.length === 0) {
      return (
        <div className="p-4 text-center text-gray-500">
          No orders found.
        </div>
      );
    }

    return (
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Order ID
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Customer
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Total
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Items
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Payment Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Date
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {orders.map((order) => (
              <tr key={order.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600">
                  {order.id.substring(0, 8)}...
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <div>{order.customer_name}</div>
                  <div className="text-xs text-gray-500">{order.customer_email}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  {formatCurrency(order.total)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-center">
                  {order.total_items}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(order.status)}`}>
                    {order.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(order.payment_status)}`}>
                    {order.payment_status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDate(order.created_at)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-right">
                  <button
                    onClick={() => fetchOrderDetail(order.id)}
                    className="text-indigo-600 hover:text-indigo-900 font-medium"
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

  // Render the order detail view
  const renderOrderDetail = () => {
    if (!selectedOrder) {
      return null;
    }

    return (
      <div className="bg-white p-6 rounded-lg shadow-lg">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-semibold text-gray-800">Order Details</h2>
          <button
            className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300"
            onClick={() => {
              setViewingOrderDetail(false);
              setSelectedOrder(null);
            }}
          >
            Back to Orders
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div>
            <h3 className="font-semibold text-gray-700 mb-2">Order Information</h3>
            <div className="bg-gray-50 p-4 rounded-md">
              <p><span className="font-medium">Order ID:</span> {selectedOrder.id}</p>
              <p><span className="font-medium">Date:</span> {formatDate(selectedOrder.created_at)}</p>
              <p><span className="font-medium">Total:</span> {formatCurrency(selectedOrder.total)}</p>
              <p className="flex items-center mt-1">
                <span className="font-medium mr-2">Status:</span>
                <span className={`px-2 py-1 text-xs leading-5 font-semibold rounded-full ${getStatusColor(selectedOrder.status)}`}>
                  {selectedOrder.status}
                </span>
              </p>
            </div>
          </div>

          <div>
            <h3 className="font-semibold text-gray-700 mb-2">Customer Information</h3>
            <div className="bg-gray-50 p-4 rounded-md">
              <p><span className="font-medium">Name:</span> {selectedOrder.customer_name}</p>
              <p><span className="font-medium">Email:</span> {selectedOrder.customer_email}</p>
              <p><span className="font-medium">Shipping Address:</span></p>
              <p className="ml-4">
                {selectedOrder.shipping_info.street}, {selectedOrder.shipping_info.ward},
                <br />
                {selectedOrder.shipping_info.district}, {selectedOrder.shipping_info.city}
              </p>
              <p><span className="font-medium">Phone:</span> {selectedOrder.shipping_info.phone}</p>
            </div>
          </div>
        </div>

        {selectedOrder.payment_info && (
          <div className="mb-6">
            <h3 className="font-semibold text-gray-700 mb-2">Payment Information</h3>
            <div className="bg-gray-50 p-4 rounded-md">
              <p><span className="font-medium">Payment ID:</span> {selectedOrder.payment_info.id}</p>
              <p><span className="font-medium">Amount:</span> {formatCurrency(selectedOrder.payment_info.amount)}</p>
              <p><span className="font-medium">Method:</span> {selectedOrder.payment_info.method}</p>
              <p className="flex items-center mt-1">
                <span className="font-medium mr-2">Status:</span>
                <span className={`px-2 py-1 text-xs leading-5 font-semibold rounded-full ${getStatusColor(selectedOrder.payment_info.status)}`}>
                  {selectedOrder.payment_info.status}
                </span>
              </p>
              {selectedOrder.payment_info.gateway && (
                <p><span className="font-medium">Gateway:</span> {selectedOrder.payment_info.gateway}</p>
              )}
              {selectedOrder.payment_info.refund_id && (
                <p><span className="font-medium">Refund ID:</span> {selectedOrder.payment_info.refund_id}</p>
              )}
            </div>
          </div>
        )}

        {/* Status change section */}
        {getAvailableStatusOptions(selectedOrder.status).length > 0 && (
          <div className="mb-6">
            <h3 className="font-semibold text-gray-700 mb-2">Change Order Status</h3>
            <div className="flex flex-wrap gap-2">
              {getAvailableStatusOptions(selectedOrder.status).map(status => (
                <LoadingButton
                  key={status}
                  className={`px-4 py-2 rounded-md text-sm font-medium ${status === 'cancelled' || status === 'refunded' ? 'bg-red-600 hover:bg-red-700 text-white' : 'bg-blue-600 hover:bg-blue-700 text-white'}`}
                  onClick={() => changeOrderStatus(selectedOrder.id, status)}
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
          <h3 className="font-semibold text-gray-700 mb-4">Order Items</h3>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Product
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Attributes
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Quantity
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Price
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {selectedOrder.products.map((product) => (
                  <tr key={product.id}>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        {product.image_url && (
                          <div className="flex-shrink-0 h-10 w-10 mr-4">
                            <img
                              className="h-10 w-10 rounded-md object-cover"
                              src={product.image_url}
                              alt={product.name}
                            />
                          </div>
                        )}
                        <div>
                          <div className="text-sm font-medium text-gray-900">
                            {product.name}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="text-sm text-gray-900">
                        {product.attributes_snapshot.map((attr, index) => (
                          <span key={index} className="px-2 py-1 mr-1 text-xs bg-gray-100 rounded">
                            {attr.name}: {attr.value}
                          </span>
                        ))}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {product.quantity}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {formatCurrency(product.line_total)}
                    </td>
                  </tr>
                ))}
              </tbody>
              <tfoot className="bg-gray-50">
                <tr>
                  <td colSpan={3} className="px-6 py-4 text-right font-medium">
                    Total:
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {formatCurrency(selectedOrder.total)}
                  </td>
                </tr>
              </tfoot>
            </table>
          </div>
        </div>
      </div>
    );
  };

  // Render pagination controls
  const renderPagination = () => {
    return (
      <div className="flex justify-between items-center px-6 py-4 bg-white border-t border-gray-200">
        <div className="text-sm text-gray-500">
          Showing page {currentPage} of {totalPages}
        </div>
        <div className="flex space-x-2">
          <button
            onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
            disabled={currentPage === 1}
            className={`px-3 py-1 rounded ${currentPage === 1 ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'}`}
          >
            Previous
          </button>
          <button
            onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
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
      <div className="mb-4 p-4 bg-white rounded-lg shadow-sm">
        <div className="flex flex-wrap items-center gap-4">
          <div>
            <label htmlFor="statusFilter" className="block text-sm font-medium text-gray-700 mb-1">
              Status Filter
            </label>
            <select
              id="statusFilter"
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value);
                setCurrentPage(1); // Reset to first page on filter change
              }}
              className="border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="all">All Statuses</option>
              <option value="pending">Pending</option>
              <option value="confirmed">Confirmed</option>
              <option value="delivering">Delivering</option>
              <option value="delivered">Delivered</option>
              <option value="completed">Completed</option>
              <option value="cancelled">Cancelled</option>
              <option value="refunded">Refunded</option>
            </select>
          </div>
          
          <button
            onClick={() => fetchOrders()}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm font-medium"
          >
            Refresh
          </button>
        </div>
      </div>
    );
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <Breadcrumb items={[
          { label: 'Admin', href: '/admin' },
          { label: 'Orders', href: '/admin/orders' },
        ]} />
        <h1 className="text-3xl font-bold text-gray-900 mt-2">Orders Management</h1>
      </div>
      
      {viewingOrderDetail ? (
        renderOrderDetail()
      ) : (
        <>
          {renderFilterControls()}
          <div className="bg-white rounded-lg shadow-sm overflow-hidden">
            {renderOrdersList()}
            {orders.length > 0 && renderPagination()}
          </div>
        </>
      )}

      {error && (
        <div className="fixed bottom-4 right-4 bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded shadow-md">
          <p className="font-bold">Error</p>
          <p>{error}</p>
          <button
            onClick={() => setError(null)}
            className="absolute top-2 right-2 text-red-500 hover:text-red-700"
          >
            &times;
          </button>
        </div>
      )}
    </div>
  );
}
