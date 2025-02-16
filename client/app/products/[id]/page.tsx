import { GenericListResponse } from '@/lib/types';

export async function generateStaticParams() {
  const products: GenericListResponse<{
    id: number;
  }> = await fetch(process.env.NEXT_API_URL + '/product/list').then((res) =>
    res.json()
  );

  return products.data.map((post) => ({
    id: post.id.toString(),
  }));
}
export default async function ProductDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const slug = (await params).id;
  const resp = await fetch(process.env.NEXT_API_URL + `/product/${slug}`).then(
    (res) => res.json()
  );
  const product = resp.data;
  return (
    <div>
      <h1>{product.name}</h1>
      <p>{product.description}</p>
    </div>
  );
}
