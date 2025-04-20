'use client';
import { use, useEffect, useState } from 'react';
import { useProductDetail } from '../../../_lib/hooks/useProductDetail';
import LoadingInline from '@/components/Common/Loadings/LoadingInline';
import { ArrowLeftCircleIcon, XCircleIcon } from '@heroicons/react/16/solid';
import Link from 'next/link';
import { Button } from '@headlessui/react';
import { ImageUploader } from '@/components/FormFields';
import { StyledMultipleComboBox } from '@/components/FormFields/StyledMultipleComboBox';
import Image from 'next/image';

export default function VariantImagesPage({
  params,
}: {
  params: Promise<{
    slug: string;
  }>;
}) {
  const { slug } = use(params);

  const { productDetail, isLoading } = useProductDetail(slug);

  const [variantImages, setVariantImages] = useState<
    { file: File; preview: string; variantIds: string[] }[]
  >([]);

  // Handle image upload
  const handleImageUpload = (files: (File & { preview: string })[]) => {
    const newImages = files.map((file) => ({
      file,
      preview: file.preview,
      variantIds: [],
    }));

    setVariantImages([...variantImages, ...newImages]);
  };

  // Handle image removal
  const handleRemoveImage = (index: number) => {
    setVariantImages((prev) => prev.filter((_, idx) => idx !== index));
  };

  // Handle variant assignment for an image
  const handleAssignVariants = (
    index: number,
    selectedVariantIds: string[]
  ) => {
    setVariantImages((prev) => {
      const updated = [...prev];
      updated[index] = {
        ...updated[index],
        variantIds: selectedVariantIds,
      };
      return updated;
    });
  };

  // Initialize with existing variant images if available
  useEffect(() => {
    if (
      productDetail?.variant_images &&
      productDetail.variant_images.length > 0 &&
      variantImages.length === 0
    ) {
      const existingImages = productDetail.variant_images.map((img: any) => ({
        id: img.id,
        file: new File([], img.external_id),
        preview: img.url,
        variantIds: img.assignments.map((assignment: any) =>
          assignment.entity_id.toString()
        ),
      }));

      setVariantImages(existingImages);
    }
  }, [productDetail, setVariantImages]);

  // Get variant options for the multiple select

  if (isLoading || !productDetail) {
    return (
      <div className='flex justify-center items-center h-full'>
        <LoadingInline />
      </div>
    );
  }

  const variantOptions = productDetail.variants.map((variant) => ({
    id: variant.id,
    name: `${variant.sku} - ${variant.attributes
      .map(
        (attr) =>
          `${attr.name}: ${attr.value.display_value || attr.value.value}`
      )
      .join(', ')}`,
  }));

  return (
    <div className='h-full px-6 py-3 overflow-auto'>
      <Link
        href={`/admin/products/${slug}`}
        className='flex items-center mb-2 space-x-2'
      >
        <ArrowLeftCircleIcon className='size-6 text-primary' />
        <span className='text-primary text-lg hover:underline'>
          Back to Product
        </span>
      </Link>
      <div className='flex justify-between items-center mb-2'>
        <h1 className='text-2xl font-semibold text-primary mb-2'>
          Variant Images
        </h1>
        <Button className='btn btn-primary text-lg'>Save</Button>
      </div>
      <p className='text-sm text-gray-500 mb-4'>
        Upload images for variants. Each image can be assigned to multiple
        variants.
      </p>
      <div className='mb-8'>
        {/* Image uploader */}
        <div className='mb-2'>
          <ImageUploader
            label='Upload Variant Images'
            multiple={true}
            onUpload={handleImageUpload}
          />
        </div>

        {/* Uploaded images with variant assignment */}
        <div className='grid grid-cols-1 gap-6 mt-4'>
          {variantImages.map((image, index) => (
            <div
              key={index}
              className='border border-gray-200 rounded-lg p-4 bg-white shadow-sm'
            >
              <div className='flex items-start space-x-4'>
                {/* Image preview */}
                <div className='relative h-32 w-32 flex-shrink-0'>
                  <Image
                    src={image.preview}
                    alt={`Variant image ${index + 1}`}
                    fill
                    className='object-cover rounded-md'
                  />
                </div>

                {/* Variant assignment */}
                <div className='flex-1'>
                  <div className='flex justify-between mb-2'>
                    <h4 className='font-medium text-gray-700'>
                      Image {index + 1}
                    </h4>
                    <button
                      type='button'
                      onClick={() => handleRemoveImage(index)}
                      className='text-red-500 hover:text-red-700 transition-colors'
                    >
                      <XCircleIcon className='size-6' />
                    </button>
                  </div>

                  <StyledMultipleComboBox
                    label='Assign to Variants'
                    selected={variantOptions?.filter((opt) =>
                      image.variantIds.includes(opt.id)
                    )}
                    options={variantOptions}
                    getDisplayValue={(opt) => opt.name}
                    setSelected={(selected) =>
                      handleAssignVariants(
                        index,
                        selected.map((s) => s.id)
                      )
                    }
                  />

                  {image.variantIds.length === 0 && (
                    <p className='text-sm text-amber-600 mt-1'>
                      This image is not assigned to any variants
                    </p>
                  )}
                </div>
              </div>
            </div>
          ))}

          {variantImages.length === 0 && (
            <div className='text-center p-4 border border-dashed border-gray-300 rounded-lg bg-gray-50'>
              <p className='text-gray-500'>No variant images uploaded yet</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );

  async function onSubmit() {
    // Handle form submission logic
  }
}
