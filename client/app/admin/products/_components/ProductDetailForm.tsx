'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import React, { useEffect, useState } from 'react';
import { FormProvider, useForm, useWatch } from 'react-hook-form';
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
  UploadImageResponseModel,
} from '@/app/lib/definitions';

import { useProductDetailFormContext } from '../_lib/contexts/ProductFormContext';
import { ProductInfoForm } from './ProductInfoForm';
import { VariantInfoForm } from './VariantInfoForm';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';
import { ConfirmDialog } from '@/app/components/Common/Dialogs/ConfirmDialog';
import { KeyedMutator } from 'swr';
import { clientSideFetch } from '@/app/lib/api/apiClient';

interface ProductEditFormProps {
  productDetail?: ProductDetailModel;
  mutate?: KeyedMutator<ProductDetailModel>;
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
          productInfo: {
            category: productDetail.category,
            attributes: productDetail.attributes,
            brand: productDetail.brand,
            collection: productDetail.collection,
            description: productDetail.description,
            name: productDetail.name,
            price: productDetail.price,
            sku: productDetail.sku,
            isActive: productDetail.isActive,
            slug: productDetail.slug,
            images: productDetail.productImages.map((image) => ({
              id: image.id,
              url: image.url,
              role: image.role,
              assignments: image.assignments.map((assignment) => assignment.entityId),
            })),
          },
        }
      : {
          variants: [],
          productInfo: {
            brand: {
              id: '',
              name: '',
            },
            category: {
              id: '',
              name: '',
            },
            attributes: [],
            collection: null,
            description: '',
            name: '',
            sku: '',
            slug: '',
            price: 1,
            isActive: true,
            images: [],
          },
        },
  });

  const {
    reset,
    control,
    formState: { isDirty, isSubmitting, dirtyFields },
  } = productForm;

  const selectedAttributes = useWatch({
    control,
    name: 'productInfo.attributes',
    defaultValue: [],
  });

  const handleDeleteProduct = async () => {
    if (!productDetail?.id) return;

    setIsDeleting(true);
    try {
      const { error } = await clientSideFetch<GenericResponse<unknown>>(
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

  useEffect(() => {
    if (productDetail) {
      reset({
        variants: productDetail.variants,
        productInfo: {
          category: productDetail.category,
          brand: productDetail.brand,
          collection: productDetail.collection,
          attributes: productDetail.attributes,
          description: productDetail.description,
          name: productDetail.name,
          shortDescription: productDetail.shortDescription,
          price: productDetail.price,
          sku: productDetail.sku,
          isActive: productDetail.isActive,
          slug: productDetail.slug,
          images: productDetail.productImages.map((image) => ({
            id: image.id,
            url: image.url,
            role: image.role,
            assignments: image.assignments.map((assignment) => assignment.entityId),
          })),
        },
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
            {/* Product Description Header */}
            <div className='mb-4'>
              <h2 className='text-xl font-semibold text-gray-700'>
                Product Details
              </h2>
              <p className='text-gray-500 mt-1'>
                {productDetail
                  ? 'Edit your product information, variants, and images below.'
                  : 'Fill in the details to create a new product in your inventory.'}
              </p>
            </div>

            <div className='mb-6'>
              <ProductInfoForm productDetail={productDetail} />
            </div>

            {/* Product Images */}
            <hr className='my-8' />
            {/* Product Variants */}
            <VariantInfoForm selectedAttributes={selectedAttributes} />
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
    let productID = productDetail?.id;
    let isAllSuccess = true;
    if (dirtyFields.productInfo) {
      const rs = await onSubmitProductDetail(data);
      productID = rs;
      isAllSuccess &&= !!rs;
    }

    if (productID && tempProductImages.length) {
      const rs = await onUploadImages(productID);
      isAllSuccess &&= rs;
    }

    if (!productDetail && productID) {
      redirect(`/admin/products/${productID}`);
    }
    if (isAllSuccess && mutate) {
      await mutate();
    }
    reset();
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

    const { error, data } = await clientSideFetch<
      GenericResponse<UploadImageResponseModel>
    >(ADMIN_API_PATHS.PRODUCT_IMAGES_UPLOAD.replaceAll(':id', productId), {
      method: 'POST',
      body: productImageFormData,
    });

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
    payload: ProductModelForm
  ): Promise<string | undefined> {
    const variants = payload.variants.map((variant) => ({
      ...variant,
      attributes: variant.attributes.map((attribute) => ({
        id: attribute.id,
        valueId: attribute.valueObject?.id,
      })),
    }));
    const { data, error } = await clientSideFetch<{ id: string }>(
      productDetail
        ? ADMIN_API_PATHS.PRODUCT_DETAIL.replace(':id', productDetail.id)
        : ADMIN_API_PATHS.PRODUCTS,
      {
        method: productDetail ? 'PUT' : 'POST',
        body: {
          ...payload.productInfo,
          variants,
          collectionId: payload.productInfo.collection?.id || null,
          attributes: payload.productInfo.attributes,
          brandId: payload.productInfo.brand?.id || null,
          categoryId: payload.productInfo.category?.id || null,
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
};
