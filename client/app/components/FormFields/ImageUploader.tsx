'use client';
import { Field, Label } from '@headlessui/react';
import { CloudArrowUpIcon } from '@heroicons/react/24/outline';
import clsx from 'clsx';
import React from 'react';
import { useDropzone } from 'react-dropzone';

interface ImageUploaderProps {
  label: string;
  multiple?: boolean;
  onUpload: (files: (File & { preview: string })[]) => void;
}
export const ImageUploader: React.FC<ImageUploaderProps> = ({
  label,
  multiple,
  onUpload,
}) => {
  const { getRootProps, getInputProps } = useDropzone({
    multiple: multiple,
    accept: {
      'image/*': ['.jpg', '.jpeg', '.png', '.gif'],
    },

    onDrop: (acceptedFiles) => {
      onUpload(
        acceptedFiles.map((file) =>
          Object.assign(file, {
            preview: URL.createObjectURL(file),
          })
        )
      );
    },
  });

  return (
    <Field>
      <Label className='text-base/6 font-semibold text-gray-500'>{label}</Label>
      <div className='w-full mt-1 flex flex-col items-center p-0'>
        <div
          {...getRootProps({
            className: clsx(
              'border-2 border-dashed border-gray-300 p-2 rounded-xl text-center cursor-pointer w-full',
              'hover:outline-none hover:border-blue-300 hover:text-blue-300 transition-all duration-500 ease-in-out'
            ),
          })}
        >
          <input {...getInputProps()} />
          <CloudArrowUpIcon
            className={clsx(
              'mx-auto text-gray-500 size-16',
              'hover:text-blue-300 transition-all duration-500 ease-in-out'
            )}
          />
          <p className='text-gray-600'>
            Drop your images here or click to browse
          </p>
        </div>
      </div>
    </Field>
  );
};
