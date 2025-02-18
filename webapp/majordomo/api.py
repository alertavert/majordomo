from dataclasses import dataclass
from pathlib import Path
from typing import List

@dataclass
class Project:
    id: str
    name: str
    description: str
    location: Path

@dataclass
class Assistant:
    id: str
    name: str
    model: str
    instructions: str

@dataclass
class Conversation:
    id: str
    name: str
    assistant: Assistant
    project: Project

def get_assistants(oai_project_id: str) -> List[Assistant]:
    """Get the list of Assistants for the OpenAI Project."""
    if oai_project_id:
        return [
            Assistant(id="1", name="PyDev", model="gpt-3.5-turbo", instructions="You are a Python developer."),
            Assistant(id="2", name="WebDev", model="gpt-4", instructions="You will provide help in building React apps."),
        ]
    return []

def get_projects() -> List[Project]:
    """Get the list of Projects."""
    return [
        Project(id="1", name="Project 1", description="Description 1", location=Path("~/dev/project_1")),
        Project(id="2", name="Project 2", description="Description 2", location=Path("/usr/local/project_2")),
        Project(id="3", name="Project 3", description="Description 3", location=Path("~/dev/python/project_3")),
    ]

import streamlit as st

def get_conversations_cache() -> dict[str, List[Conversation]]:
    """Get the conversations cache."""
    if "conversations" not in st.session_state:
        st.session_state.conversations = {}
    return st.session_state.conversations

def get_conversations(project_id: str) -> List[Conversation]:
    """Get the list of conversations for the OpenAI Project."""
    conversations = get_conversations_cache()
    if project_id not in conversations:
        conversations[project_id] = []
    return conversations[project_id]

def new_conversation(
            assistant: Assistant,
            project: Project,
            name: str = "New Conversation",
    ) -> Conversation:
        """Create a new conversation."""
        conversations = get_conversations_cache()
        for convs in conversations.values():
            print(f"---  {[conv.name for conv in convs]}")
        conversation = Conversation(id="1", name=name, assistant=assistant, project=project)
        if project.id not in conversations:
            conversations[project.id] = []
        conversations[project.id].append(conversation)
        return conversation