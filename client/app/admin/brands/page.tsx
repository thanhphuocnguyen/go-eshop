'use client';

import dayjs from 'dayjs';
import { Button } from '@headlessui/react';
import Link from 'next/link';
import { ReactNode, useState } from 'react';
import { useBrands } from '../hooks';
import Image from 'next/image';
import {
  PencilIcon,
  TrashIcon,
  EyeIcon,
  ArrowPathIcon,
} from '@heroicons/react/24/outline';
import { useRouter } from 'next/navigation';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';

export default function Page() {
  const router = useRouter();
  const { brands, isLoading, mutate } = useBrands();
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [sortOrder, setSortOrder] = useState('date-desc');

  // Function to delete a brand
  const handleDeleteBrand = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this brand?')) {
      setDeletingId(id);
      try {
        const response = await clientSideFetch<boolean>(
          ADMIN_API_PATHS.BRAND.replace(':id', id),
          {
            method: 'DELETE',
          }
        );

        if (response.data) {
          toast.success('Brand deleted successfully');
          mutate(); // Refresh brands data
        } else {
          toast.error('Failed to delete brand');
        }
      } catch (error) {
        toast.error('Error deleting brand');
        console.error(error);
      } finally {
        setDeletingId(null);
      }
    }
  };

  // Filter brands based on search and status
  const filteredBrands =
    brands?.filter((brand) => {
      const matchesSearch =
        searchTerm === '' ||
        brand.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        (brand.description &&
          brand.description.toLowerCase().includes(searchTerm.toLowerCase()));

      const matchesStatus =
        statusFilter === 'all' ||
        (statusFilter === 'published' && brand.published) ||
        (statusFilter === 'draft' && !brand.published);

      return matchesSearch && matchesStatus;
    }) || [];

  // Sort filtered brands
  const sortedBrands = [...filteredBrands].sort((a, b) => {
    switch (sortOrder) {
      case 'name-asc':
        return a.name.localeCompare(b.name);
      case 'name-desc':
        return b.name.localeCompare(a.name);
      case 'date-asc':
        return (
          new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
        );
      case 'date-desc':
      default:
        return (
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
        );
    }
  });

  if (isLoading)
    return (
      <div className='flex justify-center items-center h-full'>
        <div className='animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary'></div>
      </div>
    );

  return (
    <div className='h-full overflow-auto my-auto p-6'>
      <div className='flex justify-between pt-4 pb-8'>
        <div>
          <h2 className='text-2xl font-bold'>Brands</h2>
          <p className='text-gray-500 text-sm mt-1'>
            Manage your product brands
          </p>
        </div>
        <div className='flex items-center space-x-4'>
          <Button
            onClick={() => mutate()}
            className='inline-flex items-center px-3 py-2 text-sm font-medium rounded-md text-gray-700 hover:bg-gray-100'
          >
            <ArrowPathIcon className='h-4 w-4 mr-2' />
            Refresh
          </Button>
          <Button
            as={Link}
            href={'/admin/brands/new'}
            className='btn btn-lg btn-primary'
          >
            Add new
          </Button>
        </div>
      </div>

      {/* Search and Filter Section */}
      <div className='mb-6 flex justify-between'>
        <div className='relative flex items-center max-w-xs'>
          <input
            type='text'
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder='Search brands...'
            className='w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent'
          />
          <div className='absolute left-3'>
            <svg
              className='h-4 w-4 text-gray-400'
              fill='none'
              viewBox='0 0 24 24'
              stroke='currentColor'
            >
              <path
                strokeLinecap='round'
                strokeLinejoin='round'
                strokeWidth='2'
                d='M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z'
              />
            </svg>
          </div>
        </div>
        <div className='flex space-x-2'>
          <select
            className='border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500'
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
          >
            <option value='all'>All Status</option>
            <option value='published'>Published</option>
            <option value='draft'>Draft</option>
          </select>
          <select
            className='border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500'
            value={sortOrder}
            onChange={(e) => setSortOrder(e.target.value)}
          >
            <option value='date-desc'>Newest first</option>
            <option value='date-asc'>Oldest first</option>
            <option value='name-asc'>Name (A-Z)</option>
            <option value='name-desc'>Name (Z-A)</option>
          </select>
        </div>
      </div>

      {/* Empty state */}
      {sortedBrands.length === 0 ? (
        <div className='text-center py-12 border-2 border-dashed border-gray-300 rounded-lg'>
          <svg
            className='mx-auto h-12 w-12 text-gray-400'
            fill='none'
            viewBox='0 0 24 24'
            stroke='currentColor'
          >
            <path
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth='2'
              d='M20 7l-8-4-8 4m16 0l-8 4m-8-4l8 4m8 0l-8 4m8-4v10a2 2 0 01-2 2H6a2 2 0 01-2-2V7'
            />
          </svg>
          <h3 className='mt-2 text-sm font-medium text-gray-900'>
            No brands found
          </h3>
          <p className='mt-1 text-sm text-gray-500'>
            {searchTerm || statusFilter !== 'all'
              ? 'Try changing your search or filter criteria'
              : 'Get started by creating a new brand'}
          </p>
          <div className='mt-6'>
            {searchTerm || statusFilter !== 'all' ? (
              <Button
                onClick={() => {
                  setSearchTerm('');
                  setStatusFilter('all');
                }}
                className='btn btn-primary'
              >
                Clear filters
              </Button>
            ) : (
              <Button
                as={Link}
                href='/admin/brands/new'
                className='btn btn-primary'
              >
                Create brand
              </Button>
            )}
          </div>
        </div>
      ) : (
        <div className='bg-white overflow-hidden shadow rounded-lg'>
          <TableContainer
            header={[
              '',
              'Name',
              'Description',
              'Status',
              'Products',
              'Created At',
              'Actions',
            ]}
          >
            {sortedBrands.map((brand) => (
              <TableRow key={brand.id}>
                <TableCell className='w-16'>
                  {brand.imageUrl ? (
                    <div className='relative h-10 w-10 rounded-md overflow-hidden'>
                      <Image
                        src={brand.imageUrl}
                        alt={brand.name}
                        fill
                        className='object-cover'
                      />
                    </div>
                  ) : (
                    <div className='h-10 w-10 rounded-md bg-gray-200 flex items-center justify-center'>
                      <span className='text-gray-400 text-xs'>
                        {brand.name.charAt(0).toUpperCase()}
                      </span>
                    </div>
                  )}
                </TableCell>
                <TableCellHead>
                  <Link
                    className='text-blue-500 hover:text-blue-700 font-medium'
                    href={'/admin/brands/' + brand.id}
                  >
                    {brand.name}
                  </Link>
                </TableCellHead>
                <TableCell className='max-w-xs'>
                  <div className='truncate w-48'>
                    {brand.description || 'No description'}
                  </div>
                </TableCell>
                <TableCell>
                  <span
                    className={`inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium ${brand.published ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}`}
                  >
                    {brand.published ? 'Published' : 'Draft'}
                  </span>
                </TableCell>
                <TableCell className='text-center'>{0}</TableCell>
                <TableCell>
                  {dayjs(brand.createdAt).format('YYYY/MM/DD')}
                </TableCell>
                <TableCell>
                  <div className='flex space-x-2 justify-end'>
                    <Button
                      as={Link}
                      href={`/admin/brands/${brand.id}`}
                      className='inline-flex items-center p-2 text-blue-600 hover:bg-blue-50 rounded-full'
                      title='Edit'
                    >
                      <PencilIcon className='h-4 w-4' />
                    </Button>
                    <Button
                      onClick={() => router.push(`/admin/brands/${brand.id}`)}
                      className='inline-flex items-center p-2 text-green-600 hover:bg-green-50 rounded-full'
                      title='View details'
                    >
                      <EyeIcon className='h-4 w-4' />
                    </Button>
                    <Button
                      onClick={() => handleDeleteBrand(brand.id)}
                      className='inline-flex items-center p-2 text-red-600 hover:bg-red-50 rounded-full'
                      disabled={deletingId === brand.id}
                      title='Delete'
                    >
                      {deletingId === brand.id ? (
                        <div className='h-4 w-4 border-2 border-t-red-500 border-r-transparent border-b-transparent border-l-transparent rounded-full animate-spin' />
                      ) : (
                        <TrashIcon className='h-4 w-4' />
                      )}
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableContainer>

          {/* Pagination */}
          <div className='bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6'>
            <div className='hidden sm:flex-1 sm:flex sm:items-center sm:justify-between'>
              <div>
                <p className='text-sm text-gray-700'>
                  Showing{' '}
                  <span className='font-medium'>{sortedBrands.length}</span>{' '}
                  brands
                </p>
              </div>
              <div>
                <nav
                  className='relative z-0 inline-flex rounded-md shadow-sm -space-x-px'
                  aria-label='Pagination'
                >
                  <Button
                    className='relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50'
                    disabled={true}
                  >
                    <span className='sr-only'>Previous</span>
                    <svg
                      className='h-5 w-5'
                      xmlns='http://www.w3.org/2000/svg'
                      viewBox='0 0 20 20'
                      fill='currentColor'
                      aria-hidden='true'
                    >
                      <path
                        fillRule='evenodd'
                        d='M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z'
                        clipRule='evenodd'
                      />
                    </svg>
                  </Button>
                  <Button className='relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50'>
                    1
                  </Button>
                  <Button className='relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50'>
                    <span className='sr-only'>Next</span>
                    <svg
                      className='h-5 w-5'
                      xmlns='http://www.w3.org/2000/svg'
                      viewBox='0 0 20 20'
                      fill='currentColor'
                      aria-hidden='true'
                    >
                      <path
                        fillRule='evenodd'
                        d='M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z'
                        clipRule='evenodd'
                      />
                    </svg>
                  </Button>
                </nav>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// Table components
interface TableContainerProps {
  header: string[];
  children: ReactNode;
}

function TableContainer({ header, children }: TableContainerProps) {
  return (
    <div className='overflow-x-auto'>
      <table className='min-w-full divide-y divide-gray-200'>
        <thead className='bg-gray-50'>
          <tr>
            {header.map((item, index) => (
              <th
                key={index}
                scope='col'
                className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
              >
                {item}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className='bg-white divide-y divide-gray-200'>{children}</tbody>
      </table>
    </div>
  );
}

function TableRow({ children }: { children: ReactNode }) {
  return <tr>{children}</tr>;
}

function TableCellHead({ children }: { children: ReactNode }) {
  return (
    <td className='px-6 py-4 whitespace-nowrap font-medium text-sm'>
      {children}
    </td>
  );
}

function TableCell({
  children,
  className = '',
}: {
  children?: ReactNode;
  className?: string;
}) {
  return (
    <td
      className={`px-6 py-4 whitespace-nowrap text-sm text-gray-500 ${className}`}
    >
      {children}
    </td>
  );
}
