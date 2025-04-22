'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import React, { useEffect } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Button, Fieldset, Legend } from '@headlessui/react';
import { redirect } from 'next/navigation';
import clsx from 'clsx';
import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import {
  GenericResponse,
  ProductDetailModel,
  ProductFormSchema,
  ProductModelForm,
  VariantModelForm,
} from '@/lib/definitions';

import { useProductDetailFormContext } from '../_lib/contexts/ProductFormContext';
import { ProductInfoForm } from './ProductInfoForm';
import { VariantInfoForm } from './VariantInfoForm';
import { apiFetch } from '@/lib/api/api';
import { API_PATHS } from '@/lib/constants/api';
import { toast } from 'react-toastify';

interface ProductEditFormProps {
  productDetail?: ProductDetailModel;
  mutate?: () => void;
}

export const ProductDetailForm: React.FC<ProductEditFormProps> = ({
  productDetail,
  mutate,
}) => {
  const { tempProductImages, setTempProductImages } =
    useProductDetailFormContext();

  const productForm = useForm<ProductModelForm>({
    resolver: zodResolver(ProductFormSchema),
    defaultValues: productDetail
      ? {
          variants: productDetail.variants,
          product_info: {
            category: productDetail.category,
            brand: productDetail.brand,
            collection: productDetail.collection,
            description: productDetail.description,
            name: productDetail.name,
            price: productDetail.price,
            sku: productDetail.sku,
            is_active: productDetail.is_active,
            slug: productDetail.slug,
          },
          product_images: productDetail.product_images.map((image) => ({
            id: image.id,
            url: image.url,
            assignments: image.assignments.map(
              (assignment) => assignment.entity_id
            ),
          })),
        }
      : {
          variants: [],
          product_info: {
            brand: {
              id: '',
              name: '',
            },
            category: {
              id: '',
              name: '',
            },
            collection: null,
            description: '',
            name: '',
            sku: '',
            slug: '',
            price: 1,
            is_active: true,
          },
          product_images: [],
        },
  });

  const {
    reset,
    formState: { isDirty, isSubmitting, dirtyFields },
  } = productForm;

  useEffect(() => {
    if (productDetail) {
      reset({
        variants: productDetail.variants,
        product_info: {
          category: productDetail.category,
          brand: productDetail.brand,
          collection: productDetail.collection,
          description: productDetail.description,
          name: productDetail.name,
          price: productDetail.price,
          sku: productDetail.sku,
          is_active: productDetail.is_active,
          slug: productDetail.slug,
        },
        product_images: productDetail.product_images.map((image) => ({
          id: image.id,
          url: image.url,
          assignments: image.assignments.map(
            (assignment) => assignment.entity_id
          ),
        })),
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [productDetail]);

  return (
    <div className='h-full px-6 py-3 overflow-auto'>
      <FormProvider {...productForm}>
        <Fieldset
          onSubmit={productForm.handleSubmit(submitHandler, console.error)}
          as='form'
        >
          <Link
            href={'/admin/products'}
            className='flex items-center mb-2 space-x-2'
          >
            <ArrowLeftCircleIcon className='size-6 text-primary' />
            <span className='text-primary text-lg hover:underline'>
              Back to Products
            </span>
          </Link>
          <Legend className='text-2xl flex justify-between font-bold text-primary mb-4'>
            {productDetail ? (
              <span>Edit Product: {productDetail.name}</span>
            ) : (
              <span>Create New Product</span>
            )}
            <Button
              disabled={isSubmitting || (!isDirty && !tempProductImages.length)}
              type='submit'
              className={clsx(
                'btn text-lg btn-primary',
                isSubmitting && 'loading',
                isDirty || tempProductImages.length
                  ? 'btn-primary'
                  : 'btn-disabled'
              )}
            >
              {isSubmitting ? (
                <span>{productDetail ? 'Saving...' : 'Creating...'}</span>
              ) : (
                <span>{productDetail ? 'Save' : 'Create'}</span>
              )}
            </Button>
          </Legend>

          {/* Combined Product Details and Variants Section */}
          <div className='bg-white rounded-lg shadow-md p-6'>
            <div className='mb-6'>
              <ProductInfoForm productDetail={productDetail} />
            </div>

            {/* Product Images */}
            <hr className='my-8' />
            {/* Product Variants */}
            <VariantInfoForm
            />
          </div>
        </Fieldset>
      </FormProvider>
    </div>
  );

  async function submitHandler(data: ProductModelForm) {
    let productID = productDetail?.id;
    const { variants, ...productData } = data;
    let isHaveAction = false;
    if (dirtyFields.product_info) {
      const rs = await onSubmitProductDetail(productData);
      // Success handling
      productID = rs;
      isHaveAction = true;
    }

    if (productID && tempProductImages.length) {
      await onUploadImages(productID);
      isHaveAction = true;
    }

    if (productID && dirtyFields.variants?.length) {
      await onSubmitVariants(productID, variants);
      isHaveAction = true;
    }

    if (!productDetail && productID) {
      redirect(`/admin/products/${productID}`);
    }

    if (mutate && isHaveAction) {
      mutate();
    }
  }

  async function onUploadImages(productId: string) {
    if (tempProductImages.length) {
      const productImageFormData = new FormData();
      tempProductImages.forEach((obj) => {
        if (obj) {
          productImageFormData.append('files', obj.file);
          productImageFormData.append(
            'assignments',
            JSON.stringify(obj.variantIds)
          );
        }
      });
      const imageUploadResp = await apiFetch<GenericResponse<unknown>>(
        API_PATHS.PRODUCT_IMAGES_UPLOAD.replaceAll(':id', productId),
        {
          method: 'POST',
          body: productImageFormData,
        }
      );
      if (imageUploadResp.error) {
        toast.error('Failed to upload images');
      }
      if (imageUploadResp.data) {
        toast.success('Images uploaded successfully');
        setTempProductImages([]);
      }
    }
  }

  async function onSubmitProductDetail(
    payload: Omit<ProductModelForm, 'variants'>
  ): Promise<string | undefined> {
    // return

    const { data, error } = await apiFetch<GenericResponse<{ id: string }>>(
      productDetail
        ? API_PATHS.PRODUCT_DETAIL.replace(':id', productDetail.id)
        : API_PATHS.PRODUCTS,
      {
        method: productDetail ? 'PUT' : 'POST',
        body: {
          ...payload.product_info,
          collection_id: payload.product_info.collection?.id || null,
          brand_id: payload.product_info.brand?.id || null,
          category_id: payload.product_info.category?.id || null,
        },
      }
    );

    if (error) {
      toast.error(
        <div>
          Failed to {productDetail ? 'update' : 'create'} product
          <br />
          {error.details}
        </div>
      );
      return undefined;
    }

    return data.id;
  }

  async function onSubmitVariants(prodId: string, payload: VariantModelForm[]) {
    const body = {
      variants: payload.map((variant) => ({
        ...variant,
        attributes: variant.attributes.map((attribute) => ({
          id: attribute.id,
          value_id: attribute.value?.id,
        })),
      })),
    };
    const { data, error } = await apiFetch<
      GenericResponse<{
        updated_ids: string[];
        created_ids: string[];
      }>
    >(API_PATHS.PRODUCT_VARIANTS.replace(':id', prodId), {
      method: 'PUT',
      body: body,
    });
    if (error) {
      toast.error('Failed to update variants');
      return;
    }
    if (data) {
      toast.success('Variants updated successfully');
    }
  }
};
