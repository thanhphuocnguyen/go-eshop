'use client';

import React from 'react';
import { ProductDetailForm } from '../_components/ProductForm';
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
    files: ProductImageFile[]
  ) {
    const resp = await apiFetch(API_PATHS.PRODUCTS, {
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

    if (resp.ok) {
      const data: GenericResponse<string> = await resp.json();
      const formData = new FormData();
      const promises = files.map((file) => {
        if (file.image) formData.append('file', file.image);
        return apiFetch(
          API_PATHS.PRODUCT_VARIANT_IMAGE_UPLOAD.replaceAll(
            ':id',
            data.data.toString()
          ),
          {
            method: 'POST',
            body: formData,
          }
        );
      });
      const uploadResp = await Promise.all(promises);
      const allOk = uploadResp.every((resp) => resp.ok);
      if (allOk) {
        toast.success('Product created successfully');
        redirect('/admin/products/' + data.data);
      } else {
        toast.error('Failed to upload images');
      }
    } else {
      toast.error('Failed to create product');
    }
  }

  return (
    <ProductDetailFormProvider>
      <ProductDetailForm onSubmit={onSubmit} />
    </ProductDetailFormProvider>
  );
};

export default Page;
