import {Logger} from "./logger";

const MajordomoServerUrl = 'http://localhost:9090';
const ScenariosApiUrl = MajordomoServerUrl + '/scenarios';
const ProjectsApiUrl = MajordomoServerUrl + '/projects';
const SessionsApiUrl = (project) => { return `${ProjectsApiUrl}/${project}/sessions`;}

/**
 * Fetches scenarios from the backend and updates the component state accordingly.
 * @param {Function} setScenarios Function to update the scenarios state.
 * @param {Function} setError Function to update the error state.
 */
export const fetchScenarios = (setScenarios, setError) => {
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
            }
            Logger.debug(`Fetch Sessions(${project}): ${data.map(session => session.SessionID)}`);
            setSessions(data);
        })
        .catch(error => setError(`Could not GET sessions for ${project}: ${error.message}`));
};
