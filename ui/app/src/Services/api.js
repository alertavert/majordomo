import {Logger} from "./logger";

const MajordomoServerUrl = 'http://localhost:5005';
const ScenariosApiUrl = MajordomoServerUrl + '/scenarios';
const ProjectsApiUrl = MajordomoServerUrl + '/projects';
const SessionsApiUrl = (project) => { return `${ProjectsApiUrl}/${project}/sessions`;}
const SpeechApiUrl = MajordomoServerUrl + '/command';
const PromptApiUrl = MajordomoServerUrl + '/prompt';


/**
 * Fetches scenarios from the backend and updates the component state accordingly.
 * @param {Function} setScenarios Function to update the scenarios state.
 * @param {Function} setError Function to update the error state.
 */
export const fetchScenarios = (setScenarios, setError) => {
    Logger.debug('Fetching Scenarios from:', ScenariosApiUrl);
    fetch(ScenariosApiUrl)
        .then(response => {
            if (!response.ok) {
                throw new Error(`Failed to retrieve scenarios: ${response.statusText}`);
            }
            return response.json();
        })
        .then((data) => {
            if (data.scenarios.length === 0) {
                throw new Error('No scenarios found');
            }
            Logger.debug(`Fetch Scenarios: ${data.scenarios}`);
            setScenarios(data.scenarios);
        })
        .catch(error => setError(`Could not retrieve Scenarios: ${error.message}`));
};

export const fetchProjects = (setActiveProject, setProjects, setError) => {
    fetch(ProjectsApiUrl)
        .then(response => {
            if (!response.ok) {
                throw new Error(`Failed to retrieve projects: ${response.statusText}`);
            }
            return response.json();
        })
        .then((data) => {
            if (!data.active_project) {
                throw new Error('No active project found');
            }
            Logger.debug(``);
            Logger.debug(`Fetch Projects: ${data.projects.map(project => project.name)},
                          Active: ${data.active_project}`);
            setActiveProject(data.active_project);
            setProjects(data.projects.map(project => project.name));
        })
        .catch(error => setError(`Could not retrieve projects: ${error.message}`));
};

/**
 * Session type definition.
 * @typedef {Object} Session
 * @property {string} SessionID - The user ID.
 * @property {string} Scenario - The Scenario for the conversation (cannot be changed).
 * @property {string} Project - The Project for the conversation (cannot be changed).
 * @property {string} DisplayName - A user-friendly name for the conversation. (not used)
 * @property {string} Description - A user-friendly short description for the conversation. (not used)
 */
export const fetchSessionsForProjects = (project, setSessions, setError) => {
    if (project === '') {
        return;
    }
    let url = SessionsApiUrl(project);
    Logger.debug(`Fetching Sessions for ${project} from ${url}`);
    fetch(url)
        .then(response => {
            if (!response.ok) {
                throw new Error(`Failed to retrieve sessions: ${response.statusText}`);
            }
            return response.json();
        })
        .then((data) => {
            if (!data || data?.length === 0) {
                Logger.warn(`No sessions found for ${project}, creating dummy data`)
                data = [
                    {
                        SessionID: 12345,
                        DisplayName: "Build UI",
                        Scenario: "web_developer",
                        Project: project,
                        Description: "Build a UI for a web application",
                    },
                    {
                        SessionID: 6987,
                        DisplayName: "Manage Projects",
                        Scenario: "project_manager",
                        Project: project,
                        Description: "Manage projects for a web application",
                    },
                ];
            } else {
                Logger.debug(`Fetch Sessions(${project}): ${data.map(session => session.SessionID)}`);
            }
            setSessions(data);
        })
        .catch(error => setError(`Could not GET sessions for ${project}: ${error.message}`));
};

/**
 * Sends a text prompt to the backend for processing.
 *
 * @param {string} content The text prompt to send.
 * @param {string} assistant The AI assistant to use for generating the response.
 * @param {string} thread The conversation to append the current User prompt to.
 *
 * @returns {Promise<string>} The response from the backend.
 */
export const sendPrompt = async (content, assistant, thread) => {
    Logger.debug('Sending Query to Majordomo (' + content.length + ' chars)');
    try {
        const response = await fetch(PromptApiUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                prompt: content,
                scenario: assistant,
                session: thread,
            }),
        });
        const data = await response.json();
        Logger.debug(`Outcome from POSTing Prompt: ${data?.response}`);
        if (!response.ok || data?.response === 'error') {
            throw new Error(`${data.message ? data.message : response.statusText}`);
        }
        return data.message;
    } catch (error) {
        Logger.error(`Could not POST Prompt to Majordomo: ${JSON.stringify(error)}`);
        throw new Error(`The Bot wasn't quite there yet: ${error.message}`);
    }
}

/**
 * Sends an audio blob to the backend for processing.
 * @param {Blob} blob The audio blob to send.
 * @param {Function} setTextareaValue Function to update the text area value.
 * @param {Function} setError Function to update the error state.
 */
export const sendAudioBlob = async (blob, setTextareaValue, setError) => {
    Logger.debug('Sending Audio Blob to:', SpeechApiUrl);
    let formData = new FormData();
    formData.append('audio', blob, 'audio.mp3');
    try {
        const response = await fetch(SpeechApiUrl, {
            method: 'POST',
            body: formData,
        });
        if (response.ok) {
            const data = await response.json();
            Logger.debug('Received:', data.message.length, 'characters');
            setTextareaValue(data.message);
        } else {
            const errorData = await response.json(); // Parse the error response as JSON
            setError('Error (' + response.status + '): ' + errorData.message);
        }
    } catch (error) {
        setError('Cannot convert Audio: ' + error);
    }
}
