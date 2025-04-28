'use client';
import { API_PATHS } from '@/lib/constants/api';
import React, { use } from 'react';
import { toast } from 'react-toastify';
import useSWR from 'swr';
import Link from 'next/link';
import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import {
  GeneralCategoryModel,
  GenericResponse,
  CategoryProductModel,
} from '@/lib/definitions';
import Loading from '@/app/loading';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import LoadingInline from '@/components/Common/Loadings/LoadingInline';
import CategoryProductList from '../../_components/CategoryProductList';
import { apiFetch } from '@/lib/apis/api';

export default function AdminBrandDetail({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = use(params);
  const { data: brand, isLoading } = useSWR(
    API_PATHS.BRAND.replace(':id', slug),
    async (url) => {
      const response = await apiFetch<GenericResponse<GeneralCategoryModel>>(
        url,
        {}
      );
      return response.data;
    },
    {
      onError: (error) => {
        toast('Failed to fetch brand', { type: 'error', data: error });
      },
    }
  );

  const { data: products, isLoading: isLoadingProducts } = useSWR(
    API_PATHS.BRAND_PRODUCTS.replace(':id', slug),
    async (url) => {
      const response = await apiFetch<GenericResponse<CategoryProductModel[]>>(
        url,
        {}
      );
      return response.data;
    },
    {
      onError: (error) => {
        toast('Failed to fetch products', { type: 'error', data: error });
      },
    }
  );

  async function handleSave(data: FormData) {
    const response = await apiFetch<GenericResponse<GeneralCategoryModel>>(
      API_PATHS.BRAND.replace(':id', slug),
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
    <div>
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
      {isLoadingProducts ? (
        <LoadingInline />
      ) : (
        <CategoryProductList products={products ?? []} />
      )}
    </div>
  );
}
