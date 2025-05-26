import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { serverSideFetch } from '@/app/lib/api/apiServer';
import { ProductListModel } from '@/app/lib/definitions';
import CategoryProducts from './_components/CategoryProducts';
import CategoryFilters from './_components/CategoryFilters';
import { categoryCache } from './layout';

// Define the search params type
type SearchParams = {
  minPrice?: string;
  maxPrice?: string;
  rating?: string;
  page?: string;
};

async function getProducts(
  categoryIds: string,
  minPrice?: number,
  maxPrice?: number,
  rating?: number,
  page = 1,
  pageSize = 100
) {
  const queryParams: Record<string, string | number> = {
    categoryIds: categoryIds,
    page,
    pageSize: pageSize,
  };

  if (minPrice !== undefined) {
    queryParams.minPrice = minPrice;
  }

  if (maxPrice !== undefined) {
    queryParams.maxPrice = maxPrice;
  }

  if (rating !== undefined) {
    queryParams.minRating = rating;
  }

  const { data, error } = await serverSideFetch<ProductListModel[]>(
    PUBLIC_API_PATHS.PRODUCTS,
    { queryParams }
  );

  if (error) {
    console.error('Failed to fetch products:', error);
    return [];
  }

  return data || [];
}

export default async function CategoryPage({
  params,
  searchParams,
}: {
  params: Promise<{ slug: string }>;
  searchParams: Promise<SearchParams>;
}) {
  const { slug } = await params;
  const queries = await searchParams;
  const category = await categoryCache(slug);

  if (!category) {
    return null; // This should never happen because of the layout check
  }

  // Parse search params
  const minPrice = queries.minPrice ? Number(queries.minPrice) : undefined;
  const maxPrice = queries.maxPrice ? Number(queries.maxPrice) : undefined;
  const rating = queries.rating ? Number(queries.rating) : undefined;
  const page = queries.page ? Number(queries.page) : 1;
  // Fetch products with filters
  const products = await getProducts(
    category.id,
    minPrice,
    maxPrice,
    rating,
    page
  );

  // Calculate price range for filters
  let priceRange = { min: 0, max: 10000 };

  if (products.length > 0) {
    const minProductPrice = Math.min(...products.map((p) => p.minPrice));
    const maxProductPrice = Math.max(...products.map((p) => p.maxPrice));
    priceRange = { min: minProductPrice, max: maxProductPrice };
  }

  return (
    <div className='flex h-max flex-col lg:flex-row gap-8'>
      {/* Client-side filters component that controls URL parameters */}
      <CategoryFilters
        priceRange={priceRange}
        minPrice={minPrice}
        maxPrice={maxPrice}
        selectedRating={rating}
        categoryId={category.id}
      />

      {/* Use the CategoryProducts component correctly */}
      <CategoryProducts
        products={products}
        loadingProducts={false}
        minPrice={minPrice}
        maxPrice={maxPrice}
        selectedRating={rating}
      />
    </div>
  );
}
