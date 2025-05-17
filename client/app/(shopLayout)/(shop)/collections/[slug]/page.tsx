import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { apiFetchServerSide } from '@/app/lib/apis/apiServer';
import {
  CollectionDetailModel,
  Pagination,
  ProductDetailModel,
} from '@/app/lib/definitions';
import { Metadata } from 'next';
import Image from 'next/image';
import ProductGrid from '@/components/Product/ProductGrid';
import { notFound } from 'next/navigation';
import { cache } from 'react';

// Define the search params type
type SearchParams = {
  page?: string;
};

// Define the params type
type Params = {
  slug: string;
};

// Dynamic metadata generation
export async function generateMetadata({
  params,
}: {
  params: Params;
}): Promise<Metadata> {
  try {
    const collection = await collectionBySlugCache(params.slug);

    return {
      title: `${collection.name} Collection | eShop`,
      description:
        collection.description ||
        `Browse the ${collection.name} collection at eShop.`,
      openGraph: {
        title: `${collection.name} Collection | eShop`,
        description:
          collection.description ||
          `Browse the ${collection.name} collection at eShop.`,
        images: collection.image_url
          ? [
              {
                url: collection.image_url,
                width: 1200,
                height: 630,
                alt: `${collection.name} Collection`,
              },
            ]
          : [],
      },
    };
  } catch (error) {
    return {
      title: 'Collection | eShop',
      description: 'Explore our curated collections',
    };
  }
}

async function getCollectionBySlug(
  slug: string
): Promise<CollectionDetailModel> {
  const result = await apiFetchServerSide<CollectionDetailModel>(
    PUBLIC_API_PATHS.COLLECTION.replace(':slug', slug)
  );
  if (result.error) {
    throw new Error(result.error.details, {
      cause: result.error,
    });
  }
  if (!result.data || result.error) {
    notFound();
  }

  return result.data;
}
const collectionBySlugCache = cache(getCollectionBySlug);

async function getCollectionProducts(
  slug: string,
  page = 1,
  pageSize = 12
): Promise<{
  data: ProductDetailModel[];
  pagination: Pagination;
}> {
  const result = await apiFetchServerSide<ProductDetailModel[]>(
    `${PUBLIC_API_PATHS.COLLECTION_PRODUCTS.replace(':slug', slug)}?page=${page}&page_size=${pageSize}`
  );

  return {
    data: result.data || [],
    pagination: result.pagination ?? {
      total: 0,
      page,
      pageSize: pageSize,
      hasNextPage: false,
      hasPreviousPage: false,
      totalPages: 1,
    },
  };
}

export default async function CollectionDetailPage({
  params,
  searchParams,
}: {
  params: Promise<Params>;
  searchParams: SearchParams;
}) {
  const { slug } = await params;
  const currentPage = searchParams.page ? Number(searchParams.page) : 1;
  const pageSize = 12; // 3x4 grid

  // Fetch collection details and products
  const collection = await collectionBySlugCache(slug);
  const { data: products, pagination } = await getCollectionProducts(
    slug,
    currentPage,
    pageSize
  );

  return (
    <div className='container mx-auto px-4 py-8'>
      {/* Collection Hero Section */}
      <div className='relative mb-12 h-[40vh] min-h-[400px] rounded-2xl overflow-hidden'>
        {collection.image_url ? (
          <Image
            src={collection.image_url}
            alt={collection.name}
            fill
            priority
            className='object-cover'
            sizes='100vw'
          />
        ) : (
          <div className='w-full h-full bg-gradient-to-r from-blue-600 to-indigo-800'></div>
        )}

        {/* Overlay */}
        <div className='absolute inset-0 bg-gradient-to-t from-black/80 to-black/30'></div>

        {/* Content */}
        <div className='absolute inset-0 flex flex-col justify-end p-8 md:p-16'>
          <div className='max-w-3xl'>
            <h1 className='text-3xl md:text-5xl font-bold text-white mb-4 tracking-tight'>
              {collection.name}
            </h1>
            <p className='text-gray-200 text-lg md:text-xl mb-4 max-w-2xl'>
              {collection.description ||
                `Explore our ${collection.name} collection.`}
            </p>
          </div>
        </div>
      </div>

      {/* Product Grid */}
      <div className='my-8'>
        <h2 className='text-2xl font-semibold mb-6'>
          Products in this Collection
        </h2>

        {products.length > 0 ? (
          <ProductGrid
            products={products}
            pagination={pagination}
            basePath={`/collections/${slug}`}
          />
        ) : (
          <div className='text-center py-16 bg-gray-50 rounded-lg'>
            <h3 className='text-xl font-medium text-gray-600'>
              No products found in this collection
            </h3>
            <p className='text-gray-500 mt-2'>
              Check back later for new additions.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
