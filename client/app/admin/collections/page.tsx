'use client';
import { API_PATHS } from '@/lib/constants/api';
import Link from 'next/link';
import { Button } from '@headlessui/react';
import dayjs from 'dayjs';
import useSWR from 'swr';
import { toast } from 'react-toastify';
import Image from 'next/image';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import { useState } from 'react';
import { ConfirmDialog } from '@/components/Common/Dialogs/ConfirmDialog';
import { apiFetch } from '@/lib/api/api';

export default function Page() {
  const [selectedCollection, setSelectedCollection] =
    useState<GeneralCategoryModel | null>(null);

  const {
    data: collections,
    isLoading,
    mutate,
  } = useSWR(
    API_PATHS.COLLECTIONS,
    (url) =>
      apiFetch<GenericResponse<GeneralCategoryModel[]>>(url, {}).then(
        (data) => data.data
      ),
    {
      onError: (error) => {
        toast(<div>{error.message}</div>, { type: 'error' });
      },
    }
  );

  async function handleDelete() {
    if (selectedCollection) {
      const response = await apiFetch(
        API_PATHS.COLLECTIONS + '/' + selectedCollection.id,
        {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`,
          },
        }
      );
      if (!response.ok) {
        toast('Failed to delete collection', { type: 'error' });
        return;
      }

      toast('Collection deleted', { type: 'success' });
      setSelectedCollection(null);
      mutate();
    }
  }

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className='h-full'>
      <div className='flex justify-between pt-4 pb-8'>
        <h2 className='text-2xl font-semibold text-primary'>Collection List</h2>
        <Button
          as={Link}
          href={'/admin/collections/new'}
          className='btn btn-lg btn-primary'
        >
          Add new
        </Button>
      </div>

      <div className='relative overflow-x-auto shadow-md sm:rounded-lg'>
        <table className='w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400'>
          <thead className='text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400'>
            <tr>
              <th scope='col' className='px-6 py-3'>
                Image
              </th>
              <th scope='col' className='px-6 py-3'>
                Collection name
              </th>
              <th scope='col' className='px-6 py-3'>
                Slug
              </th>
              <th scope='col' className='px-6 py-3'>
                Description
              </th>
              <th scope='col' className='px-6 py-3'>
                Created At
              </th>
              <th scope='col' className='px-6 py-3'>
                Action
              </th>
            </tr>
          </thead>
          <tbody>
            {collections?.map((collection) => (
              <tr
                key={collection.id}
                className='odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700 border-gray-200'
              >
                <td className='px-6 py-4'>
                  {collection.image_url && (
                    <div className='h-10 w-10 relative'>
                      <Image
                        src={collection.image_url}
                        alt={collection.name}
                        fill
                        className='object-cover rounded'
                      />
                    </div>
                  )}
                </td>
                <th
                  scope='row'
                  className='px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white'
                >
                  <Link
                    href={`/admin/collections/${collection.id}`}
                    className='text-blue-500 hover:underline'
                  >
                    {collection.name}
                  </Link>
                </th>
                <td className='px-6 py-4'>{collection.slug}</td>
                <td className='px-6 py-4'>{collection.description}</td>
                <td className='px-6 py-4'>
                  {dayjs(collection.created_at).format('YYYY/MM/DD')}
                </td>
                <td className='px-6 py-4'>
                  <div className='flex space-x-2'>
                    <Link
                      href={`/admin/collections/${collection.id}`}
                      className='font-medium text-blue-600 dark:text-blue-500 hover:underline mr-3'
                    >
                      Edit
                    </Link>
                    <Button
                      onClick={() => setSelectedCollection(collection)}
                      className='font-medium text-red-600 dark:text-red-500 hover:underline'
                    >
                      Delete
                    </Button>
                  </div>
                </td>
              </tr>
            ))}
            {(!collections || collections.length === 0) && (
              <tr className='odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700 border-gray-200'>
                <td colSpan={6} className='px-6 py-4 text-center'>
                  No collections found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      <ConfirmDialog
        message='Are you sure you want to delete this collection?'
        onClose={() => {
          setSelectedCollection(null);
        }}
        open={!!selectedCollection}
        title='Delete Collection'
        onConfirm={handleDelete}
        confirmStyle='btn-danger'
      />
    </div>
  );
}
