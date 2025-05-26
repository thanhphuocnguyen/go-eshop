'use client';

import React, { useState } from 'react';
import ImageUploader from '@/app/components/ImageUploader';

export default function ImageUploaderExamplePage() {
  const [imageUrl, setImageUrl] = useState<string | null>(null);

  const handleImageChange = (file: File | null) => {
    // Show in console for demonstration
    console.log('Selected File:', file);
    
    // In a real app, you might upload to a server here
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setImageUrl(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-2xl font-bold mb-6">Image Editor Example</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div>
          <h2 className="text-lg font-medium mb-4">Square Image (1:1)</h2>
          <ImageUploader 
            name="square-image"
            label="Upload square image"
            onChange={handleImageChange}
            width={300}
            height={300}
            aspectRatio={1}
            maxFileSizeMB={5}
          />
        </div>
        
        <div>
          <h2 className="text-lg font-medium mb-4">Wide Image (16:9)</h2>
          <ImageUploader 
            name="wide-image"
            label="Upload wide image" 
            onChange={handleImageChange}
            width={300}
            height={169}
            aspectRatio={16/9}
            maxFileSizeMB={5}
          />
        </div>
      </div>
      
      {imageUrl && (
        <div className="mt-8">
          <h2 className="text-lg font-medium mb-4">Preview of Last Edited Image</h2>
          <img 
            src={imageUrl} 
            alt="Uploaded preview" 
            className="max-w-md border rounded shadow-sm"
          />
        </div>
      )}

      <div className="mt-12 p-4 bg-gray-50 rounded-lg">
        <h2 className="text-lg font-medium mb-2">Usage Instructions</h2>
        <ol className="list-decimal pl-5 space-y-2">
          <li>Click on any of the upload boxes above</li>
          <li>Select an image from your device</li>
          <li>Use the cropping tool to adjust the image</li>
          <li>Adjust zoom using the slider</li>
          <li>Click "Apply" to confirm your edit</li>
          <li>The edited image will appear in the preview section below</li>
        </ol>
      </div>
    </div>
  );
}
