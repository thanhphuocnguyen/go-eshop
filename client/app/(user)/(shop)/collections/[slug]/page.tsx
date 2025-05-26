import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { serverSideFetch } from '@/app/lib/api/apiServer';
import { Pagination, ProductDetailModel } from '@/app/lib/definitions';
import ProductGrid from '@/app/components/Product/ProductGrid';
// Import the CategoryFilter component
import { CheckboxGroup } from './components/CategoryFilter';
import { collectionBySlugCache } from './layout';

// Define the search params type
type SearchParams = {
  page?: string;
  category?: string;
  brand?: string;
  attribute?: string;
};

// Define the params type
type Params = {
  slug: string;
};

async function getCollectionProducts(
  slug: string,
  page = 1,
  pageSize = 12,
  categorySlug?: string,
  brandSlug?: string,
  attributeSlug?: string
): Promise<{
  data: ProductDetailModel[];
  pagination: Pagination;
}> {
  const queryParams: Record<string, string> = {
    collections: slug,
  };

  // Add category filter if provided
  if (categorySlug) {
    queryParams.category = categorySlug;
  }

  // Add brand filter if provided
  if (brandSlug) {
    queryParams.brand = brandSlug;
  }

  // Add attribute filter if provided
  if (attributeSlug) {
    queryParams.attribute = attributeSlug;
  }

  const result = await serverSideFetch<ProductDetailModel[]>(
    `${PUBLIC_API_PATHS.PRODUCTS}?page=${page}&pageSize=${pageSize}`,
    {
      queryParams,
    }
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
  searchParams: Promise<SearchParams>;
}) {
  const { slug } = await params;
  const queries = await searchParams;
  const currentPage = queries.page ? Number(queries.page) : 1;
  const categoryFilter = queries.category || undefined;
  const brandFilter = queries.brand || undefined;
  const attributeFilter = queries.attribute || undefined;
  const pageSize = 12; // 3x4 grid

  // Fetch collection details and products
  const { categories, brands, attributes } = await collectionBySlugCache(slug);
  console.log(attributes, 'attributes');
  const { data: products, pagination } = await getCollectionProducts(
    slug,
    currentPage,
    pageSize,
    categoryFilter,
    brandFilter,
    attributeFilter
  );

  // Prepare filter indicators
  const activeCategory = categoryFilter
    ? categories.find((c) => c.id === categoryFilter)
    : null;

  const activeBrand = brandFilter
    ? brands.find((b) => b.id === brandFilter)
    : null;

  // Parse the attribute filter which should be in format "attributeKey:valueId"
  let activeAttribute = null;
  if (attributeFilter && attributeFilter.includes(':')) {
    const [attributeKey, valueId] = attributeFilter.split(':');
    
    // Check if this attribute key exists in our attributes object
    if (attributes[attributeKey]) {
      // Find the specific value within that attribute's values
      const value = attributes[attributeKey].find(v => v.id === valueId);
      if (value) {
        activeAttribute = {
          key: attributeKey,
          value: value
        };
      }
    }
  }

  // Check if any filters are active
  const hasActiveFilters = activeCategory || activeBrand || activeAttribute;

  return (
    <div className='flex flex-col lg:flex-row gap-6'>
      {/* Filter Panel */}
      <div className='lg:w-64 shrink-0'>
        <CheckboxGroup
          categories={categories}
          brands={brands}
          attributes={attributes}
        />
      </div>

      {/* Content Area */}
      <div className='flex-1'>
        {/* Product Grid */}
        <div>
          <div className='flex flex-col sm:flex-row sm:items-center justify-between mb-6'>
            <h2 className='text-2xl font-semibold'>
              {hasActiveFilters ? 'Filtered Products' : 'All Products'}
            </h2>

            {/* Filter indicators */}
            {hasActiveFilters && (
              <div className='mt-2 sm:mt-0 flex items-center flex-wrap gap-2 text-sm text-gray-500'>
                <span>Filtered by:</span>
                {activeCategory && (
                  <span className='bg-indigo-100 text-indigo-700 px-2 py-0.5 rounded-full text-xs font-medium flex items-center'>
                    <span className='mr-1'>Category:</span>{' '}
                    {activeCategory.name}
                  </span>
                )}
                {activeBrand && (
                  <span className='bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full text-xs font-medium flex items-center'>
                    <span className='mr-1'>Brand:</span> {activeBrand.name}
                  </span>
                )}
                {activeAttribute && (
                  <span className='bg-green-100 text-green-700 px-2 py-0.5 rounded-full text-xs font-medium flex items-center'>
                    <span className='mr-1'>Attribute:</span>{' '}
                    {`${activeAttribute.key}: ${activeAttribute.value.name}`}
                  </span>
                )}
              </div>
            )}
          </div>

          {products.length > 0 ? (
            <ProductGrid
              products={products}
              pagination={pagination}
              basePath={`/collections/${slug}`}
            />
          ) : (
            <div className='text-center py-16 bg-gray-50 rounded-lg'>
              <h3 className='text-xl font-medium text-gray-600'>
                No products found with the selected filters
              </h3>
              <p className='text-gray-500 mt-2'>
                Try adjusting your filters or check back later for new
                additions.
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
