'use client';
import { ProductDetailModel } from '@/app/lib/definitions';
import React, { useState } from 'react';

interface TabsSectionProps {
  productDetail: ProductDetailModel;
  details: string[];
}

export const TabsSection: React.FC<TabsSectionProps> = ({ productDetail, details }) => {
  const [activeTab, setActiveTab] = useState('details');

  return (
    <>
      <div className="flex border-b border-gray-200">
        <button 
          onClick={() => setActiveTab('details')}
          className={`py-3 px-6 font-medium transition-colors duration-200 ${
            activeTab === 'details' 
              ? 'border-b-2 border-indigo-500 text-indigo-600' 
              : 'text-gray-500 hover:text-gray-700'
          }`}
        >
          Details
        </button>
        <button 
          onClick={() => setActiveTab('shipping')}
          className={`py-3 px-6 font-medium transition-colors duration-200 ${
            activeTab === 'shipping' 
              ? 'border-b-2 border-indigo-500 text-indigo-600' 
              : 'text-gray-500 hover:text-gray-700'
          }`}
        >
          Shipping
        </button>
        <button 
          onClick={() => setActiveTab('returns')}
          className={`py-3 px-6 font-medium transition-colors duration-200 ${
            activeTab === 'returns' 
              ? 'border-b-2 border-indigo-500 text-indigo-600' 
              : 'text-gray-500 hover:text-gray-700'
          }`}
        >
          Returns
        </button>
      </div>
      
      <div className="py-6">
        {activeTab === 'details' && (
          <div>
            <h3 className='text-base font-medium text-gray-900 mb-4'>Product Features</h3>
            <div className='space-y-4'>
              <ul role='list' className="space-y-2">
                {details.map((item) => (
                  <li
                    className='flex items-start'
                    key={item}
                  >
                    <svg className="h-5 w-5 text-green-500 mr-2 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                    </svg>
                    <span className="text-gray-600">{item}</span>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        )}
        
        {activeTab === 'shipping' && (
          <div>
            <h3 className='text-base font-medium text-gray-900 mb-4'>Shipping Information</h3>
            <div className="space-y-4 text-gray-600">
              <p>We offer the following shipping options:</p>
              <ul className="space-y-2">
                <li className="flex items-start">
                  <svg className="h-5 w-5 text-blue-500 mr-2 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <div>
                    <span className="font-medium">Standard Shipping:</span> 3-5 business days (Free on orders over $50)
                  </div>
                </li>
                <li className="flex items-start">
                  <svg className="h-5 w-5 text-blue-500 mr-2 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <div>
                    <span className="font-medium">Express Shipping:</span> 1-2 business days ($9.99)
                  </div>
                </li>
                <li className="flex items-start">
                  <svg className="h-5 w-5 text-blue-500 mr-2 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <div>
                    <span className="font-medium">Next Day Delivery:</span> Next business day ($19.99)
                  </div>
                </li>
              </ul>
              <p className="mt-4">All orders are processed within 1-2 business days. Orders placed after 2PM will be processed the next business day.</p>
            </div>
          </div>
        )}
        
        {activeTab === 'returns' && (
          <div>
            <h3 className='text-base font-medium text-gray-900 mb-4'>Return Policy</h3>
            <div className="space-y-4 text-gray-600">
              <p>We want you to be completely satisfied with your purchase. If you're not happy with your order, we offer a simple return policy:</p>
              <ul className="space-y-2">
                <li className="flex items-start">
                  <svg className="h-5 w-5 text-amber-500 mr-2 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <div>
                    <span className="font-medium">30-day return period</span> for unused items in original packaging
                  </div>
                </li>
                <li className="flex items-start">
                  <svg className="h-5 w-5 text-amber-500 mr-2 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <div>
                    <span className="font-medium">Free returns</span> on defective or damaged items
                  </div>
                </li>
                <li className="flex items-start">
                  <svg className="h-5 w-5 text-amber-500 mr-2 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <div>
                    <span className="font-medium">Refund or exchange</span> available for all returns
                  </div>
                </li>
              </ul>
              <p className="mt-4">To initiate a return, please contact our customer service team within 30 days of receiving your order.</p>
            </div>
          </div>
        )}
      </div>
    </>
  );
};
