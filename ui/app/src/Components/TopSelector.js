import React, {useState, useRef, useEffect} from 'react';

import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faPen, faTrash} from '@fortawesome/free-solid-svg-icons';

import '../styles/TopSelector.css';
import AudioRecorder from "./AudioRecorder";
import {Logger} from "../Services/logger";


/**
 * TopSelector component.
 *
 * Tracks essentially the state of the application, by showing the currently active project,
 * the available conversations (for the project) and the scenario associated with the currently
 * selected conversation.
 *
 * The user can change the active project, and select to continue an existing conversation
 * (in which case the scenario cannot be changed), or start a new one (in which case they will
 * also be allowed to choose a Scenario for the new conversation).
 *
 * @param {string[]} projects - the currently available projects. (we currently do not support creating new projects)
 * @param {string} activeProject - the currently active project.
 * @param {Session[]} sessions - the currently available conversations for the active project.
 * @param {string[]} scenarios - the currently available scenarios.
 * @param onProjectChange - the function to call when the user changes the active project.
 * @param onConversationChange - the function to call when the user changes the active conversation.
 * @param handleAudioRecording - the function to call when the user starts recording audio.
 * @param handleAudioRecordingError - the function to call when the user stops recording audio.
 * @returns {Element}
 */
const TopSelector = ({
                         projects,
                         activeProject,
                         sessions,
                         scenarios,
                         onProjectChange,
                         onConversationChange,
                         handleAudioRecording,
                         handleAudioRecordingError
                     }) => {
    // These will be filled dynamically as the conversation progress.
    // TODO: should use session.DisplayName
    const [conversations, setConversations] = useState(sessions.map((s) => {return s.SessionID}));
    const [selectedConversation, setSelectedConversation] = useState(0);
    const [selectedProject, setSelectedProject] = useState(0);

    const [isAdding, setIsAdding] = useState(false);
    const [newConversation, setNewConversation] = useState('');
    const [isBlocked, setIsBlocked] = useState(false);

    const inputRef = useRef(null);
    Logger.debug(`TopSelector: ${activeProject} from [${projects}] 
                  has [${sessions.map((s) => {return s.SessionID})}]. 
                  Scenarios: ${scenarios}`);

    useEffect(() => {
        Logger.debug("useEffect");
        if (isAdding) {
            inputRef.current.focus();
        }
        setConversations(sessions.map((s) => {return s.SessionID}));
        navigator.mediaDevices.getUserMedia({audio: true})
            .then(function (stream) {
                setIsBlocked(false);
            })
            .catch(function (err) {
                setIsBlocked(true);
            });
    }, [isAdding, sessions]);

    const handleProjectChange = (event) => {
        setSelectedProject(event.target.value);
        onProjectChange(projects[event.target.value]);
    };

    const handleScenarioChange = (event) => {
        let selectedScenario = scenarios[event.target.value - 1];
        // onScenarioChange(selectedScenario);
    };


    const handleConversationChange = (event) => {
        setSelectedConversation(event.target.value);
    }

    const handleNewConversationChange = (event) => {
        setNewConversation(event.target.value);
    }

    const addConversation = () => {
        setIsAdding(true);
    }

    const deleteConversation = () => {
        //setConversations(prevConversations => prevConversations.filter((conversation, index) => index + 1 !== selectedConversation));
        setSelectedConversation(1); // Select the first conversation after deletion
    }

    const handleKeyDown = (event) => {
        if (event.key === 'Enter') {
            // setConversations(prevConversations => [...prevConversations, newConversation]);
            // setSelectedConversation(conversations.length + 1);
            setIsAdding(false);
            setNewConversation('');
        } else if (event.key === 'Escape') {
            setIsAdding(false);
            setNewConversation('');
        }
    }

    return (
        <div className="row top-selector">
            <div className="col-md-3">
                <span className="bold-label">Active Project:&nbsp;</span>
                <select className='form-control'
                        value={selectedProject}
                        onChange={handleProjectChange}>
                    {projects.map((project, index) => (
                        <option key={index} value={index}>{project}</option>
                    ))}
                </select>
            </div>
            <div className="col-md-3">
                <span className="bold-label">Conversation:&nbsp;</span>
                {isAdding
                    ? <input ref={inputRef} className="form-control"
                             type="text"
                             onChange={handleNewConversationChange}
                             onKeyDown={handleKeyDown}
                             value={newConversation}/>
                    : <select className='form-control'
                              value={selectedConversation}
                              onChange={handleConversationChange}>
                        {conversations.map((session, index) => (
                            <option key={index} value={index}>{session}</option>
                        ))}
                    </select>
                }
            </div>
            <div className="col-md-2">
                <button
                    className="btn btn-outline-dark btn-sm conversation-btn"
                    onClick={addConversation}><FontAwesomeIcon icon={faPen}/></button>
                {conversations.length > 1 &&
                    <button className="btn btn-outline-dark btn-sm conversation-btn" onClick={deleteConversation}>
                        <FontAwesomeIcon icon={faTrash}/></button>}
            </div>
            <div className="col-md-3">
                <span className="bold-label">Scenarios:&nbsp;</span>
                <select className='form-control' onChange={handleScenarioChange}
                        disabled={!isAdding}>
                    {scenarios.map((option, index) => (
                        <option key={index} value={index + 1}>{option}</option>
                    ))}
                </select>
            </div>
            <div className="col-md-1">
                {isBlocked ? <span>Microphone access denied.</span> :
                    <AudioRecorder
                        handleAudioRecording={handleAudioRecording}
                        handleAudioRecordingError={handleAudioRecordingError}/>
                }
            </div>
        </div>
    )
};

export default TopSelector;
