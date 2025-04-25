'use client';

import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import { API_PATHS } from '@/lib/constants/api';
import { toast } from 'react-toastify';
import { getCookie } from 'cookies-next';
import { redirect } from 'next/navigation';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';

export default function Page() {
  const handleSave = async (form: FormData) => {
    const response = await apiFetch(API_PATHS.CATEGORIES, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${getCookie('access_token')}`,
      },
      body: form,
    });
    const data: GenericResponse<GeneralCategoryModel> = await response.json();
    if (!response.ok) {
      toast('Failed to save category', { type: 'error' });
      return;
    }
    toast('Category saved', { type: 'success' });
    redirect('/admin/categories/' + data.data.id);
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
