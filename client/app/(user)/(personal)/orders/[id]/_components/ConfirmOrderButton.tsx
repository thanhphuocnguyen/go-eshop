'use client';

import { useState } from 'react';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { useRouter } from 'next/navigation';
import { CheckCircleIcon } from '@heroicons/react/24/solid';
import { GenericResponse } from '@/app/lib/definitions/apiResponse';

interface ConfirmOrderButtonProps {
  orderId: string;
}

export default function ConfirmOrderButton({
  orderId,
}: ConfirmOrderButtonProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const router = useRouter();

  const confirmReceipt = async () => {
    if (loading || success) return;

    setLoading(true);
    setError(null);

    const { data, error } = await clientSideFetch<GenericResponse<boolean>>(
      PUBLIC_API_PATHS.CONFIRM_RECEIVED_ORDER.replace(':id', orderId),
      {
        method: 'PUT',
      }
    );
    if (error) {
      setError(error.details || 'Failed to confirm order receipt');
      setLoading(false);
      return;
    }

    if (data) {
      setSuccess(true);
      // Refresh the page data
      router.refresh();
    }

    setLoading(false);
  };

  if (success) {
    return (
      <div className='flex items-center justify-center text-green-600 font-medium'>
        <CheckCircleIcon className='w-5 h-5 mr-2' />
        Order confirmed as received
      </div>
    );
  }

  return (
    <div className='flex flex-col items-center'>
      <button
        onClick={confirmReceipt}
        disabled={loading}
        className={`px-6 py-3 bg-indigo-600 text-white rounded-md font-medium shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-colors ${
          loading ? 'opacity-75 cursor-not-allowed' : ''
        }`}
      >
        {loading ? 'Processing...' : 'Confirm Order Receipt'}
      </button>

      {error && <p className='mt-2 text-sm text-red-600'>{error}</p>}

      <p className='mt-3 text-sm text-gray-500 max-w-md text-center'>
        Click this button when you have received your order to mark it as
        completed. This will finalize the order process.
      </p>
    </div>
  );
}
