import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { OrderStatus } from '@/app/lib/definitions';
import { OrderModel } from '@/app/lib/definitions/order';
import { ArrowRightIcon, ArrowLeftIcon } from '@heroicons/react/16/solid';
import {
  CheckCircleIcon,
  TruckIcon,
  ClockIcon,
  ShoppingBagIcon,
} from '@heroicons/react/24/solid';
import clsx from 'clsx';
import dayjs from 'dayjs';
import { Metadata } from 'next';
import Image from 'next/image';
import Link from 'next/link';
import { cache } from 'react';
import dynamic from 'next/dynamic';
import { serverSideFetch } from '@/app/lib/api/apiServer';

// Import the client component with dynamic to avoid SSR issues
const PaymentInfoSection = dynamic(
  () => import('@/app/components/Payment/PaymentInfoSection'),
  { ssr: true }
);

const ConfirmOrderButton = dynamic(
  () => import('./_components/ConfirmOrderButton'),
  { ssr: true }
);

const OrderItemRating = dynamic(() => import('./_components/OrderItemRating'), {
  ssr: true,
});

const CancelOrderButton = dynamic(
  () => import('./_components/CancelOrderButton'),
  {
    ssr: true,
  }
);

export const getOrderDetails = cache(async (id: string) => {
  const order = await serverSideFetch<OrderModel>(
    PUBLIC_API_PATHS.ORDER.replace(':id', id),
    {
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
  params: Promise<{ id: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { id } = await params;

  const { data, error } = await getOrderDetails(id);
  if (error && !data) {
    throw new Error(error.details);
  }

  return {
    title: 'Order Details ' + data.id,
    description: 'Order detail for order with ' + data.status,
  };
}

// Define the main order progression steps
const orderProgressSteps = [
  { status: OrderStatus.Pending, label: 'Order placed', icon: ShoppingBagIcon },
  { status: OrderStatus.Confirmed, label: 'Processing', icon: ClockIcon },
  { status: OrderStatus.Delivering, label: 'Shipped', icon: TruckIcon },
  { status: OrderStatus.Delivered, label: 'Delivered', icon: CheckCircleIcon },
];

export default async function Page({ params }: Props) {
  const { id } = await params;
  const { data: orderDetail } = await getOrderDetails(id);
  // Find the current step index based on status
  console.log({ orderDetail });
  const getCurrentStepIndex = (status: OrderStatus): number => {
    // Special handling for Completed, Cancelled, and Refunded
    if (status === OrderStatus.Completed) return orderProgressSteps.length - 1;
    if (status === OrderStatus.Cancelled || status === OrderStatus.Refunded)
      return -1;
    const index = orderProgressSteps.findIndex(
      (step) => step.status === status
    );
    return index >= 0 ? index : 0;
  };

  const currentStepIndex = getCurrentStepIndex(orderDetail.status);

  const isSpecialStatus =
    orderDetail.status === OrderStatus.Cancelled ||
    orderDetail.status === OrderStatus.Refunded;

  // Check if order can be cancelled (Pending status and not paid)
  const canCancel =
    orderDetail.status === OrderStatus.Pending &&
    (!orderDetail.paymentInfo || orderDetail.paymentInfo.status !== 'paid');

  return (
    <div className='h-full container mx-auto my-20'>
      <div className='mb-6'>
        <Link
          href='/orders'
          className='inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-indigo-600 bg-indigo-50 rounded-md hover:bg-indigo-100 transition-colors'
        >
          <ArrowLeftIcon className='w-4 h-4' />
          Back to Orders
        </Link>
      </div>
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
            {dayjs(orderDetail.createdAt).format('MMM DD, YYYY')}
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
                    src={e.imageUrl}
                    alt={e.name}
                    className='object-cover border border-lime-400 rounded-md'
                    width={90}
                    height={90}
                    priority
                  />
                  <div className='flex flex-col gap-0.5'>
                    <div className='text-lg font-bold'>{e.name}</div>
                    <div className='text-base text-gray-600'>
                      ${e.lineTotal}
                    </div>
                    <div className='flex gap-3'>
                      {e?.attributeSnapshot?.map((e) => (
                        <div className='text-sm text-indigo-400' key={e.name}>
                          <span className='font-medium'>{e.name}: </span>
                          <span>{e.value}</span>
                        </div>
                      ))}
                    </div>

                    {/* Display rating component for completed orders */}
                    {orderDetail.status === OrderStatus.Completed && (
                      <div className='mt-2'>
                        <OrderItemRating
                          orderItemId={e.id}
                          ratingModel={e.rating}
                        />
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
            <div className='w-1/2 flex gap-6'>
              <div className='flex-1'>
                <div className='font-medium mb-2'>Delivery address</div>
                <div className='text-gray-500 text-sm'>
                  <div>{orderDetail.shippingInfo.name}</div>
                  <div>{orderDetail.shippingInfo.address}</div>
                  <div>Ward {orderDetail.shippingInfo.ward}</div>
                  <div>District {orderDetail.shippingInfo.district}</div>
                  <div>City {orderDetail.shippingInfo.city}</div>
                </div>
              </div>
              <div className='flex-1'>
                <div className='font-medium mb-2'>Shipping updates</div>
                <div className='text-gray-500 text-sm'>
                  <div>{orderDetail.shippingInfo.phone}</div>
                  <div>{orderDetail.customerEmail}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <hr />
        <div className='px-16'>
          <div className='font-semibold mb-3'>
            {isSpecialStatus ? (
              <span
                className={clsx(
                  'px-3 py-1 text-sm font-medium rounded-full',
                  orderDetail.status === OrderStatus.Cancelled
                    ? 'bg-red-100 text-red-800'
                    : 'bg-orange-100 text-orange-800'
                )}
              >
                Order {orderDetail.status}
              </span>
            ) : (
              <div className='flex justify-between items-center'>
                <span>
                  Last update:{' '}
                  {dayjs(orderDetail.createdAt).format('MMMM DD, YYYY')}
                </span>

                {/* Add cancel button for orders that can be cancelled */}
                {canCancel && <CancelOrderButton orderId={orderDetail.id} />}
              </div>
            )}
          </div>

          {/* Progress tracking section */}
          {!isSpecialStatus && (
            <div className='mb-2'>
              <div className='relative mb-8'>
                {/* Progress bar */}
                <div className='w-full bg-gray-200 h-1 absolute top-4'>
                  <div
                    className='bg-indigo-500 h-1 transition-all duration-500'
                    style={{
                      width: `${currentStepIndex >= 0 ? (100 * currentStepIndex) / (orderProgressSteps.length - 1) : 0}%`,
                    }}
                  />
                </div>

                {/* Step indicators */}
                <div className='flex justify-between relative z-10'>
                  {orderProgressSteps.map((step, index) => {
                    const isActive = currentStepIndex >= index;
                    const isCurrent = currentStepIndex === index;

                    return (
                      <div
                        key={step.status}
                        className='flex flex-col items-center'
                      >
                        <div
                          className={clsx(
                            'rounded-full p-2 w-8 h-8 flex items-center justify-center',
                            isActive
                              ? 'bg-indigo-500 text-white'
                              : 'bg-gray-200 text-gray-500'
                          )}
                        >
                          <step.icon className='w-4 h-4' />
                        </div>
                        <div
                          className={clsx(
                            'mt-2 text-sm font-medium text-center',
                            isCurrent
                              ? 'text-indigo-600'
                              : isActive
                                ? 'text-indigo-500'
                                : 'text-gray-500'
                          )}
                        >
                          {step.label}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            </div>
          )}

          {/* Confirm received button for delivered orders */}
          {orderDetail.status === OrderStatus.Delivered && (
            <div className='mt-6 flex justify-center'>
              <ConfirmOrderButton orderId={orderDetail.id} />
            </div>
          )}
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
              paymentInfo={orderDetail.paymentInfo || null}
              orderId={orderDetail.id}
              total={orderDetail.total}
              orderStatus={orderDetail.status}
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
