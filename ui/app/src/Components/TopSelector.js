import React, { useState, useRef, useEffect } from 'react';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faPen, faTrash } from '@fortawesome/free-solid-svg-icons'

import '../styles/TopSelector.css';

const TopSelector = ({ scenarios, onScenarioChange, onConversationChange }) => {

    // These will be filled dynamically as the conversation progress.
    const [conversations, setConversations] = useState(['Ask Majordomo']);
    const [selectedConversation, setSelectedConversation] = useState(1);
    const [isAdding, setIsAdding] = useState(false);
    const [newConversation, setNewConversation] = useState('');

    const inputRef = useRef(null);

    useEffect(() => {
        if (isAdding) {
            inputRef.current.focus();
        }
    }, [isAdding]);

    const handleScenarioChange = (event) => {
        let selectedScenario = scenarios[event.target.value-1];
        onScenarioChange(selectedScenario);
    };

    const handleConversationChange = (event) => {
        setSelectedConversation(event.target.value);
        onConversationChange(conversations[event.target.value-1]);
    }

    const handleNewConversationChange = (event) => {
        setNewConversation(event.target.value);
    }

    const addConversation = () => {
        setIsAdding(true);
    }

    const deleteConversation = () => {
        setConversations(prevConversations => prevConversations.filter((conversation, index) => index+1 !== selectedConversation));
        setSelectedConversation(1); // Select the first conversation after deletion
    }

    const handleKeyDown = (event) => {
        if(event.key === 'Enter') {
            setConversations(prevConversations => [...prevConversations, newConversation]);
            setSelectedConversation(conversations.length+1);
            setIsAdding(false);
            setNewConversation('');
        } else if(event.key === 'Escape') {
            setIsAdding(false);
            setNewConversation('');
        }
    }

    return (
        <div className="row">
            <div className="col-md-2">
                <span className="bold-label">Scenario:&nbsp;</span>
                <select className='form-control' onChange={handleScenarioChange}>
                    {scenarios.map((option, index) => (
                        <option value={index+1}>{option}</option>
                    ))}
                </select>
            </div>
            <div className="col-md-2">
                <span className="bold-label">Conversation:&nbsp;</span>
                {isAdding
                    ? <input ref={inputRef} className="form-control"
                             type="text"
                             onChange={handleNewConversationChange}
                             onKeyDown={handleKeyDown}
                             value={newConversation} />
                    : <select className='form-control'
                              value={selectedConversation}
                              onChange={handleConversationChange}>
                        {conversations.map((option, index) => (
                            <option value={index+1}>{option}</option>
                        ))}
                    </select>}
            </div>
            <div className="col-md-2 align-self-center">
                <button
                    className="btn btn-outline-dark btn-sm conversation-btn"
                    onClick={addConversation}><FontAwesomeIcon icon={faPen} /></button>
                {conversations.length > 1 && <button className="btn btn-outline-dark btn-sm conversation-btn" onClick={deleteConversation}><FontAwesomeIcon icon={faTrash} /></button>}
            </div>
        </div>
    )
}

export default TopSelector;
