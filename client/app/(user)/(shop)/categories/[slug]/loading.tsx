import Link from 'next/link';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import CategoryProductSkeleton from '@/app/components/Product/CategoryProductSkeleton';

export default function CategoryLoading() {
  return (
    <div className='container mx-auto px-4 py-8'>
      <Link
        href='/categories'
        className='inline-flex items-center text-indigo-600 hover:underline mb-6'
      >
        <ArrowLeftIcon className='h-4 w-4 mr-2' />
        Back to all categories
      </Link>

      <div className='flex items-center mb-8'>
        {/* Category image skeleton */}
        <div className='h-24 w-24 mr-6 relative overflow-hidden rounded-lg shadow-md bg-gray-200 animate-pulse'></div>
        <div>
          {/* Category name skeleton */}
          <div className='h-8 w-64 bg-gray-200 rounded-md animate-pulse mb-2'></div>
          {/* Category description skeleton */}
          <div className='h-4 w-96 bg-gray-200 rounded-md animate-pulse'></div>
        </div>
      </div>

      <div className='flex flex-col lg:flex-row gap-8'>
        {/* Filters skeleton */}
        <div className='lg:w-64 w-full'>
          <div className='bg-white rounded-lg shadow-sm p-5 border border-gray-100 sticky top-24'>
            <div className='h-6 w-32 bg-gray-200 rounded-md animate-pulse mb-4'></div>
            <div className='space-y-4'>
              <div className='h-4 w-full bg-gray-200 rounded-md animate-pulse'></div>
              <div className='h-4 w-full bg-gray-200 rounded-md animate-pulse'></div>
              <div className='h-4 w-3/4 bg-gray-200 rounded-md animate-pulse'></div>
            </div>
            
            <div className='my-6 border-t border-gray-100'></div>
            
            <div className='h-6 w-32 bg-gray-200 rounded-md animate-pulse mb-4'></div>
            <div className='space-y-4'>
              <div className='h-4 w-full bg-gray-200 rounded-md animate-pulse'></div>
              <div className='h-4 w-full bg-gray-200 rounded-md animate-pulse'></div>
            </div>
          </div>
        </div>
        
        {/* Products skeleton */}
        <div className='flex-1'>
          <CategoryProductSkeleton count={9} />
        </div>
      </div>
    </div>
  );
}
