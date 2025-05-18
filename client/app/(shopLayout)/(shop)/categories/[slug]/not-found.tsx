import Link from 'next/link';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';

export default function CategoryNotFound() {
  return (
    <div className='container mx-auto px-4 py-8'>
      <Link
        href='/categories'
        className='inline-flex items-center text-indigo-600 hover:underline mb-6'
      >
        <ArrowLeftIcon className='h-4 w-4 mr-2' />
        Back to all categories
      </Link>
      
      <div className='text-center py-12 bg-white rounded-lg shadow'>
        <div className='h-24 w-24 mx-auto mb-4 text-gray-300'>
          <svg 
            xmlns="http://www.w3.org/2000/svg" 
            fill="none" 
            viewBox="0 0 24 24" 
            stroke="currentColor"
          >
            <path 
              strokeLinecap="round" 
              strokeLinejoin="round" 
              strokeWidth={1.5} 
              d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" 
            />
          </svg>
        </div>
        <h1 className='text-2xl font-bold text-gray-800 mb-2'>
          Category Not Found
        </h1>
        <p className='text-gray-600 mb-6 max-w-md mx-auto'>
          The category you are looking for does not exist or has been removed.
        </p>
        <div className='flex justify-center space-x-4'>
          <Link
            href='/categories'
            className='bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 transition-colors'
          >
            Browse All Categories
          </Link>
          <Link
            href='/'
            className='bg-gray-100 text-gray-700 px-4 py-2 rounded-md hover:bg-gray-200 transition-colors'
          >
            Back to Home
          </Link>
        </div>
      </div>
    </div>
  );
}
