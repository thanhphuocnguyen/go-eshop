'use client';

import { useAppUser } from '@/lib/contexts/AppUserContext';

export default function AddressesPage() {
  const { user } = useAppUser();
  
  return (
    <div className="space-y-6">
      <h2 className="text-xl font-semibold text-gray-900">Shipping Addresses</h2>
      
      {/* Address list will go here */}
      {user?.addresses && user.addresses.length > 0 ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {user.addresses.map((address, index) => (
            <div key={index} className="bg-white border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow">
              {/* <p className="font-medium">{address.name}</p> */}
              <p className="text-sm text-gray-500">{address.street}</p>
              <p className="text-sm text-gray-500">{address.city}, {address.district} {address.ward}</p>
              <p className="text-sm text-gray-500">{address.phone}</p>
              <div className="mt-4 flex space-x-2">
                <button className="text-sm text-indigo-600 hover:text-indigo-800">Edit</button>
                <button className="text-sm text-red-600 hover:text-red-800">Delete</button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <p className="text-gray-500">No addresses found. Add a shipping address to get started.</p>
      )}
      
      <button className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 transition-colors">
        Add New Address
      </button>
    </div>
  );
}