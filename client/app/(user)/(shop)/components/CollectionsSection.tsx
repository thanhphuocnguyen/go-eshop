import { GeneralCategoryModel } from '@/app/lib/definitions';
import Image from 'next/image';
import Link from 'next/link';

export default function CollectionsSection({
  collections,
}: {
  collections: GeneralCategoryModel[];
}) {
  return (
    <section className='container pt-20 pb-10'>
      <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8'>
        {collections?.map((collection) => (
          <Link 
            href={`/collections/${collection.slug}`} 
            key={collection.id} 
            className="group block"
          >
            <div className='bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-xl transition-all duration-300 transform group-hover:scale-[1.02]'>
              <div className='relative h-[42rem] w-full'>
                <Image
                  className='object-cover'
                  alt={`${collection.name} collection`}
                  src={collection.imageUrl ?? '/images/product-placeholder.webp'}
                  fill
                  sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                />
                <div className='absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent opacity-60 group-hover:opacity-70 transition-opacity duration-300'></div>
                
                <div className='absolute bottom-0 left-0 right-0 p-6 text-left z-10'>
                  <h3 className='text-2xl font-bold text-white mb-2 group-hover:text-blue-200 transition-colors drop-shadow-lg'>
                    {collection.name}
                  </h3>
                  <p className='text-gray-200 line-clamp-2 group-hover:text-gray-100 transition-colors'>
                    {collection.description || 'Explore this collection'}
                  </p>
                  
                  <div className="mt-4 inline-flex items-center text-white opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-y-2 group-hover:translate-y-0">
                    <span className="font-medium mr-2">Explore Collection</span>
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path>
                    </svg>
                  </div>
                </div>
              </div>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}
