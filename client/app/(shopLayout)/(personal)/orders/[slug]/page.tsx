import { apiFetch } from '@/lib/apis/api';
import { PUBLIC_API_PATHS } from '@/lib/constants/api';
import { GenericResponse, OrderStatus } from '@/lib/definitions';
import { OrderModel } from '@/lib/definitions/order';
import { ArrowRightIcon } from '@heroicons/react/16/solid';
import clsx from 'clsx';
import dayjs from 'dayjs';
import { Metadata } from 'next';
import { cookies } from 'next/headers';
import Image from 'next/image';
import Link from 'next/link';
import { cache } from 'react';
import dynamic from 'next/dynamic';

// Import the client component with dynamic to avoid SSR issues
const PaymentInfoSection = dynamic(
  () => import('@/components/Payment/PaymentInfoSection'),
  { ssr: true }
);

export const getOrderDetails = cache(async (slug: string) => {
  const cookieStorage = await cookies();
  const order = await apiFetch<GenericResponse<OrderModel>>(
    PUBLIC_API_PATHS.ORDER_ITEM.replace(':id', slug),
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

const orderStatuses = [
  [OrderStatus.Pending],
  [OrderStatus.Confirm],
  [OrderStatus.Delivering],
  [OrderStatus.Delivered],
];
export default async function Page({ params }: Props) {
  const { slug } = await params;
  const { data: orderDetail } = await getOrderDetails(slug);

  const getPercentageByStatus = (status: OrderStatus) => {
    const index = orderStatuses.findIndex((e) => e.includes(status));
    if (index === -1) {
      return 0;
    }
    const ratio = index / orderStatuses.length;
    return (ratio < 1 ? ratio + 0.1 : ratio) * 100;
  };

  return (
    <div className='h-full container mx-auto my-20'>
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

      <div className='mt-8 py-12 rounded-md border border-gray-200 shadow-md flex flex-col gap-4'>
        <div className='px-16'>
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
        <div className='px-16'>
          <div className='font-semibold mb-3'>Shipped on March 23, 2021</div>
          <div className='w-full bg-gray-200 rounded-full h-3 mb-4 dark:bg-gray-700'>
            <div
              className='bg-sky-500 h-3 rounded-full dark:bg-blue-500'
              style={{ width: `${getPercentageByStatus(orderDetail.status)}%` }}
            />
          </div>
          <div className='w-full flex justify-between'>
            <div
              className={clsx(
                'text-sm font-medium',
                orderDetail.status === OrderStatus.Pending
                  ? 'text-indigo-500'
                  : ''
              )}
            >
              Order placed
            </div>
            <div
              className={clsx(
                'text-sm font-medium',
                orderDetail.status === OrderStatus.Confirm
                  ? 'text-indigo-500'
                  : ''
              )}
            >
              Processing
            </div>
            <div
              className={clsx(
                'text-sm font-medium',
                orderDetail.status === OrderStatus.Delivering
                  ? 'text-indigo-500'
                  : ''
              )}
            >
              Shipped
            </div>
            <div
              className={clsx(
                'text-sm font-medium',
                orderDetail.status === OrderStatus.Delivered
                  ? 'text-indigo-500'
                  : ''
              )}
            >
              Delivered
            </div>
          </div>
        </div>
      </div>

      <div className='bg-gray-100 flex border border-gray-100 justify-between gap-12 w-full px-10 py-8 mt-12 shadow-md rounded-md'>
        <div className='w-1/2 flex px-8'>
          <div className='w-1/2'>
            <div className='font-semibold'>Billing address</div>
            <div className='mt-2'>In progress</div>
          </div>
          <div className='w-1/2'>
            <PaymentInfoSection
              paymentInfo={orderDetail.payment_info || null}
              orderId={orderDetail.id}
              total={orderDetail.total}
            />
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
