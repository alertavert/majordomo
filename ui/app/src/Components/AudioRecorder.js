import React, {useState} from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {faMicrophoneLines, faMicrophoneLinesSlash} from '@fortawesome/free-solid-svg-icons'

import '../styles/TopSelector.css';

import MicRecorder from "mic-recorder-to-mp3";

const Mp3Recorder = new MicRecorder({bitRate: 128});

function AudioRecorder({handleAudioRecording, handleAudioRecordingError}) {
    const [isRecording, setIsRecording] = useState(false);

    const startRecording = async () => {
            Mp3Recorder.start().then(() => {
                setIsRecording(true);
            }).catch((e) => handleAudioRecordingError(e));
    };

    const stopRecording = async () => {
        setIsRecording(false);
        Mp3Recorder
            .stop()
            .getMp3()
            .then(async ([buffer, blob]) => {
                handleAudioRecording(blob);
            }).catch((e) => handleAudioRecordingError(e));
    };



    return (
        <div className="col-md-2 align-self-center">
        {isRecording ? (
            <button
                className="btn btn-outline-dark btn-sm conversation-btn"
                onClick={stopRecording}><FontAwesomeIcon icon={faMicrophoneLinesSlash}/></button>
        ) : (
            <button
                className="btn btn-outline-dark btn-sm conversation-btn"
                onClick={startRecording}><FontAwesomeIcon icon={faMicrophoneLines}/></button>
        )}
        </div>
    );
}

export default AudioRecorder;
