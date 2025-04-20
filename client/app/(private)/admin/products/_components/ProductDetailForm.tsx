'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import React from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Button, Fieldset, Legend } from '@headlessui/react';

import clsx from 'clsx';
import { ArrowLeftCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import {
  ProductDetailModel,
  ProductFormSchema,
  ProductModelForm,
} from '@/lib/definitions';

import { useProductDetailFormContext } from '../_lib/contexts/ProductFormContext';
import { ProductInfoForm } from './ProductInfoForm';
import { VariantInfoForm } from './VariantInfoForm';

export type SubmitResult = {
  createProductSuccess: boolean;
  uploadImagesSuccess: boolean;
  uploadVariantImagesSuccess: boolean;
};

interface ProductEditFormProps {
  onSubmit: (
    data?: ProductModelForm,
    productImages?: File[],
    variantImages?: File[],
    variantImageAssignments?: Record<string, string[]> // Map of image IDs to variant IDs they're assigned to
  ) => Promise<SubmitResult>;
  productDetail?: ProductDetailModel;
}

export const ProductDetailForm: React.FC<ProductEditFormProps> = ({
  onSubmit,
  productDetail,
}) => {
  const { productImages, variantImages, setProductImages, setVariantImages } =
    useProductDetailFormContext();

  const productForm = useForm<ProductModelForm>({
    resolver: zodResolver(ProductFormSchema),
    mode: 'onBlur',
    reValidateMode: 'onChange',
    defaultValues: productDetail
      ? {
          id: productDetail.id,
          variants: productDetail.variants.map((variant) => ({
            attributes: variant.attributes.map((attribute) => ({
              id: attribute.id,
              name: attribute.name,
              value: attribute.value,
            })),
            id: variant.id,
            price: variant.price,
            stock: variant.stock_qty,
            sku: variant.sku,
            weight: variant.weight,
            is_active: variant.is_active,
            image_url: variant.image_url,
          })),
          category: {
            id: productDetail.category.id,
            name: productDetail.category.name,
          },
          brand: productDetail.brand
            ? { id: productDetail.brand.id, name: productDetail.brand.name }
            : null,
          collection: productDetail.collection
            ? {
                id: productDetail.collection.id,
                name: productDetail.collection.name,
              }
            : null,
          images: productDetail?.images
            ? productDetail.images.map((image) => ({
                id: image.id,
                image_url: image.url,
              }))
            : [],
          description: productDetail.description,
          name: productDetail.name,
          price: productDetail.price,
          sku: productDetail.sku,
          is_active: productDetail.is_active,
          slug: productDetail.slug,
        }
      : {
          variants: [],
          brand: null,
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
  });
  const {
    formState: { isDirty, isSubmitting },
  } = productForm;

  async function submitHandler(data: ProductModelForm) {
    // Extract just the files from variant images
    const variantImageFiles = variantImages.map((img) => img.file);

    // Create assignments map
    let variantImageAssignments: Record<string, string[]> | undefined =
      undefined;
    if (variantImages.length) {
      variantImageAssignments = {};
      variantImages.forEach((img, index) => {
        // Use index as key for new images, id for existing ones
        const key = img.id || `new_${index}`;
        variantImageAssignments![key] = img.variantIds;
      });
    }

    const result = await onSubmit(
      isDirty ? data : undefined,
      productImages.length ? productImages : undefined,
      variantImageFiles.length ? variantImageFiles : undefined,
      variantImageAssignments
    );

    if (result.createProductSuccess) {
      // Success handling
    }
    if (result.uploadImagesSuccess) {
      setProductImages([]);
    }
    if (result.uploadVariantImagesSuccess) {
      setVariantImages([]);
    }
  }

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
              disabled={
                isSubmitting ||
                (!isDirty && !productImages.length && !variantImages.length)
              }
              type='submit'
              className={clsx(
                'btn text-lg btn-primary',
                isSubmitting && 'loading',
                isDirty ? 'btn-primary' : 'btn-disabled'
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
            <VariantInfoForm />
          </div>
        </Fieldset>
      </FormProvider>
    </div>
  );
};
