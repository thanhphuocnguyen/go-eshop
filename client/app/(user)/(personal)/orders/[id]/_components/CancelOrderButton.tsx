'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Dialog } from '@headlessui/react';
import { ExclamationTriangleIcon } from '@heroicons/react/24/outline';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';

export default function CancelOrderButton({ orderId }: { orderId: string }) {
  const [isOpen, setIsOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleCancelOrder = async () => {
    setIsLoading(true);
    try {
      const response = await clientSideFetch(
        PUBLIC_API_PATHS.CANCEL_ORDER.replace(':id', orderId),
        {
          method: 'POST',
        }
      );

      if (response.error) {
        throw new Error(response.error.details || 'Failed to cancel order');
      }

      toast.success('Order cancelled successfully');
      // Close the dialog
      setIsOpen(false);
      // Refresh the page to show updated order status
      router.refresh();
    } catch (error: any) {
      toast.error(error.message || 'An error occurred while cancelling the order');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <>
      <button
        onClick={() => setIsOpen(true)}
        className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
      >
        Cancel Order
      </button>

      {/* Confirmation Dialog */}
      <Dialog
        open={isOpen}
        onClose={() => !isLoading && setIsOpen(false)}
        className="relative z-50"
      >
        {/* Backdrop */}
        <div className="fixed inset-0 bg-black/30" aria-hidden="true" />

        {/* Dialog position */}
        <div className="fixed inset-0 flex items-center justify-center p-4">
          <Dialog.Panel className="w-full max-w-md rounded-lg bg-white p-6 shadow-xl">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <ExclamationTriangleIcon className="h-6 w-6 text-red-600" aria-hidden="true" />
              </div>
              <div className="ml-3">
                <Dialog.Title className="text-lg font-medium text-gray-900">
                  Cancel Order
                </Dialog.Title>
              </div>
            </div>

            <div className="mt-4">
              <p className="text-sm text-gray-500">
                Are you sure you want to cancel this order? This action cannot be undone,
                and your order will be permanently cancelled.
              </p>
            </div>

            <div className="mt-6 flex justify-end space-x-3">
              <button
                type="button"
                className="inline-flex justify-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                onClick={() => setIsOpen(false)}
                disabled={isLoading}
              >
                No, keep order
              </button>
              <button
                type="button"
                className="inline-flex justify-center px-4 py-2 text-sm font-medium text-white bg-red-600 border border-transparent rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                onClick={handleCancelOrder}
                disabled={isLoading}
              >
                {isLoading ? 'Cancelling...' : 'Yes, cancel order'}
              </button>
            </div>
          </Dialog.Panel>
        </div>
      </Dialog>
    </>
  );
}
