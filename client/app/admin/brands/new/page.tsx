'use client';

import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import { toast } from 'react-toastify';
import { redirect } from 'next/navigation';
import { GeneralCategoryModel } from '@/app/lib/definitions';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { clientSideFetch } from '@/app/lib/api/apiClient';

export default function Page() {
  const handleSave = async (form: FormData) => {
    const response = await clientSideFetch<GeneralCategoryModel>(
      ADMIN_API_PATHS.BRANDS,
      {
        method: 'POST',
        body: form,
      }
    );
    if (response.error) {
      console.error(response.error);
      toast.error(
        <div>
          <p className='text-red-500'>Failed to create brand</p>
          <p className='text-red-500'>{JSON.stringify(response.error)}</p>
        </div>
      );
      return;
    }
    if (response.data) {
      toast('Brand created', { type: 'success' });
      redirect('/admin/brands/' + response.data.id);
    }
  };
  return (
    <div className='h-full overflow-auto my-auto'>
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
