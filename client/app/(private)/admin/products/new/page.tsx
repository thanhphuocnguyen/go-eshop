'use client';

import React from 'react';
import { ProductDetailForm, SubmitResult } from '../_components/ProductForm';
import { API_PATHS } from '@/lib/constants/api';
import { toast } from 'react-toastify';
import { GenericResponse, ProductModelForm } from '@/lib/definitions';
import { redirect } from 'next/navigation';
import {
  ProductDetailFormProvider,
  ProductImageFile,
} from '../_lib/contexts/ProductFormContext';
import { apiFetch } from '@/lib/api/api';

const Page: React.FC = () => {
  async function onSubmit(
    payload: ProductModelForm,
    productImages: ProductImageFile[],
    variantImages: ProductImageFile[]
  ): Promise<SubmitResult> {
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
            value_ids: attribute.values?.map((value) => value.id),
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
        if (image.image) {
          productImageFormData.append('file', image.image);
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

      const variantImageFormData = new FormData();
      const promises = data.variants.reduce(
        (acc, variantID, index) => {
          if (variantID && variantImages[index]?.image) {
            variantImageFormData.append('file', variantImages[index].image);
            acc.push(
              apiFetch<GenericResponse<unknown>>(
                API_PATHS.PRODUCT_VARIANT_IMAGE_UPLOAD.replaceAll(
                  ':id',
                  variantID.toString()
                ),
                {
                  method: 'POST',
                  body: variantImageFormData,
                }
              )
            );
          }
          return acc;
        },
        [] as Promise<GenericResponse<unknown>>[]
      );

      const uploadResp = await Promise.all(promises);
      const allOk = uploadResp.every((resp) => !resp.error);
      if (allOk) {
        toast.success('Product created successfully');
        result.uploadVariantImagesSuccess = true;
      } else {
        toast.error('Failed to upload images');
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
