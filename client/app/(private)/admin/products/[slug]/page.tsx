'use client';
import LoadingInline from '@/components/Common/Loadings/LoadingInline';
import { API_PATHS } from '@/lib/constants/api';
import {
  GenericResponse,
  ProductImageModel,
  ProductModelForm,
} from '@/lib/definitions';
import { use } from 'react';
import {
  ProductDetailForm,
  SubmitResult,
} from '../_components/ProductDetailForm';
import { toast } from 'react-toastify';
import { ProductDetailFormProvider } from '../_lib/contexts/ProductFormContext';
import { apiFetch } from '@/lib/api/api';
import { useProductDetail } from '../../_lib/hooks/useProductDetail';

export default function ProductFormEditPage({
  params,
}: {
  params: Promise<{
    slug: string;
  }>;
}) {
  const { slug } = use(params);

  const { productDetail, isLoading, mutate } = useProductDetail(slug);

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
    payload?: ProductModelForm,
    productImages?: File[],
    variantImages?: File[]
  ): Promise<SubmitResult> {
    const result: SubmitResult = {
      createProductSuccess: false,
      uploadImagesSuccess: false,
      uploadVariantImagesSuccess: false,
    };
    if (payload) {
      const resp = await apiFetch<
        GenericResponse<{
          id: string;
          variants: string[];
        }>
      >(API_PATHS.PRODUCT_DETAIL.replace(':id', slug), {
        method: 'PUT',
        body: {
          ...payload,
          variants: payload.variants.map((va) => {
            return {
              ...va,
              attributes: va.attributes.map((att) => {
                return {
                  value_id: att.value?.id,
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
        return result;
      }

      if (resp.data.variants.length && variantImages?.length) {
        const variantImgFormData = new FormData();
        variantImages.forEach((img) => {
          variantImgFormData.append('files', img);
        });

        const resp = await apiFetch<GenericResponse<ProductImageModel[]>>(
          API_PATHS.VARIANT_IMAGES_UPLOAD.replaceAll(':id', slug),
          {
            method: 'POST',
            body: variantImgFormData,
          }
        );

        if (resp.data) {
          toast.success('Variant images uploaded successfully');
          result.uploadVariantImagesSuccess = true;
        } else {
          toast.error('Failed to upload images');
        }
      }
    }

    result.createProductSuccess = true;

    if (productImages?.length) {
      const productImgFormData = new FormData();
      productImages.forEach((img) => {
        if (img) {
          productImgFormData.append('files', img);
        }
      });
      const resp = await apiFetch<GenericResponse<ProductImageModel[]>>(
        API_PATHS.PRODUCT_IMAGES_UPLOAD.replaceAll(':id', slug),
        {
          method: 'POST',
          body: productImgFormData,
        }
      );

      if (resp.data) {
        toast.success('Product images uploaded successfully');
        result.uploadImagesSuccess = true;
      } else {
        toast.error('Failed to upload images');
      }
    }

    mutate();
    toast.success('Update product successfully');
    return result;
  }
}
