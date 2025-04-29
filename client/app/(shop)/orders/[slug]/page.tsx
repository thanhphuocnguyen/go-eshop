import { apiFetch } from '@/lib/apis/api';
import { API_PATHS } from '@/lib/constants/api';
import { GenericResponse } from '@/lib/definitions';
import { OrderModel } from '@/lib/definitions/order';
import { ArrowRightIcon } from '@heroicons/react/16/solid';
import clsx from 'clsx';
import dayjs from 'dayjs';
import { Metadata } from 'next';
import { cookies } from 'next/headers';
import Image from 'next/image';
import Link from 'next/link';
import { cache } from 'react';

export const getOrderDetails = cache(async (slug: string) => {
  const cookieStorage = await cookies();
  const order = await apiFetch<GenericResponse<OrderModel>>(
    API_PATHS.ORDER_ITEM.replace(':id', slug),
    {
      authToken: cookieStorage.get('access_token')?.value,
      nextOptions: {
        next: {
          tags: ['order'],
        },
      },
    }
  );
  if (order.error && !order.data) {
    console.error(order.error);
    throw new Error(order.error.details);
  }
  return order;
});

type Props = {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;

  const { data, error } = await getOrderDetails(slug);
  if (error && !data) {
    throw new Error(error.stack);
  }

  return {
    title: 'Order Details ' + data.id,
    description: 'Order detail for order with ' + data.status,
  };
}

export default async function Page({ params }: Props) {
  const { slug } = await params;
  const { data: orderDetail } = await getOrderDetails(slug);
  const orderStatuses = ['Order placed', 'Processing', 'Shipped', 'Delivered'];
  if (!orderDetail.payment_info) {
    orderStatuses.splice(1, 0, 'Pending payment');
  }
  const getPercentageByStatus = (status: string) => {
    const index = orderStatuses.indexOf(status);
    return (index / orderStatuses.length) * 100;
  };
  return (
    <div className='h-full container mx-auto py-10 px-8'>
      <div className='flex justify-between items-end'>
        <div className='flex items-end'>
          <h1 className='text-3xl font-bold'>Order #{orderDetail.id}</h1>
          <Link
            className='flex gap-2 font-medium ml-2 text-indigo-500 items-center'
            href={`/order/${orderDetail.id}/invoice`}
          >
            <span>View invoice</span>
            <ArrowRightIcon className='size-5' />
          </Link>
        </div>
        <div className='flex gap-1'>
          <span>Order placed</span>
          <span className='font-semibold'>
            {dayjs(orderDetail.created_at).format('MMM DD, YYYY')}
          </span>
        </div>
      </div>

      <div className='mt-8 rounded-md border border-gray-300 shadow-md  flex flex-col gap-4'>
        <div className='px-12 py-8'>
          <div className='flex gap-10'>
            <div className='w-1/2 flex flex-col gap-4'>
              {orderDetail.products.map((e) => (
                <div
                  key={e.id}
                  className='flex gap-4 border-b pb-4 border-gray-200'
                >
                  <Image
                    src={e.image_url}
                    alt={e.name}
                    className='object-cover border border-lime-400 rounded-md'
                    width={90}
                    height={90}
                    priority
                  />
                  <div className='flex flex-col gap-0.5'>
                    <div className='text-lg font-bold'>{e.name}</div>
                    <div className='text-base text-gray-600'>
                      ${e.line_total}
                    </div>
                    <div className='flex gap-3'>
                      {e.attribute_snapshot.map((e) => (
                        <div className='text-sm text-indigo-400' key={e.name}>
                          <span className='font-medium'>{e.name}: </span>
                          <span>{e.value}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              ))}
            </div>
            <div className='w-1/2 flex gap-6'>
              <div className='flex-1'>
                <div className='font-medium mb-2'>Delivery address</div>
                <div className='text-gray-500 text-sm'>
                  <div>{orderDetail.shipping_info.name}</div>
                  <div>{orderDetail.shipping_info.address}</div>
                  <div>Ward {orderDetail.shipping_info.ward}</div>
                  <div>District {orderDetail.shipping_info.district}</div>
                  <div>City {orderDetail.shipping_info.city}</div>
                </div>
              </div>
              <div className='flex-1'>
                <div className='font-medium mb-2'>Shipping updates</div>
                <div className='text-gray-500 text-sm'>
                  <div>{orderDetail.shipping_info.phone}</div>
                  <div>{orderDetail.customer_email}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <hr />
        <div className='px-6 pb-8'>
          <div className='font-semibold mb-3'>Shipped on March 23, 2021</div>
          <div className='w-full bg-gray-200 rounded-full h-4 mb-4 dark:bg-gray-700'>
            <div
              className='bg-sky-500 h-4 rounded-full dark:bg-blue-500'
              style={{ width: `${orderStatuses.length / 100}%` }}
            />
          </div>
          <div className='w-full flex justify-between'>
            {orderStatuses.map((e, i) => (
              <div
                key={i}
                className={clsx(
                  'text-sm font-medium ',
                  `w-1/${orderStatuses.length} text-center`
                )}
              >
                {e}
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className='bg-gray-100 flex justify-between gap-12 w-full px-10 py-8 mt-10 shadow-md rounded-md'>
        <div className='w-1/2 flex px-8'>
          <div className='w-1/2'>
            <div className='font-semibold'>Billing address</div>
            <div className='mt-2'>In progress</div>
          </div>
          <div className='w-1/2'>
            <div className='font-semibold'>Payment information</div>
            <div className='mt-2'>Stripe</div>
          </div>
        </div>
        <div className='w-1/2 flex flex-col px-8 gap-4'>
          <div className='flex justify-between'>
            <div className='text-gray-500'>Subtotal</div>
            <div className='font-bold text-indigo-500'>
              ${orderDetail.total}
            </div>
          </div>
          <hr />
          <div className='flex justify-between'>
            <div className='text-gray-500'>Shipping</div>
            <div className='font-bold text-indigo-500'>$0.00</div>
          </div>
          <hr />
          <div className='flex justify-between'>
            <div className='text-gray-500'>Tax</div>
            <div className='font-bold text-indigo-500'>$0.00</div>
          </div>
          <hr />
          <div className='flex justify-between'>
            <div className='text-gray-500'>Total</div>
            <div className='font-bold text-indigo-500'>
              ${orderDetail.total}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
