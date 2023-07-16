import React, {useEffect, useState} from 'react';
import '../styles/ResizableTextArea.css';


// Resizable textarea component
function ResizableTextArea({value, maxRows = 20}) {
    const [rows, setRows] = useState(6);

    // this updates the row size of the textarea to match the number of lines
    // of text in the response
    useEffect(() => {
        const lineCount = value.split('\n').length;
        if (lineCount <= maxRows) {
            setRows(lineCount);
        } else {
            setRows(maxRows);
        }
    }, [value, maxRows]);

    return (
        <textarea
            rows={rows}
            value={value}
            className="form-control resizable-textarea"
            readOnly={true}
        />
    );
}

export default ResizableTextArea;
