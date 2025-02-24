from dataclasses import dataclass
from pathlib import Path
from typing import List, Union, Any, Dict, Tuple

import requests
import streamlit as st

from constants import SERVER, PORT
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


@dataclass
class Project:
    name: str
    description: str
    location: Path

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "Project":
        return Project(
            name=data.get("name"),
            description=data.get("description"),
            location=Path(data.get("location"))
        )


@dataclass
class Assistant:
    id: str
    name: str
    model: str
    instructions: str

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "Assistant":
        return Assistant(
            id=data.get("id"),
            name=data.get("name"),
            model=data.get("model"),
            instructions=data.get("instructions")
        )


@dataclass
class Conversation:
    id: str
    name: str
    assistant: str
    project: str

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "Conversation":
        return Conversation(
            id=data.get("id"),
            name=data.get("name"),
            assistant=data.get("assistant"),
            project=data.get("project")
        )


def get_assistants() -> Union[List[Assistant], ResponseError]:
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


def get_projects() -> Union[Tuple[str, List[Project]], ResponseError]:
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


def get_conversations(project_id: str) -> Union[List[Conversation], ResponseError]:
    """Get the list of conversations for a given project."""
    url = f"{BASE_URL}/projects/{project_id}/sessions"
    try:
        response = requests.get(url)
    except requests.RequestException as e:
        return ResponseError(title="Connection Error", message=str(e))
    if response.status_code != 200:
        return ResponseError(title="API Error", message=response.text)
    try:
        data = response.json()
        conversations = [Conversation.from_dict(item) for item in data.get("threads")]
        return conversations
    except Exception as e:
        return ResponseError(title="Decoding Error", message=str(e))


def get_conversations_cache() -> dict[str, List[Conversation]]:
    """Get the conversations cache."""
    if "conversations" not in st.session_state:
        st.session_state.conversations = {}
    return st.session_state.conversations


def new_conversation(
        assistant: Assistant,
        project: Project,
        name: str = "New Conversation",
) -> Conversation:
    """Create a new conversation."""
    conversations = get_conversations_cache()
    for convs in conversations.values():
        print(f"---  {[conv.name for conv in convs]}")
    conversation = Conversation(id="1", name=name, assistant=assistant.name, project=project.name)
    if project.name not in conversations:
        conversations[project.name] = []
    conversations[project.name].append(conversation)
    return conversation


def ask_assistant(prompt: str, thread_id: str, assistant: str) -> Union[Dict[str, Any], ResponseError]:
    """
    Send a prompt to the assistant via POST '/prompt' and return the API response.

    :param prompt: The input prompt text.
    :param thread_id: The thread ID associated with the conversation.
    :param assistant: The name of the assistant.
    :return: Dictionary with API response or a ResponseError.
    """
    url = f"{BASE_URL}/prompt"
    payload = {
        "prompt": prompt,
        "thread_id": thread_id,
        "assistant": assistant
    }
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
        return resp.get("message")
    except Exception as e:
        return ResponseError(title="Decoding Error", message=str(e))
