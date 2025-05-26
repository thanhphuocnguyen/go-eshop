import React from 'react';
import CategoryProductSkeleton from '@/app/components/Product/CategoryProductSkeleton';

export default function HomeLoading() {
  return (
    <div className="block relative mb-10">
      {/* Hero section skeleton */}
      <section className="new-arrival-ads bg-gray-300 animate-pulse">
        <div className="overlay"></div>
        <div className="relative w-[600] mx-auto flex flex-col justify-center items-center h-[500px]">
          <div className="h-12 bg-gray-400 rounded-md w-2/3 mb-6"></div>
          <div className="h-4 bg-gray-400 rounded-md w-3/4 mb-4"></div>
          <div className="h-4 bg-gray-400 rounded-md w-3/4 mb-4"></div>
          <div className="h-4 bg-gray-400 rounded-md w-1/2 mb-8"></div>
          <div className="h-12 bg-gray-400 rounded-lg w-44"></div>
        </div>
      </section>
      
      <div className="px-10">
        {/* Categories section skeleton */}
        <div className="pt-24">
          <div className="flex justify-between mb-6">
            <div className="h-8 bg-gray-300 rounded-md w-48"></div>
            <div className="h-6 bg-gray-300 rounded-md w-40"></div>
          </div>
          <CategoryProductSkeleton />
        </div>
        
        {/* First feature section skeleton */}
        <section className="relative pt-24 h-[700px]">
          <div className="relative h-full bg-gray-300 animate-pulse rounded-lg">
            <div className="absolute inset-0 flex flex-col justify-center items-center">
              <div className="w-1/2 m-auto flex flex-col items-center">
                <div className="h-12 bg-gray-400 rounded-md w-3/4 mb-6"></div>
                <div className="h-4 bg-gray-400 rounded-md w-full mb-3"></div>
                <div className="h-4 bg-gray-400 rounded-md w-full mb-3"></div>
                <div className="h-4 bg-gray-400 rounded-md w-4/5 mb-8"></div>
                <div className="h-12 bg-gray-400 rounded-lg w-44"></div>
              </div>
            </div>
          </div>
        </section>
        
        {/* Collections section skeleton */}
        <div className="pt-20">
          <div className="flex justify-between mb-6">
            <div className="h-8 bg-gray-300 rounded-md w-48"></div>
            <div className="h-6 bg-gray-300 rounded-md w-40"></div>
          </div>
          <CategoryProductSkeleton />
        </div>
        
        {/* Second feature section skeleton */}
        <section className="relative my-24 mb-40 h-[600px]">
          <div className="h-full bg-gray-300 animate-pulse rounded-lg">
            <div className="absolute inset-0 flex flex-col justify-center items-center">
              <div className="w-1/2 m-auto flex flex-col items-center">
                <div className="h-10 bg-gray-400 rounded-md w-1/2 mb-6"></div>
                <div className="h-4 bg-gray-400 rounded-md w-4/5 mb-3"></div>
                <div className="h-4 bg-gray-400 rounded-md w-4/5 mb-3"></div>
                <div className="h-4 bg-gray-400 rounded-md w-3/5 mb-8"></div>
                <div className="h-10 bg-gray-400 rounded-lg w-44"></div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}