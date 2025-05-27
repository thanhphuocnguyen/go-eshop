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
  const [selectedReason, setSelectedReason] = useState('');
  const [customReason, setCustomReason] = useState('');
  const router = useRouter();

  const predefinedReasons = [
    'Changed my mind',
    'Found a better price elsewhere',
    'Ordered by mistake',
    'No longer needed',
    'Shipping takes too long',
    'Payment issues',
    'Product not as described',
    'Other'
  ];

  const handleCancelOrder = async () => {
    // Validate reason selection
    if (!selectedReason) {
      toast.error('Please select a reason for cancellation');
      return;
    }

    if (selectedReason === 'Other' && !customReason.trim()) {
      toast.error('Please provide a custom reason for cancellation');
      return;
    }

    setIsLoading(true);
    try {
      const reason = selectedReason === 'Other' ? customReason : selectedReason;
      
      const response = await clientSideFetch(
        PUBLIC_API_PATHS.CANCEL_ORDER.replace(':id', orderId),
        {
          method: 'POST',
          body: {
            reason: reason
          }
        }
      );

      if (response.error) {
        throw new Error(response.error.details || 'Failed to cancel order');
      }

      toast.success('Order cancelled successfully');
      // Close the dialog
      setIsOpen(false);
      // Reset form state
      setSelectedReason('');
      setCustomReason('');
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
          <Dialog.Panel className="w-full max-w-lg rounded-lg bg-white p-6 shadow-xl">
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
              <p className="text-sm text-gray-500 mb-4">
                Are you sure you want to cancel this order? This action cannot be undone,
                and your order will be permanently cancelled.
              </p>

              {/* Reason Selection */}
              <div className="space-y-4">
                <div>
                  <label htmlFor="reason" className="block text-sm font-medium text-gray-700 mb-2">
                    Please select a reason for cancellation:
                  </label>
                  <select
                    id="reason"
                    value={selectedReason}
                    onChange={(e) => setSelectedReason(e.target.value)}
                    className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm"
                    disabled={isLoading}
                  >
                    <option value="">Select a reason...</option>
                    {predefinedReasons.map((reason) => (
                      <option key={reason} value={reason}>
                        {reason}
                      </option>
                    ))}
                  </select>
                </div>

                {/* Custom reason input - only show when "Other" is selected */}
                {selectedReason === 'Other' && (
                  <div>
                    <label htmlFor="customReason" className="block text-sm font-medium text-gray-700 mb-2">
                      Please specify your reason:
                    </label>
                    <textarea
                      id="customReason"
                      value={customReason}
                      onChange={(e) => setCustomReason(e.target.value)}
                      rows={3}
                      className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500 sm:text-sm"
                      placeholder="Please provide your reason for cancellation..."
                      disabled={isLoading}
                    />
                  </div>
                )}
              </div>
            </div>

            <div className="mt-6 flex justify-end space-x-3">
              <button
                type="button"
                className="inline-flex justify-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                onClick={() => {
                  setIsOpen(false);
                  setSelectedReason('');
                  setCustomReason('');
                }}
                disabled={isLoading}
              >
                No, keep order
              </button>
              <button
                type="button"
                className={`inline-flex justify-center px-4 py-2 text-sm font-medium text-white border border-transparent rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 ${
                  !selectedReason || (selectedReason === 'Other' && !customReason.trim())
                    ? 'bg-gray-400 cursor-not-allowed'
                    : 'bg-red-600 hover:bg-red-700'
                }`}
                onClick={handleCancelOrder}
                disabled={isLoading || !selectedReason || (selectedReason === 'Other' && !customReason.trim())}
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
