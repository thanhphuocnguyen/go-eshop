'use client';
import LoadingInline from '@/components/Common/Loadings/LoadingInline';
import { API_PATHS } from '@/lib/constants/api';
import {
  GenericResponse,
  ProductDetailModel,
  ProductModelForm,
} from '@/lib/definitions';
import { use } from 'react';
import useSWR from 'swr';
import { ProductDetailForm } from '../_components/ProductForm';
import { getCookie } from 'cookies-next';
import { toast } from 'react-toastify';
import {
  ProductDetailFormProvider,
  ProductImageFile,
} from '../_lib/contexts/ProductFormContext';
import { apiFetch } from '@/lib/api/api';

export default function ProductFormEditPage({
  params,
}: {
  params: Promise<{
    slug: string;
  }>;
}) {
  const { slug } = use(params);
  const {
    data: productDetail,
    isLoading,
    mutate,
  } = useSWR(
    API_PATHS.PRODUCT_DETAIL.replace(':id', slug),
    (url) =>
      apiFetch<GenericResponse<ProductDetailModel>>(url).then(
        (data) => data.data
      ),
    {
      refreshInterval: 0,
      revalidateOnFocus: false,
      onError: (err) => {
        toast.error(
          <div>
            Failed to fetch product detail:
            <div>{JSON.stringify(err)}</div>
          </div>
        );
      },
    }
  );

  if (isLoading) {
    return (
      <div className='flex justify-center items-center h-full'>
        <LoadingInline />
      </div>
    );
  }

  return (
    <div className='h-full overflow-hidden'>
      <ProductDetailFormProvider>
        <ProductDetailForm onSubmit={onSubmit} productDetail={productDetail} />
      </ProductDetailFormProvider>
    </div>
  );

  async function onSubmit(
    values: ProductModelForm,
    productImages: ProductImageFile[],
    variantImages: ProductImageFile[]
  ) {
    if (Object.keys(values).length) {
      const resp = await apiFetch<
        GenericResponse<{
          id: string;
          variants: string[];
        }>
      >(API_PATHS.PRODUCT_DETAIL.replace(':id', slug), {
        method: 'PUT',
        body: {
          ...values,
          variants: values.variants.map((va) => {
            return {
              ...va,
              attributes: va.attributes.map((att) => {
                return {
                  value_ids: att.values.map((v) => {
                    return v.id;
                  }),
                  attribute_id: att.id ? Number(att.id) : null,
                };
              }),
            };
          }),
        },
      });

      if (resp.error) {
        toast.error(
          <div>
            Failed to update product:
            <div>{JSON.stringify(resp)}</div>
          </div>
        );
        return;
      }

      toast.success('Update product successfully');
      const productImgFormData = new FormData();
      if (productImages.length) {
        productImages.forEach((img) => {
          if (img.image) productImgFormData.append('files', img.image);
        });
        const resp = await apiFetch<GenericResponse<unknown>>(
          API_PATHS.PRODUCT_IMAGES_UPLOAD.replaceAll(':id', slug),
          {
            headers: {
              Authorization: `Bearer ${getCookie('token')}`,
            },
            method: 'POST',
            body: productImgFormData,
          }
        );
        if (resp.data) {
          toast.success('Product images uploaded successfully');
        } else {
          toast.error('Failed to upload images');
        }
      }
      const variantImgFormData = new FormData();
      if (resp.data.variants.length && variantImages.length) {
        const promises = resp.data.variants.reduce(
          (acc, curr, idx) => {
            if (variantImages[idx]?.image) {
              variantImgFormData.append('file', variantImages[idx].image);
              acc.push(
                apiFetch<GenericResponse<unknown>>(
                  API_PATHS.PRODUCT_VARIANT_IMAGE_UPLOAD.replaceAll(
                    ':id',
                    curr.toString()
                  ),
                  {
                    method: 'POST',
                    body: variantImgFormData,
                  }
                )
              );
            }
            return acc;
          },
          [] as Promise<GenericResponse<unknown>>[]
        );
        const uploadResp = await Promise.all(promises);
        const allOk = uploadResp.every((resp) => resp?.data);
        if (allOk) {
          toast.success('Variant images uploaded successfully');
        } else {
          toast.error('Failed to upload images');
        }
      }
      mutate();
    }
  }
}
