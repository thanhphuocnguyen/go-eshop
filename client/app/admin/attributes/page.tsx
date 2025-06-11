'use client';

import { useEffect, useState } from 'react';
import Loading from '@/app/loading';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { AttributeDetailModel } from '@/app/lib/definitions';
import { AddNewDialog } from './_components/AttribueFormDialog';
import { ConfirmDialog } from '@/app/components/Common/Dialogs/ConfirmDialog';
import { toast } from 'react-toastify';
import dayjs from 'dayjs';
import Link from 'next/link';
import { clientSideFetch } from '@/app/lib/api/apiClient';

export default function Page() {
  const [attributes, setAttributes] = useState<AttributeDetailModel[]>([]);
  const [modalMode, setModalMode] = useState<'add' | 'edit' | 'delete'>();
  const [loading, setLoading] = useState(true);
  const [selectedAttribute, setSelectedAttribute] =
    useState<AttributeDetailModel | null>(null);

  const fetchAttributes = async () => {
    const response = await clientSideFetch<AttributeDetailModel[]>(
      ADMIN_API_PATHS.ATTRIBUTES,
      {
        nextOptions: {
          next: {
            tags: ['attributes'],
          },
        },
      }
    );
    if (response.error) {
      toast('Failed to fetch attributes', { type: 'error' });
      return;
    }
    setAttributes(response.data ?? []);
    setLoading(false);
  };

  const handleDelete = async () => {
    if (selectedAttribute?.id) {
      const response = await clientSideFetch<AttributeDetailModel[]>(
        ADMIN_API_PATHS.ATTRIBUTE.replace(
          ':id',
          selectedAttribute.id.toString()
        ),
        {
          method: 'DELETE',
        }
      );
      if (response.error) {
        toast('Failed to delete attribute', { type: 'error' });
        return;
      }
      if (response.data) {
        setAttributes(attributes.filter((e) => e.id !== selectedAttribute?.id));
        setSelectedAttribute(null);
        setModalMode(undefined);
        toast('Remove attribute successfully', { type: 'success' });
      } else {
        toast('Failed to delete attribute', { type: 'error' });
      }
    }
  };

  useEffect(() => {
    fetchAttributes();
  }, []);

  if (loading) {
    return <Loading />;
  }

  return (
    <div className=''>
      <div className='flex justify-between pt-4 pb-8'>
        <h2 className='text-2xl font-semibold text-primary'>Attributes List</h2>
        <button
          onClick={() => {
            setModalMode('add');
          }}
          className='btn btn-lg btn-primary'
        >
          Add Attribute
        </button>
      </div>
      <AddNewDialog
        open={modalMode === 'add'}
        handleSubmitted={(newAttr) => {
          setAttributes((prev) => [...prev, newAttr]);
          setSelectedAttribute(null);
        }}
        onClose={() => {
          setSelectedAttribute(null);
          setModalMode(undefined);
        }}
      />
      <div className='relative overflow-x-auto shadow-md sm:rounded-lg'>
        <table className='w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400'>
          <thead className='text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400'>
            <tr>
              <th scope='col' className='px-6 py-3'>
                Id
              </th>
              <th scope='col' className='px-6 py-3'>
                Name
              </th>
              <th scope='col' className='px-6 py-3'>
                Values
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
            {attributes.map((e) => (
              <tr
                key={e.id}
                className='odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700 border-gray-200'
              >
                <td className='px-6 py-4'>{e.id}</td>
                <th
                  scope='row'
                  className='px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white'
                >
                  <Link
                    href={`/admin/attributes/${e.id}`}
                    className='text-blue-500 hover:underline'
                  >
                    {e.name}
                  </Link>
                </th>
                <td className='px-6 py-4'>
                  {e.values?.map((e) => e.name || e.code)?.join(', ')}
                </td>
                <td className='px-6 py-4'>
                  {dayjs(e.createdAt).format('YYYY/MM/DD')}
                </td>
                <td className='px-6 py-4'>
                  <div className='flex space-x-2'>
                    <Link
                      href={`/admin/attributes/${e.id}`}
                      className='font-medium text-blue-600 dark:text-blue-500 hover:underline mr-3'
                    >
                      Edit
                    </Link>
                    <button
                      onClick={() => {
                        setModalMode('delete');
                        setSelectedAttribute(e);
                      }}
                      className='font-medium text-red-600 dark:text-red-500 hover:underline'
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            ))}
            {attributes.length === 0 && (
              <tr className='odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700 border-gray-200'>
                <td colSpan={5} className='px-6 py-4 text-center'>
                  No attributes found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      <ConfirmDialog
        open={modalMode === 'delete'}
        title='Delete Attribute'
        message='Are you sure you want to delete this attribute?'
        onClose={() => {
          setModalMode(undefined);
          setSelectedAttribute(null);
        }}
        onConfirm={handleDelete}
        confirmStyle='btn-danger'
      />
    </div>
  );
}
