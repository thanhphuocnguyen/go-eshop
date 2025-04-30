'use client';

import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import { ADMIN_API_PATHS } from '@/lib/constants/api';
import { toast } from 'react-toastify';
import { getCookie } from 'cookies-next';
import { redirect } from 'next/navigation';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import { apiFetch } from '@/lib/apis/api';

export default function Page() {
  const handleSave = async (form: FormData) => {
    const { data, error } = await apiFetch<
      GenericResponse<GeneralCategoryModel>
    >(ADMIN_API_PATHS.CATEGORIES, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${getCookie('access_token')}`,
      },
      body: form,
    });
    if (error) {
      toast.error(
        <div>
          <p className='text-red-500'>Failed to save category</p>
          <p className='text-red-500'>{JSON.stringify(error)}</p>
        </div>
      );
      return;
    }
    if (!data) {
      toast.error('Failed to save category');
      return;
    }
    toast.success('Category created');
    redirect('/admin/categories/' + data.id);
  };
  return (
    <div className='h-full overflow-hidden'>
      <Link
        href='/admin/categories'
        className='flex space-x-2 items-center hover:underline text-blue-400'
      >
        <ArrowLeftCircleIcon className='size-5 ' />
        <span className='text-blue-500'>Back to Categories</span>
      </Link>
      <CategoryEditForm handleSave={handleSave} title='Create new category' />
    </div>
  );
}
