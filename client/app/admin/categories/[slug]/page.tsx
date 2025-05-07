'use client';

import { ADMIN_API_PATHS } from '@/lib/constants/api';
import { use } from 'react';
import { getCookie } from 'cookies-next';
import { toast } from 'react-toastify';
import Link from 'next/link';
import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import useSWR from 'swr';
import Loading from '@/app/loading';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import CategoryProductList from '../../_components/CategoryProductList';
import { apiFetch } from '@/lib/apis/api';

export default function AdminCategoryDetail({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = use(params);
  const { data: category, isLoading } = useSWR(
    ADMIN_API_PATHS.CATEGORY.replaceAll(':id', slug),
    async (url) => {
      const response = await apiFetch<GenericResponse<GeneralCategoryModel>>(
        url,
        {}
      );
      return response.data;
    },
    {
      onError: (error) => {
        toast.error(
          <div>
            Failed to fetch category:
            <div>{JSON.stringify(error)}</div>
          </div>
        );
      },
    }
  );

  async function handleSave(data: FormData) {
    const response = await apiFetch<GenericResponse<number>>(
      ADMIN_API_PATHS.CATEGORY.replace(':id', slug),
      {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${getCookie('access_token')}`,
        },
        body: data,
      }
    );
    if (response.error) {
      toast.error(
        <div>
          Failed to update category:
          <div>{JSON.stringify(response.error)}</div>
        </div>
      );
      return;
    }
    toast('Category updated', { type: 'success' });
  }

  if (isLoading) {
    return <Loading />;
  }

  return (
    <div className='overflow-y-auto h-full'>
      <Link
        href='/admin/categories'
        className='flex space-x-2 items-center hover:underline text-blue-400'
      >
        <ArrowLeftCircleIcon className='size-5 ' />
        <span className='text-blue-500'>Back to Categories</span>
      </Link>
      <CategoryEditForm
        data={category}
        handleSave={handleSave}
        title='Category Detail'
      />

      <CategoryProductList products={category?.products ?? []} />
    </div>
  );
}
