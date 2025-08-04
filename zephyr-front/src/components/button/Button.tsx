import React from 'react';
import clsx from 'clsx';
import styles from './Button.module.css';

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: 'primary' | 'secondary' | 'header__button';
};

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  children,
  ...rest
}) => {
  return (
    <button
      className={clsx(
        styles.button,
        {
          [styles.header__button]: variant === 'header__button',
        }
      )}
      {...rest}
    >
      {children}
    </button>
  );
};
