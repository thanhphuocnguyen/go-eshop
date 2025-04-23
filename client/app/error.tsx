'use client';

import { motion } from 'framer-motion';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className='min-h-screen flex flex-col items-center justify-center p-5 bg-gradient-to-b from-gray-50 to-gray-100'>
      <div className='w-full max-w-md'>
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className='text-center mb-8'
        >
          <h1 className='text-6xl font-bold text-red-500 mb-2'>Oops!</h1>
          <div className='animate-bounce text-8xl mb-4'>⚠️</div>
          <h2 className='text-2xl font-semibold text-gray-800 mb-2'>
            Something went wrong
          </h2>
          <p className='text-gray-600 mb-4'>
            We&apos;re sorry for the inconvenience. Our team has been notified.
          </p>

          {error.digest && (
            <p className='text-sm text-gray-500 mb-4'>
              Error reference:{' '}
              <code className='bg-gray-100 px-1 py-0.5 rounded'>
                {error.digest}
              </code>
            </p>
          )}
        </motion.div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3, duration: 0.5 }}
          className='flex flex-col space-y-4'
        >
          <button
            onClick={() => reset()}
            className='transition transform hover:scale-105 bg-gradient-to-r from-blue-600 to-blue-700 text-white font-medium py-3 px-6 rounded-lg shadow-lg hover:shadow-xl flex items-center justify-center'
          >
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-5 w-5 mr-2'
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
          </button>

          <button
            onClick={() => (window.location.href = '/')}
            className='transition hover:bg-gray-100 text-gray-800 font-medium py-3 px-6 rounded-lg border border-gray-300 flex items-center justify-center'
          >
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-5 w-5 mr-2'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M3 12l2-2m0 0l7-7 7 7m-7-7v14'
              />
            </svg>
            Return to Home
          </button>
        </motion.div>
      </div>

      {process.env.NODE_ENV === 'development' && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6, duration: 0.5 }}
          className='mt-12 w-full max-w-3xl p-6 bg-gray-800 rounded-lg shadow-xl text-white overflow-auto'
        >
          <h3 className='text-xl font-semibold mb-4 flex items-center'>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              className='h-6 w-6 mr-2'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth={2}
                d='M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
              />
            </svg>
            Developer Error Details
          </h3>
          <div className='bg-gray-900 p-4 rounded-md overflow-auto'>
            <pre className='text-gray-300 text-sm whitespace-pre-wrap'>
              {JSON.stringify(error.cause) ||
                error.message ||
                JSON.stringify(error, null, 2)}
            </pre>
          </div>
        </motion.div>
      )}
    </div>
  );
}
