import React from 'react';
import dayjs from 'dayjs';

type UsageHistory = {
  id: string;
  orderId: string;
  customerName: string;
  amount: number;
  discountAmount: number;
  date: string;
};

type UsageOverview = {
  usedCount: number;
  usageLimit?: number;
};

type UsageTableProps = {
  usageHistory?: UsageHistory[];
  overview: UsageOverview;
};

export const UsageTable: React.FC<UsageTableProps> = ({
  usageHistory,
  overview,
}) => {
  return (
    <div>
      <div className='mb-6'>
        <h3 className='text-lg font-medium mb-3'>Usage Overview</h3>
        <div className='grid grid-cols-1 md:grid-cols-3 gap-4'>
          <div className='p-4 border rounded-md'>
            <div className='text-sm font-medium text-gray-500'>
              Total Uses
            </div>
            <div className='mt-1 text-3xl font-semibold'>
              {overview.usedCount}
            </div>
          </div>
          <div className='p-4 border rounded-md'>
            <div className='text-sm font-medium text-gray-500'>
              Remaining
            </div>
            <div className='mt-1 text-3xl font-semibold'>
              {overview.usageLimit
                ? overview.usageLimit - overview.usedCount
                : '∞'}
            </div>
          </div>
          <div className='p-4 border rounded-md'>
            <div className='text-sm font-medium text-gray-500'>
              Usage Limit
            </div>
            <div className='mt-1 text-3xl font-semibold'>
              {overview.usageLimit || '∞'}
            </div>
          </div>
        </div>
      </div>

      <div>
        <h3 className='text-lg font-medium mb-3'>Recent Usage</h3>

        {usageHistory && usageHistory.length > 0 ? (
          <div className='overflow-x-auto'>
            <table className='min-w-full divide-y divide-gray-200'>
              <thead className='bg-gray-50'>
                <tr>
                  <th
                    scope='col'
                    className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
                  >
                    Order ID
                  </th>
                  <th
                    scope='col'
                    className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
                  >
                    Customer
                  </th>
                  <th
                    scope='col'
                    className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'
                  >
                    Date
                  </th>
                  <th
                    scope='col'
                    className='px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider'
                  >
                    Order Total
                  </th>
                  <th
                    scope='col'
                    className='px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider'
                  >
                    Discount Applied
                  </th>
                </tr>
              </thead>
              <tbody className='bg-white divide-y divide-gray-200'>
                {usageHistory.map((usage) => (
                  <tr key={usage.id}>
                    <td className='px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900'>
                      {usage.orderId}
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                      {usage.customerName}
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                      {dayjs(usage.date).format('MMM D, YYYY h:mm A')}
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-right'>
                      ${usage.amount.toFixed(2)}
                    </td>
                    <td className='px-6 py-4 whitespace-nowrap text-sm text-red-600 font-medium text-right'>
                      -${usage.discountAmount.toFixed(2)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className='bg-gray-50 rounded-md p-8 text-center'>
            <p className='text-gray-500'>
              This discount hasn't been used yet.
            </p>
          </div>
        )}
      </div>
    </div>
  );
};
