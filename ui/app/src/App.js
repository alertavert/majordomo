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
import {
    fetchProjects,
    fetchSessionsForProjects,
    sendAudioBlob,
    sendPrompt,
} from './Services/api'; // Import fetchScenarios function


function App() {
    const [responseValue, setResponseValue] = useState('Bot says...');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [textareaValue, setTextareaValue] = useState('');
    const refTextAreaValue = useRef(textareaValue);
    refTextAreaValue.current = textareaValue;

    // Scenarios to choose from
    const [Projects, setProjects] = useState([]);  // New state to store fetched Scenarios
    const [Sessions, setSessions] = useState([]);  // New state to store fetched Scenarios

    // TopSelector state declarations
    const [selectedScenario, setSelectedScenario] = useState(null);
    const [selectedConversation, setSelectedConversation] = useState('Ask Majordomo');
    const [activeProject, setActiveProject] = useState('');

    // Fetching scenarios on initial mount
    useEffect(() => {
        fetchProjects(setActiveProject, setProjects, setError);
        fetchSessionsForProjects(activeProject, setSessions, setError);
        // TODO: the scenario should be retrieved from the session
        setSelectedScenario('web_developer')
    }, [activeProject]);
    Logger.debug(`Projects: ${Projects}, Active: ${activeProject}, Sessions: ${Sessions}`)

    const handleProjectChange = (project) => {
        setActiveProject(project);
        // TODO: Fetch conversations for the selected project
        console.log('Fetching conversations for project:', project, Sessions)
    }

    const handleConversationChange = (conversation) => {
        setSelectedConversation(conversation);
    };

    const handleAudioRecording = async (blob) => {
        Logger.debug('Sending Audio Blob to Majordomo')
        await sendAudioBlob(blob, setTextareaValue, setError);
    }

    const handleFormSubmit = async (content) => {
        Logger.debug('Sending Query to Majordomo (' + content.length + ' chars)')
        setLoading(true);
        setError(null);
        setResponseValue('')
        try {
            const response = await sendPrompt(content, selectedScenario, selectedConversation);
            Logger.debug(`Bot response is ${response.length} characters`);
            setResponseValue(response);
        } catch (error) {
            setError("The Bot wasn't quite there yet: " + error.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="container">
            <Logo/>
            <Header/>
            <TopSelector
                projects={Projects}
                activeProject={activeProject}
                sessions={Sessions}
                onProjectChange={handleProjectChange}
                onConversationChange={handleConversationChange}
                onAudioRecording={handleAudioRecording}
                onError={setError}
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
