import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import React from 'react';
import { toast } from 'react-toastify';
import useSWR from 'swr';

type Category = {
  id: string;
  name: string;
};

export const CategoriesTable: React.FC<{ id: string }> = ({ id }) => {
  const { data: categories, isLoading } = useSWR(
    ADMIN_API_PATHS.DISCOUNT_CATEGORIES.replace(':id', id),
    (url) =>
      clientSideFetch<Category[]>(url, {
        queryParams: {
          page: 1,
          limit: 10,
        },
      }).then((res) => {
        if (res.error) {
          throw new Error(res.error.details);
        }

        return res.data;
      }),
    {
      onError: (error) => {
        console.error(`Failed to fetch categories: ${error.message}`);
        toast.error(`Failed to fetch categories: ${error.message}`, {
          position: 'top-right',
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
          progress: undefined,
        });
      },
    }
  );

  if (isLoading) {
    return (
      <div className='bg-gray-50 rounded-md p-8 text-center'>
        <p className='text-gray-500'>Loading categories...</p>
      </div>
    );
  }

  if (!categories || categories.length === 0) {
    return (
      <div className='bg-gray-50 rounded-md p-8 text-center'>
        <p className='text-gray-500'>
          No specific categories selected for this discount.
        </p>
      </div>
    );
  }

  return (
    <div className='overflow-x-auto'>
      <table className='min-w-full divide-y divide-gray-200'>
        <thead className='bg-gray-50'>
          <tr>
            <th
              scope='col'
              className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
            >
              Category Name
            </th>
            <th
              scope='col'
              className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
            >
              Category ID
            </th>
          </tr>
        </thead>
        <tbody className='bg-white divide-y divide-gray-200'>
          {categories.map((category) => (
            <tr key={category.id}>
              <td className='px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900'>
                {category.name}
              </td>
              <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                {category.id}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
