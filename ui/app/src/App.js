import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import './App.css';

import React, {useState, useEffect, useRef} from "react";
import Logo from './Components/Logo';
import ResizableTextArea from './Components/ResizableTextArea';
import Spinner from './Components/Spinner';
import ErrorBox from "./Components/ErrorBox";
import TopSelector from './Components/TopSelector';

// MicRecorder for Audio Recording
import MicRecorder from 'mic-recorder-to-mp3';
import AudioRecorder from './Components/AudioRecorder';


// FIXME: this should not be used, but the UI served directly from the host.
const MajordomoServerUrl = 'http://localhost:5005';
const SpeechApiUrl = MajordomoServerUrl + '/command';
const PromptApiUrl = MajordomoServerUrl + '/prompt';
const ScenariosApiUrl = MajordomoServerUrl + '/scenarios';

const Mp3Recorder = new MicRecorder({bitRate: 128});

function App() {
    const [responseValue, setResponseValue] = useState('Bot says...');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const [isRecording, setIsRecording] = useState(false);
    const [isBlocked, setIsBlocked] = useState(false);
    const [blobURL, setBlobURL] = useState('');

    const [textareaValue, setTextareaValue] = useState('');
    const refTextAreaValue = useRef(textareaValue);
    refTextAreaValue.current = textareaValue;

    // Scenarios to choose from
    const [Scenarios, setScenarios] = useState([]);  // New state to store fetched Scenarios

    // TopSelector state declarations
    const [selectedScenario, setSelectedScenario] = useState('');
    const [selectedConversation, setSelectedConversation] = useState('Ask Majordomo');

    // Fetching scenarios on initial mount
    useEffect(() => {
        fetch(ScenariosApiUrl)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then((data) => {
                setScenarios(data.scenarios);
                setSelectedScenario(Scenarios[0])  // Default scenario set based on fetched data
            })
            .catch(error => setError(`Could not retrieve Scenarios: ${error.message}`));
    }, []);

    const handleScenarioChange = (scenario) => {
        setSelectedScenario(scenario);
    };

    const handleConversationChange = (conversation) => {
        setSelectedConversation(conversation);
    };

    const startRecording = async () => {
        setError(null);
        if (isBlocked) {
            setError('Please give permission for using microphone');
        } else {
            Mp3Recorder.start().then(() => {
                setIsRecording(true);
            }).catch((e) => setError(e.message));
        }
    };

    const stopRecording = async () => {
        setError(null);
        setIsRecording(false);
        Mp3Recorder
            .stop()
            .getMp3()
            .then(async ([buffer, blob]) => {
                const blobURL = URL.createObjectURL(blob)
                setBlobURL(blobURL);
                // Send the blob to the server here or in separate function
                try {
                    let formData = new FormData();
                    formData.append('audio', blob, 'audio.mp3');
                    const response = await fetch(SpeechApiUrl, {
                        method: 'POST',
                        body: formData,
                    });
                    if (response.ok) {
                        const data = await response.json();
                        console.log('Received:', data.message.length, 'characters');
                        setTextareaValue(data.message);
                    } else {
                        const errorData = await response.json(); // Parse the error response as JSON
                        setError('Error (' + response.status + '): ' + errorData.message); // Show the content of the error message
                        // setError('Error (' + response.status + '): ' + response.message);
                    }
                } catch (error) {
                    setError('Cannot POST Audio: ' + error);
                }
            }).catch((e) => setError(e));
    };

    useEffect(() => {
        navigator.mediaDevices.getUserMedia({audio: true})
            .then(function (stream) {
                console.log('Permission Granted');
                setIsBlocked(false);
            })
            .catch(function (err) {
                console.log('Permission Denied');
                setIsBlocked(true);
            });
    }, []);

    const handleFormSubmit = async (content) => {
        console.log('Sending Query to Majordomo (' + content.length + ' chars)')
        setLoading(true);
        setError(null);
        setResponseValue('Bot says...')
        try {
            const response = await fetch(PromptApiUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    prompt: content,
                    scenario: selectedScenario,
                    session: selectedConversation,
                }),
            });
            setLoading(false);
            const data = await response.json();
            if (response.ok) {
                console.log('Received:', data.message.length, 'characters');
                setResponseValue(data.message);
            } else {
                let errMsg = 'Error (' + response.status + '): ' + response.statusText;
                if (data.message) {
                    errMsg = errMsg + ' - ' + data.message;
                }
                setError('Error: ' + errMsg);
            }
        } catch (error) {
            setError("The Bot wasn't quite there yet: " + error.message);
        }
    };

    return (
        <div className="container">
            <Logo/>
            <Header/>
            <TopSelector
                scenarios={Scenarios}
                onScenarioChange={handleScenarioChange}
                onConversationChange={handleConversationChange}
            />
            <AudioRecorder
                startRecording={startRecording}
                stopRecording={stopRecording}
                isRecording={isRecording}
            />
            <PromptBox
                onSubmit={handleFormSubmit}
                textareaValue={refTextAreaValue.current}
                setTextareaValue={setTextareaValue}/>
            <span className="d-block p-2"/>
            <span> {loading ? <Spinner/> : <ResponseBox responseValue={responseValue}/>}</span>
            <div className='error-box'>
                {error !== '' ? <ErrorBox message={error}/> : null}
            </div>
            <Footer/>
        </div>
    );
}

function Header() {
    return (
        <div className="App-header">
            <h1 className="title">Majordomo&nbsp;<span className="font-monospace">Your helpful code bot</span></h1>
        </div>
    );
}

function PromptBox({onSubmit, textareaValue, setTextareaValue}) {
    // const [textareaValue, setTextareaValue] = useState('');

    const handleSubmit = async () => {
        onSubmit(textareaValue);
    };

    const handleTextareaChange = (event) => {
        setTextareaValue(event.target.value);
    };

    return (
        <div className="container-fluid">
            <div className="jumbotron">
                <h6>Your request:</h6>
                <textarea className="form-control"
                          style={{height: '200px'}}
                          value={textareaValue}
                          onChange={handleTextareaChange}
                />
                <button className="btn btn-primary btn-sm ask-btn"
                        onClick={handleSubmit}>Ask Majordomo
                </button>
            </div>
        </div>
    );
}

function ResponseBox({responseValue}) {
    const [rows, setRows] = useState(5); // just an example row size
    const maxRows = 25; // assuming maximum rows limit is 10

    console.log('ResponseBox responseValue:', responseValue);
    return (
        <div className="container-fluid">
            <div className="jumbotron">
                <h6>Response</h6>
                <ResizableTextArea value={responseValue}/>
            </div>
        </div>
    );
}


function Footer() {
    return (
        <div>
            <p className="footer">&copy; 2023 AlertAvert. All rights reserved.</p>
        </div>
    )
}

export default App;
