'use client';

import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import { API_PATHS } from '@/lib/constants/api';
import { toast } from 'react-toastify';
import { redirect } from 'next/navigation';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import { apiFetch } from '@/lib/apis/api';

export default function Page() {
  const handleSave = async (form: FormData) => {
    const response = await apiFetch(API_PATHS.COLLECTIONS, {
      method: 'POST',
      body: form,
    });

    const data: GenericResponse<GeneralCategoryModel> = await response.json();
    if (!response.ok) {
      console.error(response);
      toast('Failed to save collection', { type: 'error' });
      return;
    }
    toast('Collection created', { type: 'success' });
    redirect('/admin/collections/' + data.data.id);
  };
  return (
    <div className=''>
      <Link
        href='/admin/collections'
        className='flex space-x-2 items-center hover:underline text-blue-400'
      >
        <ArrowLeftCircleIcon className='size-5 ' />
        <span className='text-blue-500'>Back to Collections</span>
      </Link>
      <CategoryEditForm handleSave={handleSave} title='Create new collection' />
    </div>
  );
}
