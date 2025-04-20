'use client';

import React from 'react';
import {
  ProductDetailForm,
  SubmitResult,
} from '../_components/ProductDetailForm';
import { API_PATHS } from '@/lib/constants/api';
import { toast } from 'react-toastify';
import {
  GenericResponse,
  ProductImageModel,
  ProductModelForm,
} from '@/lib/definitions';
import { redirect } from 'next/navigation';
import { ProductDetailFormProvider } from '../_lib/contexts/ProductFormContext';
import { apiFetch } from '@/lib/api/api';

const Page: React.FC = () => {
  async function onSubmit(
    payload: ProductModelForm,
    productImages: File[]
  ): Promise<SubmitResult> {
    console.log(payload);
    console.log(productImages);
    // return
    const result: SubmitResult = {
      createProductSuccess: false,
      uploadImagesSuccess: false,
      uploadVariantImagesSuccess: false,
    };

    const { data, error } = await apiFetch<
      GenericResponse<{
        id: string;
        variants: string[];
      }>
    >(API_PATHS.PRODUCTS, {
      method: 'POST',
      body: {
        ...payload,
        collection_id: payload.collection?.id || null,
        brand_id: payload.brand?.id || null,
        category_id: payload.category?.id || null,
        variants: payload.variants.map((variant) => ({
          ...variant,
          attributes: variant.attributes.map((attribute) => ({
            id: attribute.id,
            value_id: attribute.value?.id,
          })),
        })),
      },
    });

    if (error) {
      toast.error('Failed to create product');
      return result;
    }

    if (data) {
      result.createProductSuccess = true;
      const productImageFormData = new FormData();
      productImages.forEach((image) => {
        if (image) {
          productImageFormData.append('files', image);
        }
      });
      const imageUploadResp = await apiFetch<GenericResponse<unknown>>(
        API_PATHS.PRODUCT_IMAGES_UPLOAD.replaceAll(':id', data.id),
        {
          method: 'POST',
          body: productImageFormData,
        }
      );
      if (imageUploadResp.error) {
        toast.error('Failed to upload images');
      }
      if (imageUploadResp.data) {
        result.uploadImagesSuccess = true;
      }

      redirect('/admin/products/' + data.id);
    } else {
      toast.error('Failed to create product');
    }

    return result;
  }

  return (
    <ProductDetailFormProvider>
      <ProductDetailForm onSubmit={onSubmit} />
    </ProductDetailFormProvider>
  );
};

export default Page;
