import {
  Button,
  Dialog,
  DialogPanel,
  DialogTitle,
  Transition,
  TransitionChild,
} from '@headlessui/react';
import clsx from 'clsx';
import React, { Fragment } from 'react';

interface ConfirmDialogProps {
  title: string;
  message: string;
  open: boolean;
  confirmStyle?: string;
  onClose: () => void;
  onConfirm: () => void;
}
export const ConfirmDialog: React.FC<ConfirmDialogProps> = ({
  message,
  onClose,
  confirmStyle,
  onConfirm,
  open,
  title,
}) => {
  return (
    <Transition appear show={open} as={Fragment}>
      <Dialog
        as='div'
        className='relative z-50 focus:outline-none'
        onClose={onClose}
      >
        <TransitionChild
          as={Fragment}
          enter='ease-out duration-300'
          enterFrom='opacity-0'
          enterTo='opacity-100'
          leave='ease-in duration-200'
          leaveFrom='opacity-100'
          leaveTo='opacity-0'
        >
          <div className='fixed inset-0 bg-black/25 backdrop-blur-sm' />
        </TransitionChild>

        <div className='fixed inset-0 z-10 w-screen overflow-y-auto'>
          <div className='flex min-h-full items-center justify-center p-4'>
            <TransitionChild
              as={Fragment}
              enter='ease-out duration-300'
              enterFrom='opacity-0 scale-95'
              enterTo='opacity-100 scale-100'
              leave='ease-in duration-200'
              leaveFrom='opacity-100 scale-100'
              leaveTo='opacity-0 scale-95'
            >
              <DialogPanel className='w-full max-w-md rounded-xl bg-white border border-green-400 shadow-lg shadow-green-100/50 p-6 transform transition-all'>
                <div className='absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1/2'>
                  <div className='inline-flex h-14 w-14 items-center justify-center rounded-full bg-white border-4 border-green-100 text-red-500'>
                    <svg
                      xmlns='http://www.w3.org/2000/svg'
                      className='h-8 w-8'
                      fill='none'
                      viewBox='0 0 24 24'
                      stroke='currentColor'
                    >
                      <path
                        strokeLinecap='round'
                        strokeLinejoin='round'
                        strokeWidth={2}
                        d='M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z'
                      />
                    </svg>
                  </div>
                </div>

                <DialogTitle
                  as='h3'
                  className='mt-5 text-xl font-semibold text-gray-900 text-center'
                >
                  {title}
                </DialogTitle>
                <p className='mt-3 text-base/6 text-gray-600 text-center'>
                  {message}
                </p>
                <div className='mt-8 flex justify-center space-x-4'>
                  <Button
                    className='px-4 py-2 rounded-lg text-gray-700 bg-gray-100 hover:bg-gray-200 transition-colors duration-200 font-medium focus:outline-none focus:ring-2 focus:ring-gray-300'
                    onClick={onClose}
                  >
                    Cancel
                  </Button>
                  <Button
                    onClick={onConfirm}
                    className={clsx(
                      'px-4 py-2 rounded-lg bg-red-600 text-white hover:bg-red-700 transition-colors duration-200 font-medium focus:outline-none focus:ring-2 focus:ring-red-500',
                      confirmStyle
                    )}
                  >
                    Confirm
                  </Button>
                </div>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
};
