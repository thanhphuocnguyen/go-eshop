import axiosInstance from '@/axios/axios';

export async function generateStaticParams() {
  const posts = await axiosInstance.get('/posts');

  return posts.map((post) => ({
    slug: post.slug,
  }));
}
export default async function ProductDetailPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const slug = (await params).slug;
  const product = await axiosInstance.get(`/products/${slug}`);

  return (
    <div>
      <h1>{product.name}</h1>
      <p>{product.description}</p>
    </div>
  );
}
