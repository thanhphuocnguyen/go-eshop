'use client';

import { API_PATHS } from '@/lib/constants/api';
import { use } from 'react';
import { getCookie } from 'cookies-next';
import { toast } from 'react-toastify';
import Link from 'next/link';
import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import useSWR from 'swr';
import Loading from '@/app/loading';
import { GenericResponse, CategoryProductModel } from '@/lib/definitions';
import LoadingInline from '@/components/Common/Loadings/LoadingInline';
import CategoryProductList from '../../_components/CategoryProductList';
import { apiFetch } from '@/lib/apis/api';

export default function AdminCategoryDetail({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = use(params);
  const { data: category, isLoading } = useSWR(
    API_PATHS.CATEGORY.replaceAll(':id', slug),
    async (url) => {
      const response = await apiFetch(url, {});
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

  const { data: products, isLoading: isLoadingProducts } = useSWR(
    API_PATHS.CATEGORY_PRODUCTS.replaceAll(':id', slug),
    async (url) => {
      const response = await apiFetch<GenericResponse<CategoryProductModel[]>>(
        url,
        {}
      );
      if (response.error) {
        toast.error(
          <div>
            Failed to fetch category products:
            <div>{JSON.stringify(response.error)}</div>
          </div>
        );
        return [];
      }
      return response.data;
    }
  );

  async function handleSave(data: FormData) {
    const response = await apiFetch<GenericResponse<number>>(
      API_PATHS.CATEGORY.replace(':id', slug),
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
      {isLoadingProducts ? (
        <LoadingInline />
      ) : (
        <CategoryProductList products={products ?? []} />
      )}
    </div>
  );
}
