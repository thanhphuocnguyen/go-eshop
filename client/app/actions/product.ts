import { apiFetch } from '@/lib/api/api';
import { API_PATHS } from '@/lib/constants/api';
import { revalidateTag } from 'next/cache';

export async function addProductToCart(data: {
  productID: number;
  variantID: number;
  quantity: number;
}) {
  await apiFetch(process.env.NEXT_API_URL! + API_PATHS.CART_ITEM, {
    headers: {
      'Content-Type': 'application/json',
    },
    body: {
      product_id: data.productID,
      variant_id: data.variantID,
      quantity: data.quantity,
    },
    nextOptions: {
      next: {
        tags: ['cart'],
      },
    },
  });
  revalidateTag('cart');
}
