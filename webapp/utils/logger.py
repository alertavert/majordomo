import logging
from datetime import datetime
import os

from constants import (
    LOG_DIR,
    LOG_FILE,
    APP_NAME,
)

def get_logger() -> logging.Logger:
    """Return the application logger."""
    return logging.getLogger(APP_NAME)

def setup_logger(log_level: int = logging.INFO) -> None:
    """Setup application logging.

        To obtain the logger, use `get_logger()`.
        :param log_level: The logging level to set.
    """
    if not os.path.exists(LOG_DIR):
        os.makedirs(LOG_DIR)
    log_file = os.path.join(
        LOG_DIR, f"{LOG_FILE}_{datetime.now().strftime('%Y%m%d')}.log"
    )
    logger = get_logger()
    logger.setLevel(level=log_level)
    fmt = logging.Formatter(fmt="%(asctime)s - %(name)s - %(levelname)s - %(message)s")
    handlers = [
        logging.FileHandler(log_file),
        logging.StreamHandler(),
    ]
    for handler in handlers:
        handler.setFormatter(fmt)
        logger.addHandler(handler)
