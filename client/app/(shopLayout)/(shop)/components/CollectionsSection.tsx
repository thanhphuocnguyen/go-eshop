import { apiFetch } from '@/lib/apis/api';
import { API_PATHS } from '@/lib/constants/api';
import { GeneralCategoryModel, GenericResponse } from '@/lib/definitions';
import Image from 'next/image';

export default async function CollectionsSection() {
  const { data: collections } = await apiFetch<
    GenericResponse<GeneralCategoryModel[]>
  >(`${API_PATHS.COLLECTIONS}?page=1&page_size=3`, {
    method: 'GET',
    nextOptions: {
      next: {
        tags: ['collections'],
      },
    },
  });

  return (
    <section className='grid h-[700px] pt-20 pb-10 grid-cols-3 space-x-10'>
      {collections?.map((collection) => (
        <div className='relative h-full' key={collection.id}>
          <div className='relative h-full rounded-md'>
            <Image
              className='relative object-cover rounded-md'
              alt='Collection image'
              src={collection.image_url ?? '/images/product-placeholder.webp'}
              fill
            />
          </div>
          <div className='relative mt-4'>
            <h3 className='text-xl text-gray-900 font-bold'>
              {collection.name}
            </h3>
            <p className='text-gray-500'>{collection.description}</p>
          </div>
        </div>
      ))}
    </section>
  );
}
