'use client';
import React from 'react';
import { ExclamationCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { ErrorResponse } from '@/app/lib/definitions';

interface ErrorComponentProps {
  error?: ErrorResponse;
}

export const ErrorComponent: React.FC<ErrorComponentProps> = ({ error }) => {
  return (
    <div className='h-full flex flex-col items-center justify-center py-12 px-4 sm:px-6 lg:px-8'>
      <div className='w-full max-w-md bg-white shadow-lg rounded-lg p-6 border border-red-100'>
        <div className='flex items-center justify-center text-red-500 mb-5'>
          <ExclamationCircleIcon className='h-16 w-16' aria-hidden='true' />
        </div>

        <h2 className='mt-2 text-center text-2xl font-bold text-gray-900'>
          Product Not Found
        </h2>

        <p className='mt-3 text-center text-gray-600'>
          {error?.details ||
            error?.stack ||
            "We couldn't load this product. The product may no longer be available or there might be a temporary issue."}
        </p>

        <div className='mt-6 flex justify-center'>
          <Link
            href='/products'
            className='inline-flex justify-center items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
          >
            Back to Products
          </Link>
        </div>
      </div>
    </div>
  );
};
