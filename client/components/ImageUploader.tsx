'use client';

import React, { useState, useEffect } from 'react';
import Image from 'next/image';
import { Input } from '@headlessui/react';
import ImageEditor from './ImageEditor';
import clsx from 'clsx';

interface ImageUploaderProps {
  defaultImage?: string;
  name: string;
  label?: string;
  onChange?: (file: File | null) => void;
  className?: string;
  width?: number;
  height?: number;
  maxFileSizeMB?: number;
  aspectRatio?: number;
}

export default function ImageUploader({
  defaultImage,
  name,
  label = 'Upload image',
  onChange,
  className = '',
  width = 150,
  height = 150,
  maxFileSizeMB = 5,
  aspectRatio = 1,
}: ImageUploaderProps) {
  const [image, setImage] = useState<string | null>(defaultImage || null);
  const [file, setFile] = useState<File | null>(null);
  const [isEditorOpen, setIsEditorOpen] = useState<boolean>(false);
  const [tempImage, setTempImage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const placeholderImage = '/images/product-placeholder.webp';

  // Create a file from a base64 string
  const base64ToFile = (dataUrl: string, filename: string): File => {
    const arr = dataUrl.split(',');
    const mime = arr[0].match(/:(.*?);/)![1];
    const bstr = atob(arr[1]);
    let n = bstr.length;
    const u8arr = new Uint8Array(n);
    while (n--) {
      u8arr[n] = bstr.charCodeAt(n);
    }
    return new File([u8arr], filename, { type: mime });
  };

  // Handle file selection
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target?.files?.[0] || null;
    if (!selectedFile) return;

    // Validate file size (convert MB to bytes)
    const maxSizeBytes = maxFileSizeMB * 1024 * 1024;
    if (selectedFile.size > maxSizeBytes) {
      setError(`File size exceeds ${maxFileSizeMB}MB limit.`);
      return;
    }

    // Validate file type
    if (!selectedFile.type.startsWith('image/')) {
      setError('Selected file is not an image.');
      return;
    }

    // Reset any previous errors
    setError(null);

    // Read file as data URL to display in editor
    const reader = new FileReader();
    reader.onload = () => {
      const result = reader.result as string;
      setTempImage(result);
      setIsEditorOpen(true);
    };
    reader.readAsDataURL(selectedFile);
  };

  // Handle crop completion
  const handleCropComplete = (croppedImageData: string) => {
    try {
      // Convert base64 back to file for form submission
      const filename = file?.name || `cropped-image-${Date.now()}.jpg`;
      const croppedFile = base64ToFile(croppedImageData, filename);
      
      setFile(croppedFile);
      setImage(croppedImageData);
      setIsEditorOpen(false);
      setTempImage(null);
      
      if (onChange) {
        onChange(croppedFile);
      }
    } catch (err) {
      console.error('Error processing cropped image:', err);
      setError('Failed to process the cropped image.');
    }
  };

  // Cancel image editing
  const handleCancelEdit = () => {
    setIsEditorOpen(false);
    setTempImage(null);
  };

  // Create a hidden input field for the file
  useEffect(() => {
    // Create a DataTransfer
    if (file) {
      const dataTransfer = new DataTransfer();
      dataTransfer.items.add(file);
      
      // Find the file input element and set its files
      const fileInput = document.querySelector(`input[name="${name}"]`) as HTMLInputElement;
      if (fileInput) {
        fileInput.files = dataTransfer.files;
      }
    }
  }, [file, name]);

  return (
    <div className={className}>
      {isEditorOpen && tempImage ? (
        <div className="fixed inset-0 z-50 bg-black bg-opacity-75 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg p-4 max-w-3xl w-full">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Edit Image</h3>
            <ImageEditor
              image={tempImage}
              onComplete={handleCropComplete}
              onCancel={handleCancelEdit}
              aspect={aspectRatio}
            />
          </div>
        </div>
      ) : (
        <div>
          <div
            className={clsx(
              'relative block shadow-md shadow-lime-500 text-sm/6 font-medium text-gray-600',
              'hover:shadow-lime-600 transition-shadow cursor-pointer'
            )}
          >
            <Image
              className="rounded-sm object-contain"
              width={width}
              height={height}
              alt="Image Upload"
              priority
              src={image || placeholderImage}
            />
            <Input
              type="file"
              className="absolute cursor-pointer inset-0 w-full h-full opacity-0 z-10"
              id={name}
              name={name}
              accept="image/*"
              onChange={handleFileChange}
              aria-label={label}
            />
            <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-0 hover:bg-opacity-20 transition-all pointer-events-none">
              <span className="text-white opacity-0 hover:opacity-100 pointer-events-none">{label}</span>
            </div>
          </div>
          {error && (
            <div className="mt-1 text-xs text-red-500">{error}</div>
          )}
        </div>
      )}
    </div>
  );
}
