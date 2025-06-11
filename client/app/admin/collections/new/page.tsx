'use client';

import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { CategoryEditForm } from '../../_components/CategoryEditForm';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';
import { redirect } from 'next/navigation';
import { GeneralCategoryModel } from '@/app/lib/definitions';
import { clientSideFetch } from '@/app/lib/api/apiClient';

export default function Page() {
  const handleSave = async (form: FormData) => {
    try {
      const { data, error } = await clientSideFetch<GeneralCategoryModel>(
        ADMIN_API_PATHS.COLLECTIONS,
        {
          method: 'POST',
          body: form,
        }
      );

      if (error) {
        console.error(error);
        toast.error(
          <div>
            <p className='text-red-500'>Failed to create collection</p>
            <p className='text-red-500'>{JSON.stringify(error)}</p>
          </div>
        );
        return;
      }
      toast.success('Collection created');
      redirect('/admin/collections/' + data.id);
    } catch (error) {
      console.error(error);
    }
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
