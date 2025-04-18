import { API_PATHS } from '@/lib/constants/api';
import { revalidateTag } from 'next/cache';
import { cookies } from 'next/headers';

export async function addProductToCart(
  state: any,
  data: {
    productID: number;
    variantID: number;
    quantity: number;
  }
) {
  const cookieStore = await cookies();
  await apiFetch(process.env.NEXT_API_URL! + API_PATHS.CART_ITEM, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${cookieStore.get('token')}`,
    },
    body: JSON.stringify({
      product_id: data.productID,
      variant_id: data.variantID,
      quantity: data.quantity,
    }),
    next: {
      tags: ['cart'],
    },
  });
  revalidateTag('cart');
}
