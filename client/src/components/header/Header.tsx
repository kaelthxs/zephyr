import React from 'react';
import { Button } from '../button/Button';
import styles from './Header.module.css';

const Header: React.FC = () => {
  return (
    <header className={styles.header}>
      <div className={styles.header__container}>
        <div className={styles.header__leftside}>
          <Button variant="header__button">
            build
          </Button>
          <Button variant="header__button">
            stays
          </Button>
          <Button variant="header__button">
            rent a car
          </Button>
          <Button variant="header__button">
            transfer
          </Button>
        </div>
        <div className={styles.header__rightside}>
          <Button variant="header__button">
            sign in
          </Button>
        </div>
      </div>
    </header>
  );
};

export default Header;

