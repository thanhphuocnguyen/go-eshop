import React from 'react';

interface ButtonProps {
  variant:
    | 'primary'
    | 'secondary'
    | 'tertiary'
    | 'danger'
    | 'warning'
    | 'success'
    | 'info';
  children: React.ReactNode;
}

const Button: React.FC<ButtonProps> = ({ variant, children }) => {
  const getColor = () => {
    switch (variant) {
      case 'primary':
        return 'blue-500';
      case 'secondary':
        return 'gray-500';
      case 'tertiary':
        return 'white';
      case 'danger':
        return 'red-500';
      case 'warning':
        return 'yellow-500';
      case 'success':
        return 'green-500';
      case 'info':
        return 'blue-500';
    }
  };

  return <button className={`btn btn-${getColor()}`}>{children}</button>;
};

export default Button;
