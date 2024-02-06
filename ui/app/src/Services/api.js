const MajordomoServerUrl = 'http://localhost:5005';
const ScenariosApiUrl = MajordomoServerUrl + '/scenarios';
const ProjectsApiUrl = MajordomoServerUrl + '/projects';

/**
 * Fetches scenarios from the backend and updates the component state accordingly.
 * @param {Function} setScenarios Function to update the scenarios state.
 * @param {Function} setSelectedScenario Function to set the selected scenario.
 * @param {Function} setError Function to update the error state.
 */
export const fetchScenarios = (setScenarios, setSelectedScenario, setError) => {
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
            setScenarios(data.scenarios);
            setSelectedScenario(data.scenarios[0]);
        })
        .catch(error => setError(`Could not retrieve Scenarios: ${error.message}`));
};

export const fetchProjects = (setActiveProject, setError) => {
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
            setActiveProject(data.active_project);
        })
        .catch(error => setError(`Could not retrieve projects: ${error.message}`));
};
