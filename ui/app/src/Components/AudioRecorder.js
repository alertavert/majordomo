import React from 'react';

function AudioRecorder({isRecording, startRecording, stopRecording}) {
    return (
        <div className="container-fluid">
            <div className="d-flex justify-content-center">
            {isRecording ? (
                <button className="btn btn-primary btn-sm stop-btn"
                onClick={stopRecording}>&nbsp;&nbsp;Stop&nbsp;&nbsp;
                </button>
            ) : (
                <button className="btn btn-primary btn-sm ask-btn"
                onClick={startRecording}>Speak your command
                </button>
            )}
            </div>
        </div>
    );
}

export default AudioRecorder;
