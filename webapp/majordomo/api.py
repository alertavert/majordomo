from dataclasses import dataclass
from typing import List, Any, Dict, Tuple

import requests

from constants import SERVER, PORT
from majordomo.model import Project, Assistant, Conversation
from utils import get_logger

BASE_URL = f"http://{SERVER}:{PORT}"
LOG = get_logger()


@dataclass
class ResponseError:
    """
    Dataclass for API response errors.
    """
    title: str
    message: str


def get_assistants() -> list[Assistant] | ResponseError:
    """Get the list of Assistants for the OpenAI Project."""
    url = f"{BASE_URL}/assistants"
    try:
        response = requests.get(url)
    except requests.RequestException as e:
        return ResponseError(title="Connection Error", message=str(e))
    if response.status_code != 200:
        return ResponseError(title="API Error", message=response.text)
    try:
        data = response.json()
        assistants = [Assistant.from_dict(item) for item in data]
        return assistants
    except Exception as e:
        return ResponseError(title="Decoding Error", message=str(e))


def get_projects() -> tuple[str, list[Project]] | ResponseError:
    """Get the list of Projects."""
    url = f"{BASE_URL}/projects"
    try:
        response = requests.get(url)
    except requests.RequestException as e:
        return ResponseError(title="Connection Error", message=str(e))
    if response.status_code != 200:
        return ResponseError(title="API Error", message=response.text)
    try:
        data = response.json()
        LOG.info("Data: %s", data)
        active_project = data.get("active_project", "")
        projects = [Project.from_dict(item) for item in data.get('projects', [])]
        return active_project, projects
    except Exception as e:
        return ResponseError(title="Decoding Error", message=str(e))


def get_conversations(project_id: str) -> list[Conversation] | ResponseError:
    """Get the list of conversations for a given project."""
    url = f"{BASE_URL}/projects/{project_id}/conversations"
    try:
        response = requests.get(url)
    except requests.RequestException as e:
        return ResponseError(title="Connection Error", message=str(e))
    if response.status_code != 200:
        return ResponseError(title="API Error", message=response.text)
    try:
        data = response.json()
        conversations = [Conversation.from_dict(item) for item in data.get("threads", [])]
        return conversations
    except Exception as e:
        return ResponseError(
            title="Decoding Error in Conversations",
            message=f"Cannot decode server response for conversations for "
                    f"project {project_id}: {e}"
            )


def ask_assistant(
        prompt: str,
        assistant: str,
        thread_id: str | None = None) -> dict[str, Any] | ResponseError:
    """
    Send a prompt to the assistant via POST '/prompt' and return the API response.

    :param prompt: The input prompt text.
    :param assistant: The name of the assistant to use.
    :param thread_id: The thread ID associated with the conversation, optional
    :return: Dictionary with API response or a ResponseError.
    """
    url = f"{BASE_URL}/prompt"
    payload = {"prompt": prompt, "assistant": assistant}
    if thread_id:
        payload["thread_id"] = thread_id
    try:
        response = requests.post(url, json=payload)
    except requests.RequestException as e:
        return ResponseError(title="Connection Error", message=str(e))
    if response.status_code != 200:
        return ResponseError(title="API Error", message=response.text)
    try:
        resp = response.json()
        if resp.get("status") != "success":
            return ResponseError(title="Assistant Error", message=resp.get("message"))
        return resp
    except Exception as e:
        return ResponseError(title="Decoding Error", message=str(e))
