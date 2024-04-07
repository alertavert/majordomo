import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCopy, faEraser, faPaste } from '@fortawesome/free-solid-svg-icons';
import '../styles/PromptBox.css';

/** Represents a prompt box with a clear functionality */
function PromptBox({onSubmit, textareaValue, setTextareaValue}) {

    /** Handles clearing the textarea */
    const clearTextarea = () => {
        setTextareaValue('');
    };

    /** Handles form submit */
    const handleSubmit = async () => {
        onSubmit(textareaValue);
    };

    /** Handles textarea value change */
    const handleTextareaChange = (event) => {
        setTextareaValue(event.target.value);
    };

    const copyTextarea = () => {
        navigator.clipboard.writeText(textareaValue);
    };

    const pasteToTextarea = () => {
        navigator.clipboard.readText()
            .then(text => {
                setTextareaValue(text);
            })
            .catch(err => {
                setTextareaValue(`Failed to read clipboard contents: ${err}`);
            });
    };

    return (
        <div className="container-fluid">
            <div className="jumbotron prompt-box-container">
                <div className="prompt-box-header">
                    <h6>Your request:</h6>
                </div>
                <textarea className="form-control prompt-box"
                          value={textareaValue}
                          onChange={handleTextareaChange}
                />
                <button className="btn textarea-btn btn-copy" onClick={copyTextarea}>
                    <FontAwesomeIcon icon={faCopy} />
                </button>
                <button className="btn textarea-btn btn-paste" onClick={pasteToTextarea}>
                    <FontAwesomeIcon icon={faPaste} />
                </button>
                <button className="btn textarea-btn btn-clear" onClick={clearTextarea}>
                    <FontAwesomeIcon icon={faEraser} />
                </button>
                <div className="d-flex justify-content-end pt-2">
                    <button className="btn btn-primary btn-sm ask-btn"
                            onClick={handleSubmit}>Ask Majordomo
                    </button>
                </div>
            </div>
        </div>
    );
}

export default PromptBox;
