const log = (message, level = 'INFO') => {
    // Format the current timestamp
    const timestamp = new Date().toISOString();
    // Log to the console with a uniform format
    console.log(`[${timestamp}] ${level}: ${message}`);
};

// Exportable log level methods
export const Logger = {
    info: (message) => log(message, 'INFO'),
    warn: (message) => log(message, 'WARN'),
    error: (message) => log(message, 'ERROR'),
    debug: (message) => log(message, 'DEBUG')
};