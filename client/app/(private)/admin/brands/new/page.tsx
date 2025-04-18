'use client';

import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import { toast } from 'react-toastify';
import { redirect } from 'next/navigation';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import { API_PATHS } from '@/lib/constants/api';
import { apiFetch } from '@/lib/api/api';

export default function Page() {
  const handleSave = async (form: FormData) => {
    const response = await apiFetch<GenericResponse<GeneralCategoryModel>>(
      API_PATHS.BRANDS,
      {
        method: 'POST',
        body: form,
      }
    );
    if (response.data) {
      toast('Brand created', { type: 'success' });
      redirect('/admin/brands/' + response.data.id);
    } else {
      console.error(response);
      toast('Failed to save brand', { type: 'error', data: response.error });
    }
  };
  return (
    <div className=''>
      <Link
        href='/admin/brands'
        className='flex space-x-2 items-center hover:underline text-blue-400'
      >
        <ArrowLeftCircleIcon className='size-5 ' />
        <span className='text-blue-500'>Back to brands</span>
      </Link>
      <CategoryEditForm handleSave={handleSave} title='Create new brand' />
    </div>
  );
}
