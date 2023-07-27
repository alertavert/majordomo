import React from 'react';

function AudioRecorder({isRecording, startRecording, stopRecording}) {
    return (
        <div className="container-fluid">
            <div>
            {isRecording ? (
                <button className="btn btn-primary btn-sm stop-btn"
                onClick={stopRecording}>Stop
                </button>
            ) : (
                <button className="btn btn-primary btn-sm record-btn"
                onClick={startRecording}>Speak your command
                </button>
            )}
            </div>
        </div>
    );
}

export default AudioRecorder;
