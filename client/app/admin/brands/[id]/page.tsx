'use client';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import React, { use } from 'react';
import { toast } from 'react-toastify';
import useSWR from 'swr';
import Link from 'next/link';
import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import { GeneralCategoryModel } from '@/app/lib/definitions';
import Loading from '@/app/loading';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import { clientSideFetch } from '@/app/lib/api/apiClient';

export default function AdminBrandDetail({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const { data: brand, isLoading } = useSWR(
    ADMIN_API_PATHS.BRAND.replace(':id', id),
    async (url) => {
      const response = await clientSideFetch<GeneralCategoryModel>(url, {});
      return response.data;
    },
    {
      onError: (error) => {
        toast('Failed to fetch brand', { type: 'error', data: error });
      },
    }
  );

  async function handleSave(data: FormData) {
    const response = await clientSideFetch<GeneralCategoryModel>(
      ADMIN_API_PATHS.BRAND.replace(':id', id),
      {
        method: 'PUT',
        body: data,
      }
    );
    if (response.data) {
      toast('Category updated', { type: 'success' });
    } else if (response.error) {
      toast('Failed to update category', {
        type: 'error',
        data: response.error,
      });
    }
  }

  if (isLoading) {
    return <Loading />;
  }

  return (
    <div className='h-full overflow-auto my-auto'>
      <Link
        href='/admin/brands'
        className='flex space-x-2 items-center hover:underline text-blue-400'
      >
        <ArrowLeftCircleIcon className='size-5 ' />
        <span className='text-blue-500'>Back to brands</span>
      </Link>
      <CategoryEditForm
        data={brand}
        handleSave={handleSave}
        title='Brand Detail'
      />
      {/* {isLoadingProducts ? (
        <LoadingInline />
      ) : (
        <CategoryProductList products={products ?? []} />
      )} */}
    </div>
  );
}
