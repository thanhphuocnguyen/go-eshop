'use client';
import { Button } from '@headlessui/react';
import { XCircleIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';
import React, { useState } from 'react';
import { ImageUploader } from './ImageUploader';

// Define the UploadFile type locally since it's not exported from lib/definitions
type UploadFile = File & {
  preview: string;
};

interface ImageUploadFormProps {
  label: string;
  imageUrl?: string;
  onFileChange: (file: File | null) => void;
}

export const ImageUploadForm: React.FC<ImageUploadFormProps> = ({
  label,
  imageUrl,
  onFileChange,
}) => {
  const [showUploader, setShowUploader] = useState<boolean>(!imageUrl);
  const [file, setFile] = useState<File | null>(null);

  const handleFileUpload = (files: UploadFile[]) => {
    if (files.length > 0) {
      const uploadedFile = files[0];
      setFile(uploadedFile);
      onFileChange(uploadedFile);
      setShowUploader(false);
    }
  };

  const handleRemoveFile = () => {
    setFile(null);
    onFileChange(null);
    setShowUploader(true);
  };

  const handleChangeImage = () => {
    setShowUploader(true);
  };

  return (
    <div className="">
      {!showUploader && (imageUrl || file) ? (
        <div className="relative w-max h-max rounded-xl">
          {file && (
            <Button
              className="absolute top-3 bg-white rounded-full right-3 z-10"
              onClick={handleRemoveFile}
            >
              <XCircleIcon className="size-9 text-red-200" />
            </Button>
          )}
          <Image
            height={200}
            width={200}
            src={file ? URL.createObjectURL(file) : imageUrl || ''}
            alt={label}
            className="object-cover rounded-xl cursor-pointer"
            onClick={handleChangeImage}
          />
          <div className="mt-2 text-center">
            <Button
              type="button"
              className="text-sm text-blue-600 hover:text-blue-800"
              onClick={handleChangeImage}
            >
              Change Image
            </Button>
          </div>
        </div>
      ) : (
        <ImageUploader
          label={label}
          onUpload={handleFileUpload}
          multiple={false}
        />
      )}
    </div>
  );
};

export default ImageUploadForm;