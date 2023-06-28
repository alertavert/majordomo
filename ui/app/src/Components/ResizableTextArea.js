import React, { useState } from 'react';
import './ResizableTextArea.css';

const maxRows = 30;
const initialRows = 16;

// Resizable textarea component
function ResizableTextArea(props) {
    const [text, setText] = useState(props.initialValue);
    const [rows, setRows] = useState(initialRows);

    const handleChange = (event) => {
        setText(event.target.value);
    event.target.style.height = 'auto';  // Reset height to calculate the actual scrollHeight
    event.target.style.height = `${event.target.scrollHeight}px`;  // Set the height to scrollHeight

    }

    return (
        <textarea
            rows={rows}
            value={text}
            className="form-control resizable-textarea"
            readOnly={true}
            onChange={handleChange}/>
    );
}

export default ResizableTextArea;
