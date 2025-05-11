'use client';

import React, { useState, useCallback } from 'react';
import Cropper from 'react-easy-crop';
import { Button } from '@headlessui/react';
import { createStyles } from './image-editor-styles';

interface ImageEditorProps {
  image: string;
  onComplete: (croppedImage: string) => void;
  onCancel: () => void;
  aspect?: number;
}

interface CropArea {
  x: number;
  y: number;
  width: number;
  height: number;
}

// Function to create an image from the crop area
const createImage = (url: string): Promise<HTMLImageElement> =>
  new Promise((resolve, reject) => {
    const image = new Image();
    image.addEventListener('load', () => resolve(image));
    image.addEventListener('error', (error) => reject(error));
    image.src = url;
  });

// Function to get the cropped image as base64
const getCroppedImg = async (
  imageSrc: string,
  pixelCrop: CropArea
): Promise<string> => {
  const image = await createImage(imageSrc);
  const canvas = document.createElement('canvas');
  const ctx = canvas.getContext('2d');

  if (!ctx) {
    throw new Error('No 2d context');
  }

  // Set canvas dimensions to the cropped size
  canvas.width = pixelCrop.width;
  canvas.height = pixelCrop.height;

  // Draw the cropped image onto the canvas
  ctx.drawImage(
    image,
    pixelCrop.x,
    pixelCrop.y,
    pixelCrop.width,
    pixelCrop.height,
    0,
    0,
    pixelCrop.width,
    pixelCrop.height
  );

  // Return the canvas as a base64 string
  return canvas.toDataURL('image/jpeg');
};

export default function ImageEditor({
  image,
  onComplete,
  onCancel,
  aspect = 1,
}: ImageEditorProps) {
  // Initialize styles
  React.useEffect(() => {
    createStyles();
  }, []);

  const [crop, setCrop] = useState({ x: 0, y: 0 });
  const [zoom, setZoom] = useState(1);
  const [croppedAreaPixels, setCroppedAreaPixels] = useState<CropArea | null>(
    null
  );

  const onCropComplete = useCallback((_: any, croppedAreaPixels: CropArea) => {
    setCroppedAreaPixels(croppedAreaPixels);
  }, []);

  const handleComplete = useCallback(async () => {
    try {
      if (croppedAreaPixels) {
        const croppedImage = await getCroppedImg(image, croppedAreaPixels);
        onComplete(croppedImage);
      }
    } catch (e) {
      console.error('Error getting cropped image:', e);
    }
  }, [croppedAreaPixels, image, onComplete]);

  return (
    <div className='relative h-[500px] w-full'>
      <div className='relative h-[400px] w-full'>
        <Cropper
          image={image}
          crop={crop}
          zoom={zoom}
          aspect={aspect}
          onCropChange={setCrop}
          onCropComplete={onCropComplete}
          onZoomChange={setZoom}
        />
      </div>
      <div className='mt-4 flex items-center justify-between p-4'>
        <div className='flex-1'>
          <label
            htmlFor='zoom'
            className='block text-sm font-medium text-gray-700'
          >
            Zoom
          </label>
          <input
            type='range'
            id='zoom'
            min={1}
            max={3}
            step={0.1}
            value={zoom}
            onChange={(e) => setZoom(Number(e.target.value))}
            className='w-full'
          />
        </div>
        <div className='ml-4 flex space-x-2'>
          <Button
            type='button'
            onClick={onCancel}
            className='px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50'
          >
            Cancel
          </Button>
          <Button
            type='button'
            onClick={handleComplete}
            className='px-4 py-2 text-sm font-medium text-white bg-green-600 border border-transparent rounded-md shadow-sm hover:bg-green-700'
          >
            Apply
          </Button>
        </div>
      </div>
    </div>
  );
}
