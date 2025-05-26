'use client';

import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { use } from 'react';
import { toast } from 'react-toastify';
import Link from 'next/link';
import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import useSWR from 'swr';
import Loading from '@/app/loading';
import { GeneralCategoryModel } from '@/app/lib/definitions';
import CategoryProductList from '../../_components/CategoryProductList';
import { clientSideFetch } from '@/app/lib/api/apiClient';

export default function AdminCollectionDetail({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const { data: category, isLoading } = useSWR(
    ADMIN_API_PATHS.COLLECTION.replace(':id', id),
    async (url) => {
      const response = await clientSideFetch<GeneralCategoryModel>(url, {});
      if (response.error) {
        toast('Failed to fetch category', { type: 'error' });
        return;
      }
      return response.data;
    }
  );

  async function handleSave(data: FormData) {
    const response = await clientSideFetch<number>(
      ADMIN_API_PATHS.COLLECTION.replace(':id', id),
      {
        method: 'PUT',
        body: data,
      }
    );
    if (response.data) {
      toast('Category updated', { type: 'success' });
    } else if (response.error) {
      toast('Failed to update category', { type: 'error' });
    }
  }

  if (isLoading) {
    return <Loading />;
  }

  return (
    <div className='flex flex-col h-full w-full'>
      <Link
        href='/admin/collections'
        className='flex space-x-2 items-center hover:underline text-blue-400'
      >
        <ArrowLeftCircleIcon className='size-5 ' />
        <span className='text-blue-500'>Back to Collections</span>
      </Link>
      <CategoryEditForm
        data={category}
        handleSave={handleSave}
        title='Collection Detail'
      />
      <CategoryProductList />
    </div>
  );
}
