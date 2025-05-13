import { GeneralCategoryModel } from '@/app/lib/definitions';
import Image from 'next/image';
import Link from 'next/link';

export default function CollectionsSection({
  collections,
}: {
  collections: GeneralCategoryModel[];
}) {
  return (
    <section className='grid h-[700px] pt-20 pb-10 grid-cols-3 space-x-10'>
      {collections?.map((collection) => (
        <div className='relative h-full' key={collection.id}>
          <Link href={`/collections/${collection.slug}`}>
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
          </Link>
        </div>
      ))}
    </section>
  );
}
