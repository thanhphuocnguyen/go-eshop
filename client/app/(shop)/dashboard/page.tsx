import React from 'react';
import { Metadata } from 'next';
import RouteGuard from '@/components/RouteGuard';
import Link from 'next/link';

export const metadata: Metadata = {
  title: 'User Dashboard - eShop',
  description: 'View and manage your eShop account',
};

export default function DashboardPage() {
  return (
    <RouteGuard requireAuth>
      <div className='container mx-auto py-8'>
        <h1 className='text-3xl font-bold mb-6'>Your Dashboard</h1>

        <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'>
          <div className='bg-white p-6 rounded-lg shadow-md'>
            <h2 className='text-xl font-semibold mb-4'>Personal Information</h2>
            <p className='text-gray-600 mb-4'>
              View and edit your personal information
            </p>
            <a
              href='/personal-info'
              className='text-blue-500 hover:text-blue-700'
            >
              Manage &rarr;
            </a>
          </div>

          <div className='bg-white p-6 rounded-lg shadow-md'>
            <h2 className='text-xl font-semibold mb-4'>Your Orders</h2>
            <p className='text-gray-600 mb-4'>Track and manage your orders</p>
            <Link href='/orders' className='text-blue-500 hover:text-blue-700'>
              View Orders &rarr;
            </Link>
          </div>

          <div className='bg-white p-6 rounded-lg shadow-md'>
            <h2 className='text-xl font-semibold mb-4'>Wishlist</h2>
            <p className='text-gray-600 mb-4'>View and manage your wishlist</p>
            <a href='/wishlist' className='text-blue-500 hover:text-blue-700'>
              View Wishlist &rarr;
            </a>
          </div>
        </div>
      </div>
    </RouteGuard>
  );
}
