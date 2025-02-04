// 'use server';
import apiClient from '@/axios/axios';
import { API_PUBLIC_PATHS } from '@/lib/constants/api';
import { Category, GenericListResponse } from '@/lib/types';
import Link from 'next/link';

export default async function Home() {
  const categoriesResp = await apiClient.get<GenericListResponse<Category>>(
    API_PUBLIC_PATHS.CATEGORIES,
    {}
  );
  return (
    <div className='w-full py-3 h-full flex flex-col gap-3'>
      {categoriesResp.data.map((category) => (
        <div key={category.category_id}>
          <Link href={`/category/${category.category_id}`} className='text-xl font-bold hover:text-green-800 text-green-700 cursor-pointer' key={category.category_id}>
            {category.name}
          </Link>
        </div>
      ))}
    </div>
  );
}
