'use client';
import clsx from 'clsx';
import './index.css';
import React from 'react';

const LoadingInline = () => {
  return (
    <div className={clsx('loadership_AFHLY')}>
      {Array.from({ length: 12 }).map((_, index) => (
        <div key={index}></div>
      ))}
    </div>
  );
};

export default LoadingInline;
