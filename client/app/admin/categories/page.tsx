'use client';

import { GeneralCategoryModel } from '@/app/lib/definitions';
import { useEffect, useState } from 'react';

import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { Button } from '@headlessui/react';
import Link from 'next/link';
import { ConfirmDialog } from '@/app/components/Common/Dialogs/ConfirmDialog';
import { toast } from 'react-toastify';
import Image from 'next/image';
import { clientSideFetch } from '@/app/lib/api/apiClient';

export default function Page() {
  const [categories, setCategories] = useState<GeneralCategoryModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] =
    useState<GeneralCategoryModel | null>(null);

  async function handleDelete() {
    if (selectedCategory) {
      const { error } = await clientSideFetch<boolean>(
        ADMIN_API_PATHS.CATEGORIES + '/' + selectedCategory.id,
        {
          method: 'DELETE',
        }
      );

      if (error) {
        toast('Failed to delete category', { type: 'error' });
        return;
      }

      const newCategories = categories.filter(
        (category) => category.id !== selectedCategory.id
      );
      toast('Category deleted', { type: 'success' });
      setCategories(newCategories);
      setSelectedCategory(null);
    }
  }

  // Fetch categories
  useEffect(() => {
    (async () => {
      setLoading(true);
      const { data, error } = await clientSideFetch<GeneralCategoryModel[]>(
        ADMIN_API_PATHS.CATEGORIES,
        {
          nextOptions: {
            next: {
              tags: ['categories'],
            },
          },
        }
      );
      setLoading(false);
      if (error) {
        toast.error(
          <div>
            Failed to fetch categories:
            <div>{JSON.stringify(error)}</div>
          </div>
        );
        return;
      }
      setCategories(data ?? []);
    })();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }
  return (
    <div className='h-full'>
      <div className='flex justify-between items-center pt-4 pb-8'>
        <h2 className='text-2xl font-semibold text-primary'>Category List</h2>
        <Link href='/admin/categories/new' className='btn btn-lg btn-primary'>
          Create Category
        </Link>
      </div>

      <div className='relative overflow-x-auto shadow-md sm:rounded-lg'>
        <table className='w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400'>
          <thead className='text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400'>
            <tr>
              <th scope='col' className='px-6 py-3'>
                Image
              </th>
              <th scope='col' className='px-6 py-3'>
                Category name
              </th>
              <th scope='col' className='px-6 py-3'>
                Slug
              </th>
              <th scope='col' className='px-6 py-3'>
                Status
              </th>
              <th scope='col' className='px-6 py-3'>
                Action
              </th>
            </tr>
          </thead>
          <tbody>
            {categories.map((category) => (
              <tr
                key={category.id}
                className='odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700 border-gray-200'
              >
                <td className='px-6 py-4'>
                  {category.imageUrl && (
                    <div className='h-10 w-10 relative'>
                      <Image
                        src={category.imageUrl}
                        sizes='100%'
                        alt={category.name}
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
                    href={`/admin/categories/${category.id}`}
                    className='text-blue-500 hover:underline'
                  >
                    {category.name}
                  </Link>
                </th>
                <td className='px-6 py-4'>{category.slug}</td>
                <td className='px-6 py-4'>
                  <span
                    className={`px-2 py-1 text-xs font-medium rounded-full ${
                      category.published
                        ? 'bg-green-100 text-green-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {category.published ? 'Published' : 'Draft'}
                  </span>
                  {category.remarkable && (
                    <span className='ml-2 px-2 py-1 text-xs font-medium bg-blue-100 text-blue-800 rounded-full'>
                      Featured
                    </span>
                  )}
                </td>
                <td className='px-6 py-4'>
                  <div className='flex space-x-2'>
                    <Link
                      href={`/admin/categories/${category.id}`}
                      className='font-medium text-blue-600 dark:text-blue-500 hover:underline mr-3'
                    >
                      Edit
                    </Link>
                    <Button
                      onClick={() => setSelectedCategory(category)}
                      className='font-medium text-red-600 dark:text-red-500 hover:underline'
                    >
                      Delete
                    </Button>
                  </div>
                </td>
              </tr>
            ))}
            {(!categories || categories.length === 0) && (
              <tr className='odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700 border-gray-200'>
                <td colSpan={5} className='px-6 py-4 text-center'>
                  No categories found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      <ConfirmDialog
        message='Are you sure you want to delete this category?'
        onClose={() => {
          setSelectedCategory(null);
        }}
        open={!!selectedCategory}
        title='Delete Category'
        onConfirm={handleDelete}
        confirmStyle='btn-danger'
      />
    </div>
  );
}
