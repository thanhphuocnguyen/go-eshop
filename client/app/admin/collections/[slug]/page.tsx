'use client';

import { ADMIN_API_PATHS } from '@/lib/constants/api';
import { use } from 'react';
import { toast } from 'react-toastify';
import Link from 'next/link';
import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import useSWR from 'swr';
import Loading from '@/app/loading';
import {
  CategoryProductModel,
  GeneralCategoryModel,
  GenericResponse,
} from '@/lib/definitions';
import LoadingInline from '@/components/Common/Loadings/LoadingInline';
import CategoryProductList from '../../_components/CategoryProductList';
import { apiFetch } from '@/lib/apis/api';

export default function AdminCollectionDetail({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = use(params);
  const { data: category, isLoading } = useSWR(
    ADMIN_API_PATHS.COLLECTION.replace(':id', slug),
    async (url) => {
      const response = await apiFetch<GenericResponse<GeneralCategoryModel>>(
        url,
        {}
      );
      if (response.error) {
        toast('Failed to fetch category', { type: 'error' });
        return;
      }
      return response.data;
    }
  );

  const { data: products, isLoading: isLoadingProducts } = useSWR(
    ADMIN_API_PATHS.COLLECTION_PRODUCTS.replace(':id', slug),
    async (url) => {
      const response = await apiFetch<GenericResponse<CategoryProductModel[]>>(
        url,
        {}
      );

      return response.data;
    }
  );

  async function handleSave(data: FormData) {
    const response = await apiFetch<GenericResponse<number>>(
      ADMIN_API_PATHS.COLLECTION.replace(':id', slug),
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
    <div>
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
      {isLoadingProducts ? (
        <LoadingInline />
      ) : (
        <CategoryProductList products={products ?? []} />
      )}
    </div>
  );
}
