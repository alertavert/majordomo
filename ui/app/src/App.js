import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import './App.css';

import React, {useState, useEffect} from "react";
import Logo from './Components/Logo';
import ResizableTextArea from './Components/ResizableTextArea';
import Spinner from './Components/Spinner';
import ErrorBox from "./Components/ErrorBox";

const botUrl = 'http://localhost:5000/prompt';

function App() {
    const [responseValue, setResponseValue] = useState('Bot says...');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleFormSubmit = async (content) => {
        console.log('POSTing:', content)
        setLoading(true);
        setError(null);
        try {
            const response = await fetch(botUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({prompt: content}),
            });
            setLoading(false);
            if (response.ok) {
                const data = await response.json();
                console.log('Received:', data.message.length, 'characters');
                setResponseValue(data.message);
            } else {
                setError('Error (' + response.status + '): ' + response.message);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    };

    return (
        <div className="container">
            <Logo/>
            <Header/>
            <PromptBox onSubmit={handleFormSubmit}/>
            <span className="d-block p-2"/>
            <span> {loading ? <Spinner /> : <ResponseBox responseValue={responseValue}/>}</span>
            <div className='error-box'>
                { error != '' ? <ErrorBox message={error}/> : null }
            </div>
            <Footer/>
        </div>
    );
}

function Header() {
    return (
        <div className="page-header">
            <h1>Majordomo <span><code>Your helpful code bot</code></span></h1>
        </div>
    );
}

function PromptBox({onSubmit}) {
    const [textareaValue, setTextareaValue] = useState('');

    const handleSubmit = async () => {
        onSubmit(textareaValue);
    };

    const handleTextareaChange = (event) => {
        setTextareaValue(event.target.value);
    };

    return (
        <div className="container-fluid">
            <div className="jumbotron">
                <h6>Prompt</h6>
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
                <ResizableTextArea value={responseValue} />
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
