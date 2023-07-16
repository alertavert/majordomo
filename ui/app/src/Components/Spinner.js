import React from 'react';

import '../styles/Spinner.css';

const Spinner = () => {
    return (
        <div className="spinner">
            <div className="spinner-child"></div>
            <div className="spinner-child"></div>
            <div className="spinner-child"></div>
        </div>
    );
};

export default Spinner;
