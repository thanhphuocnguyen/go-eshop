'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import React, { useEffect, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Button, Fieldset, Legend } from '@headlessui/react';
import { redirect, useRouter } from 'next/navigation';
import clsx from 'clsx';
import { ArrowLeftCircleIcon, TrashIcon } from '@heroicons/react/24/outline';
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
import { apiFetch } from '@/lib/apis/api';
import { ADMIN_API_PATHS } from '@/lib/constants/api';
import { toast } from 'react-toastify';
import { ConfirmDialog } from '@/components/Common/Dialogs/ConfirmDialog';

interface ProductEditFormProps {
  productDetail?: ProductDetailModel;
  mutate?: () => void;
}

export const ProductDetailForm: React.FC<ProductEditFormProps> = ({
  productDetail,
  mutate,
}) => {
  const router = useRouter();
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
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
            images: productDetail.product_images.map((image) => ({
              id: image.id,
              url: image.url,
              role: image.role,
              assignments: image.assignments.map(
                (assignment) => assignment.entity_id
              ),
            })),
          },
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
            images: [],
          },
        },
  });

  const {
    reset,
    formState: { isDirty, isSubmitting, dirtyFields },
  } = productForm;

  useEffect(() => {
    if (productDetail) {
      reset(
        {
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
            images: productDetail.product_images.map((image) => ({
              id: image.id,
              url: image.url,
              role: image.role,
              assignments: image.assignments.map(
                (assignment) => assignment.entity_id
              ),
            })),
          },
        },
        {
          keepDirty: false,
        }
      );
    }
  }, [productDetail, reset]);

  const handleDeleteProduct = async () => {
    if (!productDetail?.id) return;

    setIsDeleting(true);
    try {
      const { error } = await apiFetch<GenericResponse<unknown>>(
        ADMIN_API_PATHS.PRODUCT_DETAIL.replace(':id', productDetail.id),
        {
          method: 'DELETE',
        }
      );

      if (error) {
        toast.error(`Failed to delete product: ${error.details}`);
      } else {
        toast.success(`Product "${productDetail.name}" deleted successfully`);
        router.push('/admin/products');
      }
    } catch (err) {
      toast.error('An unexpected error occurred while deleting the product');
      console.error(err);
    } finally {
      setIsDeleting(false);
      setShowDeleteConfirm(false);
    }
  };

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
            <div className='flex space-x-3'>
              {productDetail && (
                <Button
                  type='button'
                  disabled={isDeleting}
                  onClick={() => setShowDeleteConfirm(true)}
                  className={clsx(
                    'btn text-lg flex items-center',
                    isDeleting ? 'btn-disabled' : 'btn-danger'
                  )}
                >
                  <TrashIcon className='h-5 w-5 mr-1' />
                  {isDeleting ? 'Deleting...' : 'Delete'}
                </Button>
              )}
              <Button
                disabled={
                  isSubmitting || (!isDirty && !tempProductImages.length)
                }
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
            </div>
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

      {/* Delete Confirmation Dialog */}
      <ConfirmDialog
        open={showDeleteConfirm}
        title='Delete Product'
        message={`Are you sure you want to delete "${productDetail?.name}"? This action cannot be undone.`}
        onClose={() => setShowDeleteConfirm(false)}
        onConfirm={handleDeleteProduct}
        confirmStyle='bg-red-600 hover:bg-red-700'
      />
    </div>
  );
  async function submitHandler(data: ProductModelForm) {
    console.log(data);
    let productID = productDetail?.id;
    const { variants, ...productData } = data;
    let isAllSuccess = true;
    if (dirtyFields.product_info) {
      const rs = await onSubmitProductDetail(productData);
      // Success handling
      productID = rs;
      isAllSuccess &&= !!rs;
    }

    if (productID && tempProductImages.length) {
      const rs = await onUploadImages(productID);
      isAllSuccess &&= rs;
    }

    if (productID && dirtyFields.variants?.length) {
      const variantsToUpdate = new Array<VariantModelForm>(0);

      dirtyFields.variants.map((val, i) => {
        if (val) {
          variantsToUpdate.push(variants[i]);
        }
      });

      const rs = await onSubmitVariants(productID, variantsToUpdate);
      isAllSuccess &&= rs;
    }

    if (!productDetail && productID) {
      redirect(`/admin/products/${productID}`);
    }
    if (isAllSuccess && mutate) {
      mutate();
    }
  }

  async function onUploadImages(productId: string) {
    if (!tempProductImages.length) {
      return true;
    }
    const productImageFormData = new FormData();
    tempProductImages.forEach((obj) => {
      if (obj) {
        productImageFormData.append('files', obj.file);
        productImageFormData.append('roles', obj.role || 'gallery');
        productImageFormData.append(
          'assignments[]',
          JSON.stringify(obj.variantIds)
        );
      }
    });

    const { error, data } = await apiFetch<GenericResponse<unknown>>(
      ADMIN_API_PATHS.PRODUCT_IMAGES_UPLOAD.replaceAll(':id', productId),
      {
        method: 'POST',
        body: productImageFormData,
      }
    );

    if (error) {
      toast.error('Failed to upload images');
      return false;
    }
    if (data) {
      toast.success('Images uploaded successfully');
      setTempProductImages([]);
    }
    return true;
  }

  async function onSubmitProductDetail(
    payload: Omit<ProductModelForm, 'variants'>
  ): Promise<string | undefined> {
    const { data, error } = await apiFetch<GenericResponse<{ id: string }>>(
      productDetail
        ? ADMIN_API_PATHS.PRODUCT_DETAIL.replace(':id', productDetail.id)
        : ADMIN_API_PATHS.PRODUCTS,
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
    if (data) {
      toast.success(
        <div>
          {productDetail ? 'Updated' : 'Created'} product successfully
          <br />
          {data.id}
        </div>
      );
    }

    return data.id;
  }

  async function onSubmitVariants(prodId: string, payload: VariantModelForm[]) {
    const body = {
      variants: payload.map((variant) => ({
        ...variant,
        attributes: variant.attributes.map((attribute) => ({
          id: attribute.id,
          value_id: attribute.value_object?.id,
        })),
      })),
    };

    const { error } = await apiFetch<
      GenericResponse<{
        updated_ids: string[];
        created_ids: string[];
      }>
    >(ADMIN_API_PATHS.PRODUCT_VARIANTS.replace(':id', prodId), {
      method: 'PUT',
      body: body,
    });

    if (error) {
      toast.error('Failed to update variants');
      return false;
    }

    toast.success('Variants updated successfully');
    return true;
  }
};
