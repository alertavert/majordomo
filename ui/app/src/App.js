import logo from './logo.svg';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import React, {useState} from "react";

const botUrl = 'http://localhost:5000/prompt';

function App() {

    const [responseValue, setResponseValue] = useState('');

    const handleFormSubmit = async (content) => {
        console.log('content', content)
        try {
            const response = await fetch(botUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ prompt: content }),
            });

            if (response.ok) {
                const data = await response.json();
                setResponseValue(data.message);
            } else {
                console.error('Error:', response.status);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    };

    return (
        <div className="container">
            <Header/>
            <PromptBox onSubmit={handleFormSubmit}/>
            <span className="d-block p-2"/>
            <ResponseBox responseValue={responseValue}/>
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
                <button className="btn btn-primary btn-sm"
                        onClick={handleSubmit}>Submit
                </button>
            </div>
        </div>
    );
}

function ResponseBox({responseValue}) {
    return (
        <div className="container-fluid">
            <div className="jumbotron">
                <h6>Response</h6>

                <textarea className="form-control"
                          style={{height: '100px', marginTop: '10px'}}
                          value={responseValue}
                          readOnly
                />
            </div>
        </div>
    );
}


function Footer() {
    return (
        <div className="footer fixed-bottom">
            <h6>AlertAvert &copy; 2023</h6>
        </div>
    )
}

export default App;
