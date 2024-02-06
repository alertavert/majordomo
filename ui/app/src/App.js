import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import './App.css';

import React, {useState, useEffect, useRef} from "react";
import Logo from './Components/Logo';
import ResizableTextArea from './Components/ResizableTextArea';
import Spinner from './Components/Spinner';
import ErrorBox from "./Components/ErrorBox";
import TopSelector from './Components/TopSelector';
import PromptBox from './Components/PromptBox'; // Imported PromptBox component

import { Logger } from "./Services/logger";
import {fetchProjects, fetchScenarios} from './Services/api'; // Import fetchScenarios function

// FIXME: this should not be used, but the UI served directly from the host.
const MajordomoServerUrl = 'http://localhost:5005';
const SpeechApiUrl = MajordomoServerUrl + '/command';
const PromptApiUrl = MajordomoServerUrl + '/prompt';


function App() {
    const [responseValue, setResponseValue] = useState('Bot says...');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [textareaValue, setTextareaValue] = useState('');
    const refTextAreaValue = useRef(textareaValue);
    refTextAreaValue.current = textareaValue;

    // Scenarios to choose from
    const [Scenarios, setScenarios] = useState([]);  // New state to store fetched Scenarios

    // TopSelector state declarations
    const [selectedScenario, setSelectedScenario] = useState(null);
    const [selectedConversation, setSelectedConversation] = useState('Ask Majordomo');
    const [activeProject, setActiveProject] = useState('');

    // Fetching scenarios on initial mount
    useEffect(() => {
        fetchScenarios(setScenarios, setSelectedScenario, setError);
        fetchProjects(setActiveProject, setError);

    }, []);

    const handleScenarioChange = (scenario) => {
        setSelectedScenario(scenario);
    };

    const handleConversationChange = (conversation) => {
        setSelectedConversation(conversation);
    };

    const handleAudioRecording = async (blob) => {
        setError(null);
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
                setError('Error (' + response.status + '): ' + errorData.message);
            }
        } catch (error) {
            setError('Cannot convert Audio: ' + error);
        }
    };

    const handleAudioRecordingError = (error) => {
        setError('Audio Recording Error: ' + error);
    }

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
                activeProject={activeProject}
                onScenarioChange={handleScenarioChange}
                onConversationChange={handleConversationChange}
                handleAudioRecording={handleAudioRecording}
                handleAudioRecordingError={handleAudioRecordingError}
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

function ResponseBox({responseValue}) {
    Logger.debug('ResponseBox responseValue:', responseValue);
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
            <p className="footer">&copy; 2023-2024 AlertAvert. All rights reserved.</p>
        </div>
    )
}

export default App;
